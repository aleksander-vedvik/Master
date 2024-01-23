package main

import (
	"fmt"
	"net"
)

func advertise(address string) {
	pc, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	addr, err := net.ResolveUDPAddr("udp4", "192.168.7.255:8829")
	if err != nil {
		panic(err)
	}

	_, err = pc.WriteTo([]byte("data to transmit"), addr)
	if err != nil {
		panic(err)
	}
}

func handleAdvertisement() {
	pc, err := net.ListenPacket("udp4", ":8829")
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buf := make([]byte, 1024)
	n, addr, err := pc.ReadFrom(buf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s sent this: %s\n", addr, buf[:n])
}
