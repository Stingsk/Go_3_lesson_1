package storage

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/Stingsk/Go_3_lesson_1/db/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

const (
	psqlDriverName      = "pgx"
	migrationSourceName = "go-bindata"
)

type DBStore struct {
	connection *sql.DB
}

func NewDBStore(dataBaseConnectionString string) (*DBStore, error) {
	var db DBStore

	conn, err := sql.Open(psqlDriverName, dataBaseConnectionString)
	if err != nil {
		return nil, err
	}

	db = DBStore{
		connection: conn,
	}

	if err := db.migrate(); err != nil {
		return nil, err
	}

	return &db, nil
}

func (db *DBStore) UpdateMetrics(ctx context.Context, metricsBatch []*Metric) error {
	tx, err := db.connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmtInsertGauge, err := tx.Prepare("INSERT INTO gauge (metric_id, metric_value) VALUES ($1, $2) " +
		"ON CONFLICT (metric_id) DO UPDATE SET metric_value = $2")
	if err != nil {
		return err
	}
	defer func(stmtInsertGauge *sql.Stmt) {
		if err := stmtInsertGauge.Close(); err != nil {
			logrus.Error("failed to close insert statement")
		}
	}(stmtInsertGauge)

	stmtSelectCounter, err := tx.Prepare("SELECT metric_delta FROM counter WHERE metric_id = $1")
	if err != nil {
		return err
	}
	defer func(stmtSelectCounter *sql.Stmt) {
		if err := stmtSelectCounter.Close(); err != nil {
			logrus.Error("failed to close insert statement")
		}
	}(stmtSelectCounter)

	stmtInsertCounter, err := tx.Prepare("INSERT INTO counter (metric_id, metric_delta) VALUES ($1, $2) " +
		"ON CONFLICT (metric_id) DO UPDATE SET metric_delta = $2")
	if err != nil {
		return err
	}
	defer func(stmtInsertCounter *sql.Stmt) {
		if err := stmtInsertCounter.Close(); err != nil {
			logrus.Error("failed to close insert statement")
		}
	}(stmtInsertCounter)

	for _, metric := range metricsBatch {
		switch {
		case metric.MType == MetricTypeGauge:
			if _, err := stmtInsertGauge.Exec(metric.ID, *(metric.Value)); err != nil {
				if err := tx.Rollback(); err != nil {
					logrus.Error("unable to rollback transaction")
				}

				return err
			}
		case metric.MType == MetricTypeCounter:
			var counter int64
			query := stmtSelectCounter.QueryRow(metric.ID)

			err = query.Scan(&counter)
			if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
				if err := tx.Rollback(); err != nil {
					logrus.Error("unable to rollback transaction")
				}

				return err
			}

			counter += *(metric.Delta)

			if _, err := stmtInsertCounter.Exec(metric.ID, counter); err != nil {
				if err := tx.Rollback(); err != nil {
					logrus.Error("unable to rollback transaction")
				}

				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *DBStore) UpdateMetricByParameters(ctx context.Context, metricName string, metricType string, value string) error {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return errors.New("only numbers  params in request are allowed")
	}

	if strings.ToLower(metricType) == MetricTypeGauge {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		db.updateGaugeMetric(ctx, metricName, v)
	} else if strings.ToLower(metricType) == MetricTypeCounter {
		newValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		db.updateCounterMetric(ctx, metricName, newValue)
	}

	return nil
}

func (db *DBStore) UpdateMetric(ctx context.Context, metric Metric) error {
	if metric.ID == "" {
		return errors.New("id is empty")
	}

	if metric.MType == "" {
		return errors.New("type is empty")
	}

	if strings.ToLower(metric.MType) == MetricTypeGauge {
		db.updateGaugeMetric(ctx, metric.ID, *metric.Value)
	} else if strings.ToLower(metric.MType) == MetricTypeCounter {
		db.updateCounterMetric(ctx, metric.ID, *metric.Delta)
	}

	return nil
}

func (db *DBStore) Ping(_ context.Context) error {
	return db.connection.Ping()
}

func (db *DBStore) GetMetric(ctx context.Context, metricName string, metricType string) (*Metric, error) {
	metric := Metric{
		ID:    metricName,
		MType: metricType,
	}

	switch metricType {
	case MetricTypeCounter:
		var counter int64
		row := db.connection.QueryRowContext(ctx,
			"SELECT metric_delta FROM counter WHERE metric_id = $1", metricName)

		err := row.Scan(&counter)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("metric not found")
		case !errors.Is(err, nil):
			return nil, err
		}
		metric.Delta = &counter
	case MetricTypeGauge:
		var gauge float64
		row := db.connection.QueryRowContext(ctx,
			"SELECT metric_value FROM gauge WHERE metric_id = $1", metricName)

		err := row.Scan(&gauge)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("metric not found")
		case !errors.Is(err, nil):
			return nil, err
		}
		metric.Value = &gauge
	default:
		return nil, errors.New("metric not found")
	}

	return &metric, nil
}

func (db *DBStore) GetMetrics(ctx context.Context) (map[string]*Metric, error) {
	metricsMap := make(map[string]*Metric)

	counters, err := db.connection.QueryContext(ctx,
		"SELECT metric_id,metric_delta FROM counter")

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Error("Couldn't close rows")
		}
	}(counters)

	for counters.Next() {
		var counter int64
		metric := Metric{
			MType: MetricTypeCounter,
			Delta: &counter,
		}
		err = counters.Scan(&metric.ID, metric.Delta)
		if err != nil {
			return nil, err
		}

		metricsMap[metric.ID] = &metric
	}

	err = counters.Err()
	if err != nil {
		return nil, err
	}

	gauges, err := db.connection.QueryContext(ctx,
		"SELECT metric_id,metric_value FROM gauge")

	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logrus.Error("Couldn't close rows")
		}
	}(gauges)

	for gauges.Next() {
		var gauge float64
		metric := Metric{
			MType: MetricTypeGauge,
			Value: &gauge,
		}

		err = gauges.Scan(&metric.ID, metric.Value)
		if err != nil {
			return nil, err
		}

		metricsMap[metric.ID] = &metric
	}

	err = gauges.Err()
	if err != nil {
		return nil, err
	}

	return metricsMap, nil
}

func (db *DBStore) Close() error {
	return db.connection.Close()
}

func (db *DBStore) updateCounterMetric(ctx context.Context, metricName string, metricData int64) error {
	var counter int64
	row := db.connection.QueryRowContext(ctx,
		"SELECT metric_delta FROM counter WHERE metric_id = $1", metricName)

	err := row.Scan(&counter)
	if !errors.Is(err, nil) && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	counter += metricData
	_, err = db.connection.ExecContext(ctx,
		"INSERT INTO counter (metric_id, metric_delta) VALUES ($1, $2) "+
			"ON CONFLICT (metric_id) DO UPDATE SET metric_delta = $2",
		metricName, counter)

	return err
}

func (db *DBStore) updateGaugeMetric(ctx context.Context, metricName string, metricData float64) error {
	_, err := db.connection.ExecContext(ctx,
		"INSERT INTO gauge (metric_id, metric_value) VALUES ($1, $2) "+
			"ON CONFLICT (metric_id) DO UPDATE SET metric_value = $2",
		metricName, metricData)

	return err
}

func (db *DBStore) migrate() error {
	data := bindata.Resource(migrations.AssetNames(), migrations.Asset)

	sourceDriver, err := bindata.WithInstance(data)
	if err != nil {
		return err
	}

	dbDriver, err := postgres.WithInstance(db.connection, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithInstance(migrationSourceName, sourceDriver, psqlDriverName, dbDriver)
	if err != nil {
		return err
	}

	if err := migration.Up(); !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
