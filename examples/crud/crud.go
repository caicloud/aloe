package crud

import (
	"fmt"
	"os"
	"testing"

	"github.com/caicloud/aloe"
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
