package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const DefaultUrl = "https://golang.org/"

func main() {
	http.HandleFunc("/", shortUrlProcessing)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func shortUrlProcessing(w http.ResponseWriter, r *http.Request) {
	var paths interface{}
	jsonData := []byte(`
    {
        "paths" : {
            "/go-http": "https://golang.org/pkg/net/http/",
            "/go-gophers" : "https://github.com/shalakhin/gophericons/blob/master/preview.jpg"
        }
        
    }`)
	err := json.Unmarshal(jsonData, &paths)

	if err != nil {
		log.Println(err)
	}

	var url = r.URL.Path
	var redirectUrl = DefaultUrl
	m := paths.(map[string]interface{})
	for key, value := range m["paths"].(map[string]interface{}) {
		if key == url{
			redirectUrl = value.(string)
			break
		}
	}

	redirect(w, r, redirectUrl)
}

func redirect(w http.ResponseWriter, r *http.Request, redirectUrl string) {
	http.Redirect(w, r, redirectUrl, 301)
}
