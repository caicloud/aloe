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

var s = server.NewProductServer()

// RunTEST runs crud test
func RunTEST(t *testing.T) {
	f := aloe.NewFramework("localhost:8080",
		"testdata",
	)
	if err := f.RegisterCleaner(s); err != nil {
		fmt.Printf("can't register cleaner: %v", err)
		os.Exit(1)
	}
	f.Run(t)
}

var _ = ginkgo.BeforeSuite(func() {
	s.Register()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("can't begin server")
		os.Exit(1)
	}
	go http.Serve(listener, nil)
})
