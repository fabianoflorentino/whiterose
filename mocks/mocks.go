package mocks

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type MockExecutor struct {
	RunFunc   func(cmd string, args ...string) (string, error)
	WhichFunc func(cmd string) (string, error)
}

func (m *MockExecutor) Run(cmd string, args ...string) (string, error) {
	if m.RunFunc != nil {
		return m.RunFunc(cmd, args...)
	}
	return "", nil
}

func (m *MockExecutor) Which(cmd string) (string, error) {
	if m.WhichFunc != nil {
		return m.WhichFunc(cmd)
	}
	return "/usr/bin/" + cmd, nil
}

type RealExecutor struct{}

func (e *RealExecutor) Run(cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	return string(out), err
}

func (e *RealExecutor) Which(cmd string) (string, error) {
	c := exec.Command("which", cmd)
	out, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

type MockFS struct {
	Files map[string][]byte
}

func NewMockFS() *MockFS {
	return &MockFS{Files: make(map[string][]byte)}
}

func (m *MockFS) Read(path string) ([]byte, error) {
	if data, ok := m.Files[path]; ok {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFS) Write(path string, data []byte) error {
	m.Files[path] = data
	return nil
}

func (m *MockFS) Exists(path string) bool {
	_, ok := m.Files[path]
	return ok
}

func (m *MockFS) MkdirAll(path string) error {
	m.Files[path+"/.dir"] = []byte{}
	return nil
}

type MockLogger struct {
	Logs []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{Logs: []string{}}
}

func (m *MockLogger) Print(v ...interface{}) {
	m.Logs = append(m.Logs, fmt.Sprint(v...))
}

func (m *MockLogger) Println(v ...interface{}) {
	m.Logs = append(m.Logs, fmt.Sprintln(v...))
}

func (m *MockLogger) Printf(format string, v ...interface{}) {
	m.Logs = append(m.Logs, fmt.Sprintf(format, v...))
}

type MockDockerClient struct {
	BuildFunc  func(dockerfile, image string, args []string) error
	DeleteFunc func(image string) error
	ListFunc   func(pattern string) ([]string, error)
}

func (m *MockDockerClient) Build(dockerfile, image string, args []string) error {
	if m.BuildFunc != nil {
		return m.BuildFunc(dockerfile, image, args)
	}
	return nil
}

func (m *MockDockerClient) Delete(image string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(image)
	}
	return nil
}

func (m *MockDockerClient) ListImages(pattern string) ([]string, error) {
	if m.ListFunc != nil {
		return m.ListFunc(pattern)
	}
	return nil, nil
}

type MockHTTPClient struct {
	GetFunc func(url string) (*http.Response, error)
	DoFunc  func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if m.GetFunc != nil {
		return m.GetFunc(url)
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}, nil
}

type MockCommandExecutor struct {
	RunFunc func(cmd string, args ...string) (string, error)
}

func (m *MockCommandExecutor) Run(cmd string, args ...string) (string, error) {
	if m.RunFunc != nil {
		return m.RunFunc(cmd, args...)
	}
	return "", nil
}

type MockVersionLister struct {
	GoVersionsFunc     func() ([]string, error)
	PackageUpdatesFunc func(path string) ([]string, error)
}

func (m *MockVersionLister) ListGoVersions() ([]string, error) {
	if m.GoVersionsFunc != nil {
		return m.GoVersionsFunc()
	}
	return []string{"1.26.0", "1.25.0", "1.24.0"}, nil
}

func (m *MockVersionLister) ListPackageUpdates(path string) ([]string, error) {
	if m.PackageUpdatesFunc != nil {
		return m.PackageUpdatesFunc(path)
	}
	return []string{}, nil
}