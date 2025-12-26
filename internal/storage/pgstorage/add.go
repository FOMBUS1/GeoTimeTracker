package pgstorage

import (
	"context"
	"fmt"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/models"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (storage *PGstorage) UpsertUserVisits(ctx context.Context, geoInfos []*models.GeoKafkaMessage) error {
	locationNames := make([]string, 0, len(geoInfos))
	for _, geoInfo := range geoInfos {
		if geoInfo.Location != nil {
			locationNames = append(locationNames, *geoInfo.Location)
		} else {
			locationNames = append(locationNames, geoInfo.LocationAddress)
		}
	}

	locationMap, err := storage.ensureLocations(ctx, locationNames)
	if err != nil {
		return errors.Wrap(err, "ensure locations")
	}

	queries := make([]squirrel.Sqlizer, 0, len(geoInfos))
	for _, userVisit := range geoInfos {
		name := ""
		if userVisit.Location != nil {
			name = *userVisit.Location
		} else {
			name = userVisit.LocationAddress
		}

		locID, ok := locationMap[name]
		if !ok {
			return fmt.Errorf("location ID not found for name: %s", name)
		}

		queries = append(queries, storage.buildSingleVisitQuery(userVisit, locID))
	}

	tx, err := storage.db.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback(ctx)

	for _, q := range queries {
		sql, args, err := q.ToSql()
		if err != nil {
			return errors.Wrap(err, "to sql error")
		}
		if _, err := tx.Exec(ctx, sql, args...); err != nil {
			return errors.Wrap(err, "exec query error")
		}
	}

	return tx.Commit(ctx)
}

func (storage *PGstorage) buildSingleVisitQuery(userVisit *models.GeoKafkaMessage, locID uint64) squirrel.Sqlizer {
	if !userVisit.Departure {
		return squirrel.Insert(tableName).
			Columns(UserIDColumnName, LocationFKColumnName, EntryTimeColumnName, DurationTimeColumnName).
			Values(userVisit.UserId, locID, userVisit.Time.AsTime(), 0).
			PlaceholderFormat(squirrel.Dollar)
	}

	return squirrel.Update(tableName).
		Set(ExitTimeColumnName, userVisit.Time.AsTime()).
		Set(DurationTimeColumnName, squirrel.Expr("EXTRACT(EPOCH FROM (?::timestamp - "+EntryTimeColumnName+"))", userVisit.Time.AsTime())).
		Where(squirrel.Eq{
			UserIDColumnName:     userVisit.UserId,
			LocationFKColumnName: locID,
			ExitTimeColumnName:   nil,
		}).
		PlaceholderFormat(squirrel.Dollar)
}

func (storage *PGstorage) ensureLocations(ctx context.Context, locations []string) (map[string]uint64, error) {
	if len(locations) == 0 {
		return make(map[string]uint64), nil
	}

	uniqueLocations := lo.Uniq(locations)

	builder := squirrel.Insert(locationTableName).
		Columns(LocationColumnName).
		PlaceholderFormat(squirrel.Dollar)

	for _, loc := range uniqueLocations {
		builder = builder.Values(loc)
	}

	query, args, err := builder.
		Suffix("ON CONFLICT (" + LocationColumnName + ") DO UPDATE SET " + LocationColumnName + " = EXCLUDED." + LocationColumnName).
		Suffix("RETURNING " + LocationIDColumnName + ", " + LocationColumnName).
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "build ensure locations query")
	}

	rows, err := storage.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "exec ensure locations")
	}
	defer rows.Close()

	locationMap := make(map[string]uint64)
	for rows.Next() {
		var id uint64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, errors.Wrap(err, "scan location row")
		}
		locationMap[name] = id
	}

	return locationMap, nil
}
