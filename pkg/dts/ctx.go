package dts

import (
	"github.com/go-stomp/stomp/v3"
	"go.uber.org/zap"
)

type Ctx struct {
	Log               *zap.SugaredLogger
	SubscriptionNames []string
	Subscriptions     map[string]*stomp.Subscription
	Username          string
	Password          string
	Host              string
	Port              string
}
