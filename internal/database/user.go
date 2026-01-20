package database

import (
	"fmt"

	"github.com/mocbotau/mocbot-archive/internal/models"
)

// GetRecentTrackPlaysByUser retrieves recent track plays that a user listened to with limit.
func (db *DB) GetRecentTrackPlaysByUser(userID string, limit int) ([]models.TrackPlay, error) {
	var trackPlays []models.TrackPlay

	query := db.
		Joins("JOIN track_play_listeners ON track_play_listeners.track_play_id = track_plays.id").
		Where("track_plays.ended_at IS NOT NULL").
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

// GetRecentArtistsByUser retrieves all unique artists listened to by a user with their latest play time and count.
func (db *DB) GetRecentArtistsByUser(userID string, seedCount int) ([]models.ArtistWithLatestPlay, error) {
	var artists []models.ArtistWithLatestPlay

	subquery := db.
		Model(&models.TrackPlay{}).
		Select("track_plays.artist, track_plays.started_at").
		Joins("JOIN track_play_listeners ON track_play_listeners.track_play_id = track_plays.id").
		Where("track_play_listeners.user_id = ?", userID).
		Where("track_plays.ended_at IS NOT NULL").
		Order("track_plays.started_at DESC").
		Limit(seedCount)

	err := db.
		Table("(?) as recent", subquery).
		Select(`
		artist,
		MAX(started_at) AS latest_play_time,
		COUNT(*)        AS count
	`).
		Group("artist").
		Order("latest_play_time DESC").
		Scan(&artists).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get artists by user: %w", err)
	}

	return artists, nil
}

// GetNumberOfTrackPlaysByUser retrieves the total number of track plays by a user for a given time range.
func (db *DB) GetNumberOfTrackPlaysByUser(userID string, startTime, endTime int64) (int64, error) {
	var count int64

	query := db.
		Model(&models.TrackPlay{}).
		Joins("JOIN track_play_listeners ON track_play_listeners.track_play_id = track_plays.id").
		Where("track_play_listeners.user_id = ?", userID).
		Where("track_plays.ended_at IS NOT NULL")

	if startTime > 0 {
		query = query.Where("track_plays.started_at >= ?", startTime)
	}

	if endTime > 0 {
		query = query.Where("track_plays.started_at <= ?", endTime)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count user's track plays: %w", err)
	}

	return count, nil
}
