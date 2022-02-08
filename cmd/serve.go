/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/moroen/tradfri2mqtt/mqttclient"
	"github.com/moroen/tradfri2mqtt/settings"
	"github.com/moroen/tradfri2mqtt/tradfri"
	"github.com/moroen/tradfri2mqtt/webinterface"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		do_serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveCmd.Flags().StringP("server-root", "s", "./www", "Location of index.html")
	serveCmd.Flags().IntP("server-port", "p", 8321, "Web server port")

	viper.BindPFlag("interface.root", serveCmd.Flags().Lookup("server-root"))
	viper.BindPFlag("interface.port", serveCmd.Flags().Lookup("server-port"))
}

var status_channel chan (error)

func do_serve() {

	var wg sync.WaitGroup

	// log.SetReportCaller(true)
	// coap.SetCoapRetry(2, 1)

	var err error

	status_channel = make(chan error)

	if viper.GetBool("interface.enable") {
		go webinterface.Interface_Server(viper.GetString("interface.root"), status_channel)
	}

	for err == nil {

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		if viper.GetBool("mqtt.enable") {
			tradfri.MQTTSendTopic = mqttclient.SendTopic
			go mqttclient.Start(&wg, status_channel)
		} else {
			log.Info("MQTT - Disabled by config")
		}

		if viper.GetBool("tradfri.enable") {
			go tradfri.Start(&wg, status_channel)
		} else {
			log.Info("Tradfri - Disabled by config")
		}

		select {
		case <-c:
			log.Debug("Sig catched")

			go tradfri.Stop()
			go mqttclient.Stop()

			wg.Wait()

			os.Exit(1)
		case err = <-status_channel:
			if err == settings.ErrConfigIsDirty {
				log.Info("Config has changed, restarting")
				go tradfri.Stop()
				go mqttclient.Stop()
				wg.Wait()
			}
			err = nil
		}
	}

	// time.Sleep(5 * time.Second)
	// coap.ObserveRestart(false)
	// tradfri.Stop()
	// wg.Wait()
}
