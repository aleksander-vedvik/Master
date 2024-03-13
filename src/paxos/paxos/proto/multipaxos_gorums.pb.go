// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.7.0-devel
// 	protoc            v3.12.4
// source: multipaxos.proto

package proto

import (
	context "context"
	fmt "fmt"
	gorums "github.com/relab/gorums"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	encoding "google.golang.org/grpc/encoding"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	net "net"
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
	nodes []*Node
	qspec QuorumSpec
	srv   *clientServerImpl
}

// ConfigurationFromRaw returns a new Configuration from the given raw configuration and QuorumSpec.
//
// This function may for example be used to "clone" a configuration but install a different QuorumSpec:
//
//	cfg1, err := mgr.NewConfiguration(qspec1, opts...)
//	cfg2 := ConfigurationFromRaw(cfg1.RawConfig, qspec2)
func ConfigurationFromRaw(rawCfg gorums.RawConfiguration, qspec QuorumSpec) *Configuration {
	// return an error if the QuorumSpec interface is not empty and no implementation was provided.
	var test interface{} = struct{}{}
	if _, empty := test.(QuorumSpec); !empty && qspec == nil {
		panic("QuorumSpec may not be nil")
	}
	return &Configuration{
		RawConfiguration: rawCfg,
		qspec:            qspec,
	}
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Manager.
//
// NOTE: mutating the returned slice is not supported.
func (c *Configuration) Nodes() []*Node {
	if c.nodes == nil {
		c.nodes = make([]*Node, 0, c.Size())
		for _, n := range c.RawConfiguration {
			c.nodes = append(c.nodes, &Node{n})
		}
	}
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
}

// NewManager returns a new Manager for managing connection to nodes added
// to the manager. This function accepts manager options used to configure
// various aspects of the manager.
func NewManager(opts ...gorums.ManagerOption) (mgr *Manager) {
	mgr = &Manager{}
	mgr.RawManager = gorums.NewRawManager(opts...)
	return mgr
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
		return nil, fmt.Errorf("wrong number of options: %d", len(opts))
	}
	c = &Configuration{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case gorums.NodeListOption:
			c.RawConfiguration, err = gorums.NewRawConfiguration(m.RawManager, v)
			if err != nil {
				return nil, err
			}
		case net.Listener:
			err = c.RegisterClientServer(v)
			if err != nil {
				return nil, err
			}
			return c, nil
		case QuorumSpec:
			// Must be last since v may match QuorumSpec if it is interface{}
			c.qspec = v
		default:
			return nil, fmt.Errorf("unknown option type: %v", v)
		}
	}
	// return an error if the QuorumSpec interface is not empty and no implementation was provided.
	var test interface{} = struct{}{}
	if _, empty := test.(QuorumSpec); !empty && c.qspec == nil {
		return nil, fmt.Errorf("missing required QuorumSpec")
	}
	return c, nil
}

// Nodes returns a slice of available nodes on this manager.
// IDs are returned in the order they were added at creation of the manager.
func (m *Manager) Nodes() []*Node {
	gorumsNodes := m.RawManager.Nodes()
	nodes := make([]*Node, 0, len(gorumsNodes))
	for _, n := range gorumsNodes {
		nodes = append(nodes, &Node{n})
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
	View *Configuration
}

func NewServer() *Server {
	srv := &Server{
		Server: gorums.NewServer(),
	}
	b := &Broadcast{
		Broadcaster: gorums.NewBroadcaster(),
		sp:          gorums.NewSpBroadcastStruct(),
	}
	srv.RegisterBroadcastStruct(b, configureHandlers(b), configureMetadata(b))
	return srv
}

func (srv *Server) SetView(config *Configuration) {
	srv.View = config
	srv.RegisterConfig(config.RawConfiguration)
	srv.ListenForBroadcast()
}

type Broadcast struct {
	*gorums.Broadcaster
	sp       *gorums.SpBroadcast
	metadata gorums.BroadcastMetadata
}

func configureHandlers(b *Broadcast) func(bh gorums.BroadcastHandlerFunc, ch gorums.BroadcastReturnToClientHandlerFunc) {
	return func(bh gorums.BroadcastHandlerFunc, ch gorums.BroadcastReturnToClientHandlerFunc) {
		b.sp.BroadcastHandler = bh
		b.sp.ReturnToClientHandler = ch
	}
}

func configureMetadata(b *Broadcast) func(metadata gorums.BroadcastMetadata) {
	return func(metadata gorums.BroadcastMetadata) {
		b.metadata = metadata
	}
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

func (c *Configuration) RegisterClientServer(lis net.Listener, opts ...grpc.ServerOption) error {
	srvImpl := &clientServerImpl{
		grpcServer: grpc.NewServer(opts...),
	}
	srv, err := gorums.NewClientServer(lis)
	if err != nil {
		return err
	}
	srvImpl.grpcServer.RegisterService(&clientServer_ServiceDesc, srvImpl)
	go srvImpl.grpcServer.Serve(lis)
	srvImpl.ClientServer = srv
	c.srv = srvImpl
	return nil
}

func (b *Broadcast) Accept(req *AcceptMsg, opts ...gorums.BroadcastOption) {
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	b.sp.BroadcastHandler("proto.MultiPaxos.Accept", req, b.metadata, options)
}

func (b *Broadcast) Learn(req *LearnMsg, opts ...gorums.BroadcastOption) {
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	b.sp.BroadcastHandler("proto.MultiPaxos.Learn", req, b.metadata, options)
}

func _clientWrite(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Response)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(clientServer).clientWrite(ctx, in)
}

func (srv *clientServerImpl) clientWrite(ctx context.Context, resp *Response) (*Response, error) {
	err := srv.AddResponse(ctx, resp)
	return resp, err
}

func (c *Configuration) Write(ctx context.Context, in *Value) (resp *Response, err error) {
	if c.srv == nil {
		return nil, fmt.Errorf("a client server is not defined. Use configuration.RegisterClientServer() to define a client server")
	}
	if c.qspec == nil {
		return nil, fmt.Errorf("a qspec is not defined.")
	}
	doneChan, cd := c.srv.AddRequest(ctx, in, gorums.ConvertToType(c.qspec.WriteQF))
	c.RawConfiguration.Multicast(ctx, cd, gorums.WithNoSendWaiting())
	response, ok := <-doneChan
	if !ok {
		return nil, fmt.Errorf("done channel was closed before returning a value")
	}
	return response.(*Response), err
}

// clientServer is the client server API for the MultiPaxos Service
type clientServer interface {
	clientWrite(ctx context.Context, request *Response) (*Response, error)
}

var clientServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.ClientServer",
	HandlerType: (*clientServer)(nil),
	Methods: []grpc.MethodDesc{

		{
			MethodName: "ClientWrite",
			Handler:    _clientWrite,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "",
}

// Accept is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Accept(ctx context.Context, in *AcceptMsg, opts ...gorums.CallOption) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "proto.MultiPaxos.Accept",
	}

	c.RawConfiguration.Multicast(ctx, cd, opts...)
}

// Ping is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Ping(ctx context.Context, in *Heartbeat, opts ...gorums.CallOption) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "proto.MultiPaxos.Ping",
	}

	c.RawConfiguration.Multicast(ctx, cd, opts...)
}

// QuorumSpec is the interface of quorum functions for MultiPaxos.
type QuorumSpec interface {
	gorums.ConfigOption

	// WriteQF is the quorum function for the Write
	// broadcastcall call method. The in parameter is the request object
	// supplied to the Write method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *Value'.
	WriteQF(replies []*Response) (*Response, bool)

	// PrepareQF is the quorum function for the Prepare
	// quorum call method. The in parameter is the request object
	// supplied to the Prepare method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *PrepareMsg'.
	PrepareQF(in *PrepareMsg, replies map[uint32]*PromiseMsg) (*PromiseMsg, bool)
}

// Prepare is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Prepare(ctx context.Context, in *PrepareMsg) (resp *PromiseMsg, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "proto.MultiPaxos.Prepare",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*PromiseMsg, len(replies))
		for k, v := range replies {
			r[k] = v.(*PromiseMsg)
		}
		return c.qspec.PrepareQF(req.(*PrepareMsg), r)
	}

	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*PromiseMsg), err
}

// MultiPaxos is the server-side API for the MultiPaxos Service
type MultiPaxos interface {
	Write(ctx gorums.ServerCtx, request *Value, broadcast *Broadcast)
	Prepare(ctx gorums.ServerCtx, request *PrepareMsg) (response *PromiseMsg, err error)
	Accept(ctx gorums.ServerCtx, request *AcceptMsg, broadcast *Broadcast)
	Learn(ctx gorums.ServerCtx, request *LearnMsg, broadcast *Broadcast)
	Ping(ctx gorums.ServerCtx, request *Heartbeat)
}

func (srv *Server) Write(ctx gorums.ServerCtx, request *Value, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method Write not implemented"))
}
func (srv *Server) Prepare(ctx gorums.ServerCtx, request *PrepareMsg) (response *PromiseMsg, err error) {
	panic(status.Errorf(codes.Unimplemented, "method Prepare not implemented"))
}
func (srv *Server) Accept(ctx gorums.ServerCtx, request *AcceptMsg) {
	panic(status.Errorf(codes.Unimplemented, "method Accept not implemented"))
}
func (srv *Server) Learn(ctx gorums.ServerCtx, request *LearnMsg, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method Learn not implemented"))
}
func (srv *Server) Ping(ctx gorums.ServerCtx, request *Heartbeat) {
	panic(status.Errorf(codes.Unimplemented, "method Ping not implemented"))
}

func RegisterMultiPaxosServer(srv *Server, impl MultiPaxos) {
	srv.RegisterHandler("proto.MultiPaxos.Write", gorums.BroadcastHandler(impl.Write, srv.Server))
	srv.RegisterClientHandler("proto.MultiPaxos.Write", gorums.ServerClientRPC("proto.MultiPaxos.Write"))
	srv.RegisterHandler("proto.MultiPaxos.Prepare", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*PrepareMsg)
		defer ctx.Release()
		resp, err := impl.Prepare(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("proto.MultiPaxos.Accept", gorums.BroadcastHandler(impl.Accept, srv.Server))
	srv.RegisterHandler("proto.MultiPaxos.Learn", gorums.BroadcastHandler(impl.Learn, srv.Server))
	srv.RegisterHandler("proto.MultiPaxos.Ping", func(ctx gorums.ServerCtx, in *gorums.Message, _ chan<- *gorums.Message) {
		req := in.Message.(*Heartbeat)
		defer ctx.Release()
		impl.Ping(ctx, req)
	})
}

func (b *Broadcast) SendToClient(resp protoreflect.ProtoMessage, err error) {
	b.sp.ReturnToClientHandler(resp, err, b.metadata)
}

func (srv *Server) SendToClient(resp protoreflect.ProtoMessage, err error, broadcastID string) {
	srv.RetToClient(resp, err, broadcastID)
}

type internalPromiseMsg struct {
	nid   uint32
	reply *PromiseMsg
	err   error
}
