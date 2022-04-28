package item

import (
  "fmt"
  "github.com/b-turchyn/idler/util"
  "github.com/b-turchyn/idler/view"
)

type UpgradeState int

const (
  Unowned UpgradeState = iota
  Purchased
  Highlighted
)

var(
  ItemList = []ItemType{
    {
      Name: "Bot",
      Field: "Bots",
      InitialCost: 10,
      BasePoints: 1,
      Upgrades: []ItemUpgrade{
        { Name: "56K Modem", Cost: 5000, Upgrade: func(base uint64) uint64 { return base * 2 } },
        { Name: "T1 Line", Cost: 25000, Upgrade: func(base uint64) uint64 { return base * 2 } },
        { Name: "T2 Line", Cost: 100000, Upgrade: func(base uint64) uint64 { return base * 2 } },
        { Name: "T3 Line", Cost: 2500000, Upgrade: func(base uint64) uint64 { return base * 2 } },
      },
    },
    { Name: "Lurker", Field: "Lurkers", InitialCost: 200, BasePoints: 30 },
    { Name: "Viewer", Field: "Viewers", InitialCost: 2700, BasePoints: 500 },
    { Name: "Follower", Field: "Followers", InitialCost: 50000, BasePoints: 1200 },
    { Name: "Tier 1 Subscriber", Field: "Tier1Subs", InitialCost: 420000, BasePoints: 3500 },
  }
)

type ItemType struct {
  Name string
  InitialCost uint64
  BasePoints uint64
  Field string

  Upgrades []ItemUpgrade
}

type ItemUpgrade struct {
  Name string
  Description string
  Cost uint64
  Upgrade func(base uint64) uint64
}

func (u ItemUpgrade) ToString(state UpgradeState) string {
  content := fmt.Sprintf("%s: %s", u.Name, util.NumberFormatLong(u.Cost))
  if state == Purchased {
    return view.DisabledListItem(content)
  }
  return view.ListItem(content, state == Highlighted)
}
