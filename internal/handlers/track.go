package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mocbotau/mocbot-archive/internal/models"
	"github.com/mocbotau/mocbot-archive/internal/utils"
)

// StartTrack creates a new track play and adds listeners in a transaction.
func (h *Handler) StartTrack(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}

	var req struct {
		models.CreateTrackPlayRequest
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	trackPlay, err := h.db.CreateTrackPlay(sessionID, &req.CreateTrackPlayRequest)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create track play"})

		return
	}

	if len(req.ListenerIDs) > 0 {
		if err := h.db.AddListeners(trackPlay.ID, req.ListenerIDs); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add listeners"})

			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, trackPlay)
}

// EndTrack updates a track play with end metadata.
func (h *Handler) EndTrack(c *gin.Context) {
	trackPlayID := c.Param("trackPlayId")
	if trackPlayID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trackPlayId is required"})
		return
	}

	var req models.UpdateTrackPlayRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EndedAt == nil {
		now := time.Now()
		req.EndedAt = &now
	}

	trackPlay, err := h.db.UpdateTrackPlay(trackPlayID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update track play"})
		return
	}

	c.JSON(http.StatusOK, trackPlay)
}

// GetTrackPlay retrieves a track play by ID.
func (h *Handler) GetTrackPlay(c *gin.Context) {
	trackPlayID := c.Param("trackPlayId")
	if trackPlayID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trackPlayId is required"})
		return
	}

	trackPlay, err := h.db.GetTrackPlay(trackPlayID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track play not found"})
		return
	}

	c.JSON(http.StatusOK, trackPlay)
}

// GetTracksBySession retrieves all track plays for a session.
func (h *Handler) GetTracksBySession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}

	trackPlays, err := h.db.GetTrackPlaysBySession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve track plays"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trackPlays": trackPlays,
		"count":      len(trackPlays),
	})
}

// GetRecentTracksByGuild retrieves recent track plays in a guild with optional limit.
// Query param: limit (optional, default 50, max 500).
func (h *Handler) GetRecentTracksByGuild(c *gin.Context) {
	guildID := c.Param("guildId")
	if guildID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guildId is required"})
		return
	}

	limit := utils.DefaultLimit

	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			if parsedLimit > utils.MaxLimit {
				parsedLimit = utils.MaxLimit
			}

			limit = parsedLimit
		}
	}

	trackPlays, err := h.db.GetRecentTrackPlaysByGuild(guildID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve track plays"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trackPlays": trackPlays,
		"count":      len(trackPlays),
		"guildId":    guildID,
	})
}

// GetRecentTracksByUser retrieves recent track plays that a user listened to with optional limit.
// Query param: limit (optional, default 50, max 500).
func (h *Handler) GetRecentTracksByUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	limit := utils.DefaultLimit

	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			if parsedLimit > utils.MaxLimit {
				parsedLimit = utils.MaxLimit
			}

			limit = parsedLimit
		}
	}

	trackPlays, err := h.db.GetRecentTrackPlaysByUser(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve track plays"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"trackPlays": trackPlays,
		"count":      len(trackPlays),
		"userId":     userID,
	})
}

// GetRecommendedTracksByUser retrieves recommended track plays for a user.
func (h *Handler) GetRecommendedTracksByUser(c *gin.Context) {

}

// GetRecommendedTracksByGuild retrieves recommended track plays for a user.
func (h *Handler) GetRecommendedTracksByGuild(c *gin.Context) {
	
}

// GetTrackListeners retrieves all listeners for a specific track play.
func (h *Handler) GetTrackListeners(c *gin.Context) {
	trackPlayID := c.Param("trackPlayId")
	if trackPlayID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trackPlayId is required"})
		return
	}

	listeners, err := h.db.GetListeners(trackPlayID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve listeners"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"listeners": listeners,
		"count":     len(listeners),
	})
}
