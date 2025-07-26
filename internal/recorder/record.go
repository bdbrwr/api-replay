package recorder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bdbrwr/api-replay/internal/cliutils"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var url, output string

	cmd := &cobra.Command{
		Use:   "record [url] [output-path]",
		Short: "Record an API reponse and save it to an output file",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			url, err := cliutils.GetArgOrFlag(cmd, args, "url", 0, "URL to fetch")
			if err != nil {
				return err
			}
			output, err := cliutils.GetArgOrFlag(cmd, args, "output", 1, "output file path")
			if err != nil {
				return err
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return fmt.Errorf("failed creating HTTP request: %w", err)
			}

			httpClient := &http.Client{}

			resp, err := httpClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed executing HTTP request: %w", err)
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed reading response body: %w", err)
			}

			cached := map[string]any{
				"status":  resp.StatusCode,
				"headers": resp.Header,
				"body":    json.RawMessage(data),
			}

			file, err := os.Create(output)
			if err != nil {
				return fmt.Errorf("failed creating output file: %w", err)
			}
			defer file.Close()

			enc := json.NewEncoder(file)
			enc.SetIndent("", " ") // make json readable by humans
			if err := enc.Encode(cached); err != nil {
				return fmt.Errorf("failed encoding JSON to output: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&url, "url", "U", "", "URL to fetch")
	cmd.Flags().StringVarP(&output, "output", "O", "", "Path to save the response")
	cmd.MarkFlagRequired("url")
	cmd.MarkFlagRequired("output")

	return cmd
}
