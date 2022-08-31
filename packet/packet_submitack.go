package packet

import "bytes"

type SubmitAck struct {
	Id     string
	Result uint8
}

func (s *SubmitAck) Decode(pktBody []byte) error {
	s.Id = string(pktBody[:8])
	s.Result = uint8(pktBody[8])
	return nil
}

func (s *SubmitAck) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.Id[:8]), []byte{s.Result}}, nil), nil
}
