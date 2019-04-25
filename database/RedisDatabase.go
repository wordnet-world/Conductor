package database

type Cache interface {
	Ping() bool
	GetGames()
}
