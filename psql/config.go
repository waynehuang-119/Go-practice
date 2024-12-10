package psql

import "fmt"

const (
	host     = "localhost"
	port     = 5432
	user     = "test"
	password = "test"
	dbname   = "testdb"
	sslmode  = "disable"
)

func GetPsqlConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
}
