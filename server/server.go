package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"udp-dig/common"
)

var (
	port      int
	clientMap = make(map[string]*net.UDPAddr)
)

func main() {
	flag.IntVar(&port, "port", 3356, "")
	flag.Parse()

	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		panic(err)
	}

	for {
		msgInfo, err := common.GetMsgInfoFromConn(listener)
		if err != nil {
			listener.WriteToUDP([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.ErrorInfo, Data: err.Error()})), msgInfo.RemoteAddr)
			continue
		}

		if msgInfo.Type == common.ConnectReq {
			clientMap[msgInfo.RemoteAddr.String()] = msgInfo.RemoteAddr

			listener.WriteToUDP([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.ConnectRsp, Data: "ok"})), msgInfo.RemoteAddr)
		} else if msgInfo.Type == common.GetClientListReq {
			var clientList []string
			for k := range clientMap {
				clientList = append(clientList, k)
			}

			listener.WriteToUDP([]byte(common.ConvertMsgInfo(&common.MsgInfo{Type: common.GetClientListRsp,
				Data: strings.Join(clientList, "|")})), msgInfo.RemoteAddr)
		} else if msgInfo.Type == common.ConnectClientReq {
			dstRemoteAddr, ok := clientMap[msgInfo.Data]
			if !ok {
				listener.WriteToUDP([]byte(common.ConvertMsgInfo(&common.MsgInfo{Type: common.ErrorInfo,
					Data: fmt.Sprintf("addr no exist: %s", msgInfo.Data)})), msgInfo.RemoteAddr)
				continue
			}

			listener.WriteToUDP([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.ConnectClientReq, Data: msgInfo.RemoteAddr.String()})), dstRemoteAddr)

			listener.WriteToUDP([]byte(common.ConvertMsgInfo(
				&common.MsgInfo{Type: common.ConnectClientRsp, Data: "ok"})), msgInfo.RemoteAddr)
		}
	}
}
