/*
Copyright © 2023 Michael Raskansky michaelraskansky@gmail.com
*/
package cmd

import (
	"github.com/go-stomp/stomp/v3"
	"github.com/michaelraskansky/ntokin/pkg/dts"
	"github.com/michaelraskansky/ntokin/pkg/health"
	"github.com/michaelraskansky/ntokin/pkg/kinesis"
	"github.com/michaelraskansky/ntokin/pkg/nrod"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start streaming to kinesis",
	Long:  `start streaming to kinesis`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := zap.NewProduction()
		if err != nil {
			panic("could not init logger")
		}
		sugar := logger.Sugar()
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		kinesisStreamArn, _ := cmd.Flags().GetString("kinesis-stream-arn")
		subscriptionNames, _ := cmd.Flags().GetStringArray("subscriptions")
		region, _ := cmd.Flags().GetString("kinesis-region")
		healthcheckPort, _ := cmd.Flags().GetInt32("healthcheck-port")

		sugar.Infof("serve called with %v %v:%v", username, host, port)
		ctx := &dts.Ctx{
			Region:            region,
			Context:           cmd.Context(),
			KillSignal:        make(chan int8),
			Log:               sugar,
			SubscriptionNames: subscriptionNames,
			Subscriptions:     make(map[string]*stomp.Subscription, len(subscriptionNames)),
			Username:          username,
			Password:          password,
			Host:              host,
			Port:              port,
			StreamARN:         kinesisStreamArn,
			InjestChan:        make(chan *dts.ShardRecord),
			HealthcheckPort:   healthcheckPort,
		}
		health.Start(ctx)
		kinesis.Start(ctx)
		nrod.Start(ctx)
	},
}

var kinesisTestCmd = &cobra.Command{
	Use:   "test",
	Short: "test kinesis connect",
	Long:  `test kinesis connect`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := zap.NewProduction()
		kinesis.Test(&dts.Ctx{
			Log: logger.Sugar(),
		})
	},
}

func init() {
	rootCmd.AddCommand(kinesisTestCmd)
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("username", os.Getenv("NTOKIN_USERNAME"), "Username")
	serveCmd.Flags().String("password", os.Getenv("NTOKIN_PASSWORD"), "Password")
	serveCmd.Flags().String("kinesis-stream-arn", os.Getenv("NTOKIN_STREAM_ARN"), "the stream arn")
	serveCmd.Flags().String("host", "darwin-dist-44ae45.nationalrail.co.uk", "Hostname")
	serveCmd.Flags().String("port", "61613", "STOMP Port")
	serveCmd.Flags().String("kinesis-region", "eu-west-1", "the stream region")
	serveCmd.Flags().Int32("healthcheck-port", 8080, "the stream region")
	serveCmd.Flags().StringArray("subscriptions", []string{
		"darwin.pushport-v16",
	}, "subscriptions")
}
