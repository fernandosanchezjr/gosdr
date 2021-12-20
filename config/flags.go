package config

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

func printDefaults() {
	flag.PrintDefaults()
	os.Exit(1)
}

func ParseFlags(handlers ...func() error) {
	flag.Parse()
	if flag.Parsed() {
		for _, handler := range handlers {
			if err := handler(); err != nil {
				log.WithError(err).Error("Argument error")
				printDefaults()
			}
		}
	} else {
		printDefaults()
	}
}
