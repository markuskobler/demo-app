package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var cmd string

func init() {
	flag.StringVar(&cmd, "cmd", filepath.Base(os.Args[0]), "check, in, out")
}

func main() {
	flag.Parse()
	Runner{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		exit:   os.Exit,
	}.Exec(cmd, flag.Args()...)
}

type Runner struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	exit   func(code int)
}

func (r Runner) Exec(cmd string, args ...string) {
	switch cmd {
	case "check":

		var req CheckRequest

		r.decodeRequest(&req)

		resp := execCheck(&r, req)

		r.encodeResponse(resp)

	case "in":
		if len(args) != 1 {
			r.Failf("usage: in <destination>")
		}
		destination := args[0]

		var req InRequest

		r.decodeRequest(&req)

		resp := execIn(&r, req, destination)

		r.encodeResponse(&resp)

	case "out":
		r.Failf("out not supported")

	default:
		r.Failf("unexpected command %s; must be check, in", cmd)
	}
}

func (r *Runner) Logf(msg string, args ...interface{}) {
	fmt.Fprintf(r.stderr, msg, args...)
	fmt.Fprintln(r.stderr)
}

func (r *Runner) Failf(msg string, args ...interface{}) {
	r.Logf(msg, args...)
	r.exit(1)
}

func (r *Runner) decodeRequest(req interface{}) {
	err := json.NewDecoder(r.stdin).Decode(req)
	if err != nil {
		r.Failf("Invalid JSON request: %s", err)
	}
}

func (r *Runner) encodeResponse(resp interface{}) {
	err := json.NewEncoder(r.stdout).Encode(resp)
	if err != nil {
		r.Failf("Invalid JSON response: %s", err)
	}
}

type Version struct {
	Ref string `json:"ref"`
}

type CheckRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
}

type Source struct {
	Endpoint string `json:"endpoint"`
}

type InRequest struct {
	Version Version `json:"version"`
}

type InResponse struct {
	Version Version `json:"version"`
}

func execCheck(r *Runner, req CheckRequest) (resp []Version) {
	if req.Source.Endpoint == "" {
		r.Failf("Missing source.endpoint")
	}
	r.Logf("checking endpoint %s", req.Source.Endpoint)

	if req.Version.Ref != "" {
		resp = append(resp, req.Version)
	}

	endpointReq, _ := http.NewRequest("GET", req.Source.Endpoint, nil)
	endpointResp, err := http.DefaultClient.Do(endpointReq)
	if err != nil {
		r.Failf("failed to check %s: %s", req.Source.Endpoint, err)
	} else if endpointResp.StatusCode != 200 {
		r.Failf("failed to check %s: status %d", req.Source.Endpoint, endpointResp.StatusCode)
		return
	}

	var body struct {
		Name   string `json:"name"`
		Commit string `json:"commit"`
	}

	defer endpointResp.Body.Close()
	err = json.NewDecoder(endpointResp.Body).Decode(&body)
	if err != nil {
		r.Failf("failed to parse %s: %s", req.Source.Endpoint, err)
	} else if body.Commit == req.Version.Ref {
		return
	}

	return append(resp, Version{Ref: body.Commit})
}

func execIn(r *Runner, req InRequest, destination string) (resp InResponse) {
	resp.Version = req.Version
	ref := filepath.Join(destination, "ref")
	err := ioutil.WriteFile(ref, []byte(req.Version.Ref), 0666)
	if err != nil {
		r.Failf("Failed to write `ref`: %s", err)
	}
	return
}
