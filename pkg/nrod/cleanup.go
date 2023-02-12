package nrod

import (
	"os"
	"os/signal"

	"github.com/michaelraskansky/ntokin/pkg/dts"
)

func cleanUpOnInterrupt(ctx *dts.Ctx) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-ctx.KillSignal:
			cleanUpAndExit(ctx)
		case <-c:
			cleanUpAndExit(ctx)
		case <-ctx.Context.Done():
			cleanUpAndExit(ctx)
		}
	}()
}

func cleanUpAndExit(ctx *dts.Ctx) {
	ctx.Log.Infof("Cleaning-up...")
	cleanUp(ctx)
	ctx.Log.Infof("Cleaned-up; exiting.")
	os.Exit(1)
}

func cleanUp(ctx *dts.Ctx) {
	for subscriptionName, subscription := range subscriptions {
		if subscription != nil {
			ctx.Log.Infof("Unsubscribing from subscription %s...", subscriptionName)
			subscription.Unsubscribe()
			ctx.Log.Infof("Successfully unsubscribed from subscription %s.", subscriptionName)

			delete(subscriptions, subscriptionName)
			subscriptionWaitGroup.Done()
		}
	}

	if connection != nil {
		ctx.Log.Infof("Disconnecting from connection...")
		connection.Disconnect()
		ctx.Log.Infof("Sucessfully disconnected from connection.")
	}
}
