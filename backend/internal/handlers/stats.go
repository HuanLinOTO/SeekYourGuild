package handlers

import (
	"net/http"
	"strconv"
	"time"

	"seekyourguild/internal/config"
	"seekyourguild/internal/database"
	"seekyourguild/internal/models"

	"github.com/gin-gonic/gin"
)

// GetGroups 获取可用群列表
func GetGroups(c *gin.Context) {
	cfg := config.Load()
	db := database.GetDB()
	
	groups := make([]gin.H, 0)
	
	// 如果配置了允许的群ID，使用配置
	if len(cfg.Groups.AllowedIDs) > 0 {
		for _, id := range cfg.Groups.AllowedIDs {
			var group models.Group
			if err := db.First(&group, id).Error; err == nil {
				groups = append(groups, gin.H{
					"id":           group.ID,
					"name":         group.Name,
					"member_count": group.MemberCount,
				})
			} else {
				// 群不存在于数据库，仍然返回ID
				groups = append(groups, gin.H{
					"id":           id,
					"name":         "",
					"member_count": 0,
				})
			}
		}
	} else {
		// 没有配置，返回数据库中所有群
		var allGroups []models.Group
		db.Find(&allGroups)
		for _, group := range allGroups {
			groups = append(groups, gin.H{
				"id":           group.ID,
				"name":         group.Name,
				"member_count": group.MemberCount,
			})
		}
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"groups": groups,
		},
	})
}

// GetOverview 群概览数据
func GetOverview(c *gin.Context) {
	groupIDStr := c.Query("group_id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "group_id is required",
		})
		return
	}

	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "invalid group_id",
		})
		return
	}

	db := database.GetDB()

	// 获取群信息
	var group models.Group
	if err := db.First(&group, groupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    2002,
			"message": "group not found",
		})
		return
	}

	// 今日统计 - 从每日统计表或事件表获取
	today := time.Now().Truncate(24 * time.Hour)

	var todayJoined int64
	var todayLeft int64

	// 尝试从 daily_stats 获取
	var dailyStat models.DailyStat
	if err := db.Where("group_id = ? AND stat_date = ?", groupID, today).First(&dailyStat).Error; err == nil {
		todayJoined = int64(dailyStat.JoinedCount)
		todayLeft = int64(dailyStat.LeftCount)
	} else {
		// 从事件表聚合（兼容旧数据）
		tomorrow := today.Add(24 * time.Hour)
		
		db.Model(&models.MemberEvent{}).
			Where("group_id = ? AND event_type = ? AND event_time >= ? AND event_time < ?",
				groupID, "join", today, tomorrow).
			Select("COALESCE(SUM(count), 0)").
			Scan(&todayJoined)

		db.Model(&models.MemberEvent{}).
			Where("group_id = ? AND event_type IN (?, ?) AND event_time >= ? AND event_time < ?",
				groupID, "leave", "kick", today, tomorrow).
			Select("COALESCE(SUM(count), 0)").
			Scan(&todayLeft)
	}

	// 机器人状态
	var botStatus models.BotStatus
	botStatusStr := "offline"
	var lastUpdate int64
	if err := db.Where("group_id = ?", groupID).Order("last_heartbeat DESC").First(&botStatus).Error; err == nil {
		// 检查心跳是否超时（5分钟）
		if time.Since(botStatus.LastHeartbeat) < 5*time.Minute {
			botStatusStr = botStatus.Status
		}
		lastUpdate = botStatus.LastHeartbeat.Unix()
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"group_id":        group.ID,
			"current_members": group.MemberCount,
			"today_joined":    todayJoined,
			"today_left":      todayLeft,
			"bot_status":      botStatusStr,
			"last_update":     lastUpdate,
		},
	})
}

// GetMemberChanges 成员变动历史
func GetMemberChanges(c *gin.Context) {
	groupIDStr := c.Query("group_id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "group_id is required",
		})
		return
	}

	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "invalid group_id",
		})
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 日期范围
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	if startDateStr != "" {
		startDate, _ = time.Parse("2006-01-02", startDateStr)
	} else {
		startDate = time.Now().AddDate(0, 0, -30) // 默认30天
	}
	if endDateStr != "" {
		endDate, _ = time.Parse("2006-01-02", endDateStr)
		endDate = endDate.Add(24 * time.Hour)
	} else {
		endDate = time.Now().Add(24 * time.Hour)
	}

	db := database.GetDB()

	// 查询每日统计
	type DayStats struct {
		Date   string `json:"date"`
		Joined int    `json:"joined"`
		Left   int    `json:"left"`
	}

	var dailyStats []DailyStat
	db.Where("group_id = ? AND stat_date >= ? AND stat_date < ?", groupID, startDate, endDate).
		Order("stat_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&dailyStats)

	var total int64
	db.Model(&models.DailyStat{}).
		Where("group_id = ? AND stat_date >= ? AND stat_date < ?", groupID, startDate, endDate).
		Count(&total)

	items := make([]gin.H, 0)
	for _, stat := range dailyStats {
		items = append(items, gin.H{
			"date":       stat.StatDate.Format("2006-01-02"),
			"joined":     stat.JoinedCount,
			"left":       stat.LeftCount,
			"net_change": stat.JoinedCount - stat.LeftCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"items":     items,
		},
	})
}

// DailyStat 用于查询的临时结构
type DailyStat struct {
	StatDate    time.Time
	JoinedCount int
	LeftCount   int
}

// GetMessageStats 发言统计
func GetMessageStats(c *gin.Context) {
	groupIDStr := c.Query("group_id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "group_id is required",
		})
		return
	}

	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "invalid group_id",
		})
		return
	}

	period := c.DefaultQuery("period", "day")

	db := database.GetDB()

	// 获取最近7天的小时统计数据
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	var hourlyStats []models.HourlyMessageStat
	db.Where("group_id = ? AND stat_hour >= ?", groupID, sevenDaysAgo).
		Order("stat_hour ASC").
		Find(&hourlyStats)

	// 按天聚合
	dayStats := make(map[string]gin.H)
	hourlyDistribution := make([]int, 24)

	for _, stat := range hourlyStats {
		dateKey := stat.StatHour.Format("2006-01-02")
		if _, exists := dayStats[dateKey]; !exists {
			dayStats[dateKey] = gin.H{
				"time":           dateKey,
				"total_messages": 0,
				"active_users":   0,
				"by_type": gin.H{
					"text":  0,
					"image": 0,
					"file":  0,
					"other": 0,
				},
			}
		}

		ds := dayStats[dateKey]
		ds["total_messages"] = ds["total_messages"].(int) + stat.MessageCount
		ds["active_users"] = ds["active_users"].(int) + stat.ActiveUsers
		byType := ds["by_type"].(gin.H)
		byType["text"] = byType["text"].(int) + stat.TextCount
		byType["image"] = byType["image"].(int) + stat.ImageCount
		byType["file"] = byType["file"].(int) + stat.FileCount
		byType["other"] = byType["other"].(int) + stat.OtherCount

		// 小时分布
		hour := stat.StatHour.Hour()
		hourlyDistribution[hour] += stat.MessageCount
	}

	stats := make([]gin.H, 0)
	for _, v := range dayStats {
		stats = append(stats, v)
	}

	hourlyData := make([]gin.H, 24)
	for i := 0; i < 24; i++ {
		hourlyData[i] = gin.H{
			"hour":  i,
			"count": hourlyDistribution[i],
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"period":              period,
			"stats":               stats,
			"hourly_distribution": hourlyData,
		},
	})
}

// GetActiveUsers 活跃用户排行
func GetActiveUsers(c *gin.Context) {
	groupIDStr := c.Query("group_id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "group_id is required",
		})
		return
	}

	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1001,
			"message": "invalid group_id",
		})
		return
	}

	period := c.DefaultQuery("period", "day")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	db := database.GetDB()

	// 计算时间范围
	var startTime time.Time
	switch period {
	case "week":
		startTime = time.Now().AddDate(0, 0, -7)
	case "month":
		startTime = time.Now().AddDate(0, -1, 0)
	default:
		startTime = time.Now().Truncate(24 * time.Hour)
	}

	// 查询活跃用户
	type UserStat struct {
		UserID       int64
		Nickname     string
		MessageCount int64
		LastActive   time.Time
	}

	var userStats []UserStat
	// 使用子查询获取每个用户最新的昵称
	db.Raw(`
		SELECT 
			m.user_id,
			(SELECT nickname FROM messages WHERE user_id = m.user_id AND group_id = ? ORDER BY send_time DESC LIMIT 1) as nickname,
			COUNT(*) as message_count,
			MAX(m.send_time) as last_active
		FROM messages m
		WHERE m.group_id = ? AND m.send_time >= ?
		GROUP BY m.user_id
		ORDER BY message_count DESC
		LIMIT ?
	`, groupID, groupID, startTime, limit).Scan(&userStats)

	users := make([]gin.H, 0)
	for i, stat := range userStats {
		users = append(users, gin.H{
			"rank":          i + 1,
			"user_id":       stat.UserID,
			"nickname":      stat.Nickname,
			"message_count": stat.MessageCount,
			"last_active":   stat.LastActive.Unix(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"period": period,
			"users":  users,
		},
	})
}
