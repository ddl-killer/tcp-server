package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"tcp-server/frame"
	"tcp-server/packet"
)

func main() {
	var wg sync.WaitGroup
	var num int = 2

	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(i int) {
			defer wg.Done()
			startClient(i)
		}(i + 1)
	}
	wg.Wait()
}

func startClient(i int) {
	quit := make(chan struct{})
	done := make(chan struct{})
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		fmt.Printf("dial error: %s\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("[client %d]: dial ok\n", i)

	frameCodec := frame.NewMyFrameCodec()
	var counter int

	go func() {
		// handel ack
		for {
			select {
			case <-quit:
				done <- struct{}{}
				return
			default:
			}

			conn.SetReadDeadline(time.Now().Add(time.Second))
			ackFramePayload, err := frameCodec.Decode(conn)
			if err != nil {
				if e, ok := err.(net.Error); ok {
					if e.Timeout() {
						fmt.Printf("[client %d]: handle ack decode time out", i)
						continue
					}
				}
				panic(err)
			}

			p, err := packet.Decode(ackFramePayload)
			submitAck, ok := p.(*packet.SubmitAck)
			if !ok {
				panic("not submitack")
			}
			fmt.Printf("[client %d]: the result of submit ack[%s] is %d\n", i, submitAck.Id, submitAck.Result)
		}
	}()

	for {
		counter++
		id := fmt.Sprintf("%08d", counter)

		s := &packet.Submit{
			Id:      id,
			Payload: []byte("hello npc"),
		}

		framePayload, err := packet.Encode(s)
		if err != nil {
			panic(err)
		}

		fmt.Printf("[client %d]: send submit id = %s, payload=%s, frame lenth = %d\n", i, s.Id, s.Payload, len(framePayload)+4)
		err = frameCodec.Encode(conn, framePayload)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 1)
		if counter >= 10 {
			quit <- struct{}{}
			<-done
			fmt.Printf("[client %d]: exit ok\n", i)
			return
		}
	}
}
