package kinesis

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/michaelraskansky/ntokin/pkg/dts"
)

type KinesisSync struct {
	Ctx       *dts.Ctx
	StreamARN string
	Client    *kinesis.Kinesis
	killChan  chan bool
}

func (ks *KinesisSync) Kill() {
	go func() {
		ks.killChan <- true
	}()

}

func Test(ctx *dts.Ctx) bool {
	kc := newKinesisClient(ctx)
	o, err := kc.DescribeLimits(&kinesis.DescribeLimitsInput{})
	if err != nil {
		ctx.Log.Errorf("could not connect to kinesis %v", err)
	}
	ctx.Log.Infof("got %v", o)
	return *o.ShardLimit > int64(0)
}

func Start(ctx *dts.Ctx) {
	kc := newKinesisClient(ctx)
	ks := &KinesisSync{
		StreamARN: ctx.StreamARN,
		Client:    kc,
		killChan:  make(chan bool),
	}
	go func() {
		injestLoop(ks, ctx)
	}()
}

func injestLoop(kc *KinesisSync, ctx *dts.Ctx) {
	for {
		select {
		case <-kc.killChan:
			ctx.Log.Infof("killing injest loop")
			ctx.Stop()
		case shardArray := <-ctx.InjestChan:
			ctx.Log.Debugf("injest data to %v", shardArray.Shard)
			var records []*kinesis.PutRecordsRequestEntry
			for _, bytearray := range *shardArray.Records {
				records = append(records, &kinesis.PutRecordsRequestEntry{
					Data:         bytearray,
					PartitionKey: &shardArray.Shard,
				})
			}
			pri := &kinesis.PutRecordsInput{
				StreamARN: &kc.StreamARN,
				Records:   records,
			}
			_, err := kc.Client.PutRecords(pri)
			if err != nil {
				ctx.Log.Errorf("could not send to kinesis %v", pri)
				ctx.Log.Errorf("could not send to kinesis %v", err)
				ctx.Stop()
			}
		}
	}
}

func newKinesisClient(ctx *dts.Ctx) *kinesis.Kinesis {
	mySession := session.Must(session.NewSession())
	svc := kinesis.New(mySession, aws.NewConfig().WithRegion(ctx.Region))
	return svc
}
