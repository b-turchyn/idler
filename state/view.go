package state

import (
	"github.com/charmbracelet/lipgloss"
)

func aboutView(m State) string {
  return "About page"
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
