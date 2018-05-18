package inner

import (
	"fmt"
	"net"
	"github.com/golang/glog"
	"github.com/qwb2333/Pierce/idl"
	"time"
	"io"
	"github.com/qwb2333/Pierce/common"
	"github.com/qwb2333/Pierce/lib"
)

type Service struct {
	common.Service
}

func NewService(msg common.OuterMsg) *Service {
	ret := new(Service)
	ret.OuterMsg = msg
	return ret
}

func (sv *Service) Run() {
	hp := fmt.Sprintf("%s:%d", sv.OuterIp, sv.OuterPort)
	glog.Infof("InnerServiceProgram Start, outer: %s.\n", hp)

	sv.solveCenter(hp)
}

func (sv *Service) solveCenter(hp string) {
	for {
		conn, err := net.Dial("tcp", hp)
		if err != nil {
			glog.Warningf("connect %s failed. wait 5s.%v\n", hp, err)
			time.Sleep(time.Second * 5)
			continue
		}

		sv.solveCenterConn(conn)
		time.Sleep(time.Second)
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

	err = writer.WritePb(idl.NewActionPipe())
	if err != nil {
		glog.Warningf("centerConn WritePb failed.%v\n", err)
		return
	}


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

	if action.GetOption() != idl.PipeOption_ACK {
		glog.Warningf("centerConn option = %s, not allowed.\n", action.GetOption())
		return
	}

	sv.Lock()
	sv.AssignMainPipe(conn, reader, writer)
	sv.Unlock()

	glog.Infof("centerConn build MainPipe success.\n")

	sv.solveMainPipeConn(sv.MainPipeContext)
}