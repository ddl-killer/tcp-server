package packet

import "bytes"

type Submit struct {
	Id      string
	Payload []byte
}

func (s *Submit) Decode(pktBody []byte) error {
	s.Id = string(pktBody[:8])
	s.Payload = pktBody[8:]
	return nil
}

func (s *Submit) Encode() ([]byte, error) {
	return bytes.Join([][]byte{[]byte(s.Id[:8]), s.Payload}, nil), nil
}
