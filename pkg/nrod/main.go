package nrod

import "github.com/michaelraskansky/nationalrail_to_kinesis/pkg/dts"

func Start(ctx *dts.Ctx) {
	cleanUpOnInterrupt(ctx)

	connectionError := connect(ctx)
	defer connection.Disconnect()
	if connectionError != nil {
		ctx.Log.Fatalf("error establishing connection: %v", connectionError.Error())
	} else {
		ctx.Log.Infof("Connected: %v", connection.Session())

		for _, subscriptionName := range ctx.SubscriptionNames {
			subscriptionWaitGroup.Add(1)
			go workSubscription(ctx, subscriptionName)
		}
		subscriptionWaitGroup.Wait()
	}

	if len(subscriptions) > 0 {
		ctx.Log.Infof("Cleaning up...")
		cleanUp(ctx)
	}
}

func workSubscription(ctx *dts.Ctx, subscriptionName string) {
	subscription, subscriptionError := subscribe(subscriptionName)

	if subscriptionError != nil {
		ctx.Log.Infof("error subscribing: %v", subscriptionError)
		return
	}

	ctx.Subscriptions[subscriptionName] = subscription
	ctx.Log.Infof("Subscribed: %s (%v)", subscriptionName, subscription.Id())

	for {
		if subscription == nil {
			return
		}

		processMessages(ctx, subscription)
	}
}
