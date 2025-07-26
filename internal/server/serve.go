package server

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bdbrwr/api-replay/internal/cliutils"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
)

type CachedResponse struct {
	Status  int                 `json:"status"`
	Headers map[string][]string `json:"headers"`
	Body    json.RawMessage     `json:"body"`
}

func NewCommand() *cobra.Command {
	var dir, port string

	cmd := &cobra.Command{
		Use:   "serve [dir]",
		Short: "Serve cached responses over HTTP",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			dirPath, err := cliutils.GetArgOrFlag(cmd, args, "from-dir", 0, "directory to serve")
			if err != nil {
				return fmt.Errorf("resolving directory: %w", err)
			}

			router := chi.NewRouter()

			err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, walkErr error) error {
				if walkErr != nil || d.IsDir() || filepath.Ext(path) != ".json" {
					return nil
				}

				routePath := "/" + filepath.Base(path)

				router.Get(routePath, func(w http.ResponseWriter, r *http.Request) {
					file, err := os.Open(path)
					if err != nil {
						http.Error(w, "Failed to open cached file", http.StatusInternalServerError)
						return
					}
					defer file.Close()

					var resp CachedResponse
					if err := json.NewDecoder(file).Decode(&resp); err != nil {
						http.Error(w, "Invalid cached file", http.StatusInternalServerError)
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

				fmt.Printf("→ %s mapped to %s\n", routePath, path)
				return nil
			})

			if err != nil {
				return fmt.Errorf("walking cache directory: %w", err)
			}

			fmt.Printf("Serving cached responses from %s on http://localhost:%s\n", dirPath, port)
			return http.ListenAndServe(":"+port, router)
		},
	}

	cmd.Flags().StringVarP(&dir, "from-dir", "D", "", "Directory containing cached responses")
	cmd.Flags().StringVarP(&port, "port", "P", "8080", "Port to serve on")

	return cmd
}
