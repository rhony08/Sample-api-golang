package connection

import "time"

const (
	DEBUG_TYPE   = "debug"
	ERROR_TYPE   = "error"
	WARNING_TYPE = "warning"
)

type ConnectionInterface interface {
	SearchMovieData(key string, pagination int) ([]Movie, error)
	GetDetailMovie(imdbID string) (Movie, error)
	SaveLogData(typeLog string, logMessage string, logTime time.Time) error
}

type OptionConnection struct {
	BaseURL string
	ApiKey  string
	HostDB  string
}

type Connection struct {
	omdbOutcall OMDBInterface
	dbOutcall   DbInterfaceConnection
}

func GetNewConnectionInterface(optionParam OptionConnection) ConnectionInterface {
	return &Connection{
		omdbOutcall: newOmdbConnection(optionParam.ApiKey, optionParam.BaseURL),
		dbOutcall:   newDBConection(optionParam.HostDB),
	}
}

func (outcallConnection *Connection) SaveLogData(typeLog string, logMessage string, logTime time.Time) error {
	return outcallConnection.dbOutcall.SaveLogData(typeLog, logMessage, logTime)
}

func (outcallConnection *Connection) GetDetailMovie(imdbID string) (Movie, error) {
	return outcallConnection.omdbOutcall.GetDetailMovie(imdbID)
}

func (outcallConnection *Connection) SearchMovieData(key string, pagination int) ([]Movie, error) {
	return outcallConnection.omdbOutcall.SearchMovieData(key, pagination)
}
