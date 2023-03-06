package data

import (
	"context"
	"database/sql"
	"time"

	entDialectSQL "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/lib/pq" // postgres driver for Go's database/sql package
	"github.com/phlx-ru/hatchet/metrics"
	"github.com/phlx-ru/hatchet/slices"

	"storage/ent"
	"storage/internal/conf"
)

const (
	maxOpenConnections = 32
	maxIdleConnections = 30
	maxConnLifetime    = 5 * time.Minute
	sendStatsEvery     = time.Second
)

// ProviderRepoSet is data providers.
var ProviderRepoSet = wire.NewSet(NewFileRepo)

var ProviderDataSet = wire.NewSet(NewData)

// Data .
type Data struct {
	db     *sql.DB
	ent    *ent.Client
	logger *log.Helper
}

type Database interface {
	DB() *sql.DB
	Ent() *ent.Client
	MigrateSoft(ctx context.Context) error
	MigrateHard(ctx context.Context) error
	Prepare(ctx context.Context, m conf.Data_Database_Migrate) error
	CollectDatabaseMetrics(ctx context.Context, metric metrics.Metrics)
	Seed(ctx context.Context, seeding func(context.Context, *ent.Client) error) error
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (Database, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "ent/data/logger-job"))

	drv, err := entDialectSQL.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}
	// Get the underlying sql.DB object of the driver.
	db := drv.DB()
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(maxConnLifetime)
	options := []ent.Option{
		ent.Driver(drv),
	}
	if c.Database.Debug {
		options = append(options, ent.Debug())
	}
	client := ent.NewClient(options...)

	cleanup := func() {
		logHelper.Info("closing database client from cleanup() function")
		if client != nil {
			err := client.Close()
			if err != nil {
				logHelper.Errorf(`failed to close database client: %v`, err)
			}
		}
	}
	return &Data{
		db:     db,
		ent:    client,
		logger: logHelper,
	}, cleanup, nil
}

func (d *Data) DB() *sql.DB {
	return d.db
}

func (d *Data) Ent() *ent.Client {
	return d.ent
}

// MigrateSoft only creates and updates schema entities
func (d *Data) MigrateSoft(ctx context.Context) error {
	err := d.ent.Schema.Create(ctx, schema.WithForeignKeys(false))
	if err != nil {
		d.logger.WithContext(ctx).Errorf(`failed to soft migrate database schema: %v`, err)
		return err
	}
	return nil
}

// MigrateHard does same as MigrateSoft, but also drop columns and indices
func (d *Data) MigrateHard(ctx context.Context) error {
	err := d.ent.Schema.Create(ctx, schema.WithDropIndex(true), schema.WithDropColumn(true))
	if err != nil {
		d.logger.WithContext(ctx).Errorf(`failed to hard migrate database schema: %v`, err)
		return err
	}
	return nil
}

func (d *Data) Prepare(ctx context.Context, m conf.Data_Database_Migrate) error {
	var err error
	if m == conf.Data_Database_none {
		return nil
	}
	if m == conf.Data_Database_soft {
		d.logger.WithContext(ctx).Info("preparing database: running soft migrate")
		err = d.MigrateSoft(ctx)
	}
	if m == conf.Data_Database_hard {
		d.logger.WithContext(ctx).Info("preparing database: running hard migrate")
		err = d.MigrateHard(ctx)
	}
	migrateValuesAllowedSeeding := []conf.Data_Database_Migrate{
		conf.Data_Database_soft,
		conf.Data_Database_hard,
	}
	if err == nil && slices.Includes(m, migrateValuesAllowedSeeding) {
		d.logger.WithContext(ctx).Info("preparing database: running seeders")
		// err = d.Seed(ctx, SeedMainEntities) // TODO: Add or remove seeding
	}
	return err
}

func (d *Data) CollectDatabaseMetrics(ctx context.Context, metric metrics.Metrics) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		stats := d.db.Stats()

		// The number of established connections both in use and idle.
		metric.Gauge(`postgres.connections.open`, stats.OpenConnections)

		// The number of connections currently in use.
		metric.Gauge(`postgres.connections.used`, stats.InUse)

		// The number of idle connections.
		metric.Gauge(`postgres.connections.idle`, stats.Idle)

		// The total number of connections waited for.
		metric.Gauge(`postgres.connections.wait`, stats.WaitCount)

		// The total time blocked waiting for a new connection.
		// metric.Gauge(`postgres.connections.wait_duration`, stats.WaitDuration) // TODO Duration or count ms?

		// The total number of connections closed due to SetMaxIdleConns.
		metric.Gauge(`postgres.connections.max_idle_closed`, stats.MaxIdleClosed)

		// The total number of connections closed due to SetConnMaxIdleTime.
		metric.Gauge(`postgres.connections.max_idle_time_closed`, stats.MaxIdleTimeClosed)

		// The total number of connections closed due to SetConnMaxLifetime.
		metric.Gauge(`postgres.connections.max_lifetime_closed`, stats.MaxLifetimeClosed)

		time.Sleep(sendStatsEvery)
	}
}

// Seed everything you need by passing the seeding func
func (d *Data) Seed(ctx context.Context, seeding func(context.Context, *ent.Client) error) error {
	return seeding(ctx, d.ent)
}

// client return client by tx in context if it exists or default ent client
func client(data Database) func(ctx context.Context) *ent.Client {
	return func(ctx context.Context) *ent.Client {
		if client := ent.FromContext(ctx); client != nil {
			return client
		}
		return data.Ent()
	}
}
