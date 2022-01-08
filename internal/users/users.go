package users

import (
	"context"

	"gbu-telegram-bot/internal/entity"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

// Users is implementation for bot.Users interface
type Users struct {
	db *pgxpool.Pool
}

// New returns bot.Users imlementation
func New(pool *pgxpool.Pool) *Users {
	return &Users{
		db: pool,
	}
}

func (u *Users) Add(ctx context.Context, user entity.User) error {
	query := "INSERT INTO users (id, username, name, is_subscribed) VALUES ($1, $2, $3, $4)"

	_, err := u.db.Exec(ctx, query, user.ID, user.Username, user.Name, user.IsSubscribed)
	if err != nil {
		return errors.Wrap(err, "can't exec insert query")
	}

	return nil
}

func (u *Users) Get(ctx context.Context, id int64) (entity.User, error) {
	query := "SELECT username, name, is_subscribed FROM users WHERE id = $1"

	row := u.db.QueryRow(ctx, query, id)

	user := entity.User{ID: id}
	err := row.Scan(&user.Username, &user.Name, &user.IsSubscribed)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.User{}, entity.ErrUserNotFound
	}
	if err != nil {
		return entity.User{}, errors.Wrap(err, "can't scan")
	}

	return user, nil
}

func (u *Users) Subscribe(ctx context.Context, id int64) error {
	query := "UPDATE users SET is_subscribed = true WHERE id = $1"

	_, err := u.db.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "can't exec query")
	}

	return nil
}

func (u *Users) Unsubscribe(ctx context.Context, id int64) error {
	query := "UPDATE users SET is_subscribed = false WHERE id = $1"

	_, err := u.db.Exec(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "can't exec query")
	}

	return nil
}

func (u *Users) GetSubscribedIDs(ctx context.Context) ([]int64, error) {
	query := "SELECT id FROM users WHERE is_subscribed = true"

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "can't exec select query")
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, errors.Wrap(err, "can't scan row")
		}
		ids = append(ids, id)
	}

	return ids, nil
}
