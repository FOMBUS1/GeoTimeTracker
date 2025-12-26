package pgstorage

import (
	"time"
)

type UserVisits struct {
	ID         uint64     `db:"id"`
	UserID     uint64     `db:"user_id"`
	LocationFK uint64     `db:"location_fk"`
	EntryTime  time.Time  `db:"entry_time"`
	ExitTime   *time.Time `db:"exit_time"`
	Duration   int64      `db:"duration_sec"`
}

type Locations struct {
	LocationID uint64 `db:"location_id"`
	Location   string `db:"location"`
}

const (
	tableName = "userVisits"

	IDСolumnName           = "id"
	UserIDColumnName       = "user_id"
	LocationFKColumnName   = "location_id"
	EntryTimeColumnName    = "entry_time"
	ExitTimeColumnName     = "exit_time"
	DurationTimeColumnName = "duration_sec"
)

const (
	locationTableName = "locations"

	LocationIDColumnName = "location_id"
	LocationColumnName   = "location"
)
