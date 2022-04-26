package database

import (
	"database/sql"

	"github.com/b-turchyn/idler/model"
)

func GetLeaderboard(db *sql.DB) []model.User {
  stmt, err := db.Prepare(`SELECT
  ident, points
  FROM users
  ORDER BY points DESC
  LIMIT 0, 5
  `)
  errCheck(err)
  defer stmt.Close()

  rows, err := stmt.Query()
  errCheck(err)
  defer rows.Close()

  var result []model.User

  for rows.Next() {
    user := model.User{
      Stats: model.PlayerStats{},
    }
    rows.Scan(&user.Ident, &user.Stats.Points)
    result = append(result, user)
  }

  return result
}
