package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/b-turchyn/idler/database"
	"github.com/b-turchyn/idler/state"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/gliderlabs/ssh"
)

const host = "0.0.0.0"
const port = 2222

var db *sql.DB

func main() {
  var err error
  db, err = database.Open("idler.sqlite3")

  go state.GetTopUsers(db)

  s, err := wish.NewServer(
    wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
    wish.WithHostKeyPath(".ssh/term_info_ed25519"),
    wish.WithMiddleware(
      bm.Middleware(teaHandler),
      lm.Middleware(),
      state.SessionCountMiddleware(),
    ),
    wish.WithPublicKeyAuth(func (ctx ssh.Context, key ssh.PublicKey) bool {
      return true
    }),
  )

  if err != nil {
    log.Fatalln(err)
  }

  done := make(chan os.Signal, 1)
  signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
  log.Printf("Starting SSH server on %s:%d", host, port)

  go func() {
    if err = s.ListenAndServe(); err != nil {
      log.Fatalln(err)
    }
  }()

  <-done
  log.Println("Stopping SSH server")
  ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
  defer func() { cancel() }()
  if err := s.Shutdown(ctx); err != nil {
    log.Fatalln(err)
  }
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
  pty, _, active := s.Pty()
  if !active {
    fmt.Println("no active terminal, skipping")
    return nil, nil
  }
  log.Printf("%s\n", base64.StdEncoding.EncodeToString(s.PublicKey().Marshal()))

  user := database.SelectUserByPublicKey(db, s.PublicKey().Marshal())
  user.Ident = s.User()

  m := state.State{
    Term: pty.Term,
    Width: pty.Window.Width,
    Height: pty.Window.Height,

    SelectedTab: 0,
    Cursor: 0,
    User: user,
    Db: db,
  }.SetupData()

  return m, []tea.ProgramOption{tea.WithAltScreen()}

}
