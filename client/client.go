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

	dealWithUDP(localAddr.String(), string(data[:dataLen]))
	//dealWithTCP(localAddr.String(), string(data[:dataLen]))
}

func dealWithUDP(local, remote string) {
	parseUDPAddr := func(addr string) *net.UDPAddr {
		t := strings.Split(addr, ":")
		port, _ := strconv.Atoi(t[1])
		return &net.UDPAddr{
			IP:   net.ParseIP(strings.Split(addr, ":")[0]),
			Port: port,
		}
	}

	localAddr, remoteAddr := parseUDPAddr(local), parseUDPAddr(remote)

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

func dealWithTCP(local, remote string) {
	parseTCPAddr := func(addr string) *net.TCPAddr {
		t := strings.Split(addr, ":")
		port, _ := strconv.Atoi(t[1])
		return &net.TCPAddr{
			IP:   net.ParseIP(strings.Split(addr, ":")[0]),
			Port: port,
		}
	}

	localAddr, remoteAddr := parseTCPAddr(local), parseTCPAddr(remote)

	listener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, _ := listener.Accept()
			data := make([]byte, 1024)
			dataLen, _ := conn.Read(data)
			if dataLen != 0 {
				fmt.Printf("%s->%s\n", remoteAddr, data[:dataLen])
			}
		}
	}()

	conn, err := net.DialTCP("tcp", localAddr, remoteAddr)
	if err != nil {
		panic(err)
	}

	conn.Write(nil)

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
