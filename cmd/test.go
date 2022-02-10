/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type MyConfig struct {
	Tradfri struct {
		Gateway string `mapstructure:"gateway"`
	}
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// var conf settings.Config

		// fmt.Println(viper.AllSettings())

		/*
			if err := viper.Unmarshal(&conf); err == nil {
				fmt.Printf("%+v\n", conf)
			} else {
				fmt.Println(err.Error())
			}
		*/
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	testCmd.Flags().Int("port", 1138, "Port to run Application server on")
	viper.BindPFlag("mqtt.port", testCmd.Flags().Lookup("port"))
}
