package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	pb "github.com/bezzang-dev/go-grpc-chat/chatproto"
	"google.golang.org/grpc"
)


var port = flag.Int("port", 50051, "The server port")

func main() {
    flag.Parse()
    
    // TCP 리스너 생성
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    // gRPC 서버 생성
    grpcServer := grpc.NewServer()
    
    // 채팅 서비스 등록
    pb.RegisterChatServiceServer(grpcServer, NewServer())
    
    // 서버 시작 (블로킹)
    grpcServer.Serve(lis)
}

type ChatServer struct {
	pb.UnimplementedChatServiceServer
	mu sync.RWMutex
	streams map[pb.ChatService_ChatServer]struct{}
}

func NewServer() *ChatServer {
	return &ChatServer{
        streams: make(map[pb.ChatService_ChatServer]struct{}),
    }
}

func (s *ChatServer) addStream(stream pb.ChatService_ChatServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.streams[stream] = struct{}{}
}

func (s *ChatServer) removeStream(stream pb.ChatService_ChatServer) {
	s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.streams, stream)
}

func (s *ChatServer) getStreamCount() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.streams)
}

func (s *ChatServer) broadcast(msg *pb.ChatMsg) {
	s.mu.RLock()
	
	streams := make([]pb.ChatService_ChatServer, 0, len(s.streams))
	for strm := range s.streams {
		streams = append(streams, strm)
	}
	s.mu.RUnlock()

	var wg sync.WaitGroup

	for _, strm := range streams {
		wg.Add(1)

		go func(st pb.ChatService_ChatServer) {
			defer wg.Done()
			if err := st.Send(msg); err != nil {
				log.Printf("Failed to send message: %v", err)
                s.removeStream(st)  // 실패한 스트림 제거
			}
		}(strm)
	}
	wg.Wait()
}

func (s *ChatServer) Chat(stream pb.ChatService_ChatServer) error {
	ctx := stream.Context()

    s.addStream(stream)
    defer s.removeStream(stream)

	log.Printf("Client connected. Total: %d", s.getStreamCount())
    defer log.Printf("Client disconnected. Total: %d", s.getStreamCount())

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Printf("Receive error: %v", err)
                return err
			}
			log.Printf("Broadcasting message from %s: %s", msg.Sender, msg.Message)
            s.broadcast(msg)
		}
	}
}
