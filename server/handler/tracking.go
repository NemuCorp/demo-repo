package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NemuCorp/demo-repo/server/db"
	"github.com/NemuCorp/demo-repo/server/myerrors"
)

type TrackingHandler struct {
	DB *db.TrackingDB
}

func NewTrackingHandler(trackingDB *db.TrackingDB) *TrackingHandler {
	return &TrackingHandler{DB: trackingDB}
}

type TrackEventRequest struct {
	EventType string                 `json:"event_type" binding:"required"`
	EventData map[string]interface{} `json:"event_data"`
}

func (h *TrackingHandler) Track(c *gin.Context) {
	var req TrackEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.EventData == nil {
		req.EventData = map[string]interface{}{}
	}

	var userID *int
	if uid, ok := GetUserID(c); ok {
		userID = &uid
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	_, err := h.DB.RecordEvent(userID, req.EventType, req.EventData, ipAddress, userAgent)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, gin.H{"status": "recorded"})
}

func (h *TrackingHandler) Dashboard(c *gin.Context) {
	metrics, err := h.DB.GetDashboardMetrics()
	if err != nil {
		JSONError(c, http.StatusInternalServerError, myerrors.ErrInternal.Error())
		return
	}

	JSONSuccess(c, http.StatusOK, metrics)
}
