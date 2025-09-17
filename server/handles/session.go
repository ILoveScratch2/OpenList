package handles

import (
	"strconv"

	"github.com/OpenListTeam/OpenList/v4/internal/conf"
	"github.com/OpenListTeam/OpenList/v4/internal/db"
	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/OpenListTeam/OpenList/v4/server/common"
	"github.com/gin-gonic/gin"
)

type SessionResp struct {
	ID         uint   `json:"id"`
	DeviceKey  string `json:"device_key"`
	UserID     uint   `json:"user_id,omitempty"`
	Username   string `json:"username,omitempty"`
	UserAgent  string `json:"user_agent"`
	IP         string `json:"ip"`
	LastActive string `json:"last_active"`
	Status     string `json:"status"`
}

type SessionListReq struct {
	IP       string `json:"ip" form:"ip"`
	Username string `json:"username" form:"username"`
}

func ListMySessions(c *gin.Context) {
	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	var req SessionListReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	
	sessions, err := db.ListUserSessions(user.ID, req.IP)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	
	resp := make([]SessionResp, len(sessions))
	for i, s := range sessions {
		status := "ACTIVE"
		if !s.IsActive() {
			status = "INACTIVE"
		}
		resp[i] = SessionResp{
			ID:         s.ID,
			DeviceKey:  s.DeviceKey,
			UserAgent:  s.UserAgent,
			IP:         s.IP,
			LastActive: s.LastActive.Format("2006-01-02 15:04:05"),
			Status:     status,
		}
	}
	common.SuccessResp(c, resp)
}

type EvictSessionReq struct {
	DeviceKey string `json:"device_key" binding:"required"`
}

func EvictMySession(c *gin.Context) {
	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	var req EvictSessionReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	
	session, err := db.GetSession(user.ID, req.DeviceKey)
	if err != nil {
		common.ErrorStrResp(c, "Session not found", 404)
		return
	}
	
	if err := db.DeactivateSession(session.DeviceKey); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c)
}

func ListSessions(c *gin.Context) {
	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	if !user.CanManageSessions() {
		common.ErrorStrResp(c, "Permission denied", 403)
		return
	}
	
	var req SessionListReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	
	sessions, err := db.ListSessions(req.Username, req.IP)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	
	resp := make([]SessionResp, len(sessions))
	for i, s := range sessions {
		status := "ACTIVE"
		if !s.IsActive() {
			status = "INACTIVE"
		}
		resp[i] = SessionResp{
			ID:         s.ID,
			DeviceKey:  s.DeviceKey,
			UserID:     s.UserID,
			Username:   s.Username,
			UserAgent:  s.UserAgent,
			IP:         s.IP,
			LastActive: s.LastActive.Format("2006-01-02 15:04:05"),
			Status:     status,
		}
	}
	common.SuccessResp(c, resp)
}

func EvictSession(c *gin.Context) {
	user := c.Request.Context().Value(conf.UserKey).(*model.User)
	if !user.CanManageSessions() {
		common.ErrorStrResp(c, "Permission denied", 403)
		return
	}
	
	var req EvictSessionReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	
	if err := db.DeactivateSession(req.DeviceKey); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c)
}