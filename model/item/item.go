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

const (
  THOUSAND = 1000
  MILLION = THOUSAND * THOUSAND
  BILLION = MILLION * THOUSAND
  TRILLION = BILLION * THOUSAND
  QUADRILLION = TRILLION * THOUSAND
)

func standardUpgrade(base uint64) uint64 {
  return base * 2
}

var(
  ItemList = []ItemType{
    {
      Name: "Bot",
      Field: "Bots",
      InitialCost: 10,
      BasePoints: 1,
      Upgrades: []ItemUpgrade{
        { Name: "56K Modem", Cost: 5 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 10 },
        { Name: "T1 Line", Cost: 25 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 20 },
        { Name: "T2 Line", Cost: 100 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 30 },
        { Name: "T3 Line", Cost: 2500 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 50 },
      },
    },
    {
      Name: "Lurker",
      Field: "Lurkers",
      InitialCost: 200,
      BasePoints: 30,
      Upgrades: []ItemUpgrade{
        { Name: "Lurker Upgrade 1", Cost: 30 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 10 },
        { Name: "Lurker Upgrade 2", Cost: 150 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 20 },
        { Name: "Lurker Upgrade 3", Cost: 600 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 30 },
        { Name: "Lurker Upgrade 4", Cost: 15 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 40 },
        { Name: "Lurker Upgrade 5", Cost: 60 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 50 },
        { Name: "Lurker Upgrade 6", Cost: 300 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 60 },
        { Name: "Lurker Upgrade 7", Cost: 1200 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 70 },
      },
    },
    {
      Name: "Viewer",
      Field: "Viewers",
      InitialCost: 2700,
      BasePoints: 500,
      Upgrades: []ItemUpgrade{
        { Name: "Viewer Upgrade 1", Cost: 450 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 10 },
        { Name: "Viewer Upgrade 2", Cost: 1800 * THOUSAND, Upgrade: standardUpgrade, MinimumQuantity: 20 },
        { Name: "Viewer Upgrade 3", Cost: 7 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 30 },
        { Name: "Viewer Upgrade 4", Cost: 28 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 40 },
        { Name: "Viewer Upgrade 5", Cost: 1 * BILLION, Upgrade: standardUpgrade, MinimumQuantity: 50 },
        { Name: "Viewer Upgrade 6", Cost: 4 * BILLION, Upgrade: standardUpgrade, MinimumQuantity: 60 },
        { Name: "Viewer Upgrade 7", Cost: 16 * BILLION, Upgrade: standardUpgrade, MinimumQuantity: 70 },
      },
    },
    {
      Name: "Follower",
      Field: "Followers",
      InitialCost: 50000,
      BasePoints: 1200,
      Upgrades: []ItemUpgrade{
        { Name: "Follower Upgrade 1", Cost: 2 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 10 },
        { Name: "Follower Upgrade 2", Cost: 8 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 20 },
        { Name: "Follower Upgrade 3", Cost: 32 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 30 },
        { Name: "Follower Upgrade 4", Cost: 128 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 40 },
        { Name: "Follower Upgrade 5", Cost: 512 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 50 },
        { Name: "Follower Upgrade 6", Cost: 2048 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 60 },
        { Name: "Follower Upgrade 7", Cost: 8192 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 70 },
      },
    },
    {
      Name: "Tier 1 Subscriber",
      Field: "Tier1Subs",
      InitialCost: 420000,
      BasePoints: 3500,
      Upgrades: []ItemUpgrade{
        { Name: "Tier 1 Sub Upgrade 1", Cost: 16 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 10 },
        { Name: "Tier 1 Sub Upgrade 2", Cost: 64 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 20 },
        { Name: "Tier 1 Sub Upgrade 3", Cost: 256 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 30 },
        { Name: "Tier 1 Sub Upgrade 4", Cost: 1024 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 40 },
        { Name: "Tier 1 Sub Upgrade 5", Cost: 4096 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 50 },
        { Name: "Tier 1 Sub Upgrade 6", Cost: 16384 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 60 },
        { Name: "Tier 1 Sub Upgrade 7", Cost: 65536 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 70 },
      },
    },
    {
      Name: "Tier 2 Subscriber",
      Field: "Tier2Subs",
      InitialCost: 2.5 * MILLION,
      BasePoints: 8700,
      Upgrades: []ItemUpgrade{
        { Name: "Tier 2 Sub Upgrade 1", Cost: 128 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 10 },
        { Name: "Tier 2 Sub Upgrade 2", Cost: 512 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 20 },
        { Name: "Tier 2 Sub Upgrade 3", Cost: 2048 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 30 },
        { Name: "Tier 2 Sub Upgrade 4", Cost: 8196 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 40 },
        { Name: "Tier 2 Sub Upgrade 5", Cost: 32768 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 50 },
        { Name: "Tier 2 Sub Upgrade 6", Cost: 130 * BILLION, Upgrade: standardUpgrade, MinimumQuantity: 60 },
        { Name: "Tier 2 Sub Upgrade 7", Cost: 520 * BILLION, Upgrade: standardUpgrade, MinimumQuantity: 70 },
      },
    },
    {
      Name: "Tier 3 Subscriber",
      Field: "Tier3Subs",
      InitialCost: 12 * MILLION,
      BasePoints: 15000,
      Upgrades: []ItemUpgrade{
        { Name: "Tier 3 Sub Upgrade 1", Cost: 2048 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 10 },
        { Name: "Tier 3 Sub Upgrade 2", Cost: 8196 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 20 },
        { Name: "Tier 3 Sub Upgrade 3", Cost: 32768 * MILLION, Upgrade: standardUpgrade, MinimumQuantity: 30 },
        { Name: "Tier 3 Sub Upgrade 4", Cost: 130 * BILLION, Upgrade: standardUpgrade, MinimumQuantity: 40 },
        { Name: "Tier 3 Sub Upgrade 5", Cost: 520 * BILLION, Upgrade: standardUpgrade, MinimumQuantity: 50 },
        { Name: "Tier 3 Sub Upgrade 6", Cost: 2 * TRILLION, Upgrade: standardUpgrade, MinimumQuantity: 60 },
        { Name: "Tier 3 Sub Upgrade 7", Cost: 8 * TRILLION, Upgrade: standardUpgrade, MinimumQuantity: 70 },
      },
    },
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
  MinimumQuantity uint
}

func (u ItemUpgrade) ToString(state UpgradeState) string {
  content := fmt.Sprintf("%s: %s", u.Name, util.NumberFormatLong(u.Cost))
  if state == Purchased {
    return view.DisabledListItem(content)
  }
  return view.ListItem(content, state == Highlighted)
}

