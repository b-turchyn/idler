package state

var(
  ItemList = []ItemType{
    { Name: "Bot", Field: "Bots", InitialCost: 10, BasePoints: 1 },
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
}
