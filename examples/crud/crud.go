package crud

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/caicloud/aloe"
	"github.com/caicloud/aloe/examples/crud/server"
	"github.com/onsi/ginkgo"
)

func init() {
	aloe.Init(nil)
}

var s = server.NewProductServer()

// RunTEST runs crud test
func RunTEST(t *testing.T) {
	aloe.AppendDataDirs("testdata")
	if err := aloe.Env("host", "localhost:8080"); err != nil {
		fmt.Printf("can't set env host: %v", err)
		os.Exit(1)
	}
	if err := aloe.RegisterCleaner(s); err != nil {
		fmt.Printf("can't register cleaner: %v", err)
		os.Exit(1)
	}
	aloe.Run(t)
}

var _ = ginkgo.BeforeSuite(func() {
	s.Register()
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("can't begin server")
		os.Exit(1)
	}
	go http.Serve(listener, nil)
})
