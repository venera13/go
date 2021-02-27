package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const DefaultUrl = "https://golang.org/"
const Port = "8080"

func main() {
	startHttpServer()
}

func startHttpServer(){
	configFile := getConfigFile()
	var configFileContent interface{}
	if configFile != "" {
		configFileContent = getConfigFileContent(configFile)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		shortUrlProcessing(w, r, configFileContent)
	})
	err := http.ListenAndServe(":" + Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getConfigFile() string {
	configFile := flag.String("f", "", "Config file")
	flag.Parse()
	return *configFile
}

func getConfigFileContent(configFile string) interface{} {
	file, err := os.Open(configFile)
	if err != nil{
		fmt.Println(err)
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
	}

	var paths interface{}
	err = json.Unmarshal(fileContent, &paths)
	if err != nil {
		fmt.Println(err)
	}

	return paths
}

func shortUrlProcessing(w http.ResponseWriter, r *http.Request, paths interface{}) {
	var url = r.URL.Path
	var redirectUrl = DefaultUrl
	if paths != nil{
		m := paths.(map[string]interface{})
		for key, value := range m["paths"].(map[string]interface{}) {
			if key == url{
				redirectUrl = value.(string)
				break
			}
		}
	}

	redirect(w, r, redirectUrl)
}

func redirect(w http.ResponseWriter, r *http.Request, redirectUrl string) {
	http.Redirect(w, r, redirectUrl, 301)
}
