package global

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DBEngine *sqlx.DB

var DBEngineShop *sqlx.DB
