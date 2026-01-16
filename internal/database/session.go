package database

import (
	"fmt"
	"time"

	"github.com/lithammer/shortuuid/v4"

	"github.com/mocbotau/mocbot-archive/internal/models"
)

// GetSession retrieves a listening session by session ID.
func (db *DB) GetSession(sessionID string) (*models.ListeningSession, error) {
	session := &models.ListeningSession{}

	err := db.
		Where("id = ?", sessionID).
		First(session).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find listening session: %w", err)
	}

	return session, nil
}

// CreateSession creates a new listening session for a guild.
func (db *DB) CreateSession(guildID int64) (*models.ListeningSession, error) {
	if guildID == 0 {
		return nil, fmt.Errorf("guild ID is required")
	}

	session := &models.ListeningSession{
		ID:        shortuuid.New(),
		GuildID:   guildID,
		StartedAt: time.Now(),
	}

	if err := db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create listening session: %w", err)
	}

	return session, nil
}

// UpdateSession updates an existing listening session for a guild.
func (db *DB) UpdateSession(sessionID string, endedAt *time.Time) (*models.ListeningSession, error) {
	session := &models.ListeningSession{}

	err := db.
		Where("id = ?", sessionID).
		First(session).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find listening session: %w", err)
	}

	if endedAt == nil {
		return nil, fmt.Errorf("endedAt is required")
	}

	session.EndedAt = endedAt

	if err := db.Save(session).Error; err != nil {
		return nil, fmt.Errorf("failed to update listening session: %w", err)
	}

	return session, nil
}

// GetSessionsByGuildIDs retrieves listening sessions by guild IDs with limit.
func (db *DB) GetSessionsByGuildIDs(guildIDs []string, limit int) ([]models.ListeningSession, error) {
	var sessions []models.ListeningSession

	query := db.Where("guild_id IN ?", guildIDs).
		Order("started_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	return sessions, nil
}
