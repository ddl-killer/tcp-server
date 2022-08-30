package frame

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"testing"
)

func TestEncode(t *testing.T) {
	codec := NewMyFrameCodec()
	buf := make([]byte, 0, 128)
	rw := bytes.NewBuffer(buf)

	err := codec.Encode(rw, []byte("hello"))
	if err != nil {
		t.Errorf("error: %s", err)
	}

	var totalLen int32
	err = binary.Read(rw, binary.BigEndian, &totalLen)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	if totalLen != 9 {
		t.Errorf("want 9, actual %d", totalLen)
	}

	left := rw.Bytes()
	if string(left) != "hello" {
		t.Errorf("want hello, acutal %s", string(left))
	}
}

func TestDecode(t *testing.T) {
	codec := NewMyFrameCodec()
	data := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}

	framePayload, err := codec.Decode(bytes.NewReader(data))
	if err != nil {
		t.Errorf("error: %s", err)
	}

	if string(framePayload) != "hello" {
		t.Errorf("want hello, acutal %s", string(framePayload))
	}
}

type ReturnErrorWriter struct {
	W  io.Writer
	Wn int // 第几次调用Write返回错误
	wc int // 写操作次数技术
}

func (w *ReturnErrorWriter) Write(p []byte) (int, error) {
	w.wc++
	if w.wc >= w.Wn {
		return 0, errors.New("write error")
	}
	return w.W.Write(p)
}

type ReturnErrorReader struct {
	R  io.Reader
	Rn int
	rc int
}

func (r *ReturnErrorReader) Read(p []byte) (int, error) {
	r.rc++
	if r.rc >= r.Rn {
		return 0, errors.New("read error")
	}
	return r.R.Read(p)
}

func TestEncodeWithWriteFail(t *testing.T) {
	codec := NewMyFrameCodec()
	buf := make([]byte, 0, 128)
	w := bytes.NewBuffer(buf)

	err := codec.Encode(&ReturnErrorWriter{
		W:  w,
		Wn: 1,
	}, []byte("hello"))
	if err == nil {
		t.Errorf("want non-nil, actual nil")
	}

	err = codec.Encode(&ReturnErrorWriter{
		W:  w,
		Wn: 2,
	}, []byte("hello"))
	if err == nil {
		t.Errorf("want non-nil, actual nil")
	}
}

func TestEncodeWithReadFail(t *testing.T) {
	codec := NewMyFrameCodec()
	data := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}

	_, err := codec.Decode(&ReturnErrorReader{
		R:  bytes.NewReader(data),
		Rn: 1,
	})
	if err == nil {
		t.Errorf("want non-nil, actual nil")
	}

	_, err = codec.Decode(&ReturnErrorReader{
		R:  bytes.NewReader(data),
		Rn: 2,
	})
	if err == nil {
		t.Errorf("want non-nil, actual nil")
	}
}
