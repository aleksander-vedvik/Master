// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.7.0-devel
// 	protoc            v3.12.4
// source: storage.proto

package __

import (
	context "context"
	fmt "fmt"
	uuid "github.com/google/uuid"
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

// NewBroadcastConfiguration returns a configuration based on the provided list of nodes (required)
// and an optional quorum specification. The QuorumSpec is necessary for call types that
// must process replies. For configurations only used for unicast or multicast call types,
// a QuorumSpec is not needed. The QuorumSpec interface is also a ConfigOption.
// Nodes can be supplied using WithNodeMap or WithNodeList, or WithNodeIDs.
// A new configuration can also be created from an existing configuration,
// using the And, WithNewNodes, Except, and WithoutNodes methods.
func (m *Manager) NewBroadcastConfiguration(nodeOpt gorums.NodeListOption, qSpec QuorumSpec, lis net.Listener) (c *Configuration, err error) {
	c, err = m.NewConfiguration(nodeOpt, qSpec)
	if err != nil {
		return nil, err
	}
	c.RegisterClientServer(lis)
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
}

func NewServer() *Server {
	srv := &Server{
		gorums.NewServer(),
	}
	b := &Broadcast{
		Broadcaster: gorums.NewBroadcaster(),
		sp:              gorums.NewSpBroadcastStruct(),
	}
	srv.RegisterBroadcastStruct(b, configureHandlers(b), configureMetadata(b))
	return srv
}

func (srv *Server) SetView(srvAddrs []string, opts ...gorums.ManagerOption) error {
	err := srv.RegisterView(srvAddrs, opts...)
	srv.ListenForBroadcast()
	return err
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

func (b *Broadcast) SaveStudents(req *States, opts ...gorums.BroadcastOption) {
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	b.sp.BroadcastHandler("protos.UniformBroadcast.SaveStudents", req, b.metadata, options)
}

func (b *Broadcast) Broadcast(req *State, opts ...gorums.BroadcastOption) {
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	b.sp.BroadcastHandler("protos.UniformBroadcast.Broadcast", req, b.metadata, options)
}

func (b *Broadcast) Deliver(req *State, opts ...gorums.BroadcastOption) {
	options := gorums.NewBroadcastOptions()
	for _, opt := range opts {
		opt(&options)
	}
	b.sp.BroadcastHandler("protos.UniformBroadcast.Deliver", req, b.metadata, options)
}

func _clientSaveStudent(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(clientServer).clientSaveStudent(ctx, in)
}

func (srv *clientServerImpl) clientSaveStudent(ctx context.Context, resp *ClientResponse) (*ClientResponse, error) {
	err := srv.AddResponse(ctx, resp)
	return resp, err
}

func (c *Configuration) SaveStudent(ctx context.Context, in *State) (resp *ClientResponse, err error) {
	if c.srv == nil {
		return nil, fmt.Errorf("a client server is not defined. Use configuration.RegisterClientServer() to define a client server")
	}
	if c.qspec == nil {
		return nil, fmt.Errorf("a qspec is not defined.")
	}
	doneChan, cd := c.srv.AddRequest(ctx, in, gorums.ConvertToType(c.qspec.SaveStudentQF))
	c.RawConfiguration.Multicast(ctx, cd, gorums.WithNoSendWaiting())
	response, ok := <-doneChan
	if !ok {
		return nil, fmt.Errorf("done channel was closed before returning a value")
	}
	return response.(*ClientResponse), err
}

func _clientSaveStudents(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(clientServer).clientSaveStudents(ctx, in)
}

func (srv *clientServerImpl) clientSaveStudents(ctx context.Context, resp *ClientResponse) (*ClientResponse, error) {
	err := srv.AddResponse(ctx, resp)
	return resp, err
}

func (c *Configuration) SaveStudents(ctx context.Context, in *States) (resp *ClientResponse, err error) {
	if c.srv == nil {
		return nil, fmt.Errorf("a client server is not defined. Use configuration.RegisterClientServer() to define a client server")
	}
	if c.qspec == nil {
		return nil, fmt.Errorf("a qspec is not defined.")
	}
	doneChan, cd := c.srv.AddRequest(ctx, in, gorums.ConvertToType(c.qspec.SaveStudentsQF))
	c.RawConfiguration.Multicast(ctx, cd, gorums.WithNoSendWaiting())
	response, ok := <-doneChan
	if !ok {
		return nil, fmt.Errorf("done channel was closed before returning a value")
	}
	return response.(*ClientResponse), err
}

// clientServer is the client server API for the UniformBroadcast Service
type clientServer interface {
	clientSaveStudent(ctx context.Context, request *ClientResponse) (*ClientResponse, error)
	clientSaveStudents(ctx context.Context, request *ClientResponse) (*ClientResponse, error)
}

var clientServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.ClientServer",
	HandlerType: (*clientServer)(nil),
	Methods: []grpc.MethodDesc{

		{
			MethodName: "ClientSaveStudent",
			Handler:    _clientSaveStudent,
		},
		{
			MethodName: "ClientSaveStudents",
			Handler:    _clientSaveStudents,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "",
}

// QuorumSpec is the interface of quorum functions for UniformBroadcast.
type QuorumSpec interface {
	gorums.ConfigOption

	// SaveStudentQF is the quorum function for the SaveStudent
	// broadcastcall call method. The in parameter is the request object
	// supplied to the SaveStudent method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *State'.
	SaveStudentQF(replies []*ClientResponse) (*ClientResponse, bool)

	// SaveStudentsQF is the quorum function for the SaveStudents
	// broadcastcall call method. The in parameter is the request object
	// supplied to the SaveStudents method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *States'.
	SaveStudentsQF(replies []*ClientResponse) (*ClientResponse, bool)

	// BroadcastQF is the quorum function for the Broadcast
	// broadcast call method. The in parameter is the request object
	// supplied to the Broadcast method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *State'.
	BroadcastQF(in *State, replies map[uint32]*View) (*View, bool)
}

// Broadcast is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Broadcast(ctx context.Context, in *State) (resp *View, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "protos.UniformBroadcast.Broadcast",

		BroadcastID: uuid.New().String(),
		Sender:      gorums.BroadcastClient,
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*View, len(replies))
		for k, v := range replies {
			r[k] = v.(*View)
		}
		return c.qspec.BroadcastQF(req.(*State), r)
	}

	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*View), err
}

// UniformBroadcast is the server-side API for the UniformBroadcast Service
type UniformBroadcast interface {
	SaveStudent(ctx gorums.ServerCtx, request *State, broadcast *Broadcast)
	SaveStudents(ctx gorums.ServerCtx, request *States, broadcast *Broadcast)
	Broadcast(ctx gorums.ServerCtx, request *State, broadcast *Broadcast)
	Deliver(ctx gorums.ServerCtx, request *State, broadcast *Broadcast)
}

func (srv *Server) SaveStudent(ctx gorums.ServerCtx, request *State, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method SaveStudent not implemented"))
}
func (srv *Server) SaveStudents(ctx gorums.ServerCtx, request *States, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method SaveStudents not implemented"))
}
func (srv *Server) Broadcast(ctx gorums.ServerCtx, request *State, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method Broadcast not implemented"))
}
func (srv *Server) Deliver(ctx gorums.ServerCtx, request *State, broadcast *Broadcast) {
	panic(status.Errorf(codes.Unimplemented, "method Deliver not implemented"))
}

func RegisterUniformBroadcastServer(srv *Server, impl UniformBroadcast) {
	srv.RegisterHandler("protos.UniformBroadcast.SaveStudent", gorums.BroadcastHandler(impl.SaveStudent, srv.Server))
	srv.RegisterClientHandler("protos.UniformBroadcast.SaveStudent", gorums.ServerClientRPC("protos.UniformBroadcast.SaveStudent"))
	srv.RegisterHandler("protos.UniformBroadcast.SaveStudents", gorums.BroadcastHandler(impl.SaveStudents, srv.Server))
	srv.RegisterClientHandler("protos.UniformBroadcast.SaveStudents", gorums.ServerClientRPC("protos.UniformBroadcast.SaveStudents"))
	srv.RegisterHandler("protos.UniformBroadcast.Broadcast", gorums.BroadcastHandler(impl.Broadcast, srv.Server))
	srv.RegisterHandler("protos.UniformBroadcast.Deliver", gorums.BroadcastHandler(impl.Deliver, srv.Server))
}

func (b *Broadcast) SendToClient(resp protoreflect.ProtoMessage, err error) {
	b.sp.ReturnToClientHandler(resp, err, b.metadata)
}

func (srv *Server) SendToClient(resp protoreflect.ProtoMessage, err error, broadcastID string) {
	srv.RetToClient(resp, err, broadcastID)
}

type internalView struct {
	nid   uint32
	reply *View
	err   error
}
