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

type Broadcast struct {
	*gorums.BroadcastStruct
}

func (b *Broadcast) Broadcast(req *State) {
	b.SetBroadcastValues("protos.UniformBroadcast.Broadcast", req)
}

func (b *Broadcast) Deliver(req *State) {
	b.SetBroadcastValues("protos.UniformBroadcast.Deliver", req)
}

// QuorumSpec is the interface of quorum functions for UniformBroadcast.
type QuorumSpec interface {
	gorums.ConfigOption

	// BroadcastQF is the quorum function for the Broadcast
	// broadcast call method. The in parameter is the request object
	// supplied to the Broadcast method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *State'.
	BroadcastQF(in *State, replies map[uint32]*ClientResponse) (*ClientResponse, bool)

	// DeliverQF is the quorum function for the Deliver
	// broadcast call method. The in parameter is the request object
	// supplied to the Deliver method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *State'.
	DeliverQF(in *State, replies map[uint32]*Empty) (*Empty, bool)
}

// Broadcast is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Broadcast(ctx context.Context, in *State) (resp *ClientResponse, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "protos.UniformBroadcast.Broadcast",

		BroadcastID: uuid.New().String(),
		Sender:      "client",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*ClientResponse, len(replies))
		for k, v := range replies {
			r[k] = v.(*ClientResponse)
		}
		return c.qspec.BroadcastQF(req.(*State), r)
	}
	//ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("BroadcastID", "broadcast"))
	//ctx = context.WithValue(ctx, "testKey", "testVal")
	//ctxMd, _ := metadata.FromOutgoingContext(ctx)
	//log.Println("Init CTX:", ctxMd)
	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*ClientResponse), err
}

// Deliver is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Deliver(ctx context.Context, in *State) (resp *Empty, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "protos.UniformBroadcast.Deliver",

		BroadcastID: uuid.New().String(),
		Sender:      "client",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*Empty, len(replies))
		for k, v := range replies {
			r[k] = v.(*Empty)
		}
		return c.qspec.DeliverQF(req.(*State), r)
	}

	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*Empty), err
}

// UniformBroadcast is the server-side API for the UniformBroadcast Service
type UniformBroadcast interface {
	Broadcast(ctx gorums.ServerCtx, request *State, broadcast *Broadcast) (err error)
	Deliver(ctx gorums.ServerCtx, request *State, broadcast *Broadcast) (err error)
}

func RegisterUniformBroadcastServer(srv *Server, impl UniformBroadcast) {
	srv.RegisterHandler("protos.UniformBroadcast.Broadcast", gorums.BroadcastHandler(impl.Broadcast, srv.Server))
	srv.RegisterHandler("protos.UniformBroadcast.Deliver", gorums.BroadcastHandler(impl.Deliver, srv.Server))
}

func (srv *Server) RegisterConfiguration(srvAddrs []string, opts ...gorums.ManagerOption) error {
	err := srv.RegisterConfig(srvAddrs, opts...)
	srv.ListenForBroadcast()
	return err
}

func (b *Broadcast) ReturnToClient(resp *ClientResponse, err error) {
	b.SetReturnToClient(resp, err)
}

func (srv *Server) ReturnToClient(resp *ClientResponse, err error, broadcastID string) {
	go srv.RetToClient(resp, err, broadcastID)
}
