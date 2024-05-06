# Problem statement

### What problem should Gorums solve?

- It is hard to implement reliable communication between servers. It is important to handle retries, faulty configurations, faulty nodes, network separation, reconfiguration, byzantine nodes (spamming), etc.
- Gorums should make this very simple. A protocol implementer should only focus on implementing his own protocol. E.g. if he implement Paxos, he should only need to create the initial configuration of Gorums (such as server addresses/ranges) and the rest of the communication between the servers should be handled accordingly.
- It can be easy to implement communication between nodes, but it is hard to optimize and make it very efficient.

### Problem

The advantage of running tasks in a single thread is to avoid thread-safety issues.
To DepFast, an extra benefit is to eliminate the possibility of
fail-slow faults caused by thread locks, a known suspect of
performance issues. In reality, running tasks concurrently with
a single thread can give high enough performance in most
cases, as shown in our evaluation (see §6.2 and §6.3).
Moving to multiple threads, to write thread-safe code the
users need to either shard the system into different threads
and regulate inter-thread communication, or write memory-
sharing code and use mutexes for mutual exclusion. DepFast
encourages the former as it minimizes the chances of perfor-
mance issues caused by waiting on mutexes. [p.7, DepFast]
