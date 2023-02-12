package nrod

import (
	"bytes"
	"compress/gzip"
	"fmt"

	"github.com/antchfx/xmlquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/go-stomp/stomp/v3"
	"github.com/michaelraskansky/nationalrail_to_kinesis/pkg/dts"
)

func newKinesisClient() *kinesis.Kinesis {
	mySession := session.Must(session.NewSession())
	svc := kinesis.New(mySession, aws.NewConfig().WithRegion("eu-west-1"))
	return svc
}

func processMessages(ctx *dts.Ctx, subscription *stomp.Subscription) error {
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

func processMessageGzip(ctx *dts.Ctx, subscription *stomp.Subscription, subscriptionMessage *stomp.Message) {
	reader, err := gzip.NewReader(bytes.NewReader(subscriptionMessage.Body))
	if err != nil {
		ctx.Log.Panicf("could not read %v", err)
	}
	defer reader.Close()
	doc, err := xmlquery.Parse(reader)
	if err != nil {
		ctx.Log.DPanicf("could not parse", err)
	}

	proceessSchedule(ctx, doc)
	proccessLocations(ctx, doc)
}

func proceessSchedule(ctx *dts.Ctx, doc *xmlquery.Node) {
	for _, node := range xmlquery.Find(doc, "//Pport/uR/schedule") {
		rid := node.SelectAttr("rid")         // darwin unique id
		uid := node.SelectAttr("uid")         // schedule unique id
		trainId := node.SelectAttr("trainId") // the train id
		ssd := node.SelectAttr("ssd")         // day when train starts
		toc := node.SelectAttr("toc")         // company running the train
		for _, pt := range []string{"OR", "OPOR", "PP", "IP", "OPIP", "DT", "OPDT"} {
			prosessSchedulePoint(ctx, node, pt, rid, uid, trainId, ssd, toc)
		}
	}
}

func proccessLocations(ctx *dts.Ctx, doc *xmlquery.Node) {
	for _, node := range xmlquery.Find(doc, "//Pport/uR/TS") {
		locations := node.SelectElements("ns5:Location")
		rid := node.SelectAttr("rid") // darwin unique id
		uid := node.SelectAttr("uid") // schedule unique id
		ssd := node.SelectAttr("ssd") // day when train starts
		for _, location := range locations {
			tpl := location.SelectAttr("tpl") // timing point location
			pta := location.SelectAttr("pta") // public timetable arraival
			ptd := location.SelectAttr("ptd") // public timetable departure
			wta := location.SelectAttr("wta") // working timetable arraival
			wtd := location.SelectAttr("wtd") // working timetable departure
			wtp := location.SelectAttr("wtp") // working timetable pass
			ctx.Log.Infof("ts,%v,%v,%v,%v,%v,%v,%v,%v,%v", rid, uid, ssd, tpl, pta, ptd, wta, wtd, wtp)
		}
	}
}

func prosessSchedulePoint(ctx *dts.Ctx, node *xmlquery.Node, schedulePointType string, rid string, uid string, trainId string, ssd string, toc string) {
	opdts := node.SelectElements(fmt.Sprintf("ns2:%v", schedulePointType))
	for _, x := range opdts {
		tpl := x.SelectAttr("tpl") // timing point location
		pta := x.SelectAttr("pta") // public timetable arraival
		ptd := x.SelectAttr("ptd") // public timetable departure
		wta := x.SelectAttr("wta") // working timetable arraival
		wtd := x.SelectAttr("wtd") // working timetable departure
		wtp := x.SelectAttr("wtp") // working timetable pass
		act := x.SelectAttr("act") // ???
		ctx.Log.Infof("s,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v", rid, uid, ssd, tpl, pta, ptd, wta, wtd, wtp, act, schedulePointType, trainId, toc)

	}
}
