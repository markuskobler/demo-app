//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -prefix assets assets/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
)

var (
	commit = "DEV"

	addr = ":8888"
	name = "demo1"
)

func main() {
	flag.StringVar(&addr, "addr", addr, "Address to listen on (default ':8888')")
	flag.StringVar(&name, "name", name, "Server name")
	flag.Parse()

	http.HandleFunc("/", Index)
	http.HandleFunc("/version", Version)

	fmt.Printf("%s (%s)\n", name, commit)
	fmt.Printf("Listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}

var index *template.Template = template.Must(template.New("index").Parse(string(MustAsset("index.html"))))

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Header().Set("Set-Cookie", "JSESSIONID="+name)
	if refresh := r.FormValue("r"); refresh != "" {
		w.Header().Set("Refresh", refresh)
	}
	w.WriteHeader(http.StatusOK)

	index.Execute(w, struct {
		Name    string
		Headers http.Header
	}{
		Name:    name,
		Headers: r.Header,
	})
}

func Version(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	v := struct {
		Name   string `json:"name"`
		Commit string `json:"Commit"`
	}{
		Name:   name,
		Commit: commit,
	}
	json.NewEncoder(w).Encode(v)
}
