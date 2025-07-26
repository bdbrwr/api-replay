package main

import (
	"log"

	"github.com/bdbrwr/api-replay/internal/config"
	"github.com/bdbrwr/api-replay/internal/recorder"
	"github.com/bdbrwr/api-replay/internal/replayer"
	"github.com/bdbrwr/api-replay/internal/server"
	"github.com/spf13/cobra"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed loading config: %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "api-replay",
		Short: "API Replay is a lightweight tool that helps to record and replay API calls for demo projects",
	}

	rootCmd.AddCommand(recorder.NewCommand(cfg))
	rootCmd.AddCommand(replayer.NewCommand(cfg))
	rootCmd.AddCommand(server.NewCommand(cfg))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
