package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "grpc-stream-demo/streampb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewStreamServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ServerStream(ctx, &pb.StreamRequest{
		Name: "Alice",
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Received:", resp.GetMessage())
	}
}
