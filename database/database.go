package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Open(filename string) (*sql.DB, error) {
  db, err := sql.Open("sqlite3", filename)
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

  return nil
}

func errCheck(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
