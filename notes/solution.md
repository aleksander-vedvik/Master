# Solution

### Preliminary

Gorums can be seen as an implementation of a Regular Register. All-to-all should potentially define the other broadcast schemes.

- m = message
- **The most important task is to make sure that gRPC function is called exactly once by each node. Otherwise, the system will incur infinite messages.**

### Assumptions

- All servers know about each other. The configuration is created in advance (both for broadcast and gossiping).
  - Change the implementation for how the manager and the configuration is created.
  - Should be a notion on node? The node acts as both server and client? E.g. node.Client and node.Server?
    - This way, the server can know about the configuration.
- The configuration has to be stored in Metadata if configuration is not known in advance.

### Types

- Basic Broadcast (Best effort broadcast: (BEB))
- Lazy Reliable broadcast
- Eager Reliable broadcast
  - BEB Broadcast first time m is received
- All-Ack Uniform Relibale Broadcast
- Majority-Ack Uniform Relibale Broadcast

### Propositions

- Best effort broadcast:
  - Gorums handles the broadcast in its entirety.
  - The gRPC function is run only once.
  - It will invoke the same gRPC function on all nodes.
  - It will ignore subsequent calls to the same gRPC function for the specific message.
- Uniform broadcast:
  - The user needs to decide a value after reaching a majority.
  - Gorums could handle the entire broadcast:
    - It will only return on the decided value.
    - The gRPC function needs to satisfy a few criterias:
      - A call to the gRPC function needs to take in one value (e.g. state)
      - It should not make sense to call the function on a single node? **Check this.**
- User defined broadcast. A user can define a function on how to handle the result:
  - The user should specify a criteria for when to finish the call. For example:
    - Ignore subsequent calls with the same message.
    - Collect all messages and save the newest/highest round number.
    - Collect all messages and return the majority (e.g. accept in Paxos).
    - Specify whether the gRPC function should be run first or after or both. Could also define another function handling the result after a specific criteria is met.

### Challenges

- A node must have only one configuration, from now called view. This must be created at the start.
  - To change the view, the nodes must instantiate a "configure view" protocol. This is not yet implemented in Gorums.
- This view needs to be part of the Gorums server.
- Optimialization:
  - Gorums will create a connection (stream) to each node in the configuration. This will be done for every node.
    - Is it possible to use the same connection for both Gorums client and server?
- A message gets a new ID for each request. How to distinguish messages in broadcast?
  - Use the same ID for broadcast messages?
    - How to distinguish chained methods? E.g. pBFT: PrePrepare -> Prepare -> Commit
  - How to check who the message was received from? Uses peer.GetAddr from context. IP changes constantly.
