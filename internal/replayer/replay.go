package replayer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bdbrwr/api-replay/internal/cliutils"
	"github.com/spf13/cobra"
)

type CachedResponse struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    json.RawMessage     `json:"body"`
}

func NewCommand() *cobra.Command {
	var inputFlag string

	cmd := &cobra.Command{
		Use:   "replay [input file]",
		Short: "Replay a cached API response from a file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile, err := cliutils.GetArgOrFlag(cmd, args, "input", 0, "input file path")
			if err != nil {
				return err
			}

			file, err := os.Open(inputFile)
			if err != nil {
				return fmt.Errorf("failed opening input file: %w", err)
			}
			defer file.Close()

			var resp CachedResponse
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&resp); err != nil {
				return fmt.Errorf("failed decoding cached response: %w", err)
			}

			fmt.Printf("Status: %d\n", resp.Status)

			fmt.Println("Headers:")
			for k, v := range resp.Headers {
				fmt.Printf("  %s: %v\n", k, v)
			}

			fmt.Println("Body:")
			var formattedBody map[string]any
			if err := json.Unmarshal(resp.Body, &formattedBody); err == nil {
				pretty, _ := json.MarshalIndent(formattedBody, "", "  ")
				fmt.Println(string(pretty))
			} else {
				fmt.Println(string(resp.Body))
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&inputFlag, "input", "I", "", "Path to cached response file")

	return cmd
}
