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

}
