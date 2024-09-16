/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GaikwadPratik/signoztest/internal/appservice"
	"github.com/GaikwadPratik/signoztest/internal/entity"
	"github.com/GaikwadPratik/signoztest/internal/webserver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("server called")

		appConfig := entity.AppConfig{
			Webserver: entity.Webserver{
				Port: viper.GetInt("webserver.port"),
				Timeouts: entity.WebserverTimeouts{
					StartWait: viper.GetDuration("webserver.timeouts.startwait"),
					Graceful:  viper.GetDuration("webserver.timeouts.graceful"),
					Write:     viper.GetDuration("webserver.timeouts.write"),
					Read:      viper.GetDuration("webserver.timeouts.read"),
					Idle:      viper.GetDuration("webserver.timeouts.idle"),
				},
			},
		}

		slog.Info(
			"start config",
			slog.Any("config", appConfig),
		)

		quitChannel := make(chan os.Signal, 1)
		signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

		appCtx, cancelAppCtxFn := context.WithCancel(cmd.Context())
		defer cancelAppCtxFn()

		cleanup := initTracer()
		defer cleanup(appCtx)

		var err error

		srv := appservice.New()

		webserverConfig := webserver.WebserverDependencies{
			Conf: &webserver.WebserverConf{
				Port:         appConfig.Webserver.Port,
				StartWait:    appConfig.Webserver.Timeouts.StartWait,
				Graceful:     appConfig.Webserver.Timeouts.Graceful,
				WriteTimeout: appConfig.Webserver.Timeouts.Write,
				ReadTimeout:  appConfig.Webserver.Timeouts.Read,
				IdleTimeout:  appConfig.Webserver.Timeouts.Idle,
			},
			AppService: srv,
			LogLevel:   logLevel,
		}

		webSrvErrChan := make(chan error)
		go webserver.Initiate(appCtx, webserverConfig, webSrvErrChan)

		var ok bool
		err, ok = <-webSrvErrChan
		if ok && err != nil {
			quitChannel <- syscall.SIGINT

			return
		}

		//listen for interrupt
		<-quitChannel
		slog.Info("shutting down gracefully, press ctrl+c again to force")
		cancelAppCtxFn()

		time.Sleep(100 * time.Millisecond)

		slog.Info("Adios!")

		if err != nil {
			os.Exit(1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
