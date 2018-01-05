package crud

import (
	"fmt"
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
	go http.ListenAndServe(":8080", nil)
})
