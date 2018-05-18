package idl

import "github.com/golang/protobuf/proto"

func NewActionAck() *PipeAction {
	ret := PipeAction {
		Option: PipeOption_ACK.Enum(),
	}
	return &ret
}

func NewRemoteInfo(ip string, port int) *RemoteInfo {
	info := RemoteInfo {
		Ip: proto.String(ip),
		Port: proto.Int(port),
	}
	return &info
}

func NewActionConnect(uniqueId uint32, info *RemoteInfo) *PipeAction {
	ret := PipeAction {
		Id:         proto.Uint32(uniqueId),
		RemoteInfo: info,
		Option:     PipeOption_CONNECT.Enum(),
	}
	return &ret
}

func NewActionData(uniqueId uint32, data []byte) *PipeAction {
	ret := PipeAction {
		Id:     proto.Uint32(uniqueId),
		Data:   data,
		Option: PipeOption_DATA.Enum(),
	}
	return &ret
}

func NewActionPipe() *PipeAction {
	ret := PipeAction {
		Option: PipeOption_PIPE.Enum(),
	}
	return &ret
}

func NewActionDisconnect(uniqueId uint32) *PipeAction {
	ret := PipeAction {
		Id:     proto.Uint32(uniqueId),
		Option: PipeOption_DISCONNECT.Enum(),
	}
	return &ret
}