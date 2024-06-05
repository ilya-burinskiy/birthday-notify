package storage

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/ilya-burinskiy/birthday-notify/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBStorage struct {
	pool *pgxpool.Pool
}

func NewDBStorage(dsn string) (*DBStorage, error) {
	if err := runMigrations(dsn); err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %w", err)
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create a connection pool: %w", err)
	}

	return &DBStorage{
		pool: pool,
	}, nil
}

func (db *DBStorage) CreateUser(ctx context.Context, email string, encryptedPassword []byte, birthDate time.Time) (models.User, error) {
	row := db.pool.QueryRow(
		ctx,
		`INSERT INTO "users" ("email", "encrypted_password", "birthdate") VALUES ($1, $2, $3) RETURNING "id"`,
		email,
		encryptedPassword,
		birthDate,
	)
	user := models.User{Email: email, EncryptedPassword: encryptedPassword, BirthDate: birthDate}
	err := row.Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return user, ErrUserNotUniq{User: user}
		}
		return user, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (db *DBStorage) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	row := db.pool.QueryRow(
		ctx,
		`SELECT "id", "encrypted_password", "birthdate"
		 FROM "users"
		 WHERE "email" = $1`,
		email,
	)
	user := models.User{Email: email}
	err := row.Scan(&user.ID, &user.EncryptedPassword, &user.BirthDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, ErrUserNotFound{User: user}
		}

		return user, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

func (db *DBStorage) CreateSubscription(ctx context.Context, subscribedUserID, subscribingUserID int) (models.Subscription, error) {
	row := db.pool.QueryRow(
		ctx,
		`INSERT INTO "subscriptions" ("subscribed_user_id", "subscribing_user_id") VALUES ($1, $2) RETURNING "id"`,
		subscribedUserID,
		subscribingUserID,
	)
	subscription := models.Subscription{SubscribedUserID: subscribedUserID, SubscribingUserID: subscribingUserID}
	err := row.Scan(&subscription.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return subscription, fmt.Errorf("failed to create subscription: %w", err)
		}

		// TODO: use proper error types instead of fmt.Errof
		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return subscription, fmt.Errorf("user with id=%d not found", subscribedUserID)
		case pgerrcode.UniqueViolation:
			return subscription, fmt.Errorf("user with id=%d already subscribed to user with id=%d", subscribingUserID, subscribedUserID)
		default:
			return subscription, pgErr
		}
	}

	return subscription, nil
}

func (db *DBStorage) FindSubscription(ctx context.Context, subscribedUserID, subscribingUserID int) (models.Subscription, error) {
	row := db.pool.QueryRow(
		ctx,
		`SELECT "id" FROM "subscriptions" WHERE "subscribed_user_id" = $1 AND "subscribing_user_id" = $2`,
		subscribedUserID,
		subscribingUserID,
	)
	subscription := models.Subscription{SubscribedUserID: subscribedUserID, SubscribingUserID: subscribingUserID}
	err := row.Scan(&subscription.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return subscription, ErrSubscriptionNotFound{Subscription: subscription}
		}
		return subscription, fmt.Errorf("failed to find subscription: %w", err)
	}
	return subscription, nil
}

func (db *DBStorage) DeleteSubscription(ctx context.Context, subscriptionID int) error {
	_, err := db.pool.Exec(ctx, `DELETE FROM "subscriptions" WHERE "id" = $1`, subscriptionID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription with id=%d: %w", subscriptionID, err)
	}
	return nil
}

//go:embed db/migrations/*.sql
var migrationsDir embed.FS

func runMigrations(dsn string) error {
	d, err := iofs.New(migrationsDir, "db/migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
	}

	return nil
}
