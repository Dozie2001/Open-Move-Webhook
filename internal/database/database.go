package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/open-move/intercord/internal/config"
	"github.com/open-move/intercord/internal/models"
)

func New(cfg *config.DatabaseConfig) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	return db
}

func RunMigrations(db *bun.DB) error {
	ctx := context.Background()

	models := []interface{}{
		(*models.User)(nil),
		(*models.Team)(nil),
		(*models.TeamMembership)(nil),
		(*models.Subscription)(nil),
		(*models.Channel)(nil),
		(*models.SubscriptionChannel)(nil),
		(*models.Notification)(nil),
		(*models.PasswordReset)(nil),
		(*models.EmailVerification)(nil),
	}

	for _, model := range models {
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			log.Printf("Error creating table for model %T: %v", model, err)
			return err
		}
	}

	return nil
}
