package state

import (
	"database/sql"
	"fmt"

	"github.com/b-turchyn/idler/database"
	"github.com/b-turchyn/idler/model"
	"github.com/b-turchyn/idler/model/item"
	"github.com/b-turchyn/idler/util"
	"github.com/b-turchyn/idler/view"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
  CURSOR_COLUMN = 0
  CURSOR_ROW = 1

  TAB_GAME = 0
  TAB_SETTINGS = 1
  TAB_LEADERBOARD = 2
  TAB_ABOUT = 3

  GAME_COLUMN_ITEMS = 0
  GAME_COLUMN_UPGRADES = 1
)

/**
 * To use:
 * - 0: Which tab you are on
 * - 1: Which column you are on
 */
var maxLengthColumns = [][]int{
  // TAB_GAME
  []int{ len(item.ItemList), len(item.ItemList[0].Upgrades) },
}

type State struct {
  Term string
  Width int
  Height int
  PerSecond uint64
  windowReady bool

  Cursor [2]int
  SelectedTab int
  SelectedItem int

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
                    lipgloss.Height(view.Title(m.User.Ident, m.User.StatsV02.Points, m.PerSecond))

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
      m = m.SetCursorRow(m.CursorRow() - 1)
    case "down", "j":
      m = m.SetCursorRow(m.CursorRow() + 1)
    case "left", "h":
      m = m.SetCursor(m.CursorColumn() - 1, 0)
    case "right", "l":
      m = m.SetCursor(m.CursorColumn() + 1, 0)
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
  m.User.StatsV02.Points += m.PerSecond

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
    m.User.StatsV02.Points,
    m.PerSecond,
    f(m),
  )
}

func (m State) purchase() State {
  if m.SelectedTab != TAB_GAME {
    return m
  }

  if m.CursorColumn() == GAME_COLUMN_ITEMS {
    i := item.ItemList[m.SelectedItem]
    field := m.User.StatsV02.GetItem(m.SelectedItem)

    price := field.CalculateNextCost(i)

    if m.User.StatsV02.Points < price {
      return m
    }

    m.User.StatsV02.Points -= price
    field.Quantity++
    m.User.StatsV02.Items[m.CursorRow()] = field


    m = m.recalculatePerSecond()
  } else if m.CursorColumn() == GAME_COLUMN_UPGRADES {
    upgradeIndex := m.CursorRow()
    i := item.ItemList[m.SelectedItem]
    field := m.User.StatsV02.GetItem(m.SelectedItem)
    upgrade := i.Upgrades[upgradeIndex]

    if m.User.StatsV02.Points < upgrade.Cost || field.IsUpgraded(upgradeIndex) {
      return m
    }

    m.User.StatsV02.Points -= upgrade.Cost

    if len(field.Upgrades) < len(i.Upgrades) {
      // Create a new array that's the correct size
      tempupgrades := make([]bool, len(i.Upgrades))
      for i, v := range field.Upgrades {
        tempupgrades[i] = v
      }
      field.Upgrades = tempupgrades
    }

    field.Upgrades[upgradeIndex] = true
    m.User.StatsV02.Items[m.SelectedItem] = field

    m = m.recalculatePerSecond()
  }

  return m
}

func (m State) recalculatePerSecond() State {
  var result uint64

  for i, v := range item.ItemList {
    field := m.User.StatsV02.Items[i]

    result += field.CalculateItemPerSecond(v)
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

  m.Cursor = [2]int{0, 0}
  m.SelectedItem = 0
  m.RecalculateCursorDisplay()

  return m
}

func (m State) ViewerCount() string {
  var listItems []string

  for i, v := range item.ItemList {
    field := m.User.StatsV02.GetItem(i)
    formattednumber := util.NumberFormatLong(field.Quantity)
    listItems = append(listItems, view.ListItem(fmt.Sprintf("%ss: %s", v.Name, formattednumber), false))
  }

  return view.List(
    "Your Viewer List",
    listItems,
  )
}

func (m State) CostList() string {
  var listItems []string

  for i, v := range item.ItemList {
    field := m.User.StatsV02.GetItem(i)
    var itemstring string

    formattednumber := util.NumberFormatLong(field.CalculateNextCost(v))
    itemstring = view.ListItem(
      fmt.Sprintf("%ss: %s", v.Name, formattednumber),
      m.Cursor[CURSOR_COLUMN] == 0 && i == m.Cursor[CURSOR_ROW],
    )

    listItems = append(listItems, itemstring)
  }

  return view.List(
    "Buy Viewers",
    listItems,
  )
}

func (m State) UpgradeList() string {
  var upgradeItems []string
  selectedItem := item.ItemList[m.SelectedItem]

  for i, v := range selectedItem.Upgrades {
    state := item.Unowned
    if m.CursorColumn() == GAME_COLUMN_UPGRADES && m.CursorRow() == i {
      state = item.Highlighted
    } else if m.User.StatsV02.GetItem(m.SelectedItem).IsUpgraded(i) {
      state = item.Purchased
    }
    upgradeItems = append(upgradeItems, v.ToString(state))
  }

  return view.List(
    "Upgrades",
    upgradeItems,
  )
}

func (m State) SetCursorColumn(col int) State {
  max := len(maxLengthColumns[m.SelectedTab])
  m.Cursor[CURSOR_COLUMN] = util.Within(0, col, max - 1)

  return m.RecalculateCursorDisplay()
}

func (m State) CursorColumn() int {
  return m.Cursor[CURSOR_COLUMN]
}

func (m State) SetCursorRow(row int) State {
  max := maxLengthColumns[m.SelectedTab][m.CursorColumn()]
  m.Cursor[CURSOR_ROW] = util.Within(0, row, util.Max(0, max - 1))

  return m.RecalculateCursorDisplay()
}

func (m State) CursorRow() int {
  return m.Cursor[CURSOR_ROW]
}

func (m State) SetCursor(col, row int) State {
  max := len(maxLengthColumns[m.SelectedTab])
  m.Cursor[CURSOR_COLUMN] = util.Within(0, col, max - 1)

  max = maxLengthColumns[m.SelectedTab][m.CursorColumn()]
  m.Cursor[CURSOR_ROW] = util.Within(0, row, util.Max(0, max - 1))

  return m.RecalculateCursorDisplay()
}

func (m State) RecalculateCursorDisplay() State {
  switch m.SelectedTab {
  case TAB_GAME:
    if m.CursorColumn() == GAME_COLUMN_ITEMS {
      m.SelectedItem = m.CursorRow()
      maxLengthColumns[TAB_GAME][GAME_COLUMN_UPGRADES] = len(item.ItemList[m.SelectedItem].Upgrades)
    }
  }
  return m
}
