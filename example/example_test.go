package example

import (
	"fmt"
	"os"

	"github.com/caicloud/aloe"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"testing"
)

func TestAPI(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	f := framework.NewFramework("https://api.github.com/", cleanUp,
		"testdata",
	)
	if err := f.Run(); err != nil {
		fmt.Printf("can't run framework: %v", err)
		os.Exit(1)
	}
	ginkgo.RunSpecs(t, "API Suite")
}

func cleanUp() {

}

var _ = ginkgo.BeforeSuite(func() {

})
