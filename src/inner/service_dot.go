package inner

import (
	"idl"
	"common"
	"fmt"
	"net"
	"github.com/golang/glog"
	"io"
)

func handleError(connName string, conn net.Conn, err error) {
	if err != nil {
		if err == io.EOF {
			glog.Warningf("dotCenterConn %s closed.\n", connName)
		} else {
			glog.Warningf("dotCenterConn %s closed.%v\n", connName, err)
		}
		if conn != nil {
			conn.Close()
		}
	}
}

func (sv *Service) solveDotCenterConn(ac idl.PipeAction, mainPipeCtx common.MainPipeContext) {
	id := ac.GetId()

	remoteInfo := ac.GetRemoteInfo()
	hp := fmt.Sprintf("%s:%d", remoteInfo.GetIp(), remoteInfo.GetPort())

	ctx := NewContext(id)

	dotConn, err := net.Dial("tcp", hp)
	if err != nil {
		glog.Infof("dotCenterConn id: %d connect %s failed.\n", id, hp)
		handleError("DotConn", dotConn, err)

		err := mainPipeCtx.MainPipeWriter.WritePb(idl.NewActionDisconnect(id))
		if err != nil {
			glog.Warningf("dotCenterConn id: %d write disconnect action failed.\n", id)
			handleError("MainPipeConn", mainPipeCtx.MainPipeConn, err)
		}
		return
	}
	ctx.AssignDot(dotConn, dotConn, dotConn)

	hp = fmt.Sprintf("%s:%d", sv.OuterIp, sv.OuterPort)
	pipeConn, err := net.Dial("tcp", hp)
	if err != nil {
		glog.Warningf("dotCenterConn id: %d connect to outer failed.\n", id)
		handleError("PipeConn", pipeConn, err)
		return
	}

	pipeReader, pipeWriter := idl.NewPbRW(pipeConn)
	glog.Infof("dotCenterConn id: %d prepare to write Connect to pipe.\n", id)
	err = pipeWriter.WritePb(idl.NewActionConnect(id, nil))
	if err != nil {
		glog.Warningf("dotCenterConn id: %d write connect action failed.\n", id)
		handleError("PipeConn", pipeConn, err)
	}
	ctx.AssignPipe(pipeConn, pipeReader, pipeWriter)
	ctx.SolvePipeAndDot()
	glog.Infof("dotCenterConn id: %d SolvePipeAndDot.\n", id)
}
