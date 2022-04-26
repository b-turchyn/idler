package view

import _ "embed"
import "github.com/charmbracelet/glamour"

//go:embed CHANGELOG.md
var changelog string

func ChangelogView() string {
  out, _ := glamour.Render(changelog, "dark")
  return out
}
