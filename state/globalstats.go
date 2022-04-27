package state

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/b-turchyn/idler/database"
	"github.com/b-turchyn/idler/model"
	"github.com/b-turchyn/idler/util"
	"github.com/b-turchyn/idler/view"
	"github.com/charmbracelet/wish"
	"github.com/gliderlabs/ssh"
)

var SessionCount = 0
var TopUsersList []model.User

func SessionCountMiddleware() wish.Middleware {
  return func (sh ssh.Handler) ssh.Handler {
    return func(s ssh.Session) {
      SessionCount++
      sh(s)
      SessionCount--
    }
  }
}

func GetTopUsers(db *sql.DB) {
  TopUsersList = database.GetLeaderboard(db)

  for range time.Tick(2 * time.Minute) {
    TopUsersList = database.GetLeaderboard(db)
  }
}

func TopUsers() []string {
  var result []string
  for i, v := range TopUsersList {
    result = append(result, view.ListItem(
      fmt.Sprintf("%d. %s (%s points)", i + 1, v.Ident, util.NumberFormatLong(v.Stats.Points)),
      false,
    ))
  }

  return result
}
