/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

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
  "github.com/b-turchyn/idler/state/server"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/wish"
  bm "github.com/charmbracelet/wish/bubbletea"
  lm "github.com/charmbracelet/wish/logging"
  "github.com/gliderlabs/ssh"
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
)

var db *sql.DB

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
  Use:   "serve",
  Short: "A brief description of your command",
  Long: `A longer description that spans multiple lines and likely contains examples
  and usage of using your command. For example:

  Cobra is a CLI library for Go that empowers applications.
  This application is a tool to generate the needed files
  to quickly create a Cobra application.`,
  Run: func(cmd *cobra.Command, args []string) {
    var err error
    db, err = database.Open(viper.GetString("database.filename"))

    go state.GetTopUsers(db)

    s, err := wish.NewServer(
      wish.WithAddress(fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))),
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
    log.Printf("Starting SSH server on %s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))

    go func() {
      if err = s.ListenAndServe(); err != nil {
        log.Fatalln(err)
      }
    }()


    if viper.GetBool("server.non-interactive") {
      <-done
    } else {
      p := tea.NewProgram(
        server.InitialModel(),
        tea.WithAltScreen(),
      )
      if err = p.Start(); err != nil {
        log.Fatalln(err)
      }
    }
    log.Println("Stopping SSH server")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer func() { cancel() }()
    if err := s.Shutdown(ctx); err != nil {
      log.Fatalln(err)
    }
  },
}

func init() {
  rootCmd.AddCommand(serveCmd)

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // serveCmd.PersistentFlags().String("foo", "", "A help for foo")
  serveCmd.PersistentFlags().String("host", "0.0.0.0", "Address to bind to")
  serveCmd.PersistentFlags().IntP("port", "p", 2222, "Port to listen on")
  serveCmd.PersistentFlags().BoolP("no-interactive", "i", false, "Disable interactive host view")
  viper.BindPFlag("server.host", serveCmd.PersistentFlags().Lookup("host"))
  viper.BindPFlag("server.port", serveCmd.PersistentFlags().Lookup("port"))
  viper.BindPFlag("server.non-interactive", serveCmd.PersistentFlags().Lookup("no-interactive"))

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
    Cursor: [2]int{0, 0},
    User: user,
    Db: db,
  }.SetupData()

  return m, []tea.ProgramOption{tea.WithAltScreen()}

}
