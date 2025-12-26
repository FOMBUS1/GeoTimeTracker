package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGstorage struct {
	db *pgxpool.Pool
}

func NewPGStorge(connString string) (*PGstorage, error) {

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}
	storage := &PGstorage{
		db: db,
	}
	err = storage.initTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGstorage) initTables() error {
	create_locations := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
			%v BIGSERIAL PRIMARY KEY,
			%v TEXT UNIQUE NOT NULL
		);
	`, locationTableName, LocationIDColumnName, LocationColumnName)

	create_user_visits := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		%v BIGSERIAL PRIMARY KEY,
		%v BIGSERIAL,
		%v BIGINT REFERENCES %v(%v),
		%v TIMESTAMPTZ NOT NULL,
		%v TIMESTAMPTZ,
		%v BIGINT DEFAULT 0
		);
	`, tableName, IDСolumnName, UserIDColumnName, LocationFKColumnName, locationTableName, LocationIDColumnName, EntryTimeColumnName, ExitTimeColumnName, DurationTimeColumnName)
	_, err := s.db.Exec(context.Background(), create_locations)
	if err != nil {
		return errors.Wrap(err, "initition locations table")
	}

	_, err = s.db.Exec(context.Background(), create_user_visits)
	if err != nil {
		return errors.Wrap(err, "initiation userVisits table")
	}

	return nil
}
