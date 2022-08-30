package frame

import (
	"encoding/binary"
	"io"
)

type myFrameCodec struct {
}

func NewMyFrameCodec() StreamFrameCodec {
	return &myFrameCodec{}
}

func (m myFrameCodec) Encode(writer io.Writer, framePayload FramePayload) error {
	var f = framePayload
	var totalLen int32 = int32(len(framePayload)) + 4

	err := binary.Write(writer, binary.BigEndian, &totalLen)
	if err != nil {
		return err
	}

	n, err := writer.Write(f)
	if err != nil {
		return err
	}
	if n < len(framePayload) {
		return ErrShortWrite
	}

	return nil
}

func (m myFrameCodec) Decode(reader io.Reader) (FramePayload, error) {
	var totalLen int32
	err := binary.Read(reader, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	n, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}

	return FramePayload(buf), nil
}
