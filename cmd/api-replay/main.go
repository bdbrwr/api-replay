package main

import (
	"log"

	"github.com/bdbrwr/api-replay/internal/recorder"
	"github.com/bdbrwr/api-replay/internal/replayer"
	"github.com/bdbrwr/api-replay/internal/server"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "api-replay",
		Short: "API Replay is a lightweight tool that helps to record and replay API calls for demo projects",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	// Do Stuff here
		// }
	}

	rootCmd.AddCommand(recorder.NewCommand())
	rootCmd.AddCommand(replayer.NewCommand())
	rootCmd.AddCommand(server.NewCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
