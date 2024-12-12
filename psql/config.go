package psql

import "fmt"

const (
	host     = "localhost"
	dbname   = "mydb"
	user     = "myuser"
	password = "myuser"
	port     = "5432"
	sslmode  = "disable"
)

func ConnectStr() string {
	return fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%s sslmode=%s", host, dbname, user, password, port, sslmode)
}
