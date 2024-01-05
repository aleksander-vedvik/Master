#! /bin/bash
set -e

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "1" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "2" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "3" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "4" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "5" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "6" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "7" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "8" &

#./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "9" &

#./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "10" &

#./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "11" &

#./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=1000 -clientId "12" &

read && killall paxosclient
