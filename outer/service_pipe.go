package outer

import (
	"github.com/qwb2333/Pierce/idl"
	"github.com/golang/glog"
	"io"
	"github.com/qwb2333/Pierce/common"
)

func (sv *Service) solveMainPipeConn(mainPipeCtx common.MainPipeContext) {
	var err error
	action := &idl.PipeAction{}
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

	for {
		err = reader.ReadPb(action)
		if err != nil {
			return
		}

		switch action.GetOption() {
		case idl.PipeOption_ACK:
			glog.Infof("mainPipeConn get ACK.\n")
			err = writer.WritePb(idl.NewActionAck())
			if err != nil {
				return
			}
		case idl.PipeOption_DISCONNECT:
			id := action.GetId()
			intf, ok := sv.ctxMap.Load(id)

			glog.Infof("mainPipeConn id: %d disconnect.\n", id)

			if ok {
				ctx := intf.(*Context)
				ctx.Lock()
				if _, ok := sv.ctxMap.Load(id); ok {
					ctx.WriteCh(false)
				}
				ctx.Unlock()
			} else {
				glog.Warningf("id: %d not exists.\n", id)
			}
		default:
			glog.Warningf("mainPipeConn option %v not allowed.\n", action.GetOption())
		}
	}

	return
}