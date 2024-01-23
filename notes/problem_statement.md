# Problem statement

### What problem should Gorums solve?

- It is hard to implement reliable communication between servers. It is important to handle retries, faulty configurations, faulty nodes, network separation, reconfiguration, byzantine nodes (spamming), etc.
- Gorums should make this very simple. A protocol implementer should only focus on implementing his own protocol. E.g. if he implement Paxos, he should only need to create the initial configuration of Gorums (such as server addresses/ranges) and the rest of the communication between the servers should be handled accordingly.
- It can be easy to implement communication between nodes, but it is hard to optimize and make it very efficient.
