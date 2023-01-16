package main

import (
	"github.com/fernandosanchezjr/gosdr/config"
	"github.com/sirupsen/logrus"
)

func main() {
	config.SetLogLevel(logrus.ErrorLevel)
	config.ParseFlags()
	handleInvocation()
}
