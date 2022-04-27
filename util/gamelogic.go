package util

import (
  "time"

  tea "github.com/charmbracelet/bubbletea"
)

const defaultFps = 2

type ViewTickMsg struct {
  Time time.Time
}

type GameTickMsg struct {}

type SaveDataMsg struct {}

func TickNow() tea.Msg {
  return ViewTickMsg{Time: time.Now()}
}

func Tick() tea.Cmd {
  return tea.Tick(time.Second / defaultFps, func (t time.Time) tea.Msg {
    return ViewTickMsg{Time: t}
  })
}

func StartGameLoop() tea.Msg {
  return GameTickMsg{}
}

func GameLoop() tea.Cmd {
  return tea.Tick(time.Second, func (t time.Time) tea.Msg {
    return GameTickMsg{}
  })
}

func SaveDataLoop() tea.Cmd {
  return tea.Tick(time.Minute, func (t time.Time) tea.Msg {
    return SaveDataMsg{}
  })
}
