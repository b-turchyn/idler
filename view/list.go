package view

import (
	"github.com/charmbracelet/lipgloss"
)

var (
  list = lipgloss.NewStyle().
    Border(lipgloss.NormalBorder(), false, true, false, false).
    BorderForeground(special).
    PaddingLeft(2).
    PaddingRight(2)


  listHeader = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderBottom(true).
    BorderForeground(subtle).
    MarginRight(2).
    Render

  listItem = lipgloss.NewStyle().Render

  activeListItem = lipgloss.NewStyle().Foreground(highlight).Render
)

func List(title string, items []string) string {
  items = prepend(items, listHeader(title))

  return list.Render(
    lipgloss.JoinVertical(
      lipgloss.Left,
      items...,
    ),
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
