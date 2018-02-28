package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	stdin  bytes.Buffer
	stdout bytes.Buffer
	stderr bytes.Buffer
	runner Runner
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupTest() {
	s.stdin.Reset()
	s.stdout.Reset()
	s.stderr.Reset()
	s.runner = Runner{
		stdin:  &s.stdin,
		stdout: &s.stdout,
		stderr: &s.stderr,
		exit: func(c int) {
			panic(fmt.Sprintf("exit(%d)", c))
		},
	}
}

func (s *Suite) TestCheckMissingSourceEndpoint() {

	s.stdin.WriteString(`{}`)

	assert.PanicsWithValue(s.T(), "exit(1)", func() {
		s.runner.Exec("check")
	})

	assert.Equal(s.T(), s.stdout.String(), "")
	assert.Contains(s.T(), s.stderr.String(), "Missing source.endpoint")
}

func (s *Suite) TestCheckInvalidJSON() {
	s.stdin.WriteString(`{"invalid"}`)

	assert.PanicsWithValue(s.T(), "exit(1)", func() {
		s.runner.Exec("check")
	})

	assert.Contains(s.T(), s.stderr.String(), "Invalid JSON")
}

func (s *Suite) TestCheckInvalidSourceEndpoint() {
	s.stdin.WriteString(`{"source":{"endpoint":"__invalid__"}}`)

	assert.PanicsWithValue(s.T(), "exit(1)", func() {
		s.runner.Exec("check")
	})

	assert.Contains(s.T(), s.stderr.String(), "failed to check __invalid__")
}

func (s *Suite) TestCheckValidEndpointInvalidResponse() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `invalid json`)
	}))
	defer ts.Close()

	fmt.Fprintf(&s.stdin, `{"source":{"endpoint": "%s"}}`, ts.URL)

	assert.PanicsWithValue(s.T(), "exit(1)", func() {
		s.runner.Exec("check")
	})

	assert.Equal(s.T(), s.stdout.String(), "")
	assert.Contains(s.T(), s.stderr.String(), "failed to parse")
}

func (s *Suite) TestCheckValidEndpointValidResponse() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"commit":"3e7b1416cf4198ebefad700e45f48315dcd44652"}`)
	}))
	defer ts.Close()

	fmt.Fprintf(&s.stdin, `{"source":{"endpoint": "%s"}}`, ts.URL)

	s.runner.Exec("check")

	assert.Contains(s.T(), s.stdout.String(), `{"ref":"3e7b1416cf4198ebefad700e45f48315dcd44652"}`)
	assert.Contains(s.T(), s.stderr.String(), "checking endpoint")
}

func (s *Suite) TestCheckWith500Service() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		fmt.Fprintf(w, "500 error")
	}))
	defer ts.Close()

	fmt.Fprintf(&s.stdin, `{
  "source":{"endpoint": "%s"},
  "version":{"ref":"3e7b1416cf4198ebefad700e45f48315dcd44652"}
}`, ts.URL)

	assert.PanicsWithValue(s.T(), "exit(1)", func() {
		s.runner.Exec("check")
	})

	assert.Contains(s.T(), s.stdout.String(), ``)
	assert.Contains(s.T(), s.stderr.String(), "checking endpoint")
}

func (s *Suite) TestCheckValidEndpointValidResponseSecondCommit() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
  "commit":"a7361d758952cd6028602e268ce302fe82fc7b5a"
}`)
	}))
	defer ts.Close()

	fmt.Fprintf(&s.stdin, `{
  "source":{"endpoint": "%s"},
  "version":{"ref":"3e7b1416cf4198ebefad700e45f48315dcd44652"}
}`, ts.URL)

	s.runner.Exec("check")

	assert.Contains(s.T(), s.stdout.String(), `[{"ref":"3e7b1416cf4198ebefad700e45f48315dcd44652"},{"ref":"a7361d758952cd6028602e268ce302fe82fc7b5a"}]`)
	assert.Contains(s.T(), s.stderr.String(), "checking endpoint")
}

func (s *Suite) TestCheckValidEndpointValidResponseSameCommit() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
  "commit":"3e7b1416cf4198ebefad700e45f48315dcd44652"
}`)
	}))
	defer ts.Close()

	fmt.Fprintf(&s.stdin, `{
  "source":{"endpoint": "%s"},
  "version":{"ref":"3e7b1416cf4198ebefad700e45f48315dcd44652"}
}`, ts.URL)

	s.runner.Exec("check")

	assert.Contains(s.T(), s.stdout.String(), `[{"ref":"3e7b1416cf4198ebefad700e45f48315dcd44652"}]`)
	assert.Contains(s.T(), s.stderr.String(), "checking endpoint")
}

func (s *Suite) TestInVersion() {

	s.stdin.WriteString(`{
  "version":{"ref":"3e7b1416cf4198ebefad700e45f48315dcd44652"}
}`)

	dir, err := ioutil.TempDir("", "version-resource")
	if !assert.NoError(s.T(), err) {
		return
	}
	s.runner.Exec("in", dir)

	assert.Contains(s.T(), s.stdout.String(), `{"version":{"ref":"3e7b1416cf4198ebefad700e45f48315dcd44652"}}`)
	assert.Equal(s.T(), s.stderr.String(), "")

	ref, err := ioutil.ReadFile(filepath.Join(dir, "ref"))
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), string(ref), "3e7b1416cf4198ebefad700e45f48315dcd44652")
}

func (s *Suite) TestOutUnsupported() {
	assert.PanicsWithValue(s.T(), "exit(1)", func() {
		s.runner.Exec("out")
	})

	assert.Contains(s.T(), s.stderr.String(), "out not supported")
}
