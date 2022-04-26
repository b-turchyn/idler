package view

import (
  "fmt"
  "strings"

	"github.com/b-turchyn/idler/util"
	"github.com/charmbracelet/lipgloss"
)

var (
  subtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
  highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

  descStyle = lipgloss.NewStyle().MarginTop(1)
  infoStyle = lipgloss.NewStyle().
    BorderStyle(lipgloss.NormalBorder()).
    BorderTop(true).
    BorderForeground(subtle)

  divider = lipgloss.NewStyle().
    SetString("•").
    Padding(0, 1).
    Foreground(subtle).
    String()

  // Tabs
  activeTabBorder = lipgloss.Border{
    Top:         "─",
    Bottom:      " ",
    Left:        "│",
    Right:       "│",
    TopLeft:     "╭",
    TopRight:    "╮",
    BottomLeft:  "┘",
    BottomRight: "└",
  }

  tabBorder = lipgloss.Border{
    Top:         "─",
    Bottom:      "─",
    Left:        "│",
    Right:       "│",
    TopLeft:     "╭",
    TopRight:    "╮",
    BottomLeft:  "┴",
    BottomRight: "┴",
  }

  tab = lipgloss.NewStyle().
    Border(tabBorder, true).
    BorderForeground(highlight).
    Padding(0, 1)

  activeTab = tab.Copy().Border(activeTabBorder, true)

  tabGap = tab.Copy().
    BorderTop(false).
    BorderLeft(false).
    BorderRight(false)

  TabList = []string{
    "Game",
    "Settings",
    "About",
    "Test",
  }
)

const (
)

func MainLayout(width int, selectedTab int, ident string, points uint64, perSecond uint64, body string) string {
  return lipgloss.JoinVertical(
    lipgloss.Top,
    Tabs(width, selectedTab),
    Title(ident, points, perSecond),
    body,
  )
}

func Title(ident string, points uint64, perSecond uint64) string {
  return lipgloss.JoinVertical(
    lipgloss.Left,
    descStyle.Render("Idler v0.1.0"),
    infoStyle.MarginBottom(1).Render(
      fmt.Sprintf(
        "%s%sPoints: %d%s%d/s",
        ident,
        divider,
        points,
        divider,
        perSecond,
      ),
    ),
  )
}

func Tabs(width int, selectIndex int) string {
  var tablist []string

  for i, v := range TabList {
    tablist = append(tablist, Tab(v, i, selectIndex))
  }

  row := lipgloss.JoinHorizontal(
    lipgloss.Top,
    tablist...,
  )
  gap := tabGap.Render(strings.Repeat(" ", util.Max(0, width - lipgloss.Width(row) - 2)))
  row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

  return row
}

func Tab(title string, index int, selectedIndex int) string {
  f := tab
  if selectedIndex == index {
    f = activeTab
  }

  return f.Render(title)
}
