package inner

import (
	"common"
	"idl"
	"io"
	"github.com/golang/glog"
	"time"
	"lib"
)

func (sv *Service) solveMainPipeConn(mainPipeCtx common.MainPipeContext) {
	var err error
	action := idl.PipeAction{}
	conn := mainPipeCtx.MainPipeConn
	reader := mainPipeCtx.MainPipeReader
	writer := mainPipeCtx.MainPipeWriter

	defer func() {
		if err == io.EOF {
			glog.Warningf("mainPipeConn closed.\n")
		} else {
			glog.Warningf("mainPipeConn closed.%v\n", err)
		}
		if conn != nil {
			conn.Close()
		}
	}()

	go func() {
		for {
			time.Sleep(time.Millisecond * lib.HEART_INTERVAL)
			err = writer.WritePb(idl.NewActionAck())
			if err != nil {
				break
			}
		}
	}()

	for {
		err = reader.ReadPb(&action)
		if err != nil {
			return
		}

		switch action.GetOption() {
		case idl.PipeOption_ACK:
			glog.Infof("mainPipeConn get ACK.\n")
			continue

		case idl.PipeOption_CONNECT:
			glog.Infof("mainPipeConn get connect id: %d\n", action.GetId())
			go sv.solveDotCenterConn(action, mainPipeCtx)

		default:
			glog.Warningf("mainPipeConn option %v not allowed.\n", action.GetOption())
		}
	}
}