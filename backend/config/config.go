package config

import (
	"picture_storage/db"
	"strconv"

	ec "github.com/kiririx/easy-config/ec"
)

var H ec.Handler

func init() {
	port, _ := strconv.Atoi(db.MYSQL_PORT)
	H = ec.Initialize(ec.NewMySQLStorage(db.MYSQL_HOST, port, db.MYSQL_USERNAME, db.MYSQL_PASSWORD, db.MYSQL_DATABASE), "main")
}
