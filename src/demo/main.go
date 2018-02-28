//go:generate go:generate go-bindata -pkg $GOPACKAGE -o assets.go assets/

package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
)

var (
	addr string
	name string
)

func main() {
	flag.StringVar(&addr, "addr", ":8888", "Address to listen on (default ':8888')")
	flag.StringVar(&name, "name", "demo1", "Server name")
	flag.Parse()

	http.HandleFunc("/", Index)

	fmt.Printf("Starting %s listening on %s\n", name, addr)
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
