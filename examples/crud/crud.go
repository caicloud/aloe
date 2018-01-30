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
	"github.com/onsi/gomega"
)

// RunTEST runs crud test
func RunTEST(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	f := aloe.NewFramework("localhost:8080", cleanUp,
		"testdata",
	)
	if err := f.Run(); err != nil {
		fmt.Printf("can't run framework: %v", err)
		os.Exit(1)
	}
	ginkgo.RunSpecs(t, "CRUD Suite")
}

func cleanUp() {
	// clean up databases
}

var _ = ginkgo.BeforeSuite(func() {
	server.Product{}.Register()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("unable to listen on 8080")
	}
	go http.Serve(listener, nil)
})
