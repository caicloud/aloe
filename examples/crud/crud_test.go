package crud

import (
	"testing"

	"github.com/caicloud/aloe/examples/crud/server"
	"github.com/onsi/ginkgo"
)

func TestCRUD(t *testing.T) {
	RunTEST(t)
}

var _ = ginkgo.BeforeSuite(func() {
	stopCh := make(chan struct{})
	defer close(stopCh)
	go server.RunServer(stopCh)
})
