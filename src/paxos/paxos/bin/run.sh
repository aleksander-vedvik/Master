#! /bin/bash
set -e

./paxos --type 0 --id 0 &
./paxos --type 0 --id 1 &
./paxos --type 0 --id 2 &
./paxos --type 0 --id 3 &
./paxos --type 0 --id 4 &

./paxos --type=1 --reqs 500 --id 0 &
./paxos --type=1 --reqs 500 --id 1 &
./paxos --type=1 --reqs 500 --id 2 &
./paxos --type=1 --reqs 500 --id 3 &
./paxos --type=1 --reqs 500 --id 4 &
./paxos --type=1 --reqs 500 --id 5 &
./paxos --type=1 --reqs 500 --id 6 &
./paxos --type=1 --reqs 500 --id 7 &
./paxos --type=1 --reqs 500 --id 8 &
./paxos --type=1 --reqs 500 --id 9 &
./paxos --type=1 --reqs 500 --id 10 &
./paxos --type=1 --reqs 500 --id 11 &
./paxos --type=1 --reqs 500 --id 12 &
./paxos --type=1 --reqs 500 --id 13 &
./paxos --type=1 --reqs 500 --id 14 &
./paxos --type=1 --reqs 500 --id 15 &

echo "running servers, enter to stop"

read && killall paxos
