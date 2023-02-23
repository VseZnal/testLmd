package errors

import (
	"database/sql"
	"log"

	"google.golang.org/grpc/status"
)

func HandleDatabaseError(err error) error {
	log.Println(err)

	switch {
	case err == sql.ErrNoRows:
		return status.Error(404, "Record not found")
	case err == sql.ErrTxDone:
		return status.Error(500, "Internal server error has occured")
	case err == sql.ErrConnDone:
		return status.Error(500, "Internal server error has occured")
	default:
		return status.Error(500, "Internal server error has occured")
	}
}

func HandleFatalError(err error, message string) {
	log.Println(err)

	log.Fatal(message)
}

func LogError(err error) error {
	log.Println(err)

	return err
}
