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
	stmt, err := r.db.Prepare("INSERT INTO url_info (original_url, short_code, click_count, qr_click_count) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("preparing insert operator error: %v", err)
	}

	res, err := stmt.Exec(url.OriginalURL, url.ShortCode, url.ClickCount, url.QRClickCount)
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

func (r *UrlRepository) FindUrlByCode(code string) (*models.URL, error) {
	stmt, err := r.db.Prepare("SELECT original_url, click_count, qr_click_count FROM url_info WHERE short_code = ?")
	if err != nil {
		return nil, fmt.Errorf("preparing select operator error: %v",  err)
	}

	url := &models.URL{}
    
    err = stmt.QueryRow(code).Scan(&url.OriginalURL, &url.ClickCount, &url.QRClickCount)
	if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("url not found")
        }
        return nil, fmt.Errorf("scan row: %v", err)
    }

	return url, nil
}

func (r *UrlRepository) UpdateUrlByCode(url *models.URL) error {
	stmt, err := r.db.Prepare("UPDATE url_info SET original_url = ?, click_count = ?, qr_click_count = ? WHERE short_code = ?")
	if err != nil {
		return fmt.Errorf("preparing update operator error: %v", err)
	}

	_, err = stmt.Exec(url.OriginalURL, url.ClickCount, url.QRClickCount, url.ShortCode)
	if err != nil {
		return fmt.Errorf("update value error: %v", err)
	}

	return nil
}

func (r *UrlRepository) UpdateUrlClicksByCode(code string) error {
	stmt, err := r.db.Prepare("UPDATE url_info SET click_count = click_count + 1 WHERE short_code = ?")
	if err != nil {
		return fmt.Errorf("preparing update clicks operator error: %v", err)
	}

	_, err = stmt.Exec(code)
	if err != nil {
		return fmt.Errorf("update clicks value error: %v", err)
	}

	return nil
}

func (r *UrlRepository) DeleteUrlByCode(code string) error {
	stmt, err := r.db.Prepare("DELETE FROM url_info WHERE short_code = ?")
	if err != nil {
		return fmt.Errorf("preparing delete operator error: %v", err)
	}

	_, err = stmt.Exec(code)
	if err != nil {
		return fmt.Errorf("delete value error: %v", err)
	}

	return nil
}