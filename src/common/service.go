package common

import (
	"sync"
)

type Service struct {
	sync.RWMutex
	OuterMsg
	MainPipeContext
}