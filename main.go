package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Api struct {
	Write    bool `json:write`
	Packages []struct {
		Url           string  `json:"url"`
		Last_modified float32 `json:last_modified`
		Name          string  `json:name`
		Version       string  `json:version`
		Filename      string  `json:filename`
	} `json:packages`
}

func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.RequestURI, "/")
	name := parts[2]
	url := fmt.Sprintf("%s/api/package/%s/", baseurl, name)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(raw, &api)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, p := range api.Packages {
		versions = append(versions, p.Version)
	}

	jsonBody, err := json.Marshal(versions)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBody)
}

var (
	bind     string
	baseurl  string
	versions []string
	api      Api
)

func init() {
	flag.StringVar(&bind, "listen", "127.0.0.1:8001", "bind to")
	flag.StringVar(&baseurl, "url", "", "base url (with credentials if any) for Pypicloud")

	flag.Parse()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/package/{package}/", ProxyRequest).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(bind, nil)
}
