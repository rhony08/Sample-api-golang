package connection

import (
	"errors"
	"time"
)

type DbInterfaceConnection interface {
	SaveLogData(typeLog string, logMessage string, logTime time.Time) error
}

type dbConnection struct{}

func newDBConection(hostDBURL string) DbInterfaceConnection {
	return &dbConnection{}
}

func (outcallDB *dbConnection) SaveLogData(typeLog string, logMessage string, logTime time.Time) error {
	return errors.New("Unimplemented")
}
