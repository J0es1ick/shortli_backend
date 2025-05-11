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

func (r *UserRepository) SaveUser(user *models.User) (int64, error) {
	query := `
		INSERT INTO User_info 
			(username, email, password_hash) 
		VALUES ($1, $2, $3)
		RETURNING user_id`

	var id int64
	err := r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&id)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("user with this code already exists")
		}
		return 0, fmt.Errorf("insert value error: %v", err)
	}

	return id, nil
}

func (r *UserRepository) FindUserById(id int) (*models.User, error) {
	query := `
		SELECT 
			id,
			username, 
			email, 
			password_hash 
		FROM User_info 
		WHERE id = $1`

	user := &models.User{}   
    err := r.db.QueryRow(query, id).Scan(
		&user.ID, 
		&user.Username, 
		&user.Email, 
		&user.PasswordHash,
	)
	
	if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("scan row: %v", err)
    }

	return user, nil
}

func (r *UserRepository) FindUserByUsername(username string) (*models.User, error) {
    query := `
		SELECT 
			id,
			username, 
			email, 
			password_hash 
		FROM User_info 
		WHERE username = $1`

	user := &models.User{}  
    err := r.db.QueryRow(query, username).Scan(
		&user.ID, 
		&user.Username, 
		&user.Email, 
		&user.PasswordHash,
	)

	if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("scan row: %v", err)
    }

	return user, nil
}

func (r *UserRepository) UpdateUserById(user *models.User) error {
	query := `
		UPDATE User_info 
		SET 
			username = $1, 
			email = $2, 
			password_hash = $3
		WHERE id = $4`


	result, err := r.db.Exec(
        query,
        user.Username, 
		user.Email, 
		user.PasswordHash, 
		user.ID,
    )	

	if err != nil {
		return fmt.Errorf("update value error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("no rows updated - user with id %d not found", user.ID)
    }
    
    return nil
}

func (r *UserRepository) DeleteUserById(id int) error {
	query := `DELETE FROM User_info WHERE id = ?`

	var deletedID int64
    err := r.db.QueryRow(query, id).Scan(&deletedID)

	if err != nil {
		if err == sql.ErrNoRows {
            return fmt.Errorf("url with code '%d' not found", deletedID)
        }
		return fmt.Errorf("delete value error: %v", err)
	}

	return nil
}