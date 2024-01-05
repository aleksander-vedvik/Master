#! /bin/bash
set -e

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "1" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "2" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "3" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "4" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "5" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "6" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "7" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "8" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "9" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "10" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "11" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "12" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "13" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "14" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "15" &

./paxosclient -addrs="localhost:50081,localhost:50082,localhost:50083,localhost:50084,localhost:50085" -clientRequest=500 -clientId "16" &

read && killall paxosclient
