package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/J0es1ick/shortli/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) SaveUser(User *models.User) (int64, error) {
	stmt, err := r.db.Prepare("INSERT INTO User_info (username, email, password_hash) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("preparing insert operator error: %v", err)
	}

	res, err := stmt.Exec(User.Username, User.Email, User.PasswordHash)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, fmt.Errorf("unique constraint violation: %v", err)
			}
		}
		return 0, fmt.Errorf("insert value error: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %v", err)
	}

	return id, nil
}

func (r *UserRepository) FindUserById(id int) (*models.User, error) {
	stmt, err := r.db.Prepare("SELECT id, username, email, password_hash FROM users WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("preparing select operator error: %v",  err)
	}

	User := &models.User{}
    
    err = stmt.QueryRow(id).Scan(&User.ID, &User.Username, &User.Email, &User.PasswordHash)
	if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("scan row: %v", err)
    }

	return User, nil
}

func (r *UserRepository) FindUserByUsername(username string) (*models.User, error) {
    stmt, err := r.db.Prepare("SELECT id, username, email, password_hash FROM users WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("preparing select operator error: %v",  err)
	}

	User := &models.User{}
    
    err = stmt.QueryRow(username).Scan(&User.ID, &User.Username, &User.Email, &User.PasswordHash)
	if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("scan row: %v", err)
    }

	return User, nil
}

func (r *UserRepository) UpdateUserById(User *models.User) error {
	stmt, err := r.db.Prepare("UPDATE User_info SET username = ?, email = ?, password_hash = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("preparing update operator error: %v", err)
	}

	_, err = stmt.Exec(User.Username, User.Email, User.PasswordHash, User.ID)
	if err != nil {
		return fmt.Errorf("update value error: %v", err)
	}

	return nil
}

func (r *UserRepository) DeleteUserById(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return fmt.Errorf("preparing delete operator error: %v", err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("delete value error: %v", err)
	}

	return nil
}