package cmd

import (
    mapper2 "github.com/Eldius/mock-server-go/internal/mapper"
    "os"

    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
    Use:   "template",
    Short: "Outputs a template for YAML config file",
    Long:  `Outputs a template for YAML config file.`,
    Run: func(cmd *cobra.Command, args []string) {
        r := mapper2.Router{}
        r.Add(mapper2.RequestMapping{
            Path:   "/v1/contract",
            Method: "POST",
            Response: mapper2.MockResponse{
                Headers: mapper2.MockHeader{
                    "Content-Type": "application/json",
                },
                StatusCode: 202,
            },
        })
        body1 := `{"id": 123, "name": "My Contract"}`
        r.Add(mapper2.RequestMapping{
            Path:   "/v1/contract",
            Method: "GET",
            Response: mapper2.MockResponse{
                Headers: mapper2.MockHeader{
                    "Content-Type": "application/json",
                },
                StatusCode: 200,
                Body:       &body1,
            },
        })
        script2 := `
var res = {
	"body": "{\"id\": 1, \"status\": \"OK\"}",
	code: 200
};
`
        r.Add(mapper2.RequestMapping{
            Path:   "/v1/contract/1",
            Method: "GET",
            Response: mapper2.MockResponse{
                Headers: mapper2.MockHeader{
                    "Content-Type": "application/json",
                },
                StatusCode: 200,
                Script:     &script2,
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
