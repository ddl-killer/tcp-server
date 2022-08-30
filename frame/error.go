package frame

import "errors"

var (
	ErrShortWrite = errors.New("short write")
	ErrShortRead  = errors.New("short read")
)
