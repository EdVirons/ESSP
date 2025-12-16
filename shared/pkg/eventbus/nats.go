package eventbus

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

type Publisher struct{ nc *nats.Conn }

func NewPublisher(nc *nats.Conn) *Publisher { return &Publisher{nc: nc} }

func (p *Publisher) PublishJSON(subject string, v any) error {
	b, err := json.Marshal(v)
	if err != nil { return err }
	return p.nc.Publish(subject, b)
}
