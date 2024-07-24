package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"mockambo/oapi"
	"mockambo/server"
)

// specFilePath is the path to the specification file
var specFilePath string

var mergerFilePath string

// port is the port the server should run on
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
		doc, err := oapi.NewDoc(specFilePath, mergerFilePath)
		if err != nil {
			log.Fatalln("The specified file `", specFilePath, "` is not a valid OpenAPI 3 specification")
		}
		if err := doc.Watch(); err != nil {
			log.Fatalln(err)
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
	RunCmd.PersistentFlags().StringVarP(&mergerFilePath, "merger", "m", "", "path to a YAML file to be merged to the original OpenAPI file")
}
