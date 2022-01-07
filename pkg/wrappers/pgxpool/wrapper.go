package pgxpool

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Connect is wrapper for original pgxpool.Connect function but with
// credential-parameters instead of connection string
func Connect(ctx context.Context, host, user, pass, dbName string) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s/%s", user, pass, host, dbName)
	return pgxpool.Connect(ctx, connString)
}

// Forward type
type Pool = pgxpool.Pool
