package tools

import "time"

type DelayMutex struct {
	ch  chan struct{}
}

func NewDelayMutex()DelayMutex{
	return DelayMutex{
		ch: make(chan struct{},1),
	}
}

func (m DelayMutex) Lock(tm int64)bool{
	select {
	case m.ch <- struct{}{}:
		return true
	case <-time.After(time.Duration(tm) * time.Second):
		return false
	}
}
func (m DelayMutex) Unlock(){
	<-  m.ch
}

