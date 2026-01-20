package database

import (
	"fmt"

	"github.com/mocbotau/mocbot-archive/internal/models"
)

// GetRecentTrackPlaysByGuild retrieves recent track plays in a guild with limit.
// Only completed plays (with EndedAt not null) are considered.
func (db *DB) GetRecentTrackPlaysByGuild(guildID string, limit int) ([]models.TrackPlay, error) {
	var trackPlays []models.TrackPlay

	query := db.
		Where("guild_id = ?", guildID).
		Where("ended_at IS NOT NULL").
		Order("started_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&trackPlays).Error; err != nil {
		return nil, fmt.Errorf("failed to get recent track plays: %w", err)
	}

	return trackPlays, nil
}

// GetRecentArtistsByGuild retrieves all unique artists played in a guild with their latest play time and count.
func (db *DB) GetRecentArtistsByGuild(guildID string, seedCount int) ([]models.ArtistWithLatestPlay, error) {
	var artists []models.ArtistWithLatestPlay

	// Subquery to get unique artists from recent plays
	subquery := db.
		Model(&models.TrackPlay{}).
		Select("artist, started_at").
		Where("guild_id = ?", guildID).
		Where("ended_at IS NOT NULL").
		Order("started_at DESC").
		Limit(seedCount)

	err := db.
		Table("(?) as recent", subquery).
		Select(`
        artist,
        MAX(started_at) AS latest_play_time,
        COUNT(*)        AS count
    `).Group("artist").
		Order("latest_play_time DESC").
		Scan(&artists).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get artists by guild: %w", err)
	}

	return artists, nil
}

// GetNumberOfTrackPlaysByGuild retrieves the total number of track plays in a guild for a given time range.
func (db *DB) GetNumberOfTrackPlaysByGuild(guildID string, startTime, endTime int64) (int64, error) {
	var count int64

	query := db.
		Model(&models.TrackPlay{}).
		Where("guild_id = ?", guildID).
		Where("ended_at IS NOT NULL")

	if startTime > 0 {
		query = query.Where("started_at >= to_timestamp(?)", startTime)
	}

	if endTime > 0 {
		query = query.Where("started_at <= to_timestamp(?)", endTime)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count track plays by guild: %w", err)
	}

	return count, nil
}
