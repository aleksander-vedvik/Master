// Code generated by protoc-gen-gorums. DO NOT EDIT.
// Source files can be found in: ./cmd/protoc-gen-gorums/dev

package gengorums

// pkgIdentMap maps from package name to one of the package's identifiers.
// These identifiers are used by the Gorums protoc plugin to generate import statements.
var pkgIdentMap = map[string]string{"fmt": "Errorf", "github.com/relab/gorums": "BroadcastMetadata", "google.golang.org/grpc": "Server", "google.golang.org/grpc/encoding": "GetCodec", "google.golang.org/protobuf/reflect/protoreflect": "ProtoMessage", "net": "Addr"}

// reservedIdents holds the set of Gorums reserved identifiers.
// These identifiers cannot be used to define message types in a proto file.
var reservedIdents = []string{"Broadcast", "Configuration", "Manager", "Node", "QuorumSpec", "Server"}

var staticCode = `// A Configuration represents a static set of nodes on which quorum remote
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

// AddClientServer starts a lightweight client-side server. This server only accepts responses
// to broadcast requests sent by the client.
//
// It is important to provide the listenAddr because this will be used to advertise the IP the
// servers should reply back to.
func (mgr *Manager) AddClientServer(lis net.Listener, clientAddr net.Addr, opts ...gorums.ServerOption) error {
	options := []gorums.ServerOption{gorums.WithListenAddr(clientAddr)}
	options = append(options, opts...)
	srv := gorums.NewClientServer(lis, options...)
	srvImpl := &clientServerImpl{
		ClientServer: srv,
	}
	registerClientServerHandlers(srvImpl)
	go func() {
		_ = srvImpl.Serve(lis)
	}()
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
func (b *Broadcast) SendToClient(resp protoreflect.ProtoMessage, err error) error {
	return b.orchestrator.SendToClientHandler(b.metadata.BroadcastID, resp, err, b.enqueueBroadcast)
}

// Cancel is a non-destructive method call that will transmit a cancellation
// to all servers in the view. It will not stop the execution but will cause
// the given ServerCtx to be cancelled, making it possible to listen for
// cancellations.
//
// Could be used together with either SendToClient() or Done().
func (b *Broadcast) Cancel() error {
	return b.orchestrator.CancelHandler(b.metadata.BroadcastID, b.srvAddrs, b.enqueueBroadcast)
}

// SendToClient sends a message back to the calling client. It also terminates
// the broadcast request, meaning subsequent messages related to the broadcast
// request will be dropped. Either SendToClient() or Done() should be used at
// the end of a broadcast request in order to free up resources.
func (srv *Server) SendToClient(resp protoreflect.ProtoMessage, err error, broadcastID uint64) error {
	return srv.SendToClientHandler(resp, err, broadcastID, nil)
}

`
