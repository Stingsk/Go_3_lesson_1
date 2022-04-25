package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

const (
	psqlDriverName = "pgx"
	psqlTimeout    = 1 * time.Second
)

type DBStore struct {
	connection   *sql.DB
	context      context.Context
	metricsCache *Metric
}

func NewDBStore(ctx context.Context, dataBaseConnectionString string) (*DBStore, error) {
	var db DBStore

	conn, err := sql.Open(psqlDriverName, dataBaseConnectionString)
	if err != nil {
		return nil, err
	}

	db = DBStore{
		connection:   conn,
		context:      ctx,
		metricsCache: &Metric{},
	}

	return &db, nil
}

func (db *DBStore) Ping() error {
	ctx, cancel := context.WithTimeout(db.context, psqlTimeout)
	defer cancel()

	return db.connection.PingContext(ctx)
}

func (db *DBStore) NewMetric(value string, metricType string, name string) (Metric, error) {
	return Metric{}, nil
}
func (db *DBStore) UpdateMetric(metricResourceMap *MetricResourceMap, metric Metric) (Metric, error) {
	return Metric{}, nil
}
func (db *DBStore) UpdateMetricByParameters(metricResourceMap *MetricResourceMap, metricName string, metricType string, value string) (Metric, error) {
	return Metric{}, nil
}

func (db *DBStore) Close() error {
	logrus.Info("Close database connection")

	return db.connection.Close()
}
