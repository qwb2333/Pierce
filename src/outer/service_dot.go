package outer

import (
	"fmt"
	"net"
	"github.com/golang/glog"
	"idl"
	"time"
	"common"
	"lib"
)

func (sv *Service) solveDotCenter(msg common.InnerOuterMsg) {
	hp := fmt.Sprintf("%s:%d", msg.OuterIp, msg.OuterPort)
	glog.Infof("listen %s\n", hp)

	listener, err := net.Listen("tcp", hp)
	if err != nil {
		glog.Fatalf("dotCenter listen failed, hp: %s.%v.\n", hp, err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		glog.Infof("dotCenter have new accept.\n")
		go sv.solveDotCenterConn(&msg, conn)
	}
}

func (sv *Service) solveDotCenterConn(msg *common.InnerOuterMsg, conn net.Conn) {
	ctx := NewContext(lib.GetNextId())
	glog.Infof("dotCenterConn id: %d connect.\n", ctx.Id)

	sv.RLock()
	mpc := sv.MainPipeContext
	sv.RUnlock()

	if mpc.MainPipeVersion == 0 {
		glog.Infof("dotCenterConn id: %d there is no MainPipe.\n", ctx.Id)
		conn.Close()
		return
	}

	sv.ctxMap.Store(ctx.Id, ctx)
	glog.Infof("dotCenterConn id: %d send connect to MainPipe.\n", ctx.Id)
	err := mpc.MainPipeWriter.WritePb(idl.NewActionConnect(ctx.Id, idl.NewRemoteInfo(msg.InnerIp, msg.InnerPort)))

	if err != nil {
		glog.Warningf("dotCenterConn id: %d writePb error, %v.\n", ctx.Id, err)
		sv.ctxMap.Delete(ctx.Id)
		mpc.MainPipeConn.Close()
		conn.Close()

		sv.Lock()
		sv.MainPipeVersion = 0
		sv.Unlock()
		return
	}

	var result bool
	select {
	case result = <-ctx.Ch:
		if result {
			glog.Infof("dotCenterConn id: %d SolvePipeAndDot.\n", ctx.Id)
		} else {
			glog.Infof("dotCenterConn id: %d inner refuse this connect.\n", ctx.Id)
		}
	case <-time.After(time.Second):
		glog.Infof("dotCenterConn id: %d timeout when assigning.\n", ctx.Id)
	}

	ctx.Lock()
	sv.ctxMap.Delete(ctx.Id)
	ctx.Unlock()

	if result {
		ctx.AssignDot(conn, conn, conn)
		ctx.SolvePipeAndDot()
	} else {
		conn.Close()
	}
}