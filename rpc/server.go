package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "grpc-stream-demo/streampb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedStreamServiceServer
}

// 服务端流式
func (s *server) ServerStream(req *pb.StreamRequest, stream pb.StreamService_ServerStreamServer) error {
	name := req.GetName()

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Hello %s, message %d", name, i)

		err := stream.Send(&pb.StreamResponse{
			Message: msg,
		})
		if err != nil {
			return err
		}

		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStreamServiceServer(grpcServer, &server{})

	fmt.Println("Server started at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
