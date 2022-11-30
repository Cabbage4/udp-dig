package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"udp-dig/common"
)

const hintInfo = `========== input command number ==========
===        [0]: link server            ===
===        [1]: show all client        ===
===        [2]: connect client         ===
==========================================`

var (
	clientAddr         *net.UDPAddr
	serverAddr         *net.UDPAddr
	clientToServerConn *net.UDPConn
)

func main() {
	var port int
	var serverIp string
	var serverPort int
	flag.IntVar(&port, "port", -1, "")
	flag.StringVar(&serverIp, "serverIp", "127.0,0.1", "")
	flag.IntVar(&serverPort, "serverPort", 3356, "")
	flag.Parse()

	clientAddr = &net.UDPAddr{IP: net.IPv4zero, Port: port}
	serverAddr = &net.UDPAddr{IP: net.ParseIP(serverIp), Port: serverPort}

	fmt.Println(hintInfo)
	for {
		var command int
		fmt.Scan(&command)

		if command < 0 || command > 2 {
			continue
		}

		if command == 0 {
			conn, err := net.DialUDP("udp", clientAddr, serverAddr)
			if err != nil {
				panic(err)
			}

			if _, err = conn.Write([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.ConnectReq}))); err != nil {
				panic(err)
			}

			clientToServerConn = conn

			go func() {
				for {
					msgInfo, err := common.GetMsgInfoFromConn(clientToServerConn)
					if err != nil {
						return
					}

					if msgInfo.Type == common.GetClientListRsp {
						fmt.Printf("->%s\n", strings.Join(strings.Split(msgInfo.Data, "|"), " "))
					} else if msgInfo.Type == common.ConnectClientRsp {
						// todo
					} else if msgInfo.Type == common.ConnectClientReq {
						clientToServerConn.Close()

						conn, err := net.DialUDP("udp", clientAddr, common.ParseAddr(msgInfo.Data))
						if err != nil {
							panic(err)
						}

						if _, err = conn.Write([]byte(common.ConvertMsgInfo(
							&common.MsgInfo{Type: common.IgnoreInfo}))); err != nil {
							panic(err)
						}

						go msgFunc(conn)

						if _, err = conn.Write([]byte(common.ConvertMsgInfo(
							&common.MsgInfo{Type: common.DataInfo, Data: "1"}))); err != nil {
							panic(err)
						}
					}
				}
			}()
		} else if command == 1 {
			if _, err := clientToServerConn.Write([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.GetClientListReq}))); err != nil {
				panic(err)
			}
		} else if command == 2 {
			var dstAddrStr string
			fmt.Scan(&dstAddrStr)

			if _, err := clientToServerConn.Write([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.ConnectClientReq, Data: dstAddrStr}))); err != nil {
				panic(err)
			}
			clientToServerConn.Close()

			conn, err := net.DialUDP("udp", clientAddr, common.ParseAddr(dstAddrStr))
			if err != nil {
				panic(err)
			}

			if _, err = conn.Write([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.IgnoreInfo}))); err != nil {
				panic(err)
			}

			go msgFunc(conn)

			if _, err = conn.Write([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.DataInfo, Data: "2"}))); err != nil {
				panic(err)
			}
		}
	}
}

func msgFunc(conn *net.UDPConn) {
	defer conn.Close()

	for {
		msgInfo, err := common.GetMsgInfoFromConn(conn)
		if err != nil {
			panic(err)
		}

		if msgInfo.Type == common.CloseConnection {
			return
		}

		if msgInfo.Type == common.DataInfo {
			fmt.Printf("[%s]->%s\n", msgInfo.RemoteAddr.String(), msgInfo.Data)
		}
	}
}
