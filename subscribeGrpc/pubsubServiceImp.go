package subscribeGrpc

import (
	"context"
	"fmt"
	"github.com/docker/docker/pkg/pubsub"
	"strings"
	"time"
)

type PubSubService struct {
	pub *pubsub.Publisher
}

func (p *PubSubService) Publish(ctx context.Context, s *String) (*String, error) {
	fmt.Println("v:", s.GetValue())
	p.pub.Publish(s.GetValue())
	return &String{}, nil
}

func (p *PubSubService) Subscribe(s *String, server PubSubService_SubscribeServer) error {
	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, s.GetValue()) {
				return true
			}
		}
		return false
	})

	for v := range ch {
		if err := server.Send(&String{Value: v.(string)}); err != nil {
			return err
		}
	}
	return nil
}

func (p *PubSubService) mustEmbedUnimplementedPubSubServiceServer() {
	panic("implement me")
}

func NewPubsubService() *PubSubService {
	return &PubSubService{
		pub: pubsub.NewPublisher(100*time.Millisecond, 10),
	}
}
