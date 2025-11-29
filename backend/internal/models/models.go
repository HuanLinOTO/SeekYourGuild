package models

import (
	"time"
)

// Group 群信息表
type Group struct {
	ID          int64     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255"`
	MemberCount int       `json:"member_count" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Group) TableName() string {
	return "groups"
}

// MemberEvent 成员变动事件表（仅记录数量变化）
type MemberEvent struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID       int64     `json:"group_id" gorm:"index"`
	EventType     string    `json:"event_type" gorm:"size:20"` // join/leave/kick
	Count         int       `json:"count" gorm:"default:1"`    // 变动数量
	CurrentTotal  int       `json:"current_total"`             // 变动后的总人数
	EventTime     time.Time `json:"event_time" gorm:"index"`
	CreatedAt     time.Time `json:"created_at"`
}

func (MemberEvent) TableName() string {
	return "member_events"
}

// Message 消息统计表
type Message struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID       int64     `json:"group_id" gorm:"index"`
	UserID        int64     `json:"user_id" gorm:"index"`
	Nickname      string    `json:"nickname" gorm:"size:255"`
	MessageID     string    `json:"message_id" gorm:"size:64;unique"`
	MessageType   string    `json:"message_type" gorm:"size:20"`
	ContentLength int       `json:"content_length"`
	SendTime      time.Time `json:"send_time" gorm:"index"`
	CreatedAt     time.Time `json:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}

// DailyStat 每日统计表
type DailyStat struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID      int64     `json:"group_id" gorm:"index"`
	StatDate     time.Time `json:"stat_date" gorm:"index;type:date"`
	MemberCount  int       `json:"member_count"`
	JoinedCount  int       `json:"joined_count"`
	LeftCount    int       `json:"left_count"`
	MessageCount int       `json:"message_count"`
	ActiveUsers  int       `json:"active_users"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (DailyStat) TableName() string {
	return "daily_stats"
}

// HourlyMessageStat 小时消息统计表
type HourlyMessageStat struct {
	ID           int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID      int64     `json:"group_id" gorm:"index"`
	StatHour     time.Time `json:"stat_hour" gorm:"index"`
	MessageCount int       `json:"message_count"`
	ActiveUsers  int       `json:"active_users"`
	TextCount    int       `json:"text_count"`
	ImageCount   int       `json:"image_count"`
	FileCount    int       `json:"file_count"`
	OtherCount   int       `json:"other_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (HourlyMessageStat) TableName() string {
	return "hourly_message_stats"
}

// BotStatus 机器人状态表
type BotStatus struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	GroupID       int64     `json:"group_id" gorm:"index"`
	BotID         int64     `json:"bot_id"`
	Status        string    `json:"status" gorm:"size:20"` // online/offline
	LastHeartbeat time.Time `json:"last_heartbeat"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (BotStatus) TableName() string {
	return "bot_status"
}
