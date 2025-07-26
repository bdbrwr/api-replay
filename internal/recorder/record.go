package recorder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bdbrwr/api-replay/internal/cliutils"
	"github.com/bdbrwr/api-replay/internal/config"
	"github.com/spf13/cobra"
)

func NewCommand(cfg *config.Config) *cobra.Command {
	var urlFlag, outputFlag string

	cmd := &cobra.Command{
		Use:   "record [url] [output-path]",
		Short: "Record an API response and save it to an output file",
		Args:  cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := os.Stat(cfg.Dir)
			if err == nil && !info.IsDir() {
				return fmt.Errorf("output_dir %q exists and is not a directory", cfg.Dir)
			}
			if os.IsNotExist(err) {
				if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
					return fmt.Errorf("creating output_dir %q: %w", cfg.Dir, err)
				}
			}
			url, err := cliutils.GetArgOrFlag(cmd, args, "url", 0, "URL to fetch")
			if err != nil {
				return err
			}
			output, err := cliutils.GetArgOrFlag(cmd, args, "output", 1, "output file path")
			if err != nil {
				return err
			}

			outputPath := filepath.Join(cfg.Dir, output)
			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				return fmt.Errorf("creating output directory: %w", err)
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return fmt.Errorf("creating HTTP request: %w", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("executing HTTP request: %w", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("reading response body: %w", err)
			}

			cached := map[string]any{
				"status":  resp.StatusCode,
				"headers": resp.Header,
				"body":    json.RawMessage(body),
			}

			file, err := os.Create(outputPath)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer file.Close()

			enc := json.NewEncoder(file)
			enc.SetIndent("", "  ")
			if err := enc.Encode(cached); err != nil {
				return fmt.Errorf("encoding response: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&urlFlag, "url", "U", "", "URL to fetch")
	cmd.Flags().StringVarP(&outputFlag, "output", "O", "", "Path to save the response (relative to output_dir)")

	return cmd
}
