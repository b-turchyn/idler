package state

var(
  ItemList = []ItemType{
    { Name: "Bot", Field: "Bots", InitialCost: 10, BasePoints: 1 },
    { Name: "Lurker", Field: "Lurkers", InitialCost: 200, BasePoints: 30 },
    { Name: "Viewer", Field: "Viewers", InitialCost: 2700, BasePoints: 500 },
  }
)

type ItemType struct {
  Name string
  InitialCost uint64
  BasePoints uint64
  Field string
}
