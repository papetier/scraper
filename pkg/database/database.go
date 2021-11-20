package database

import (
	"context"
	"errors"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/papetier/crawler/pkg/config"
	log "github.com/sirupsen/logrus"
)

type Connection struct {
	Config *config.DBConfig
	Pool   *pgxpool.Pool
}

var dbConnection *Connection
var ErrNotFound = errors.New("item not found")

func Connect() {
	// parse db config
	poolConfig, err := pgxpool.ParseConfig(config.DB.ConnectionString())
	if err != nil {
		log.Fatalf("unable to parse the postgres connection string: %s", err)
	}

	// register types
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtype.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}

	// connect to db
	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("creating the postgres connection pool: %s", err)
	}

	// set global connection pool
	dbConnection = &Connection{
		Config: config.DB,
		Pool:   pool,
	}

	log.Infof("successfully connected to the postgres database")
}

func CloseConnection() {
	dbConnection.Pool.Close()
	log.Infof("closed connection to the postgres database")
}
