package users

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"persons/internal/domain/models"
	"persons/internal/storage/postgres"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(user models.User) error {
	const op = "postgres.users.Repository.Save"
	query := `INSERT INTO 
    			persons 
    			(id, email, login, password_hash) 
				VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(query, user.Id, user.Email, user.Login, user.PasswordHash)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, postgres.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) GetById(id uuid.UUID) (*models.User, error) {
	const op = "postgres.users.Repository.GetById"
	var user models.User
	query := `SELECT 
    			id, email, login, password_hash, about_me, image 
				FROM
				    persons 
				WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, postgres.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &user, nil
}

func (r *Repository) GetAll() ([]models.User, error) {
	const op = "postgres.users.Repository.GetAll"
	var persons []models.User
	query := `SELECT 
				id, email, login, password_hash, about_me, image
				FROM 
				    persons`
	err := r.db.Select(&persons, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return persons, nil
}
