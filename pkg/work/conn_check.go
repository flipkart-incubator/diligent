package work

import (
	"database/sql"
)

func ConnCheck(db *sql.DB) error {
	_, err := db.Exec("SELECT 1;")
	if err != nil {
		return err
	}
	return nil
}
