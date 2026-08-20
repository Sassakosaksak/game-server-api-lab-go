package player

import (
	"context"
	"errors"
	"time"
)

// ErrNotFoundは指定されたPlayerが存在しないことを表す。
var ErrNotFound = errors.New("player not found")

type Player struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// RepositoryはPlayerの保存先に関係なく、Handlerが必要とする操作を定義する。
type Repository interface {
	Create(ctx context.Context, name string) (Player, error)
	FindByID(ctx context.Context, id string) (Player, error)
}
