package main

import (
	"log"
	"mockambo/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
	//data, _ := os.ReadFile("test_data/github.yaml")
	//doc, err := oapi.NewDoc(data)
	//fmt.Println(err)
	//sx := server.NewServer(8080, doc)
	//if err := sx.Run(); err != nil {
	//	log.Panic(err)
	//}
}
