package leaderdetector

// A MonLeaderDetector represents a Monarchical Eventual Leader Detector as
// described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.
type MonLeaderDetector struct {
	// DONE(student): Add needed fields
	nodeIds     []int
	leader      int
	suspected   map[int]bool
	subscribers []chan int
}

// NewMonLeaderDetector returns a new Monarchical Eventual Leader Detector
// given a list of node ids.
func NewMonLeaderDetector(nodeIDs []int) *MonLeaderDetector {
	// DONE(student): Add needed implementation
	m := &MonLeaderDetector{
		nodeIds:   nodeIDs,
		suspected: make(map[int]bool),
	}
	m.leader = m.maxRank()
	return m
}

// Leader returns the current leader. Leader will return UnknownID if all nodes
// are suspected.
func (m *MonLeaderDetector) Leader() int {
	// DONE(student): Implement
	return m.leader
}

// Suspect instructs the leader detector to consider the node with matching
// id as suspected. If the suspect indication result in a leader change
// the leader detector should publish this change to its subscribers.
func (m *MonLeaderDetector) Suspect(id int) {
	// DONE(student): Implement
	m.suspected[id] = true
	if id == m.leader {
		m.changeLeader()
	}
}

// Restore instructs the leader detector to consider the node with matching
// id as restored. If the restore indication result in a leader change
// the leader detector should publish this change to its subscribers.
func (m *MonLeaderDetector) Restore(id int) {
	// DONE(student): Implement
	s, ok := m.suspected[id]
	suspected := s && ok
	if suspected {
		m.suspected[id] = false
		delete(m.suspected, id) // Optional
	}
	if m.leader < id {
		m.changeLeader()
	}
}

// Subscribe returns a buffered channel which will be used by the leader
// detector to publish the id of the highest ranking node.
// The leader detector will publish UnknownID if all nodes become suspected.
// Subscribe will drop publications to slow subscribers.
// Note: Subscribe returns a unique channel to every subscriber;
// it is not meant to be shared.
func (m *MonLeaderDetector) Subscribe() <-chan int {
	// DONE(student): Implement
	subChan := make(chan int, len(m.nodeIds))
	m.subscribers = append(m.subscribers, subChan)
	return subChan
}

// DONE(student): Add other unexported functions or methods if needed.
func (m *MonLeaderDetector) maxRank() int {
	maxRank := UnknownID
	for _, nodeId := range m.nodeIds {
		s, ok := m.suspected[nodeId]
		suspected := s && ok
		if (nodeId > maxRank) && !suspected {
			maxRank = nodeId
		}
	}
	return maxRank
}

func (m *MonLeaderDetector) changeLeader() {
	m.leader = m.maxRank()
	for _, subChan := range m.subscribers {
		subChan <- m.leader
	}
}
