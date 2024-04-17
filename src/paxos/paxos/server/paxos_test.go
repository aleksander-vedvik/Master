package server

import "testing"

var testValues = []struct{ id, currentRnd, nextRnd uint32 }{
	{
		id:         0,
		currentRnd: 0,
		nextRnd:    3,
	},
	{
		id:         0,
		currentRnd: 1,
		nextRnd:    3,
	},
	{
		id:         0,
		currentRnd: 3,
		nextRnd:    6,
	},
	{
		id:         0,
		currentRnd: 4,
		nextRnd:    6,
	},
	{
		id:         0,
		currentRnd: 5,
		nextRnd:    6,
	},
	{
		id:         1,
		currentRnd: 0,
		nextRnd:    4,
	},
	{
		id:         1,
		currentRnd: 1,
		nextRnd:    4,
	},
	{
		id:         1,
		currentRnd: 3,
		nextRnd:    7,
	},
	{
		id:         1,
		currentRnd: 4,
		nextRnd:    7,
	},
	{
		id:         1,
		currentRnd: 5,
		nextRnd:    7,
	},
}

func TestSetNewRound(t *testing.T) {
	srvAddrs := map[int]string{
		0: "127.0.0.1:5000",
		1: "127.0.0.1:5001",
		2: "127.0.0.1:5002",
	}
	srv1 := NewPaxosServer(0, srvAddrs)
	srv2 := NewPaxosServer(1, srvAddrs)

	for _, testValue := range testValues {
		if testValue.id == 0 {
			srv1.rnd = testValue.currentRnd
			//srv1.setNewRnd()
			if testValue.nextRnd != srv1.rnd {
				t.Fatalf("wrong rnd for srv%v, \n\rgot: \t%v \n\twant: \t%v \n\tcurrentRnd: %v", testValue.id+1, srv1.rnd, testValue.nextRnd, testValue.currentRnd)
			}
		} else if testValue.id == 1 {
			srv2.rnd = testValue.currentRnd
			//srv2.setNewRnd()
			if testValue.nextRnd != srv2.rnd {
				t.Fatalf("wrong rnd for srv%v, \n\tgot: \t%v \n\twant: \t%v \n\tcurrentRnd: %v", testValue.id+1, srv2.rnd, testValue.nextRnd, testValue.currentRnd)
			}
		}
	}
}

func BenchmarkPaxos(b *testing.B) {
	srvAddrs := map[int]string{
		0: "127.0.0.1:5000",
		1: "127.0.0.1:5001",
		2: "127.0.0.1:5002",
	}
	srv1 := NewPaxosServer(1, srvAddrs, true)
	srv2 := NewPaxosServer(2, srvAddrs, true)
	srv3 := NewPaxosServer(3, srvAddrs, true)
	srv1.Start()
	srv2.Start()
	srv3.Start()

	// create 3 servers
	// create 5 servers
	// create 7 servers

	// run 1000 reqs with
	// 	- 1 client
	// 	- 5 clients
	// 	- 10 clients

	// run 10 000 reqs with
	// 	- 1 client
	// 	- 5 clients
	// 	- 10 clients

	// run 50 000 reqs with
	// 	- 1 client
	// 	- 5 clients
	// 	- 10 clients

	// measure average, mean, min, and max processing time
	// measure throughput: reqs/sec
	// measure successful and failed reqs
}
