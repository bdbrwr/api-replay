package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bdbrwr/api-replay/internal/cliutils"
	"github.com/bdbrwr/api-replay/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

type CachedResponse struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    json.RawMessage     `json:"body"`
}

func NewCommand(cfg *config.Config) *cobra.Command {
	var dir, port string

	cmd := &cobra.Command{
		Use:   "serve [dir]",
		Short: "Serve cached responses over HTTP",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dirPath, err := cliutils.GetArgOrFlag(cmd, args, "from-dir", 0, "directory to serve")
			if err != nil || dirPath == "" {
				dirPath = cfg.Dir
			}

			router := chi.NewRouter()

			err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil || d.IsDir() || filepath.Ext(path) != ".json" {
					return nil
				}

				relPath, err := filepath.Rel(dirPath, path)
				if err != nil {
					return fmt.Errorf("failed to calculate relative path: %w", err)
				}

				raw := strings.TrimSuffix(filepath.ToSlash(relPath), ".json")

				var routePath string
				var expectedQuery string

				if at := strings.Index(raw, "@"); at != -1 {
					decodedQuery, err := url.QueryUnescape(raw[at+1:])
					if err != nil {
						return fmt.Errorf("failed to decode query from %q: %w", raw, err)
					}
					routePath = "/" + raw[:at]
					expectedQuery = decodedQuery
				} else {
					routePath = "/" + raw
				}

				handlerPath := path
				router.Get(routePath, func(w http.ResponseWriter, r *http.Request) {
					// If we have expected query parameters, validate them
					if expectedQuery != "" {
						if r.URL.RawQuery != expectedQuery {
							http.NotFound(w, r)
							return
						}
					}

					file, err := os.Open(handlerPath)
					if err != nil {
						http.Error(w, "failed opening cached file", http.StatusInternalServerError)
						return
					}
					defer file.Close()

					var resp CachedResponse
					if err := json.NewDecoder(file).Decode(&resp); err != nil {
						http.Error(w, "invalid cached file", http.StatusInternalServerError)
						return
					}

					for k, v := range resp.Headers {
						for _, h := range v {
							w.Header().Add(k, h)
						}
					}
					w.WriteHeader(resp.Status)
					w.Write(resp.Body)
				})

				displayPath := routePath
				if expectedQuery != "" {
					displayPath += "?" + expectedQuery
				}
				fmt.Printf("→ %s mapped to %s\n", displayPath, path)
				return nil
			})

			if err != nil {
				return fmt.Errorf("failed walking cache directory: %w", err)
			}

			fmt.Printf("Serving cached responses from %s on http://localhost:%s\n", dirPath, port)
			return http.ListenAndServe(":"+port, router)
		},
	}

	cmd.Flags().StringVarP(&dir, "from-dir", "D", "", "Directory containing cached responses")
	cmd.Flags().StringVarP(&port, "port", "P", cfg.Port, "Port to serve on")

	return cmd
}
