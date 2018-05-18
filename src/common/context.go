package common

import (
	"net"
	"idl"
	"lib"
	"github.com/golang/glog"
	"io"
)

type MainPipeContext struct {
	MainPipeVersion uint32
	MainPipeConn    net.Conn
	MainPipeReader  *idl.PbReader
	MainPipeWriter  *idl.PbWriter
}

func (ctx *MainPipeContext) AssignMainPipe(conn net.Conn, reader *idl.PbReader, writer *idl.PbWriter) {
	ctx.MainPipeVersion = lib.GetNextId()
	ctx.MainPipeConn = conn
	ctx.MainPipeReader = reader
	ctx.MainPipeWriter = writer
}

type PipeContext struct {
	PipeConn   net.Conn
	PipeReader *idl.PbReader
	PipeWriter *idl.PbWriter
}

func (ctx *PipeContext) AssignPipe(conn net.Conn, reader *idl.PbReader, writer *idl.PbWriter) {
	ctx.PipeConn = conn
	ctx.PipeReader = reader
	ctx.PipeWriter = writer
}

type DotContext struct {
	DotConn   net.Conn
	DotReader io.Reader
	DotWriter io.Writer
}

func (ctx *DotContext) AssignDot(conn net.Conn, reader io.Reader, writer io.Writer) {
	ctx.DotConn = conn
	ctx.DotReader = reader
	ctx.DotWriter = writer
}


type Context struct {
	PipeContext
	DotContext

	Id uint32
}

func handleError(connName string, ctx *Context, err error) {
	if err == io.EOF {
		glog.Infof("dot id: %d %s closed.\n", ctx.Id, connName)
	} else {
		glog.Infof("dot id: %d %s closed.%v\n", ctx.Id, connName, err)
	}

	if ctx.DotConn != nil {
		ctx.DotConn.Close()
	}
	if ctx.PipeConn != nil {
		ctx.PipeConn.Close()
	}
}

func (ctx *Context) SolvePipeAndDot() {

	go func() {
		var err error

		buff := make([]byte, lib.BUFF_SIZE)
		for {
			var l int
			l, err = ctx.DotReader.Read(buff)
			if err != nil {
				handleError("dot", ctx, err)
				return
			} else {
				glog.Infof("PipeAndDot id: %d read data, len: %d\n", ctx.Id, l)
			}

			err = ctx.PipeWriter.WritePb(idl.NewActionData(ctx.Id, buff[0:l]))
			if err != nil {
				handleError("pipe", ctx, err)
				return
			}
		}
	}()

	go func() {
		var err error

		action := idl.PipeAction{}
		for {
			err = ctx.PipeReader.ReadPb(&action)
			if err != nil {
				handleError("pipe", ctx, err)
				return
			}

			switch action.GetOption() {
			case idl.PipeOption_DATA:

				data := action.GetData()
				glog.Infof("PipeAndDot id: %d read pipe data, len: %d\n", ctx.Id, len(data))

				err = idl.WriteBytes(ctx.DotWriter, data)
				if err != nil {
					handleError("dot", ctx, err)
					return
				}

			default:
				glog.Warningf("pipe option %v not allowed.\n", action.GetOption())
			}
		}
	}()
}


type InnerMsg struct {
	InnerIp   string
	InnerPort int
}

type OuterMsg struct {
	OuterIp   string
	OuterPort int
}
type InnerOuterMsg struct {
	InnerMsg
	OuterMsg
}