package common

import (
	"bytes"
	"net"
	"strconv"
	"strings"
)

const (
	ConnectReq MsgType = iota
	ConnectRsp

	GetClientListReq
	GetClientListRsp

	ConnectClientReq
	ConnectClientRsp

	ErrorInfo
	IgnoreInfo

	DataInfo
	CloseConnection
)

type MsgType int

type MsgInfo struct {
	Type       MsgType
	Data       string
	RemoteAddr *net.UDPAddr
}

func GetMsgInfoFromConn(conn *net.UDPConn) (*MsgInfo, error) {
	data := make([]byte, 1024)
	dataLen, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		return nil, err
	}

	dataInfo := strings.Split(string(data[:dataLen]), "\n")

	msgType, err := strconv.Atoi(dataInfo[0])
	if err != nil {
		return nil, err
	}

	r := &MsgInfo{
		Type:       MsgType(msgType),
		Data:       dataInfo[1],
		RemoteAddr: remoteAddr,
	}
	return r, nil
}

func ConvertMsgInfo(msgInfo *MsgInfo) string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString(strconv.Itoa(int(msgInfo.Type)) + "\n")
	buf.WriteString(msgInfo.Data)
	return buf.String()
}

func ParseAddr(addr string) *net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return &net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}
