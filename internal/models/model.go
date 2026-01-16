package models

import (
	"time"
)

// ListeningSession represents a listening session in a guild.
type ListeningSession struct {
	ID        string     `gorm:"type:varchar(36);primaryKey"`
	GuildID   int64      `gorm:"not null;index:idx_listening_sessions_guild"`
	StartedAt time.Time  `gorm:"type:timestamp;not null;index:idx_listening_sessions_started_at,sort:desc"`
	EndedAt   *time.Time `gorm:"type:timestamp"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

// TrackPlay represents an atomic unit of a track played during a listening session.
type TrackPlay struct {
	ID           string           `gorm:"type:varchar(36);primaryKey"`
	SessionID    string           `gorm:"type:varchar(36);not null;index:idx_track_plays_session"`
	Session      ListeningSession `json:"-" gorm:"foreignKey:SessionID;constraint:OnDelete:NO ACTION;"`
	GuildID      int64            `gorm:"not null;index:idx_track_plays_guild_time"`
	Source       string           `gorm:"type:varchar(255);not null;index:idx_track_plays_source"`
	SourceID     string           `gorm:"type:varchar(255);not null;index:idx_track_plays_source"`
	URL          string           `gorm:"type:varchar(1000);not null"`
	Title        string           `gorm:"type:varchar(500);not null"`
	Artist       string           `gorm:"type:varchar(500);not null"`
	StartedAt    time.Time        `gorm:"type:timestamp;not null;index:idx_track_plays_guild_time,sort:desc"`
	EndedAt      *time.Time       `gorm:"type:timestamp"`
	DurationMs   int              `gorm:"type:int"`
	QueuedByUser int64            `gorm:"not null;index:idx_track_plays_queued_by"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

// TrackPlayListener represents the users who listened to a particular track play.
type TrackPlayListener struct {
	TrackPlayID string    `gorm:"type:varchar(36);not null;primaryKey"`
	TrackPlay   TrackPlay `json:"-" gorm:"foreignKey:TrackPlayID;constraint:OnDelete:NO ACTION;"`
	UserID      int64     `gorm:"not null;primaryKey"`
}

// CreateTrackPlayRequest represents a request to create a track play.
type CreateTrackPlayRequest struct {
	Source       string  `json:"source" binding:"required"`
	SourceID     string  `json:"sourceId" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Artist       string  `json:"artist" binding:"required"`
	URL          string  `json:"url" binding:"required"`
	DurationMs   int     `json:"durationMs" binding:"required"`
	QueuedByUser int64   `json:"queuedByUser" binding:"required"`
	ListenerIDs  []int64 `json:"listenerIds" binding:"required"`
}

// UpdateTrackPlayRequest represents a request to update a track play.
type UpdateTrackPlayRequest struct {
	EndedAt *time.Time `json:"endedAt"`
}

// CreateSessionRequest represents a request to create a listening session.
type CreateSessionRequest struct {
	GuildID int64 `json:"guildId" binding:"required"`
}

// UpdateSessionRequest represents a request to update a listening session.
type UpdateSessionRequest struct {
	EndedAt *time.Time `json:"endedAt"`
}
