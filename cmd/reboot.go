/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	coap "github.com/moroen/gocoap/v5"
	"github.com/moroen/tradfri2mqtt/settings"
	"github.com/moroen/tradfri2mqtt/tradfri"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rebootCmd represents the reboot command
var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// tradfri.
		//	tradfri.RebootGateway()

		conf := settings.GetConfig(false)
		connection := coap.CoapDTLSConnection{}
		connection.Host = conf.Tradfri.Gateway
		connection.Port = 5684
		connection.Ident = conf.Tradfri.Identity
		connection.Key = conf.Tradfri.Passkey
		connection.OnConnect = func() {
			tradfri.SetConnecion(connection)
			tradfri.RebootGateway()
		}
		connection.OnConnectionFailed = func() {
			log.Error("Unable to connect to gateway")
		}

		connection.Connect()
	},
}

func init() {
	rootCmd.AddCommand(rebootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rebootCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rebootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
