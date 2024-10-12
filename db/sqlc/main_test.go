package db

import (
	"context"
	"github.com/WooDMaNbtw/BankApp/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")

	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
