// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.7.0-devel
// 	protoc            v3.12.4
// source: node.proto

package __

import (
	context "context"
	fmt "fmt"
	net "net"

	uuid "github.com/google/uuid"
	gorums "github.com/relab/gorums"
	grpc "google.golang.org/grpc"
	encoding "google.golang.org/grpc/encoding"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
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
	srv.RegisterBroadcastStruct(&Broadcast{gorums.NewBroadcastStruct()})
	return srv
}

func (srv *Server) RegisterConfiguration(ownAddr string, srvAddrs []string, opts ...gorums.ManagerOption) error {
	err := srv.RegisterConfig(ownAddr, srvAddrs, opts...)
	srv.ListenForBroadcast()
	return err
}

type Broadcast struct {
	*gorums.BroadcastStruct
}

func (b *Broadcast) Write(req *WriteRequest, serverAddresses ...string) {
	b.SetBroadcastValues("protos.PBFTNode.Write", req, serverAddresses...)
}

func (b *Broadcast) PrePrepare(req *PrePrepareRequest, serverAddresses ...string) {
	b.SetBroadcastValues("protos.PBFTNode.PrePrepare", req, serverAddresses...)
}

func (b *Broadcast) Prepare(req *PrepareRequest, serverAddresses ...string) {
	b.SetBroadcastValues("protos.PBFTNode.Prepare", req, serverAddresses...)
}

func (b *Broadcast) Commit(req *CommitRequest, serverAddresses ...string) {
	b.SetBroadcastValues("protos.PBFTNode.Commit", req, serverAddresses...)
}

// QuorumSpec is the interface of quorum functions for PBFTNode.
type QuorumSpec interface {
	gorums.ConfigOption

	// WriteQF is the quorum function for the Write
	// quorum call method. The in parameter is the request object
	// supplied to the Write method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *WriteRequest'.
	WriteQF(in *WriteRequest, replies map[uint32]*ClientResponse) (*ClientResponse, bool)
}

type QuroumType int
const (
	ByzantineQuorum QuroumType = iota
	MajorityQuorum
	All
)

// Write is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Write(ctx context.Context, in *WriteRequest, addr string, resultHandler func(resps []*ClientResponse), returnWhen QuroumType) (resp *ClientResponse, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "protos.PBFTNode.Write",

		BroadcastID: uuid.New().String(),
		Sender:      "client",
		OriginAddr: addr,
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*ClientResponse, len(replies))
		for k, v := range replies {
			r[k] = v.(*ClientResponse)
		}
		return c.qspec.WriteQF(req.(*WriteRequest), r)
	}
	numServers := 3
	tmpSrv, err := createTmpServer(addr, resultHandler, numServers)
	if err != nil {
		return nil, err
	}

	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	tmpSrv.handleClient(int(returnWhen))
	return res.(*ClientResponse), err
}

type tmpServer interface {
	client(context.Context, *ClientResponse) (any, error)
}

type tmpServerImpl struct {
	grpcServer *grpc.Server
	handler func(resps []*ClientResponse)
	resps []*ClientResponse
	respChan chan *ClientResponse
}

func (srv tmpServerImpl) client(ctx context.Context, resp *ClientResponse) (any, error) {
	srv.respChan <- resp
	return nil, nil
}

func (srv tmpServerImpl) handleClient(returnWhen int) {
	for resp := range srv.respChan {
		srv.resps = append(srv.resps, resp)
		if len(srv.resps) > returnWhen {
			break
		}
	}
	srv.handler(srv.resps)
	srv.grpcServer.GracefulStop()
}

func createTmpServer(addr string, handler func(resps []*ClientResponse), maxNumResponses int) (*tmpServerImpl, error) {
	var opts []grpc.ServerOption
	srv := tmpServerImpl{
		grpcServer: grpc.NewServer(opts...),
		respChan: make(chan *ClientResponse, maxNumResponses),
		resps: make([]*ClientResponse, maxNumResponses),
		handler: handler,
	}
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		return nil, err
	}
	registerTmpServer(srv.grpcServer, srv)
	go srv.grpcServer.Serve(lis)
	return &srv, nil
}

func registerTmpServer(s grpc.ServiceRegistrar, srv tmpServer) {
	s.RegisterService(&TmpServer_ServiceDesc, srv)
}

func _TmpServer_Client_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(tmpServer).client(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.TmpServer/Client",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(tmpServer).client(ctx, req.(*ClientResponse))
	}
	return interceptor(ctx, in, info, handler)
}

var TmpServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.TmpServer",
	HandlerType: (*tmpServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Client",
			Handler:    _TmpServer_Client_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "",
}

// PBFTNode is the server-side API for the PBFTNode Service
type PBFTNode interface {
	Write(ctx gorums.BroadcastCtx, request *WriteRequest, broadcast *Broadcast) (err error)
	PrePrepare(ctx gorums.BroadcastCtx, request *PrePrepareRequest, broadcast *Broadcast) (err error)
	Prepare(ctx gorums.BroadcastCtx, request *PrepareRequest, broadcast *Broadcast) (err error)
	Commit(ctx gorums.BroadcastCtx, request *CommitRequest, broadcast *Broadcast) (err error)
}

func RegisterPBFTNodeServer(srv *Server, impl PBFTNode) {
	srv.RegisterHandler("protos.PBFTNode.Write", gorums.BroadcastHandler(impl.Write, srv.Server))
	srv.RegisterHandler("protos.PBFTNode.PrePrepare", gorums.BroadcastHandler(impl.PrePrepare, srv.Server))
	srv.RegisterHandler("protos.PBFTNode.Prepare", gorums.BroadcastHandler(impl.Prepare, srv.Server))
	srv.RegisterHandler("protos.PBFTNode.Commit", gorums.BroadcastHandler(impl.Commit, srv.Server))
}

func (b *Broadcast) ReturnToClient(resp *ClientResponse, err error) {
	b.SetReturnToClient(resp, err)
}

func (srv *Server) ReturnToClient(resp *ClientResponse, err error, broadcastID string) {
	go srv.RetToClient(resp, err, broadcastID)
}

type internalClientResponse struct {
	nid   uint32
	reply *ClientResponse
	err   error
}
