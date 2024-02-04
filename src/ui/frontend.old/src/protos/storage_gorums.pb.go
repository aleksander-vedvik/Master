// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.7.0-devel
// 	protoc            v3.12.4
// source: storage.proto

package __

import (
	context "context"
	fmt "fmt"
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

// QuorumSpec is the interface of quorum functions for QCStorage.
type QuorumSpec interface {
	gorums.ConfigOption

	// ReadQF is the quorum function for the Read
	// quorum call method. The in parameter is the request object
	// supplied to the Read method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *ReadRequest'.
	ReadQF(in *ReadRequest, replies map[uint32]*State) (*State, bool)

	// WriteQF is the quorum function for the Write
	// quorum call method. The in parameter is the request object
	// supplied to the Write method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *State'.
	WriteQF(in *State, replies map[uint32]*WriteResponse) (*WriteResponse, bool)
}

// Read is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Read(ctx context.Context, in *ReadRequest) (resp *State, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "protos.QCStorage.Read",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*State, len(replies))
		for k, v := range replies {
			r[k] = v.(*State)
		}
		return c.qspec.ReadQF(req.(*ReadRequest), r)
	}

	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*State), err
}

// Write is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) Write(ctx context.Context, in *State) (resp *WriteResponse, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "protos.QCStorage.Write",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*WriteResponse, len(replies))
		for k, v := range replies {
			r[k] = v.(*WriteResponse)
		}
		return c.qspec.WriteQF(req.(*State), r)
	}

	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*WriteResponse), err
}

// Status is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Status(ctx context.Context, in *StatusRequest) (resp *StatusResponse, err error) {
	cd := gorums.CallData{
		Message: in,
		Method:  "protos.QCStorage.Status",
	}

	res, err := n.RawNode.RPCCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*StatusResponse), err
}

// QCStorage is the server-side API for the QCStorage Service
type QCStorage interface {
	Read(ctx gorums.ServerCtx, request *ReadRequest) (response *State, err error)
	Write(ctx gorums.ServerCtx, request *State) (response *WriteResponse, err error)
	Status(ctx gorums.ServerCtx, request *StatusRequest) (response *StatusResponse, err error)
}

func RegisterQCStorageServer(srv *gorums.Server, impl QCStorage) {
	srv.RegisterHandler("protos.QCStorage.Read", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*ReadRequest)
		defer ctx.Release()
		resp, err := impl.Read(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("protos.QCStorage.Write", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*State)
		defer ctx.Release()
		resp, err := impl.Write(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("protos.QCStorage.Status", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*StatusRequest)
		defer ctx.Release()
		resp, err := impl.Status(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
}

type internalState struct {
	nid   uint32
	reply *State
	err   error
}

type internalWriteResponse struct {
	nid   uint32
	reply *WriteResponse
	err   error
}