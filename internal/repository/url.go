package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/J0es1ick/shortli/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UrlRepository struct {
	db *sqlx.DB
}

func NewUrlRepository(db *sqlx.DB) *UrlRepository {
	return &UrlRepository{
		db: db,
	}
}

func (r *UrlRepository) SaveUrl(url *models.URL) (int64, error) {
    query := `
        INSERT INTO url_info 
            (original_url, short_code, click_count, qr_click_count, created_at) 
        VALUES ($1, $2, $3, $4, $5)
        RETURNING url_id
    `
    
    var id int64
    err := r.db.QueryRow(
        query,
        url.OriginalURL,
        url.ShortCode,
        url.ClickCount,
        url.QRClickCount,
        url.CreatedAt,
    ).Scan(&id)
    
    if err != nil {
        var pgErr *pq.Error
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return 0, fmt.Errorf("url with this code already exists")
        }
        return 0, fmt.Errorf("insert value error: %v", err)
    }
    
    return id, nil
}

func (r *UrlRepository) FindUrlByCode(code string) (*models.URL, error) {
    query := `
        SELECT 
            url_id, 
            original_url, 
            short_code, 
            click_count, 
            qr_click_count, 
            created_at 
        FROM url_info 
        WHERE short_code = $1
    `
    
    url := &models.URL{}
    err := r.db.QueryRow(query, code).Scan(
        &url.ID,
        &url.OriginalURL,
        &url.ShortCode,
        &url.ClickCount,
        &url.QRClickCount,
        &url.CreatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("url not found")
        }
        return nil, fmt.Errorf("select error: %v", err)
    }
    
    return url, nil
}

func (r *UrlRepository) UpdateUrlByCode(url *models.URL) error {
    query := `
        UPDATE url_info 
        SET 
            original_url = $1, 
            click_count = $2, 
            qr_click_count = $3,
            created_at = $4
        WHERE short_code = $5
    `
    
    result, err := r.db.Exec(
        query,
        url.OriginalURL,
        url.ClickCount,
        url.QRClickCount,
        url.CreatedAt,
        url.ShortCode,
    )
    
    if err != nil {
        return fmt.Errorf("update value error: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("no rows updated - url with code '%s' not found", url.ShortCode)
    }
    
    return nil
}

func (r *UrlRepository) DeleteUrlByCode(code string) error {
    query := `
        DELETE FROM url_info 
        WHERE short_code = $1
        RETURNING url_id
    `
    
    var deletedID int64
    err := r.db.QueryRow(query, code).Scan(&deletedID)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("url with code '%s' not found", code)
        }
        return fmt.Errorf("delete value error: %w", err)
    }
    
    return nil
}