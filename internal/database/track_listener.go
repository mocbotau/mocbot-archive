package database

import (
	"fmt"

	"github.com/mocbotau/mocbot-archive/internal/models"
)

// AddListeners adds multiple listeners to a track play in a batch.
func (db *DB) AddListeners(trackPlayID string, userIDs []int64) error {
	if trackPlayID == "" {
		return fmt.Errorf("track play ID is required")
	}

	if len(userIDs) == 0 {
		return fmt.Errorf("no user IDs provided")
	}

	listeners := make([]models.TrackPlayListener, len(userIDs))
	for i, userID := range userIDs {
		listeners[i] = models.TrackPlayListener{
			TrackPlayID: trackPlayID,
			UserID:      userID,
		}
	}

	if err := db.Create(&listeners).Error; err != nil {
		return fmt.Errorf("failed to add listeners: %w", err)
	}

	return nil
}

// GetListeners retrieves all listeners for a specific track play.
func (db *DB) GetListeners(trackPlayID string) ([]models.TrackPlayListener, error) {
	var listeners []models.TrackPlayListener

	err := db.
		Where("track_play_id = ?", trackPlayID).
		Find(&listeners).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get listeners: %w", err)
	}

	return listeners, nil
}
