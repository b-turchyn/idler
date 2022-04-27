package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func Open(filename string) (*sql.DB, error) {
  db, err := sql.Open(viper.GetString("database.type"), filename)
  if err != nil {
    return nil, err
  }

  err = initTables(db)

  return db, nil
}

func initTables(db *sql.DB) error {
  stmt, err := db.Prepare(`
  CREATE TABLE IF NOT EXISTS
  users(id INTEGER PRIMARY KEY, ident VARCHAR, pubkey BLOB, points INT, bots INT, lurkers INT, viewers INT)
  `)
  errCheck(err)

  _, err = stmt.Exec()
  errCheck(err)
  stmt.Close()

  stmt, err = db.Prepare(`
  CREATE UNIQUE INDEX IF NOT EXISTS users_pubkey ON users(pubkey)
  `)
  _, err = stmt.Exec()
  errCheck(err)
  stmt.Close()

  addColumn(db, "users", "followers", "INT")
  addColumn(db, "users", "tier1subs", "INT")
  addColumn(db, "users", "tier2subs", "INT")
  addColumn(db, "users", "tier3subs", "INT")
  addColumn(db, "users", "data", "BLOB")

  return nil
}

func addColumn(db *sql.DB, table string, column string, columndef string) {
  stmt, err := db.Prepare(fmt.Sprintf("SELECT %s FROM %s LIMIT 0, 1", column, table))

  if err == nil {
    stmt.Close()
    return
  }

  stmt, err = db.Prepare(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, columndef))
  errCheck(err)
  defer stmt.Close()
  _, err = stmt.Exec()
  errCheck(err)
}

func errCheck(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
