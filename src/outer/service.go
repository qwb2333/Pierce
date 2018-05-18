package outer

import (
	"fmt"
	"net"
	"github.com/golang/glog"
	"idl"
	"sync"
	"io"
	"time"
	"common"
	"lib"
)

type Service struct {
	common.Service

	ctxMap sync.Map
}

func NewService(msg common.OuterMsg) (ret *Service) {
	ret = new(Service)
	ret.OuterMsg = msg
	return ret
}

func (sv *Service) Run(msgLists []common.InnerOuterMsg) {
	hp := fmt.Sprintf("%s:%d", sv.OuterIp, sv.OuterPort)
	glog.Infof("OuterService Program Start, outer: %s.\n", hp)

	for _, msg := range msgLists {
		go sv.solveDotCenter(msg)
	}
	sv.solveCenter(hp)
}

func (sv *Service) solveCenter(hp string) {
	listener, err := net.Listen("tcp", hp)
	if err != nil {
		glog.Fatalf("center listen failed, hp: %s.%v\n", hp, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		glog.Infof("center have new accept.\n")
		go sv.solveCenterConn(conn)
	}
}

func (sv *Service) solveCenterConn(conn net.Conn) {
	var err error
	defer func() {
		if err != nil {
			if err == io.EOF {
				glog.Warningf("centerConn closed before assign.\n", )
			} else {
				glog.Warningf("centerConn closed before assign.%v\n", err)
			}
			if conn != nil {
				conn.Close()
			}
		}
	}()

	action := &idl.PipeAction{}
	reader, writer := idl.NewPbRW(conn)

	var runStatus lib.RunStatus
	runStatus, err = lib.RunWithTimeout (
		func() (bool, error) {
			e := reader.ReadPb(action)
			if e != nil {
				return false, e
			}
			return true, nil
		}, time.Second)

	switch runStatus {
	case lib.RUN_STATUS_FAILED:
		return
	case lib.RUN_STATUS_TIMEOUT:
		glog.Warningf("centerConn readPb timeout.\n")
		conn.Close()
		return
	}

	switch action.GetOption() {
	case idl.PipeOption_CONNECT:
		id := action.GetId()
		intf, ok := sv.ctxMap.Load(id)

		if !ok {
			glog.Infof("centerConn id: %d have assigned.\n", id)
			conn.Close()
			return
		}

		ctx := intf.(*Context)
		ctx.Lock()
		if _, ok := sv.ctxMap.Load(id); ok {
			ctx.AssignPipe(conn, reader, writer)
			ctx.WriteCh(true)
		}
		ctx.Unlock()

	case idl.PipeOption_PIPE:
		err = writer.WritePb(idl.NewActionAck())
		if err != nil {
			return
		}

		sv.Lock()
		if sv.MainPipeVersion != 0 &&
			time.Since(*sv.MainPipeReader.LastReadTime) < time.Millisecond * lib.HEART_TIMEOUT {
			sv.Unlock()
			glog.Warningf("centerConn has MainPipeConn already.\n")
			conn.Close()
			return
		} else {
			sv.AssignMainPipe(conn, reader, writer)
		}
		sv.Unlock()

		glog.Infof("centerConn build MainPipe success.\n")
		go sv.solveMainPipeConn(sv.MainPipeContext)

	default:
		glog.Warningf("centerConn option %v not allowed.\n", action.GetOption())
	}
}