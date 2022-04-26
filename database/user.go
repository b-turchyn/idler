package database

import (
	"database/sql"
	"log"

	"github.com/b-turchyn/idler/model"
)

func SelectUserByPublicKey(db *sql.DB, pubkey []byte) model.User {
  stmt, err := db.Prepare(`SELECT
  ident, points, bots, lurkers, viewers
  FROM users
  WHERE pubkey = ?
  `)
  errCheck(err)
  defer stmt.Close()

  row, err := stmt.Query(pubkey)
  errCheck(err)
  defer row.Close()

  result := model.User{
    PublicKey: pubkey,
    Stats: model.PlayerStats{
      Points: 0,
      Bots: 1,
      Lurkers: 0,
      Viewers: 0,
    },
  }

  if row.Next() {
    log.Println("Found user data")
    row.Scan(&result.Ident, &result.Stats.Points, &result.Stats.Bots, &result.Stats.Lurkers, &result.Stats.Viewers)
  } else {
    log.Println("Creating user data")
    stmt, err = db.Prepare(`INSERT INTO users 
      (ident, pubkey, points, bots, lurkers, viewers)
      VALUES (?, ?, ?, ?, ?, ?)
    `)
    errCheck(err)
    defer stmt.Close()
    _, err = stmt.Exec(
      result.Ident,
      pubkey,
      result.Stats.Points,
      result.Stats.Bots,
      result.Stats.Lurkers,
      result.Stats.Viewers,
    )
    errCheck(err)
  }

  return result
}

func SaveUserByPublicKey(db *sql.DB, user model.User) {
  stmt, err := db.Prepare(`UPDATE users
  SET ident = ?,
  points = ?,
  bots = ?,
  lurkers = ?,
  viewers = ?
  WHERE pubkey = ?
  `)
  errCheck(err)
  defer stmt.Close()

  _, err = stmt.Exec(
    user.Ident,
    user.Stats.Points,
    user.Stats.Bots,
    user.Stats.Lurkers,
    user.Stats.Viewers,
    user.PublicKey,
  )
  errCheck(err)
}
