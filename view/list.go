package view

import (
	"github.com/charmbracelet/lipgloss"
)

var (
  list = lipgloss.NewStyle().
    Border(lipgloss.NormalBorder(), false, true, false, false).
    MarginRight(2).
    Height(8).
    Width(30)

  listHeader = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderBottom(true).
    BorderForeground(subtle).
    MarginRight(2).
    Render

  listItem = lipgloss.NewStyle().PaddingLeft(2).Render

  activeListItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(highlight).Render
)

func List(title string, items []string) string {
  items = prepend(items, listHeader(title))

  return lipgloss.JoinVertical(
    lipgloss.Top,
    items...,
  )
}

func ListItem(text string, selected bool) string {
  if selected {
    return activeListItem(text)
  }

  return listItem(text)
}

func ListHeader(text string) string {
  return listHeader(text)
}

func prepend(items []string, prepend string) []string {
  return append([]string{prepend}, items...)
}
