package pgstorage

import (
	"context"
	"fmt"

	"github.com/FOMBUS1/GeoTimeTracker/internal/pb/geo_stats_api"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func (storage *PGstorage) GetStatsByLocations(ctx context.Context, reqs []*geo_stats_api.UserLocationRequest) ([]*geo_stats_api.TimeSpentResponse, error) {
	type key struct {
		userID uint64
		period geo_stats_api.Period
	}

	grouped := lo.GroupBy(reqs, func(r *geo_stats_api.UserLocationRequest) key {
		return key{r.UserId, r.Period}
	})

	var responses []*geo_stats_api.TimeSpentResponse

	for k, requests := range grouped {
		locIDs := lo.Map(requests, func(r *geo_stats_api.UserLocationRequest, _ int) uint64 {
			return r.LocationId
		})

		query, args, err := selectStatsQuery(k.userID, locIDs, k.period).ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "build sql")
		}

		rows, err := storage.db.Query(ctx, query, args...)
		if err != nil {
			return nil, errors.Wrap(err, "query db")
		}
		defer rows.Close()

		resp := &geo_stats_api.TimeSpentResponse{UserId: k.userID}
		for rows.Next() {
			stat := &geo_stats_api.LocationStat{Period: k.period}
			if err := rows.Scan(&stat.LocationId, &stat.TimeSeconds); err != nil {
				return nil, errors.Wrap(err, "scan row")
			}
			resp.Stats = append(resp.Stats, stat)
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

func (storage *PGstorage) GetTopStats(ctx context.Context, reqs []*geo_stats_api.TimePeriodRequest) ([]*geo_stats_api.TimeSpentResponse, error) {
	var responses []*geo_stats_api.TimeSpentResponse

	for _, req := range reqs {
		limit := uint32(5)
		if req.TopK != nil {
			limit = *req.TopK
		}

		query, args, err := selectTopStatsQuery(req.UserId, limit, req.Period).ToSql()
		if err != nil {
			return nil, errors.Wrap(err, "build top query")
		}

		rows, err := storage.db.Query(ctx, query, args...)
		if err != nil {
			return nil, errors.Wrap(err, "exec top query")
		}
		defer rows.Close()

		userResp := &geo_stats_api.TimeSpentResponse{UserId: req.UserId}
		for rows.Next() {
			stat := &geo_stats_api.LocationStat{Period: req.Period}
			if err := rows.Scan(&stat.LocationId, &stat.TimeSeconds); err != nil {
				return nil, errors.Wrap(err, "scan top row")
			}
			userResp.Stats = append(userResp.Stats, stat)
		}
		responses = append(responses, userResp)
	}

	return responses, nil
}

func switch_period(q squirrel.SelectBuilder, period geo_stats_api.Period) squirrel.SelectBuilder {
	switch period {
	case geo_stats_api.Period_DAY:
		q = q.Where(fmt.Sprintf("%v > now() - interval '1 day'", EntryTimeColumnName))
	case geo_stats_api.Period_WEEK:
		q = q.Where(fmt.Sprintf("%v > now() - interval '1 week'", EntryTimeColumnName))
	case geo_stats_api.Period_MONTH:
		q = q.Where(fmt.Sprintf("%v > now() - interval '1 month'", EntryTimeColumnName))
	case geo_stats_api.Period_YEAR:
		q = q.Where(fmt.Sprintf("%v > now() - interval '1 year'", EntryTimeColumnName))
	case geo_stats_api.Period_LAST:
		return q.OrderBy(fmt.Sprintf("%s DESC", EntryTimeColumnName)).Limit(1)
	default:
		return q
	}
	return q
}

func selectStatsQuery(userID uint64, locationIDs []uint64, period geo_stats_api.Period) squirrel.SelectBuilder {
	q := squirrel.Select(LocationFKColumnName, fmt.Sprintf("SUM(%v) as total_time", DurationTimeColumnName)).
		From(tableName).
		Where(squirrel.Eq{UserIDColumnName: userID, LocationFKColumnName: locationIDs})

	q = switch_period(q, period)

	return q.GroupBy(LocationFKColumnName).PlaceholderFormat(squirrel.Dollar)
}

func selectTopStatsQuery(userID uint64, limit uint32, period geo_stats_api.Period) squirrel.SelectBuilder {
	q := squirrel.Select(LocationFKColumnName, fmt.Sprintf("SUM(%v) as total_time", DurationTimeColumnName)).
		From(tableName).
		Where(squirrel.Eq{UserIDColumnName: userID})

	q = switch_period(q, period)

	q = q.GroupBy(LocationFKColumnName).
		OrderBy("total_time DESC").
		PlaceholderFormat(squirrel.Dollar)

	if limit > 0 {
		q = q.Limit(uint64(limit))
	}
	return q
}
