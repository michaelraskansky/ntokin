package dts

import (
	"context"

	"github.com/go-stomp/stomp/v3"
	"go.uber.org/zap"
)

type ShardRecord struct {
	Shard   string
	Records *[][]byte
}
type Ctx struct {
	Region            string
	Context           context.Context
	KillSignal        chan int8
	StreamARN         string
	Log               *zap.SugaredLogger
	SubscriptionNames []string
	Subscriptions     map[string]*stomp.Subscription
	Username          string
	Password          string
	Host              string
	Port              string
	InjestChan        chan *ShardRecord
}

func (c *Ctx) Stop() {
	c.KillSignal <- 1
}
