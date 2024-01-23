# Structure

### Thesis - Proposition:

- Write about distributed systems and their applications
- Write about a few consensus protocols, such as Paxos, PBFT, and other state-of-the-art protocols.
  - Mention how they work
    - How nodes communicate
    - How consensus is achieved
  - How they are implemented and wheter all-to-all communcation would be beneficial
- Write about the implementation (most important)
- Evaluate the performance of a few protocols (such as Paxos and PBFT) using old and new gorums library
- Maybe all-to-all communication would be beneficial for distributed workloads such as un Apache Spark or training of AI models? Maybe implement streaming?
- Write about gRPC and protocol buffers in general. Mention other solutions and briefly explain similarities/differences.
- Compare with event bus? Look at .NET MassTransit?
- Comparte with implementations in other languages?
- Gorums take a bigger part of the implementation.
  - An all-to-all function needs to be idempotent?
  - The quorum spec needs to be configured in advance on all servers.
  - There is a separation of a node: Client and Server
    - Would it be possible to avoid this separation? E.g. create the definition of a node, which handles both client and server operations.

### Questions:

- A node going offline while creating the QuorumSpecification will cause the application to crash. Could this be automatically handled by gorums?

  - _Solution propsition:_ Do not return early from the loop in config_opts.go:32. Instead, replace the line with continue. Additionally, create a custom error type that collects all errors and return this to the caller. The client can then decide themselves whether to still call the quorum function (ignore the errors) or exit the application/do something else (handle the errors).

- Could **auto discovery** of nodes be added when creating the quorum configuration?

  - _Solution propsition:_ E.g. perform a scan of the local network or implement a default function on servers listening for broadcast packets using UDP.

### Gorums - Proposition

- Autoconfiguration?
  - E.g. use QuorumSpecification to define both manager and configuration.
- Autodiscovery?
  - See question in section above.
- Handle quorum directly, instead of leaving it to the user?
  - The user now needs to return false or true depending on whether a quorum has been achieved.
  - The gorums library could handle everything and only return when a quorum is achieved.
  - Create a config flag for this? E.g. gorums.AdvancedQuorumHandling()

### Multiparty - Notes

Paxos is used as example protocol.

- Accept messages should be sent to all replicas.
- Attach an ID (server creates this and puts it in the metadata) to the accept message.
- A multiparty function would be based on the initial quorum specification, meaning all servers are added to the quorum spec at the start and does not change.
- Gorums will keep track of handled messages (by ID) and only accept messages not seen before. This should prevent infinite message propagation.
  - Maybe gorums could punish spammers?

### TODO

- Start on the actual implementation:
  - Review solutions
  - Figure out how source generation work with protocol buffers
  - Define Github issues (break down the implementation into steps)
- Define goals and milestones for each period/week (gantt)
- Create Paxos and PBFT on current gorums implementation
- Find relevant papers about distributed systems, specifically:
  - All-to-all communication
  - Similar solutions in other code languages
