package model


type User struct {
  Ident string
  PublicKey []byte

  Stats PlayerStats
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
