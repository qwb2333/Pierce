package outer

import (
	"common"
	"lib"
	"sync"
)

type Context struct {
	common.Context

	Ch  chan bool
	sync.Mutex
}

func NewContext(id uint32) (ret *Context) {
	ret = new(Context)

	ret.Id = id
	ret.Ch = make(chan bool)
	return ret
}

func (t *Context) WriteCh(x bool) {
	lib.WriteCh(t.Ch, x)
}