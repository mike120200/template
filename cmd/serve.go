/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"my_template/cmn"

	"github.com/spf13/cobra"
)

var (
	configFile string
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
		fmt.Println("serve called")

		// 初始化配置文件
		if err := cmn.ViperInit(configFile); err != nil {
			msg := fmt.Sprintf("初始化配置文件失败:%s", err.Error())
			panic(msg)
		}

		//初始化日志
		if err := cmn.LoggerInit(); err != nil {
			msg := fmt.Sprintf("初始化日志失败:%s", err.Error())
			panic(msg)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&configFile, "config", ".conf_linux.json", "config file (default is .conf_linux.json)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
