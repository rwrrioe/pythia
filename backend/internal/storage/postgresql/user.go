package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rwrrioe/pythia/backend/internal/auth"
	"github.com/rwrrioe/pythia/backend/internal/domain/entities"
	"github.com/rwrrioe/pythia/backend/internal/storage/models"
)

type UserStorage struct {
	conn *pgx.Conn
}

func NewUserStorage(conn *pgx.Conn) (*UserStorage, error) {
	const op = "storage.postgres.New"

	return &UserStorage{conn: conn}, nil
}

func (s *UserStorage) GetUser(ctx context.Context) (*entities.User, error) {
	const op = "postgresql.UserStorage.GetUser"

	var user models.User
	userId, ok := auth.UIDFromContext(ctx)
	if !ok {
		return nil, ErrUserNotFound
	}

	if err := s.conn.QueryRow(ctx,
		`SELECT * 
									FROM users u 
									JOIN languages l ON u.lang_id = l.id
									JOIN levels lv ON u.level_id = lv.id
									WHERE u.user_id=$1
									`, userId).Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return &entities.User{
		Email:       user.Email,
		Name:        user.Name,
		Level:       user.Level,
		Lang:        user.Lang,
		WordsPerDay: user.WordsPerDay,
	}, nil
}
