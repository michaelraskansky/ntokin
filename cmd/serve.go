/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/go-stomp/stomp/v3"
	nrod "github.com/michaelraskansky/nationalrail_to_kinesis/pkg/nrod"
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
		subscriptionNames, _ := cmd.Flags().GetStringArray("subscriptions")

		sugar.Infof("serve called with %v %v:%v", username, host, port)
		ctx := &nrod.Ctx{
			Log:               sugar,
			SubscriptionNames: subscriptionNames,
			Subscriptions:     make(map[string]*stomp.Subscription, len(subscriptionNames)),
			Username:          username,
			Password:          password,
			Host:              host,
			Port:              port,
		}
		nrod.Start(ctx)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().String("host", "darwin-dist-44ae45.nationalrail.co.uk", "Hostname")
	serveCmd.Flags().String("username", "", "Username")
	serveCmd.Flags().String("password", "", "Password")
	serveCmd.Flags().String("port", "61613", "STOMP Port")
	serveCmd.Flags().StringArray("subscriptions", []string{
		"darwin.pushport-v16",
	}, "subscriptions")
}