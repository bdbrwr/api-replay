package recorder

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bdbrwr/api-replay/internal/cliutils"
	"github.com/bdbrwr/api-replay/internal/config"
	"github.com/spf13/cobra"
)

func NewCommand(cfg *config.Config) *cobra.Command {
	var outputFlag, baseURLFlag string
	var headerFlags []string

	cmd := &cobra.Command{
		Use:   "record [url] [output-path]",
		Short: "Record an API response and save it to an output file",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := os.Stat(cfg.Dir)
			if err == nil && !info.IsDir() {
				return fmt.Errorf("output_dir %q exists and is not a directory", cfg.Dir)
			}
			if os.IsNotExist(err) {
				if err := os.MkdirAll(cfg.Dir, 0755); err != nil {
					return fmt.Errorf("failed creating output_dir %q: %w", cfg.Dir, err)
				}
			}

			urlStr, err := cliutils.GetArgOrFlag(cmd, args, "url", 0, "URL to fetch")
			if err != nil {
				return err
			}

			baseURL, _ := cmd.Flags().GetString("base-url")

			parsedURL, err := url.Parse(urlStr)
			if err != nil {
				return fmt.Errorf("invalid URL: %w", err)
			}

			relPath := parsedURL.Path
			if baseURL != "" {
				baseParsed, err := url.Parse(baseURL)
				if err != nil {
					return fmt.Errorf("invalid base-url: %w", err)
				}
				basePath := baseParsed.Path
				if basePath != "/" && strings.HasPrefix(relPath, basePath) {
					relPath = strings.TrimPrefix(relPath, basePath)
				}
			}

			if relPath == "" || relPath == "/" {
				relPath = "index"
			}

			if parsedURL.RawQuery != "" {
				encodedQuery := url.QueryEscape(parsedURL.RawQuery)
				relPath += "@" + encodedQuery
			}

			if filepath.Ext(relPath) != ".json" {
				relPath += ".json"
			}

			baseDir := cfg.Dir
			if outputFlag != "" {
				baseDir = outputFlag
			}

			outputPath := filepath.Join(baseDir, filepath.FromSlash(relPath))

			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				return fmt.Errorf("failed creating output directory: %w", err)
			}

			req, err := http.NewRequest("GET", urlStr, nil)
			if err != nil {
				return fmt.Errorf("failed creating HTTP request: %w", err)
			}

			for _, h := range headerFlags {
				parts := strings.SplitN(h, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid header format, expected 'Key: Value': %s", h)
				}
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				req.Header.Add(key, value)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("failed executing HTTP request: %w", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed reading response body: %w", err)
			}

			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("HTTP error: %d\nResponse body:\n%s", resp.StatusCode, body)
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

			fmt.Printf("Recorded %s -> %s\n", urlStr, outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFlag, "output", "O", "", "Override path to save the response (relative to dir)")
	cmd.Flags().StringVarP(&baseURLFlag, "base-url", "B", "", "Base URL to strip from request path when saving")
	cmd.Flags().StringSliceVarP(&headerFlags, "header", "H", nil, "Custom headers to include (eg -H 'Authorization: Bearer ...')")

	return cmd
}
