package comm

import (
	"testing"
)

func TestLogger(t *testing.T) {
	Config("E:\\Logs\\goTransferLog", 3)
	Debugf("test log debug")
	Debugf("test log debug", "ccc")
	Infof("Love", "cccc")
	Warnf("gogo", "cccc")
}
