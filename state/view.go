package state

import (
	"fmt"

	"github.com/b-turchyn/idler/util"
	"github.com/b-turchyn/idler/view"
	"github.com/charmbracelet/lipgloss"
)

func aboutView(m State) string {
  return lipgloss.JoinVertical(
    lipgloss.Left,
    m.ChangelogViewport.View(),
  )
}

func settingsView(m State) string {
  return "Settings page"
}

func gameView(m State) string {
  return lipgloss.JoinHorizontal(
    lipgloss.Top,
    m.ViewerCount(),
    m.CostList(),
  )
}

func leaderboardView(m State) string {
  content := append(
    []string{},
    view.ListHeader("Leaderboard"),
    view.ListHeader(fmt.Sprintf("Online users: %s", util.NumberFormatLong(uint64(SessionCount)))),
    view.ListHeader("Top Users"),
  )
  content = append(content, TopUsers()...)

  return lipgloss.JoinVertical(
    lipgloss.Center,
    content...,
  )
}
