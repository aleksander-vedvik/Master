# Meeting

## What have I done?

- Revisited abstractions in the Distributed programming book
- Created a simple storage server without gorums (trying to explore what challenge we want to solve)
- Dived into the gorums code
- Looked at some other frameworks in Golang
- Looked at Github issues
- Tried to come up with different solutions:
  - Extend and implement several types of broadcast abstractions?

## Problems

- A configuration will not be created if a node in the nodesList is offline.
  - Check storage_gorums.pb.go:111
    - gorums.NewRawConfiguration() will return an error if dial fails. This will cause the configuration being returned to be nil (subsequent line in storage_gorums.pb.go).
      - See config_opts.go:64 and mgr.go:121.
- How to generate the source code? Would like a brief introduction on how it is done in gorums.
- I need a configuration on the server side:
  - It would be nice to specify this configuration at the start of the server. Advantage: Gossiping
  - Could include this in the Metadata, but would maybe be unefficient for large configurations?
  - Implement Views?

## What will I do until next time?

- Implement a solution in gorums
- Research: Find and read papers
- Write: Start on the introduction
