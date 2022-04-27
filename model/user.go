package model

import "log"


type User struct {
  Ident string
  PublicKey []byte

  StatsV01 PlayerStats
  StatsV02 PlayerStatsV2
}

type PlayerStats struct {
  Points uint64
  Bots uint64
  Lurkers uint64
  Viewers uint64
  Followers uint64
  Tier1Subs uint64
  Tier2Subs uint64
  Tier3Subs uint64
}

type PlayerStatsV2 struct {
  Points uint64
  Items []PurchasedItem
}
type PurchasedItem struct {
  Quantity uint64
  Upgrades []bool
}

func (u User) Migrate() User {
  if u.StatsV01 != (PlayerStats{}) {
    log.Println("User is on v1 stats, migrating to v2")
    u = u.migrateStatsV01ToV02()
  }

  return u
}

func (u User) migrateStatsV01ToV02() User {
  u.StatsV02 = PlayerStatsV2{
    Points: u.StatsV01.Points,
    Items: []PurchasedItem{},
  }

  u.StatsV02.Items = append(u.StatsV02.Items, PurchasedItem{Quantity: u.StatsV01.Bots})
  u.StatsV02.Items = append(u.StatsV02.Items, PurchasedItem{Quantity: u.StatsV01.Lurkers})
  u.StatsV02.Items = append(u.StatsV02.Items, PurchasedItem{Quantity: u.StatsV01.Viewers})
  u.StatsV02.Items = append(u.StatsV02.Items, PurchasedItem{Quantity: u.StatsV01.Followers})
  u.StatsV02.Items = append(u.StatsV02.Items, PurchasedItem{Quantity: u.StatsV01.Tier1Subs})
  u.StatsV02.Items = append(u.StatsV02.Items, PurchasedItem{Quantity: u.StatsV01.Tier2Subs})
  u.StatsV02.Items = append(u.StatsV02.Items, PurchasedItem{Quantity: u.StatsV01.Tier3Subs})
  u.StatsV01.Points = 0
  u.StatsV01.Bots = 0
  u.StatsV01.Lurkers = 0
  u.StatsV01.Viewers = 0
  u.StatsV01.Followers = 0
  u.StatsV01.Tier1Subs = 0
  u.StatsV01.Tier2Subs = 0
  u.StatsV01.Tier3Subs = 0

  return u
}
