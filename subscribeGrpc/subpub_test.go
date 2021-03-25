package subscribeGrpc

import (
	"context"
	"fmt"
	"github.com/docker/docker/pkg/pubsub"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"testing"
	"time"
)

func TestDockerPub(t *testing.T) {
	p := pubsub.NewPublisher(100*time.Millisecond, 10)
	p.Publish("aaa")
}

func TestSubpubService(t *testing.T) {
	grpcServer := grpc.NewServer()
	RegisterPubSubServiceServer(grpcServer, NewPubsubService())
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}

func TestPublish(t *testing.T) {
	conn, err := grpc.Dial(":1234", grpc.WithInsecure())

	defer conn.Close()
	client := NewPubSubServiceClient(conn)
	_, err = client.Publish(context.Background(), &String{Value: "golang:hello Go"})
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.Publish(context.Background(), &String{Value: "docker: hello Docker"})

	if err != nil {
		log.Fatal(err)
	}
}

func TestSubscribe(t *testing.T) {
	conn, err := grpc.Dial(":1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := NewPubSubServiceClient(conn)
	stream, err := client.Subscribe(context.Background(), &String{Value: "golang:"})
	if err != nil {
		log.Fatal(err)
	}

	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}
}
