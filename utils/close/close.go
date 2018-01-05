package close

import (
	"io"

	"github.com/golang/glog"
)

// Level defines log level if error occurred
type Level int

const (
	// INFO level
	INFO Level = iota
	// WARNING level
	WARNING
	// FATAL level
	FATAL
)

// Close will deal with error if io.Closer returns an error
// It will log in warning level defaultly
func Close(c io.Closer) {
	WARNING.Close(c)
}

// Close will call io.Closer with log level
func (v Level) Close(c io.Closer) {
	if err := c.Close(); err != nil {
		switch v {
		case INFO:
			glog.Infof("close error: %v", err)
		case WARNING:
			glog.Warningf("close error: %v", err)
		case FATAL:
			glog.Fatalf("close error: %v", err)
		}
	}
}
