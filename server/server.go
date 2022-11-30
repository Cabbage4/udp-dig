package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9980, "")
	flag.Parse()

	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		panic(err)
	}

tag:
	for {
		peers := make([]*net.UDPAddr, 2)
		data := make([]byte, 1024)
		for i := 0; i < 2; i++ {
			_, remoteAddr, err := listener.ReadFromUDP(data)
			if err != nil {
				fmt.Println(err)
				continue tag
			}
			peers[i] = remoteAddr
		}

		fmt.Printf("[%s]<--->[%s]\n", peers[0].String(), peers[1].String())
		listener.WriteToUDP([]byte(peers[1].String()), peers[0])
		listener.WriteToUDP([]byte(peers[0].String()), peers[1])
	}
}
