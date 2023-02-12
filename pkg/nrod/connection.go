package nrod

import (
	"fmt"
	"net"
	"time"

	"github.com/michaelraskansky/ntokin/pkg/dts"

	"github.com/go-stomp/stomp/v3"
)

const STOMP_SERVER_TIMEOUT = 10 * time.Second

var connection *stomp.Conn

func connect(ctx *dts.Ctx) error {
	stompServer := fmt.Sprintf("%v:%v", ctx.Host, ctx.Port)
	ctx.Log.Infof("Connecting: stomp://%v ...\n", stompServer)
	networkConnection, networkConnectionError := net.DialTimeout("tcp", stompServer, STOMP_SERVER_TIMEOUT)
	if networkConnectionError != nil {
		return networkConnectionError
	}

	login := stomp.ConnOpt.Login(ctx.Username, ctx.Password)
	newConnection, connectionError := stomp.Connect(networkConnection, login)
	connection = newConnection
	if connectionError != nil {
		return connectionError
	}

	return nil
}
