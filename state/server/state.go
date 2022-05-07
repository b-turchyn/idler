package server

import (
  tea "github.com/charmbracelet/bubbletea"
)

type ServerState struct {
  SelectedTab int
}

func InitialModel() ServerState {
  return ServerState{}
}

func (m ServerState) Init() tea.Cmd {
  return nil
}

func (m ServerState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "q", "ctrl+c":
      return m, tea.Quit
    }
  }
  return m, nil
}

func (m ServerState) View() string {
  return "Hello, world"
}
