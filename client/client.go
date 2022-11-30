package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	var port int
	var serverIp string
	var serverPort int
	flag.IntVar(&port, "port", 9981, "")
	flag.StringVar(&serverIp, "serverIp", "127.0.0.1", "")
	flag.IntVar(&serverPort, "serverPort", 9980, "")
	flag.Parse()

	localAddr := &net.UDPAddr{IP: net.IPv4zero, Port: port}

	serverConn, err := net.DialUDP("udp", localAddr, &net.UDPAddr{IP: net.ParseIP(serverIp), Port: serverPort})
	if err != nil {
		panic(err)
	}
	serverConn.Write(nil)

	data := make([]byte, 1024)
	dataLen, _, err := serverConn.ReadFromUDP(data)
	if err != nil {
		panic(err)
	}
	serverConn.Close()

	remoteAddr := parseAddr(string(data[:dataLen]))
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		panic(err)
	}

	conn.Write(nil)

	go func() {
		for {
			data := make([]byte, 1024)
			dataLen, _, _ := conn.ReadFromUDP(data)
			if dataLen != 0 {
				fmt.Printf("%s->%s\n", remoteAddr, data[:dataLen])
			}
		}
	}()

	for {
		var msg string
		fmt.Scanln(&msg)

		conn.Write([]byte(msg))

		if msg == "exit" {
			conn.Close()
			return
		}
	}
}

func parseAddr(addr string) *net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return &net.UDPAddr{
		IP:   net.ParseIP(strings.Split(addr, ":")[0]),
		Port: port,
	}
}
