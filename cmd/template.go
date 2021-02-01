package cmd

import (
	"os"

	"github.com/Eldius/mock-server-go/mapper"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Outputs a template for YAML config file",
	Long:  `Outputs a template for YAML config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		r := mapper.Router{}
		r.Add(mapper.RequestMapping{
			Path:   "/v1/contract",
			Method: "POST",
			Response: mapper.MockResponse{
				Headers: mapper.MockHeader{
					"Content-Type": []string{"application/json"},
				},
				StatusCode: 202,
			},
		})
		r.Add(mapper.RequestMapping{
			Path:   "/v1/contract",
			Method: "GET",
			Response: mapper.MockResponse{
				Headers: mapper.MockHeader{
					"Content-Type": []string{"application/json"},
				},
				StatusCode: 200,
				Body:       `{"id": 123, "name": "My Contract"}`,
			},
		})

		_ = yaml.NewEncoder(os.Stdout).Encode(r)
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
