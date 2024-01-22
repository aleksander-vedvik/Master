package paxos

type Acceptor struct {
	N      int
	Val    string
	Leader int
}

func NewAcceptor() *Acceptor {
	return &Acceptor{}
}

func (a *Acceptor) handleAccept() {

}
