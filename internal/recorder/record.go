package recorder

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var url, output string

	cmd := &cobra.Command{
		Use:   "record",
		Short: "Record an API reponse and save it to an output file",
		RunE: func(cmd *cobra.Command, args []string) error {

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}

			httpClient := &http.Client{}

			resp, err := httpClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			cached := map[string]interface{}{
				"status":  resp.StatusCode,
				"headers": resp.Header,
				"body":    json.RawMessage(data),
			}

			file, err := os.Create(output)
			if err != nil {
				return err
			}
			defer file.Close()

			enc := json.NewEncoder(file)
			enc.SetIndent("", " ") // make json readable by humans
			return enc.Encode(cached)
		},
	}

	cmd.Flags().StringVarP(&url, "url", "U", "", "URL to fetch")
	cmd.Flags().StringVarP(&output, "output", "O", "", "Path to save the response")
	cmd.MarkFlagRequired("url")
	cmd.MarkFlagRequired("output")

	return cmd
}
