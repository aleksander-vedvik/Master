package defs

type SuspectRestorer interface {
	Suspect(int)
	Restore(int)
}

type Sr struct {
	suspect chan int
	restore chan int
}

func NewSr(suspectChan, restoreChan chan int) *Sr {
	return &Sr{
		suspect: suspectChan,
		restore: restoreChan,
	}
}

func (sr *Sr) Suspect(nodeId int) {
	sr.suspect <- nodeId
}

func (sr *Sr) Restore(nodeId int) {
	sr.restore <- nodeId
}
