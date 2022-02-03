package redis

import (
	"fmt"
	"sync/atomic"

	"github.com/up9inc/mizu/tap/api"
)

var matchCounter uint64

func handleClientStream(tcpID *api.TcpID, counterPair *api.CounterPair, superTimer *api.SuperTimer, emitter api.Emitter, request *RedisPacket) error {
	counterPair.Lock()
	counterPair.Request++
	requestCounter := counterPair.Request
	counterPair.Unlock()
	ident := fmt.Sprintf(
		"%d_%s:%s_%s:%s_%d",
		counterPair.StreamId,
		tcpID.SrcIP,
		tcpID.DstIP,
		tcpID.SrcPort,
		tcpID.DstPort,
		requestCounter,
	)
	// fmt.Printf("ident: %v\n", ident)
	item := reqResMatcher.registerRequest(ident, request, superTimer.CaptureTime)
	if item != nil {
		atomic.AddUint64(&matchCounter, 1)
		fmt.Printf("matchCounter: %v\n", matchCounter)
		item.ConnectionInfo = &api.ConnectionInfo{
			ClientIP:   tcpID.SrcIP,
			ClientPort: tcpID.SrcPort,
			ServerIP:   tcpID.DstIP,
			ServerPort: tcpID.DstPort,
			IsOutgoing: true,
		}
		emitter.Emit(item)
	}
	return nil
}

func handleServerStream(tcpID *api.TcpID, counterPair *api.CounterPair, superTimer *api.SuperTimer, emitter api.Emitter, response *RedisPacket) error {
	counterPair.Lock()
	counterPair.Response++
	responseCounter := counterPair.Response
	counterPair.Unlock()
	ident := fmt.Sprintf(
		"%d_%s:%s_%s:%s_%d",
		counterPair.StreamId,
		tcpID.DstIP,
		tcpID.SrcIP,
		tcpID.DstPort,
		tcpID.SrcPort,
		responseCounter,
	)
	// fmt.Printf("ident: %v\n", ident)
	item := reqResMatcher.registerResponse(ident, response, superTimer.CaptureTime)
	if item != nil {
		atomic.AddUint64(&matchCounter, 1)
		fmt.Printf("matchCounter: %v\n", matchCounter)
		item.ConnectionInfo = &api.ConnectionInfo{
			ClientIP:   tcpID.DstIP,
			ClientPort: tcpID.DstPort,
			ServerIP:   tcpID.SrcIP,
			ServerPort: tcpID.SrcPort,
			IsOutgoing: false,
		}
		emitter.Emit(item)
	}
	return nil
}
