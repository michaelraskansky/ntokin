package nrod

import (
	"bytes"
	"compress/gzip"
	"github.com/go-stomp/stomp/v3"
	messageTypes "github.com/michaelraskansky/nationalrail_to_kinesis/pkg/nrod/messages"
	"io/ioutil"
)

func processMessages(ctx *Ctx, subscription *stomp.Subscription) error {
	message, messageError := getMessage(subscription)
	if messageError != nil {
		return messageError
	}

	if message != nil {
		processMessageGzip(ctx, subscription, message)
	}

	return nil
}

func getMessage(subscription *stomp.Subscription) (*stomp.Message, error) {
	message := <-subscription.C
	if message == nil {
		return nil, nil
	} else {
		return message, message.Err
	}
}

func processMessageGzip(ctx *Ctx, subscription *stomp.Subscription, subscriptionMessage *stomp.Message) {
	reader, err := gzip.NewReader(bytes.NewReader(subscriptionMessage.Body))
	if err != nil {
		ctx.Log.Panicf("could not read %v", err)
	}
	reader.Close()
	s, err := ioutil.ReadAll(reader)
	if err != nil {
		ctx.Log.Panicf("could not read %v", err)
	}
	xml := string(s)
	ctx.Log.Debugf("%v", xml)
}

func processMessage(ctx *Ctx, subscription *stomp.Subscription, subscriptionMessage *stomp.Message) {
	messages := messageTypes.Detect(subscriptionMessage.Body)
	for _, message := range messages {
		output := message.ToString()
		if output != "" {
			ctx.Log.Infof("[%v] %v", subscription.Id(), output)
		}
	}
}
