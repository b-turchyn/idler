package database

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/b-turchyn/idler/model"
)

func SelectUserByPublicKey(db *sql.DB, pubkey []byte) model.User {
  stmt, err := db.Prepare(`SELECT
  ident, points, bots, lurkers, viewers, followers, tier1subs, tier2subs, tier3subs, data
  FROM users
  WHERE pubkey = ?
  `)
  errCheck(err)

  row, err := stmt.Query(pubkey)
  errCheck(err)

  result := model.User{
    PublicKey: pubkey,
    StatsV01: model.PlayerStats{
      Points: 0,
      Bots: 1,
      Lurkers: 0,
      Viewers: 0,
      Followers: 0,
      Tier1Subs: 0,
      Tier2Subs: 0,
      Tier3Subs: 0,
    },
  }

  if row.Next() {
    log.Println("Found user data")
    var data []byte
    row.Scan(
      &result.Ident,
      &result.StatsV01.Points,
      &result.StatsV01.Bots,
      &result.StatsV01.Lurkers,
      &result.StatsV01.Viewers,
      &result.StatsV01.Followers,
      &result.StatsV01.Tier1Subs,
      &result.StatsV01.Tier2Subs,
      &result.StatsV01.Tier3Subs,
      &data,
    )

    if len(data) >= 2 {
      err = json.Unmarshal(data, &result.StatsV02)
      errCheck(err)
    }
  } else {
    log.Println("Creating user data")
    stmt.Close()
    stmt, err = db.Prepare(`INSERT INTO users 
      (ident, pubkey, points, bots, lurkers, viewers, followers, tier1subs, tier2subs, tier3subs)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
    errCheck(err)
    _, err = stmt.Exec(
      result.Ident,
      pubkey,
      result.StatsV01.Points,
      result.StatsV01.Bots,
      result.StatsV01.Lurkers,
      result.StatsV01.Viewers,
      result.StatsV01.Followers,
      result.StatsV01.Tier1Subs,
      result.StatsV01.Tier2Subs,
      result.StatsV01.Tier3Subs,
    )
    errCheck(err)
    stmt.Close()
  }
  row.Close()

  result = result.Migrate()

  SaveUserByPublicKey(db, result)

  return result
}

func SaveUserByPublicKey(db *sql.DB, user model.User) {
  stmt, err := db.Prepare(`UPDATE users
  SET ident = ?,
  points = ?,
  bots = ?,
  lurkers = ?,
  viewers = ?,
  followers = ?,
  tier1subs = ?,
  tier2subs = ?,
  tier3subs = ?,
  data = ?
  WHERE pubkey = ?
  `)
  errCheck(err)
  defer stmt.Close()

  data, err := json.Marshal(user.StatsV02)

  _, err = stmt.Exec(
    user.Ident,
    user.StatsV02.Points,
    user.StatsV01.Bots,
    user.StatsV01.Lurkers,
    user.StatsV01.Viewers,
    user.StatsV01.Followers,
    user.StatsV01.Tier1Subs,
    user.StatsV01.Tier2Subs,
    user.StatsV01.Tier3Subs,
    data,
    user.PublicKey,
  )
  errCheck(err)
}
