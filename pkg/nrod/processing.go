package nrod

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/antchfx/xmlquery"
	"github.com/go-stomp/stomp/v3"
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
	defer reader.Close()
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		ctx.Log.DPanicf("could not parse", err)
	}
	for _, node := range xmlquery.Find(doc, "//Pport/uR/TS") {
		locations := node.SelectElements("ns5:Location")
		for _, location := range locations {
			x := fmt.Sprintf("%v", location.OutputXML(true))
			ctx.Log.Infof("%v", x)
		}
	}

	for _, node := range xmlquery.Find(doc, "//Pport/uR/schedule") {
		ors := node.SelectElements("ns2:OR")
		for _, or := range ors {
			x := fmt.Sprintf("%v", or.OutputXML(true))
			ctx.Log.Debugf("Schedule OR: %v", x)
		}
		ips := node.SelectElements("ns2:IP")
		for _, x := range ips {
			x := fmt.Sprintf("%v", x.OutputXML(true))
			ctx.Log.Debugf("Schedule IP: %v", x)
		}
		dts := node.SelectElements("ns2:DT")
		for _, x := range dts {
			x := fmt.Sprintf("%v", x.OutputXML(true))
			ctx.Log.Debugf("Schedule DT: %v", x)
		}
	}
}
