package database

import (
	"fmt"
	"time"

	"github.com/lithammer/shortuuid/v4"

	"github.com/mocbotau/mocbot-archive/internal/models"
)

// CreateTrackPlay creates a new track play record.
func (db *DB) CreateTrackPlay(sessionID string, req *models.CreateTrackPlayRequest) (*models.TrackPlay, error) {
	session, err := db.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	trackPlay := &models.TrackPlay{
		ID:           shortuuid.New(),
		SessionID:    sessionID,
		GuildID:      session.GuildID,
		Source:       req.Source,
		SourceID:     req.SourceID,
		Title:        req.Title,
		Artist:       req.Artist,
		URL:          req.URL,
		StartedAt:    time.Now(),
		DurationMs:   req.DurationMs,
		QueuedByUser: req.QueuedByUser,
		CreatedAt:    time.Now(),
	}

	if err := db.Create(trackPlay).Error; err != nil {
		return nil, fmt.Errorf("failed to create track play: %w", err)
	}

	return trackPlay, nil
}

// UpdateTrackPlay updates an existing track play record.
func (db *DB) UpdateTrackPlay(id string, req *models.UpdateTrackPlayRequest) (*models.TrackPlay, error) {
	trackPlay := &models.TrackPlay{}

	err := db.
		Where("id = ?", id).
		First(trackPlay).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find track play: %w", err)
	}

	if req.EndedAt != nil {
		trackPlay.EndedAt = req.EndedAt
	}

	if err := db.Save(trackPlay).Error; err != nil {
		return nil, fmt.Errorf("failed to update track play: %w", err)
	}

	return trackPlay, nil
}

// GetTrackPlay retrieves a track play by ID.
func (db *DB) GetTrackPlay(id string) (*models.TrackPlay, error) {
	trackPlay := &models.TrackPlay{}

	err := db.
		Preload("Session").
		Where("id = ?", id).
		First(trackPlay).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find track play: %w", err)
	}

	return trackPlay, nil
}

// GetTrackPlaysBySession retrieves all track plays for a session.
func (db *DB) GetTrackPlaysBySession(sessionID string) ([]models.TrackPlay, error) {
	var trackPlays []models.TrackPlay

	err := db.
		Where("session_id = ?", sessionID).
		Order("started_at ASC").
		Find(&trackPlays).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get track plays for session: %w", err)
	}

	return trackPlays, nil
}

// GetRecentTrackPlaysByGuild retrieves recent track plays in a guild with limit.
func (db *DB) GetRecentTrackPlaysByGuild(guildID string, limit int) ([]models.TrackPlay, error) {
	var trackPlays []models.TrackPlay

	query := db.
		Where("guild_id = ?", guildID).
		Order("started_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&trackPlays).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent track plays: %w", err)
	}

	return trackPlays, nil
}

// GetRecentTrackPlaysByUser retrieves recent track plays that a user listened to with limit.
func (db *DB) GetRecentTrackPlaysByUser(userID string, limit int) ([]models.TrackPlay, error) {
	var trackPlays []models.TrackPlay

	query := db.
		Joins("JOIN track_play_listeners ON track_play_listeners.track_play_id = track_plays.id").
		Where("track_play_listeners.user_id = ?", userID).
		Order("track_plays.started_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&trackPlays).Error; err != nil {
		return nil, fmt.Errorf("failed to get user's track plays: %w", err)
	}

	return trackPlays, nil
}
