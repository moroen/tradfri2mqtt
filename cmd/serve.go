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
	"fmt"
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

var _server_port int
var _server_root string

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
	cobra.OnInitialize(setConf)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveCmd.Flags().StringP("server-root", "s", "", "Location of index.html")
	serveCmd.Flags().IntP("server-port", "p", 0, "Web server port")

	// viper.BindPFlag("interface.root", serveCmd.Flags().Lookup("server-root"))
	// viper.BindPFlag("interface.port-cmd", serveCmd.Flags().Lookup("server-port"))
}

func setConf() {
	if port, err := serveCmd.Flags().GetInt("server-port"); err == nil {
		if port == 0 {
			_server_port = viper.GetInt("interface.port")
		} else {
			_server_port = port
		}
	} else {
		fmt.Println(err.Error())
	}

	if root, err := serveCmd.Flags().GetString("server-root"); err == nil {
		if root == "" {
			_server_root = viper.GetString("interface.root")
		} else {
			_server_root = root
		}
	} else {
		fmt.Println(err.Error())
	}
}

var status_channel chan (error)

func do_serve() {

	var wg sync.WaitGroup

	// log.SetReportCaller(true)
	// coap.SetCoapRetry(2, 1)

	var err error

	status_channel = make(chan error)

	if viper.GetBool("interface.enable") {
		go webinterface.Interface_Server(_server_root, _server_port, status_channel)
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
