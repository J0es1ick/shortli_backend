package urlHandlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	response "github.com/J0es1ick/shortli/internal/app/httputils"
	"github.com/J0es1ick/shortli/internal/config"
	"github.com/J0es1ick/shortli/internal/models"
	"github.com/J0es1ick/shortli/internal/repository"
	"github.com/J0es1ick/shortli/pkg/shortener"
	"github.com/skip2/go-qrcode"
)

type Handler struct {
	cfg *config.Config
	urlRepository *repository.UrlRepository
}

func NewHandler(cfg *config.Config, urlRepository *repository.UrlRepository) *Handler {
	return &Handler{
		cfg: cfg,
		urlRepository: urlRepository,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.Redirect(w, r)
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"message": "URL Shortener API",
		"version": "1.0",
	})
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req UrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.OriginalURL == "" {
		response.Error(w, http.StatusBadRequest, "Required original_url")
		return
	}

	userID := 0

    existingURL, err := h.urlRepository.FindUrlByOriginalUrl(req.OriginalURL)
    if err == nil {
        qrCode, err := qrcode.Encode(existingURL.OriginalURL, qrcode.Medium, 256)
        if err != nil {
            response.Error(w, http.StatusInternalServerError, "Failed to generate QR code")
            return
        }
        
        response.JSON(w, http.StatusOK, UrlResponse{
            OriginalURL:  existingURL.OriginalURL,
            ShortCode:    existingURL.ShortCode,
            ShortURL:     fmt.Sprintf("http://%s/%s", h.cfg.ServerPort, existingURL.ShortCode),
            QRCodeBase64: fmt.Sprintf("data:image/png;base64,%s", qrCode),
        })
        return
    }
	shortCode := shortener.GenerateShortCode(req.OriginalURL)
	for {
		existingURL, err := h.urlRepository.FindUrlByCode(shortCode)
		if err != nil && strings.Contains(err.Error(), "url not found") {
			break 
		}
		if existingURL != nil && existingURL.OriginalURL != req.OriginalURL {
			shortCode = shortener.GenerateShortCode(req.OriginalURL + time.Now().String())
			continue
		}
		break
	}

	qrCode, err := qrcode.Encode(req.OriginalURL, qrcode.Medium, 256) 
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate QR code")
		return
	}

	url := &models.URL{
		OriginalURL: req.OriginalURL,
		ShortCode: shortCode,
		UserId: userID,
		ClickCount: 0,
		QRClickCount: 0,
		CreatedAt: time.Now(),
	}

	if _, err := h.urlRepository.SaveUrl(url); err != nil {
		if strings.Contains(err.Error(), "unique constraint violation") {
			response.Error(w, http.StatusConflict, "URL already exists")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to save URL")
		return
	}

	response.JSON(w, http.StatusCreated, UrlResponse{
		OriginalURL: req.OriginalURL,
		ShortCode: shortCode,
		ShortURL: fmt.Sprintf("http://%s/%s", h.cfg.ServerPort, shortCode),
		QRCodeBase64: fmt.Sprintf("data:image/png;base64,%s", qrCode),
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")
	url, err := h.urlRepository.FindUrlByCode(shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "URL not found")
		} else {
			response.Error(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	if strings.Contains(r.UserAgent(), "QR") {
		url.QRClickCount++
	} else {
		url.ClickCount++
	}

	if err := h.urlRepository.UpdateUrlByCode(url); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update click count")
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusMovedPermanently)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	shortCode := strings.TrimPrefix(r.URL.Path, "/api/stats/")
	url, err := h.urlRepository.FindUrlByCode(shortCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(w, http.StatusNotFound, "URL not found")
		} else {
			response.Error(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	response.JSON(w, http.StatusOK, StatsResponse{
		URL: *url,
		TotalClicks: url.ClickCount + url.QRClickCount,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }

    shortCode := strings.TrimPrefix(r.URL.Path, "/urls/")

    if err := h.urlRepository.DeleteUrlByCode(shortCode); err != nil {
        if strings.Contains(err.Error(), "not found") {
            response.Error(w, http.StatusNotFound, "URL not found")
        } else {
            response.Error(w, http.StatusInternalServerError, "Failed to delete URL")
        }
        return
    }

    response.JSON(w, http.StatusOK, map[string]string{
        "status":  "success",
        "message": "URL deleted successfully",
        "code":    shortCode,
    })
}
