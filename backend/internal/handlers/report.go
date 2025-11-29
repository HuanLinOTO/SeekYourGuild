package handlers

import (
	"net/http"
	"time"

	"seekyourguild/internal/database"
	"seekyourguild/internal/models"

	"github.com/gin-gonic/gin"
)

// MemberChangeRequest 成员数量变动上报请求（简化版）
type MemberChangeRequest struct {
	GroupID      int64  `json:"group_id" binding:"required"`
	EventType    string `json:"event_type" binding:"required,oneof=join leave kick sync"`
	Count        int    `json:"count" binding:"required"`         // 变动数量（正数）
	CurrentTotal int    `json:"current_total" binding:"required"` // 当前群总人数
	Timestamp    int64  `json:"timestamp" binding:"required"`
}

// MessageRequest 消息上报请求
type MessageRequest struct {
	GroupID       int64  `json:"group_id" binding:"required"`
	UserID        int64  `json:"user_id" binding:"required"`
	Nickname      string `json:"nickname"`
	MessageID     string `json:"message_id" binding:"required"`
	MessageType   string `json:"message_type" binding:"required"`
	ContentLength int    `json:"content_length"`
	Timestamp     int64  `json:"timestamp" binding:"required"`
}

// HeartbeatRequest 心跳上报请求
type HeartbeatRequest struct {
	GroupID     int64  `json:"group_id" binding:"required"`
	BotID       int64  `json:"bot_id" binding:"required"`
	Status      string `json:"status" binding:"required"`
	MemberCount int    `json:"member_count"` // 可选：心跳时也可以同步人数
	Timestamp   int64  `json:"timestamp" binding:"required"`
}

// ReportMemberChange 成员数量变动上报
func ReportMemberChange(c *gin.Context) {
	var req MemberChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": err.Error(),
		})
		return
	}

	db := database.GetDB()
	eventTime := time.Unix(req.Timestamp, 0)

	// 确保群存在并更新人数
	var group models.Group
	result := db.First(&group, req.GroupID)
	if result.Error != nil {
		// 群不存在，创建
		group = models.Group{
			ID:          req.GroupID,
			MemberCount: req.CurrentTotal,
		}
		db.Create(&group)
	} else {
		// 更新群人数
		db.Model(&group).Update("member_count", req.CurrentTotal)
	}

	// 记录变动事件（sync 类型不记录事件，仅同步人数）
	if req.EventType != "sync" {
		event := models.MemberEvent{
			GroupID:      req.GroupID,
			EventType:    req.EventType,
			Count:        req.Count,
			CurrentTotal: req.CurrentTotal,
			EventTime:    eventTime,
		}
		if err := db.Create(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    2001,
				"message": "database error",
			})
			return
		}

		// 更新每日统计
		updateDailyStat(req.GroupID, eventTime, req.EventType, req.Count)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"member_count": req.CurrentTotal,
		},
	})
}

// updateDailyStat 更新每日统计
func updateDailyStat(groupID int64, eventTime time.Time, eventType string, count int) {
	db := database.GetDB()
	statDate := time.Date(eventTime.Year(), eventTime.Month(), eventTime.Day(), 0, 0, 0, 0, eventTime.Location())

	var dailyStat models.DailyStat
	result := db.Where("group_id = ? AND stat_date = ?", groupID, statDate).First(&dailyStat)

	if result.Error != nil {
		// 创建新记录
		dailyStat = models.DailyStat{
			GroupID:  groupID,
			StatDate: statDate,
		}
		if eventType == "join" {
			dailyStat.JoinedCount = count
		} else {
			dailyStat.LeftCount = count
		}
		db.Create(&dailyStat)
	} else {
		// 更新记录
		if eventType == "join" {
			db.Model(&dailyStat).Update("joined_count", dailyStat.JoinedCount+count)
		} else {
			db.Model(&dailyStat).Update("left_count", dailyStat.LeftCount+count)
		}
	}
}

// ReportMessage 消息上报
func ReportMessage(c *gin.Context) {
	var req MessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": err.Error(),
		})
		return
	}

	db := database.GetDB()
	sendTime := time.Unix(req.Timestamp, 0)

	// 记录消息
	message := models.Message{
		GroupID:       req.GroupID,
		UserID:        req.UserID,
		Nickname:      req.Nickname,
		MessageID:     req.MessageID,
		MessageType:   req.MessageType,
		ContentLength: req.ContentLength,
		SendTime:      sendTime,
	}

	if err := db.Create(&message).Error; err != nil {
		// 可能是重复消息，忽略
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
		})
		return
	}

	// 更新小时统计
	hourStart := time.Date(sendTime.Year(), sendTime.Month(), sendTime.Day(), sendTime.Hour(), 0, 0, 0, sendTime.Location())

	var hourStat models.HourlyMessageStat
	result := db.Where("group_id = ? AND stat_hour = ?", req.GroupID, hourStart).First(&hourStat)

	if result.Error != nil {
		// 创建新统计
		hourStat = models.HourlyMessageStat{
			GroupID:      req.GroupID,
			StatHour:     hourStart,
			MessageCount: 1,
			ActiveUsers:  1,
		}
		switch req.MessageType {
		case "text":
			hourStat.TextCount = 1
		case "image":
			hourStat.ImageCount = 1
		case "file":
			hourStat.FileCount = 1
		default:
			hourStat.OtherCount = 1
		}
		db.Create(&hourStat)
	} else {
		// 更新统计
		updates := map[string]interface{}{
			"message_count": hourStat.MessageCount + 1,
		}
		switch req.MessageType {
		case "text":
			updates["text_count"] = hourStat.TextCount + 1
		case "image":
			updates["image_count"] = hourStat.ImageCount + 1
		case "file":
			updates["file_count"] = hourStat.FileCount + 1
		default:
			updates["other_count"] = hourStat.OtherCount + 1
		}
		db.Model(&hourStat).Updates(updates)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}

// ReportHeartbeat 心跳上报
func ReportHeartbeat(c *gin.Context) {
	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": err.Error(),
		})
		return
	}

	db := database.GetDB()
	heartbeatTime := time.Unix(req.Timestamp, 0)

	// 如果心跳带了人数，更新群人数
	if req.MemberCount > 0 {
		db.Model(&models.Group{}).Where("id = ?", req.GroupID).Update("member_count", req.MemberCount)
	}

	var botStatus models.BotStatus
	result := db.Where("group_id = ? AND bot_id = ?", req.GroupID, req.BotID).First(&botStatus)

	if result.Error != nil {
		// 创建新记录
		botStatus = models.BotStatus{
			GroupID:       req.GroupID,
			BotID:         req.BotID,
			Status:        req.Status,
			LastHeartbeat: heartbeatTime,
		}
		db.Create(&botStatus)
	} else {
		// 更新记录
		db.Model(&botStatus).Updates(map[string]interface{}{
			"status":         req.Status,
			"last_heartbeat": heartbeatTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
