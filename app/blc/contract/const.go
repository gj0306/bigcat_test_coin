package contract

type ConnType int64

const (
	ConnTypeNull int64 = iota
	ConnTypeProfile
	ConnTypeNews
	ConnTypeSource
)
