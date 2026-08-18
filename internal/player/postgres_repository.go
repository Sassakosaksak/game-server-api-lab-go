package player

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &postgresRepository{pool: pool}
}

func (repository *postgresRepository) Create(ctx context.Context, name string) (Player, error) {
	const query = `
		INSERT INTO players (name)
		VALUES ($1)
		RETURNING id, name, created_at
	`

	var player Player
	if err := repository.pool.QueryRow(ctx, query, name).Scan(&player.ID, &player.Name, &player.CreatedAt); err != nil {
		return Player{}, fmt.Errorf("Playerの保存に失敗しました: %w", err)
	}

	return player, nil
}

func (repository *postgresRepository) FindByID(ctx context.Context, id string) (Player, error) {
	const query = `
		SELECT id, name, created_at
		FROM players
		WHERE id = $1
	`

	var player Player
	if err := repository.pool.QueryRow(ctx, query, id).Scan(&player.ID, &player.Name, &player.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Player{}, ErrNotFound
		}

		return Player{}, fmt.Errorf("Playerの取得に失敗しました: %w", err)
	}

	return player, nil
}
