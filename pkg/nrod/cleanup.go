package nrod

import (
	"os"
	"os/signal"
)

func cleanUpOnInterrupt(ctx *Ctx) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			ctx.Log.Infof("Cleaning-up...")
			cleanUp(ctx)
			ctx.Log.Infof("Cleaned-up; exiting.")
			os.Exit(1)
		}
	}()
}

func cleanUp(ctx *Ctx) {
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