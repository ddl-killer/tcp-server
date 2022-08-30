package frame

import "io"

type FramePayload []byte

type StreamFrameCodec interface {
	Encode(writer io.Writer, framePayload FramePayload) error
	Decode(reader io.Reader) (FramePayload, error)
}
