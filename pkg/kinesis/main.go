package kinesis

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func Start() {

}
func newKinesisClient() *kinesis.Kinesis {
	mySession := session.Must(session.NewSession())
	svc := kinesis.New(mySession, aws.NewConfig().WithRegion("eu-west-1"))
	return svc
}
