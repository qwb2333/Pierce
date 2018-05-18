package inner

import (
	"common"
)

type Context struct {
	common.Context
}

func NewContext(id uint32) (ret *Context) {
	ret = new(Context)

	ret.Id = id
	return ret
}