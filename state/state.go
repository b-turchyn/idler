package state

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/b-turchyn/idler/database"
	"github.com/b-turchyn/idler/model"
	"github.com/b-turchyn/idler/util"
	"github.com/b-turchyn/idler/view"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State struct {
  Term string
  Width int
  Height int
  PerSecond uint64
  windowReady bool

  Cursor int
  SelectedTab int

  User model.User

  Db *sql.DB

  ChangelogViewport viewport.Model
}

func (m State) Init() tea.Cmd {
  m.recalculatePerSecond()
  return tea.Batch(util.TickNow, util.StartGameLoop, util.SaveDataLoop())
}

func (m State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var (
    cmd tea.Cmd
    cmds []tea.Cmd
  )

  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    m.Height = msg.Height
    m.Width = msg.Width

    headerHeight := lipgloss.Height(view.Tabs(m.Width, 0)) +
                    lipgloss.Height(view.Title(m.User.Ident, m.User.Stats.Points, m.PerSecond))

    if !m.windowReady {
      m.ChangelogViewport = viewport.New(m.Width, m.Height - headerHeight)
      m.ChangelogViewport.YPosition = headerHeight
      m.ChangelogViewport.HighPerformanceRendering = false
      m.ChangelogViewport.SetContent(view.ChangelogView())

      m.windowReady = true
    } else {
      m.ChangelogViewport.Width = m.Width
      m.ChangelogViewport.Height = m.Height - headerHeight
    }

  case tea.KeyMsg:
    switch msg.String() {
    case "q", "ctrl+c":
      database.SaveUserByPublicKey(m.Db, m.User)
      return m, tea.Quit
    case "tab":
      m = m.IncrementTab(true)
    case "shift+tab":
      m = m.IncrementTab(false)
    case "enter":
      m = m.purchase()
    case "up", "k":
      m.Cursor--
    case "down", "j":
      m.Cursor++
    case "ctrl+l":
      tea.EnterAltScreen()
    }
  case util.ViewTickMsg:
    cmds = append(cmds, util.Tick())
  case util.GameTickMsg:
    m = m.GameTick()
    cmds = append(cmds, util.GameLoop())
  case util.SaveDataMsg:
    database.SaveUserByPublicKey(m.Db, m.User)
    cmds = append(cmds, util.SaveDataLoop())
  }

  if m.SelectedTab == 3 {
    m.ChangelogViewport, cmd = m.ChangelogViewport.Update(msg)
    cmds = append(cmds, cmd)
  }

  return m, tea.Batch(cmds...)
}

func (m State) SetupData() State {
  return m.recalculatePerSecond()
}

func (m State) GameTick() State {
  m.User.Stats.Points += m.PerSecond

  return m
}

func (m State) View() string {
  var f func(m State) string

  switch m.SelectedTab {
  case 3:
    f = aboutView
  case 2:
    f = settingsView
  case 1:
    f = leaderboardView
  case 0:
    fallthrough
  default:
    f = gameView
  }

  return view.MainLayout(
    m.Width,
    m.SelectedTab,
    m.User.Ident,
    m.User.Stats.Points,
    m.PerSecond,
    f(m),
  )
}

func (m State) purchase() State {
  if m.Cursor < 0 || m.Cursor >= len(ItemList) {
    return m
  }

  item := ItemList[m.Cursor]
  refmodel := reflect.ValueOf(&m.User.Stats).Elem()
  field := refmodel.FieldByName(item.Field)
  price := util.Cost(item.InitialCost, field.Uint())

  if m.User.Stats.Points < price {
    return m
  }

  m.User.Stats.Points -= price
  field.SetUint(field.Uint() + 1)

  m = m.recalculatePerSecond()

  return m
}

func (m State) recalculatePerSecond() State {
  var result uint64

  for _, v := range ItemList {
    refmodel := reflect.ValueOf(m.User.Stats)
    field := refmodel.FieldByName(v.Field).Uint()

    result += v.BasePoints * field
  }

  m.PerSecond = result

  return m
}

func (m State) IncrementTab(up bool) State {
  if up {
    if m.SelectedTab + 1 >= len(view.TabList) {
      m.SelectedTab = 0
    } else {
      m.SelectedTab = m.SelectedTab + 1
    }
  } else {
    if m.SelectedTab <= 0 {
      m.SelectedTab = len(view.TabList) - 1
    } else {
      m.SelectedTab = m.SelectedTab - 1
    }
  }

  m.Cursor = 0

  return m
}

func (m State) ViewerCount() string {
  var listItems []string

  for _, v := range ItemList {
    refmodel := reflect.ValueOf(m.User.Stats)
    formattednumber := util.NumberFormatLong(refmodel.FieldByName(v.Field).Uint())
    listItems = append(listItems, view.ListItem(fmt.Sprintf("%ss: %s", v.Name, formattednumber), false))
  }

  return view.List(
    "Your Viewer List",
    listItems,
  )
}

func (m State) CostList() string {
  var listItems []string

  for i, v := range ItemList {
    refmodel := reflect.ValueOf(m.User.Stats)
    var itemstring string

    formattednumber := util.NumberFormatLong(util.Cost(v.InitialCost, refmodel.FieldByName(v.Field).Uint()))
    itemstring = view.ListItem(
      fmt.Sprintf("%ss: %s", v.Name, formattednumber),
      i == m.Cursor,
    )

    listItems = append(listItems, itemstring)
  }

  return view.List(
    "Buy Viewers",
    listItems,
  )
}
