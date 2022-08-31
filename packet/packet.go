package packet

import (
	"bytes"
	"fmt"

	_const "tcp-server/const"
)

type Packet interface {
	Decode([]byte) error
	Encode() ([]byte, error)
}

func Decode(packet []byte) (Packet, error) {
	commandId := packet[0]
	pktBody := packet[1:]

	switch commandId {
	case _const.CommandConn:
		return nil, nil
	case _const.CommandConnAck:
		return nil, nil
	case _const.CommandSubmit:
		s := Submit{}
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	case _const.CommandSubmitAck:
		s := SubmitAck{}
		err := s.Decode(pktBody)
		if err != nil {
			return nil, err
		}
		return &s, nil
	default:
		return nil, fmt.Errorf("unknown commandId [%d]", commandId)
	}
}

func Encode(p Packet) ([]byte, error) {
	var commandId uint8
	var pktBody []byte
	var err error
	switch t := p.(type) {
	case *Submit:
		commandId = _const.CommandSubmit
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	case *SubmitAck:
		commandId = _const.CommandSubmitAck
		pktBody, err = p.Encode()
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown type [%s]", t)
	}
	return bytes.Join([][]byte{[]byte{commandId}, pktBody}, nil), nil
}
