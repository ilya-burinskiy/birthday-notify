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

func (db *DBStorage) FetchNotificationsForCurrentDate(ctx context.Context) ([]models.Notification, error) {
	rows, err := db.pool.Query(
		ctx,
		`SELECT "subscribing_users"."email" AS "subscribing_user_email",
		        COALESCE("days_before_notify", 1) AS "days_before_notify",
				"subscribed_users"."email" AS "subscribed_user_email"
		 FROM "subscriptions"
		 INNER JOIN "users" AS "subscribed_users" ON "subscriptions"."subscribed_user_id" = "subscribed_users"."id"
		 INNER JOIN "users" AS "subscribing_users" ON "subscriptions"."subscribing_user_id" = "subscribing_users"."id"
		 LEFT JOIN "notify_settings" ON "subscribing_users"."id" = "notify_settings"."user_id"
		 WHERE EXTRACT(DAY FROM CURRENT_DATE + COALESCE("days_before_notify", 1)) = EXTRACT(DAY FROM "subscribed_users"."birthdate")
		   AND EXTRACT(MONTH FROM CURRENT_DATE + COALESCE("days_before_notify", 1)) = EXTRACT(MONTH FROM "subscribed_users"."birthdate")`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}

	result, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.Notification, error) {
		var notification models.Notification
		err := row.Scan(
			&notification.SubscribingUserEmail,
			&notification.DaysBeforeNotify,
			&notification.SubscribedUserEmail,
		)
		return notification, err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch user secrets: %w", err)
	}

	return result, nil
}

func (db *DBStorage) CreateNotificationSetting(ctx context.Context, userID, daysBeforeNotify int) (models.NotifySetting, error) {
	row := db.pool.QueryRow(
		ctx,
		`INSERT INTO "notify_settings"("user_id", "days_before_notify") VALUES ($1, $2) RETURNING "id"`,
		userID,
		daysBeforeNotify,
	)
	notifySetting := models.NotifySetting{UserID: userID, DaysBeforeNotify: daysBeforeNotify}
	err := row.Scan(&notifySetting.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) {
			return notifySetting, fmt.Errorf("failed to create notify setting: %w", err)
		}

		// TODO: use proper error types instead of fmt.Errorf
		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return notifySetting, fmt.Errorf("user with id=%d does not exists", userID)
		case pgerrcode.UniqueViolation:
			return notifySetting, fmt.Errorf("user with id=%d already has notify setting", userID)
		default:
			return notifySetting, nil
		}
	}

	return notifySetting, nil
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
