// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.7.0-devel
// 	protoc            v3.12.4
// source: simple.proto

package __

import (
	context "context"
	fmt "fmt"
	empty "github.com/golang/protobuf/ptypes/empty"
	gorums "github.com/relab/gorums"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	encoding "google.golang.org/grpc/encoding"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	net "net"
	time "time"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = gorums.EnforceVersion(7 - gorums.MinVersion)
	// Verify that the gorums runtime is sufficiently up-to-date.
	_ = gorums.EnforceVersion(gorums.MaxVersion - 7)
)

// A Configuration represents a static set of nodes on which quorum remote
// procedure calls may be invoked.
type Configuration struct {
	gorums.RawConfiguration
	qspec     QuorumSpec
	srv       *clientServerImpl
	snowflake gorums.Snowflake
	nodes     []*Node
}

// ConfigurationFromRaw returns a new Configuration from the given raw configuration and QuorumSpec.
//
// This function may for example be used to "clone" a configuration but install a different QuorumSpec:
//
//	cfg1, err := mgr.NewConfiguration(qspec1, opts...)
//	cfg2 := ConfigurationFromRaw(cfg1.RawConfig, qspec2)
func ConfigurationFromRaw(rawCfg gorums.RawConfiguration, qspec QuorumSpec) (*Configuration, error) {
	// return an error if the QuorumSpec interface is not empty and no implementation was provided.
	var test interface{} = struct{}{}
	if _, empty := test.(QuorumSpec); !empty && qspec == nil {
		return nil, fmt.Errorf("config: missing required QuorumSpec")
	}
	newCfg := &Configuration{
		RawConfiguration: rawCfg,
		qspec:            qspec,
	}
	// initialize the nodes slice
	newCfg.nodes = make([]*Node, newCfg.Size())
	for i, n := range rawCfg {
		newCfg.nodes[i] = &Node{n}
	}
	return newCfg, nil
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Manager.
//
// NOTE: mutating the returned slice is not supported.
func (c *Configuration) Nodes() []*Node {
	return c.nodes
}

// And returns a NodeListOption that can be used to create a new configuration combining c and d.
func (c Configuration) And(d *Configuration) gorums.NodeListOption {
	return c.RawConfiguration.And(d.RawConfiguration)
}

// Except returns a NodeListOption that can be used to create a new configuration
// from c without the nodes in rm.
func (c Configuration) Except(rm *Configuration) gorums.NodeListOption {
	return c.RawConfiguration.Except(rm.RawConfiguration)
}

func init() {
	if encoding.GetCodec(gorums.ContentSubtype) == nil {
		encoding.RegisterCodec(gorums.NewCodec())
	}
}

// Manager maintains a connection pool of nodes on
// which quorum calls can be performed.
type Manager struct {
	*gorums.RawManager
	srv *clientServerImpl
}

// NewManager returns a new Manager for managing connection to nodes added
// to the manager. This function accepts manager options used to configure
// various aspects of the manager.
func NewManager(opts ...gorums.ManagerOption) *Manager {
	return &Manager{
		RawManager: gorums.NewRawManager(opts...),
	}
}

func (mgr *Manager) Close() {
	if mgr.RawManager != nil {
		mgr.RawManager.Close()
	}
	if mgr.srv != nil {
		mgr.srv.stop()
	}
}

func (mgr *Manager) AddClientServer(lis net.Listener, opts ...gorums.ServerOption) error {
	srv := gorums.NewClientServer(lis, opts...)
	srvImpl := &clientServerImpl{
		ClientServer: srv,
	}
	registerClientServerHandlers(srvImpl)
	go srvImpl.Serve(lis)
	mgr.srv = srvImpl
	return nil
}

// NewConfiguration returns a configuration based on the provided list of nodes (required)
// and an optional quorum specification. The QuorumSpec is necessary for call types that
// must process replies. For configurations only used for unicast or multicast call types,
// a QuorumSpec is not needed. The QuorumSpec interface is also a ConfigOption.
// Nodes can be supplied using WithNodeMap or WithNodeList, or WithNodeIDs.
// A new configuration can also be created from an existing configuration,
// using the And, WithNewNodes, Except, and WithoutNodes methods.
func (m *Manager) NewConfiguration(opts ...gorums.ConfigOption) (c *Configuration, err error) {
	if len(opts) < 1 || len(opts) > 2 {
		return nil, fmt.Errorf("config: wrong number of options: %d", len(opts))
	}
	c = &Configuration{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case gorums.NodeListOption:
			c.RawConfiguration, err = gorums.NewRawConfiguration(m.RawManager, v)
			if err != nil {
				return nil, err
			}
		case QuorumSpec:
			// Must be last since v may match QuorumSpec if it is interface{}
			c.qspec = v
		default:
			return nil, fmt.Errorf("config: unknown option type: %v", v)
		}
	}
	// register the client server if it exists.
	// used to collect responses in BroadcastCalls
	if m.srv != nil {
		c.srv = m.srv
	}
	c.snowflake = m.Snowflake()
	//var test interface{} = struct{}{}
	//if _, empty := test.(QuorumSpec); !empty && c.qspec == nil {
	//	return nil, fmt.Errorf("config: missing required QuorumSpec")
	//}
	// initialize the nodes slice
	c.nodes = make([]*Node, c.Size())
	for i, n := range c.RawConfiguration {
		c.nodes[i] = &Node{n}
	}
	return c, nil
}

// Nodes returns a slice of available nodes on this manager.
// IDs are returned in the order they were added at creation of the manager.
func (m *Manager) Nodes() []*Node {
	gorumsNodes := m.RawManager.Nodes()
	nodes := make([]*Node, len(gorumsNodes))
	for i, n := range gorumsNodes {
		nodes[i] = &Node{n}
	}
	return nodes
}

// Node encapsulates the state of a node on which a remote procedure call
// can be performed.
type Node struct {
	*gorums.RawNode
}

type Server struct {
	*gorums.Server
	broadcast *Broadcast
	View      *Configuration
}

func NewServer(opts ...gorums.ServerOption) *Server {
	srv := &Server{
		Server: gorums.NewServer(opts...),
	}
	b := &Broadcast{
		orchestrator: gorums.NewBroadcastOrchestrator(srv.Server),
	}
	srv.broadcast = b
	srv.RegisterBroadcaster(newBroadcaster)
	return srv
}

func newBroadcaster(m gorums.BroadcastMetadata, o *gorums.BroadcastOrchestrator, e gorums.EnqueueBroadcast) gorums.Broadcaster {
	return &Broadcast{
		orchestrator:     o,
		metadata:         m,
		srvAddrs:         make([]string, 0),
		enqueueBroadcast: e,
	}
}

func (srv *Server) SetView(config *Configuration) {
	srv.View = config
	srv.RegisterConfig(config.RawConfiguration)
}

type Broadcast struct {
	orchestrator     *gorums.BroadcastOrchestrator
	metadata         gorums.BroadcastMetadata
	srvAddrs         []string
	enqueueBroadcast gorums.EnqueueBroadcast
}

// Returns a readonly struct of the metadata used in the broadcast.
//
// Note: Some of the data are equal across the cluster, such as BroadcastID.
// Other fields are local, such as SenderAddr.
func (b *Broadcast) GetMetadata() gorums.BroadcastMetadata {
	return b.metadata
}

type clientServerImpl struct {
	*gorums.ClientServer
	grpcServer *grpc.Server
}

func (c *clientServerImpl) stop() {
	c.ClientServer.Stop()
	if c.grpcServer != nil {
		c.grpcServer.Stop()
	}
}

func (b *Broadcast) To(addrs ...string) *Broadcast {
	if len(addrs) <= 0 {
		return b
	}
	b.srvAddrs = append(b.srvAddrs, addrs...)
	return b
}

func (b *Broadcast) Forward(req protoreflect.ProtoMessage, addr string) error {
	if addr == "" {
		return fmt.Errorf("cannot forward to empty addr, got: %s", addr)
	}
	if !b.metadata.IsBroadcastClient {
		return fmt.Errorf("can only forward client requests")
	}
	go b.orchestrator.ForwardHandler(req, b.metadata.OriginMethod, b.metadata.BroadcastID, addr, b.metadata.OriginAddr)
	return nil
}

// Done signals the end of a broadcast request. It is necessary to call
// either Done() or SendToClient() to properly terminate a broadcast request
// and free up resources. Otherwise, it could cause poor performance.
func (b *Broadcast) Done() {
	b.orchestrator.DoneHandler(b.metadata.BroadcastID, b.enqueueBroadcast)
}

// SendToClient sends a message back to the calling client. It also terminates
// the broadcast request, meaning subsequent messages related to the broadcast
// request will be dropped. Either SendToClient() or Done() should be used at
// the end of a broadcast request in order to free up resources.
func (b *Broadcast) SendToClient(resp protoreflect.ProtoMessage, err error) {
	b.orchestrator.SendToClientHandler(b.metadata.BroadcastID, resp, err, b.enqueueBroadcast)
}

// Cancel is a non-destructive method call that will transmit a cancellation
// to all servers in the view. It will not stop the execution but will cause
// the given ServerCtx to be cancelled, making it possible to listen for
// cancellations.
//
// Could be used together with either SendToClient() or Done().
func (b *Broadcast) Cancel() {
	b.orchestrator.CancelHandler(b.metadata.BroadcastID, b.srvAddrs, b.enqueueBroadcast)
}

// SendToClient sends a message back to the calling client. It also terminates
// the broadcast request, meaning subsequent messages related to the broadcast
// request will be dropped. Either SendToClient() or Done() should be used at
// the end of a broadcast request in order to free up resources.
func (srv *Server) SendToClient(resp protoreflect.ProtoMessage, err error, broadcastID uint64) {
	srv.SendToClientHandler(resp, err, broadcastID, nil)
}

func (b *Broadcast) BroadcastIntermediate(req *BroadcastRequest, opts ...gorums.BroadcastOption) {
	if b.metadata.BroadcastID == 0 {
		panic("broadcastID cannot be empty. Use srv.BroadcastBroadcastIntermediate instead")
	}
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	options.ServerAddresses = append(options.ServerAddresses, b.srvAddrs...)
	b.orchestrator.BroadcastHandler("protosSimple.Simple.BroadcastIntermediate", req, b.metadata.BroadcastID, b.enqueueBroadcast, options)
}

func (b *Broadcast) Broadcast(req *BroadcastRequest, opts ...gorums.BroadcastOption) {
	if b.metadata.BroadcastID == 0 {
		panic("broadcastID cannot be empty. Use srv.BroadcastBroadcast instead")
	}
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	options.ServerAddresses = append(options.ServerAddresses, b.srvAddrs...)
	b.orchestrator.BroadcastHandler("protosSimple.Simple.Broadcast", req, b.metadata.BroadcastID, b.enqueueBroadcast, options)
}

func (srv *clientServerImpl) clientBroadcastCall1(ctx context.Context, resp *WriteResponse1, broadcastID uint64) (*WriteResponse1, error) {
	err := srv.AddResponse(ctx, resp, broadcastID)
	return resp, err
}

func (c *Configuration) BroadcastCall1(ctx context.Context, in *WriteRequest1, cancelOnTimeout ...bool) (resp *WriteResponse1, err error) {
	if c.srv == nil {
		return nil, fmt.Errorf("config: a client server is not defined. Use mgr.AddClientServer() to define a client server")
	}
	if c.qspec == nil {
		return nil, fmt.Errorf("a qspec is not defined")
	}
	var (
		timeout  time.Duration
		ok       bool
		response protoreflect.ProtoMessage
	)
	// use the same timeout as defined in the given context.
	// this is used for cancellation.
	deadline, ok := ctx.Deadline()
	if ok {
		timeout = deadline.Sub(time.Now())
	} else {
		timeout = 5 * time.Second
	}
	broadcastID := c.snowflake.NewBroadcastID()
	doneChan, cd := c.srv.AddRequest(broadcastID, ctx, in, gorums.ConvertToType(c.qspec.BroadcastCall1QF), "protosSimple.Simple.BroadcastCall1")
	c.RawConfiguration.BroadcastCall(ctx, cd, gorums.WithNoSendWaiting())
	select {
	case response, ok = <-doneChan:
	case <-ctx.Done():
		if len(cancelOnTimeout) > 0 && cancelOnTimeout[0] {
			go func() {
				bd := gorums.BroadcastCallData{
					Method:      gorums.Cancellation,
					BroadcastID: broadcastID,
				}
				cancelCtx, cancelCancel := context.WithTimeout(context.Background(), timeout)
				defer cancelCancel()
				c.RawConfiguration.BroadcastCall(cancelCtx, bd)
			}()
		}
		return nil, fmt.Errorf("context cancelled")
	}
	if !ok {
		return nil, fmt.Errorf("done channel was closed before returning a value")
	}
	resp, ok = response.(*WriteResponse1)
	if !ok {
		return nil, fmt.Errorf("wrong proto format")
	}
	return resp, nil
}

func (srv *clientServerImpl) clientBroadcastCall2(ctx context.Context, resp *WriteResponse2, broadcastID uint64) (*WriteResponse2, error) {
	err := srv.AddResponse(ctx, resp, broadcastID)
	return resp, err
}

func (c *Configuration) BroadcastCall2(ctx context.Context, in *WriteRequest2, cancelOnTimeout ...bool) (resp *WriteResponse2, err error) {
	if c.srv == nil {
		return nil, fmt.Errorf("config: a client server is not defined. Use mgr.AddClientServer() to define a client server")
	}
	if c.qspec == nil {
		return nil, fmt.Errorf("a qspec is not defined")
	}
	var (
		timeout  time.Duration
		ok       bool
		response protoreflect.ProtoMessage
	)
	// use the same timeout as defined in the given context.
	// this is used for cancellation.
	deadline, ok := ctx.Deadline()
	if ok {
		timeout = deadline.Sub(time.Now())
	} else {
		timeout = 5 * time.Second
	}
	broadcastID := c.snowflake.NewBroadcastID()
	doneChan, cd := c.srv.AddRequest(broadcastID, ctx, in, gorums.ConvertToType(c.qspec.BroadcastCall2QF), "protosSimple.Simple.BroadcastCall2")
	c.RawConfiguration.BroadcastCall(ctx, cd, gorums.WithNoSendWaiting())
	select {
	case response, ok = <-doneChan:
	case <-ctx.Done():
		if len(cancelOnTimeout) > 0 && cancelOnTimeout[0] {
			go func() {
				bd := gorums.BroadcastCallData{
					Method:      gorums.Cancellation,
					BroadcastID: broadcastID,
				}
				cancelCtx, cancelCancel := context.WithTimeout(context.Background(), timeout)
				defer cancelCancel()
				c.RawConfiguration.BroadcastCall(cancelCtx, bd)
			}()
		}
		return nil, fmt.Errorf("context cancelled")
	}
	if !ok {
		return nil, fmt.Errorf("done channel was closed before returning a value")
	}
	resp, ok = response.(*WriteResponse2)
	if !ok {
		return nil, fmt.Errorf("wrong proto format")
	}
	return resp, nil
}

func registerClientServerHandlers(srv *clientServerImpl) {

	srv.RegisterHandler("protosSimple.Simple.BroadcastCall1", gorums.ClientHandler(srv.clientBroadcastCall1))
	srv.RegisterHandler("protosSimple.Simple.BroadcastCall2", gorums.ClientHandler(srv.clientBroadcastCall2))
}

// QuorumSpec is the interface of quorum functions for Simple.
type QuorumSpec interface {
	gorums.ConfigOption

	// BroadcastCall1QF is the quorum function for the BroadcastCall1
	// broadcastcall call method. The in parameter is the request object
	// supplied to the BroadcastCall1 method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *WriteRequest1'.
	BroadcastCall1QF(in *WriteRequest1, replies []*WriteResponse1) (*WriteResponse1, bool)

	// BroadcastCall2QF is the quorum function for the BroadcastCall2
	// broadcastcall call method. The in parameter is the request object
	// supplied to the BroadcastCall2 method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *WriteRequest2'.
	BroadcastCall2QF(in *WriteRequest2, replies []*WriteResponse2) (*WriteResponse2, bool)

	// BenchmarkQF is the quorum function for the Benchmark
	// quorum call method. The in parameter is the request object
	// supplied to the Benchmark method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *empty.Empty'.
	BenchmarkQF(in *empty.Empty, replies map[uint32]*Result) (*Result, bool)
}

// Benchmark is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Benchmark(ctx context.Context, in *empty.Empty) (resp *Result, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "protosSimple.Simple.Benchmark",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Result, len(replies))
		for k, v := range replies {
			r[k] = v.(*Result)
		}
		return c.qspec.BenchmarkQF(req.(*empty.Empty), r)
	}

	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Result), err
}

// Simple is the server-side API for the Simple Service
type Simple interface {
	BroadcastCall1(ctx gorums.ServerCtx, request *WriteRequest1, broadcast *Broadcast)
	BroadcastCall2(ctx gorums.ServerCtx, request *WriteRequest2, broadcast *Broadcast)
	BroadcastIntermediate(ctx gorums.ServerCtx, request *BroadcastRequest, broadcast *Broadcast)
	Broadcast(ctx gorums.ServerCtx, request *BroadcastRequest, broadcast *Broadcast)
	Benchmark(ctx gorums.ServerCtx, request *empty.Empty) (response *Result, err error)
}

func (srv *Server) BroadcastCall1(ctx gorums.ServerCtx, request *WriteRequest1, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method BroadcastCall1 not implemented"))
}
func (srv *Server) BroadcastCall2(ctx gorums.ServerCtx, request *WriteRequest2, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method BroadcastCall2 not implemented"))
}
func (srv *Server) BroadcastIntermediate(ctx gorums.ServerCtx, request *BroadcastRequest, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method BroadcastIntermediate not implemented"))
}
func (srv *Server) Broadcast(ctx gorums.ServerCtx, request *BroadcastRequest, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method Broadcast not implemented"))
}
func (srv *Server) Benchmark(ctx gorums.ServerCtx, request *empty.Empty) (response *Result, err error) {
	panic(status.Errorf(codes.Unimplemented, "method Benchmark not implemented"))
}

func RegisterSimpleServer(srv *Server, impl Simple) {
	srv.RegisterHandler("protosSimple.Simple.BroadcastCall1", gorums.BroadcastHandler(impl.BroadcastCall1, srv.Server))
	srv.RegisterClientHandler("protosSimple.Simple.BroadcastCall1")
	srv.RegisterHandler("protosSimple.Simple.BroadcastCall2", gorums.BroadcastHandler(impl.BroadcastCall2, srv.Server))
	srv.RegisterClientHandler("protosSimple.Simple.BroadcastCall2")
	srv.RegisterHandler("protosSimple.Simple.BroadcastIntermediate", gorums.BroadcastHandler(impl.BroadcastIntermediate, srv.Server))
	srv.RegisterHandler("protosSimple.Simple.Broadcast", gorums.BroadcastHandler(impl.Broadcast, srv.Server))
	srv.RegisterHandler("protosSimple.Simple.Benchmark", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*empty.Empty)
		defer ctx.Release()
		resp, err := impl.Benchmark(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler(gorums.Cancellation, gorums.BroadcastHandler(gorums.CancelFunc, srv.Server))
}

func (srv *Server) BroadcastBroadcastIntermediate(req *BroadcastRequest, opts ...gorums.BroadcastOption) {
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	if options.RelatedToReq > 0 {
		srv.broadcast.orchestrator.BroadcastHandler("protosSimple.Simple.BroadcastIntermediate", req, options.RelatedToReq, nil, options)
	} else {
		srv.broadcast.orchestrator.ServerBroadcastHandler("protosSimple.Simple.BroadcastIntermediate", req, options)
	}
}

func (srv *Server) BroadcastBroadcast(req *BroadcastRequest, opts ...gorums.BroadcastOption) {
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	if options.RelatedToReq > 0 {
		srv.broadcast.orchestrator.BroadcastHandler("protosSimple.Simple.Broadcast", req, options.RelatedToReq, nil, options)
	} else {
		srv.broadcast.orchestrator.ServerBroadcastHandler("protosSimple.Simple.Broadcast", req, options)
	}
}

const (
	SimpleBroadcastCall1        string = "protosSimple.Simple.BroadcastCall1"
	SimpleBroadcastCall2        string = "protosSimple.Simple.BroadcastCall2"
	SimpleBroadcastIntermediate string = "protosSimple.Simple.BroadcastIntermediate"
	SimpleBroadcast             string = "protosSimple.Simple.Broadcast"
)

type internalResult struct {
	nid   uint32
	reply *Result
	err   error
}
