package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mocbotau/mocbot-archive/internal/models"
)

// GetSession retrieves a listening session by ID.
func (h *Handler) GetSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}

	session, err := h.db.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetSessionsByGuilds retrieves sessions by guild IDs with optional limit.
// Query params: guildIds (comma-separated), limit (optional, default 50).
func (h *Handler) GetSessionsByGuilds(c *gin.Context) {
	guildIDsParam := c.Query("guildIds")
	if guildIDsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guildIds query parameter is required"})
		return
	}

	guildIDs := strings.Split(guildIDsParam, ",")
	for i := range guildIDs {
		guildIDs[i] = strings.TrimSpace(guildIDs[i])
	}

	limit := 50

	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	sessions, err := h.db.GetSessionsByGuildIDs(guildIDs, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// StartSession creates a new listening session for a guild.
func (h *Handler) StartSession(c *gin.Context) {
	var req models.CreateSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.db.CreateSession(req.GuildID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// EndSession updates a session with an end time.
func (h *Handler) EndSession(c *gin.Context) {
	sessionID := c.Param("sessionId")

	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sessionId is required"})
		return
	}

	var req models.UpdateSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	endedAt := req.EndedAt
	if endedAt == nil {
		now := time.Now()
		endedAt = &now
	}

	session, err := h.db.UpdateSession(sessionID, endedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to end session"})
		return
	}

	c.JSON(http.StatusOK, session)
}
