package main

import (
	"fmt"
	"log"
	"net"

	"tcp-server/frame"
	"tcp-server/packet"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("listen error: ", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			break
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	codec := frame.NewMyFrameCodec()

	for {
		// decode the frame to get the payload
		framePayload, err := codec.Decode(conn)
		if err != nil {
			log.Println("handleConn: frame decode error ", err)
			break
		}

		// do something with the packet
		ackFramePayload, err := handlePacket(framePayload)
		if err != nil {
			log.Println("handleConn: handlePacket error ", err)
			break
		}

		// write ack frame to connection
		err = codec.Encode(conn, ackFramePayload)
		if err != nil {
			log.Println("handleConn: frame encode error ", err)
			break
		}
	}
}

func handlePacket(framePayload []byte) ([]byte, error) {
	var p packet.Packet
	p, err := packet.Decode(framePayload)
	if err != nil {
		log.Println("handle Packet: packet decode error ", err)
		return nil, err
	}

	switch p.(type) {
	case *packet.Submit:
		submit := p.(*packet.Submit)
		fmt.Printf("recv packet: id is %s, body is %s\n", submit.Id, string(submit.Payload))
		submitAck := &packet.SubmitAck{
			Id:     submit.Id,
			Result: 0,
		}
		ackFramePayload, err := packet.Encode(submitAck)
		if err != nil {
			return nil, err
		}
		return ackFramePayload, nil
	default:
		return nil, fmt.Errorf("unknow packet type")
	}
}
