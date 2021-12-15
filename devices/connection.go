package devices

import log "github.com/sirupsen/logrus"

type Connection interface {
	Close() error
	IsOpen() bool
	Refresh() error
	Fields() log.Fields
	GetInfo() *Info
}
