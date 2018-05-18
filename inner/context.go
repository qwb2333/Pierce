package inner

import (
	"github.com/qwb2333/Pierce/common"
)

type Context struct {
	common.Context
}

func NewContext(id uint32) (ret *Context) {
	ret = new(Context)

	ret.Id = id
	return ret
}