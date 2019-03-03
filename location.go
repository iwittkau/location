package location

import "time"

// CheckIn is the main object for abstracting a check-in
type CheckIn struct {
	ID        string
	Name      string
	Time      time.Time
	Latitude  float64
	Longitude float64
}

// DB is the main database interface
type DB interface {
	Open(path string) error
	Close() error
}

// CheckInStore is the storage interface for CheckIns
type CheckInStore interface {
	CreateCheckIn(c CheckIn) error
	DeleteCheckIn(id string) error
	CheckIn(id string) (CheckIn, error)
	ListCheckIns(size int, offset int) ([]CheckIn, error)
}

// Storage combines all storage related interfaces
type Storage interface {
	DB
	CheckInStore
}
