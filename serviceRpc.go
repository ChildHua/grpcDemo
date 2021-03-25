package main

import (
	"context"
	"io"
)

type HelloServiceImpl struct {
}

func (h *HelloServiceImpl) Hello(ctx context.Context, s *String) (*String, error) {
	reply := &String{Value: "hello:" + s.GetValue()}
	return reply, nil
}

func (h *HelloServiceImpl) Channel(server HelloService_ChannelServer) error {
	for {
		args, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
		}
		reply := &String{Value: "hello stream:" + args.GetValue()}
		err = server.Send(reply)
		if err != nil {
			return err
		}
	}
}

func (h *HelloServiceImpl) mustEmbedUnimplementedHelloServiceServer() {
	panic("implement me")
}
