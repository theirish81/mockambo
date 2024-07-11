package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"mockambo/oapi"
	"mockambo/server"
	"os"
)

var specFilePath string
var port int

var RootCmd = &cobra.Command{
	Use:   "mockambo",
	Short: "Mockambo is an OpenAPI-based REST API mocking system with gateway, recording, and testing capabilities",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the Mockambo server",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(specFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		doc, err := oapi.NewDoc(data)
		if err != nil {
			log.Fatalln("The specified file `", specFilePath, "` is not a valid OpenAPI 3 specification")
		}
		sx := server.NewServer(port, doc)
		if err := sx.Run(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(RunCmd)
	RunCmd.PersistentFlags().StringVarP(&specFilePath, "spec", "s", "", "path to an OpenAPI 3 specification file")
	_ = RunCmd.MarkPersistentFlagRequired("spec")
	RunCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "the port the mocking server should use")
}
