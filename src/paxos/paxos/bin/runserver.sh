#! /bin/bash
set -e

./paxos --type 0 --id 0 &
./paxos --type 0 --id 1 &
./paxos --type 0 --id 2 &
./paxos --type 0 --id 3 &
./paxos --type 0 --id 4 &


echo "running servers, enter to stop"

read && killall paxos
