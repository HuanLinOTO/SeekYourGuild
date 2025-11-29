# API 设计文档

## 基础信息

- 基础路径: `/api`
- 数据格式: JSON
- 字符编码: UTF-8

## 认证

所有上报接口需要在请求头中携带 API Key：

```
Authorization: Bearer <API_KEY>
```

---

## 上报 API

### 1. 群成员变动上报

**接口地址**: `POST /api/report/member`

**请求体**:

```json
{
  "group_id": 123456789,
  "event_type": "join", // "join" | "leave" | "kick"
  "member": {
    "user_id": 987654321,
    "nickname": "用户昵称",
    "card": "群名片"
  },
  "operator_id": 111111111, // 操作者ID（踢人时有值）
  "timestamp": 1700000000
}
```

**响应**:

```json
{
  "code": 0,
  "message": "success"
}
```

### 2. 群成员列表同步

**接口地址**: `POST /api/report/members`

**请求体**:

```json
{
  "group_id": 123456789,
  "members": [
    {
      "user_id": 987654321,
      "nickname": "用户昵称",
      "card": "群名片",
      "role": "member", // "owner" | "admin" | "member"
      "join_time": 1699999999
    }
  ],
  "timestamp": 1700000000
}
```

**响应**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "synced": 100,
    "added": 5,
    "updated": 10
  }
}
```

### 3. 消息上报

**接口地址**: `POST /api/report/message`

**请求体**:

```json
{
  "group_id": 123456789,
  "user_id": 987654321,
  "message_id": "msg_123456",
  "message_type": "text", // "text" | "image" | "file" | "other"
  "content_length": 100, // 消息长度（不存储具体内容，仅统计）
  "timestamp": 1700000000
}
```

**响应**:

```json
{
  "code": 0,
  "message": "success"
}
```

### 4. 心跳上报

**接口地址**: `POST /api/report/heartbeat`

**请求体**:

```json
{
  "group_id": 123456789,
  "bot_id": 111111111,
  "status": "online",
  "timestamp": 1700000000
}
```

**响应**:

```json
{
  "code": 0,
  "message": "success"
}
```

---

## 数据查询 API

### 1. 群概览数据

**接口地址**: `GET /api/stats/overview`

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| group_id | int64 | 是 | 群号 |

**响应**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "group_id": 123456789,
    "current_members": 500,
    "today_joined": 10,
    "today_left": 3,
    "recent_members": [
      {
        "user_id": 987654321,
        "nickname": "新成员1",
        "card": "群名片",
        "join_time": 1700000000
      }
    ],
    "bot_status": "online",
    "last_update": 1700000000
  }
}
```

### 2. 成员变动历史

**接口地址**: `GET /api/stats/member-changes`

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| group_id | int64 | 是 | 群号 |
| start_date | string | 否 | 开始日期 (YYYY-MM-DD) |
| end_date | string | 否 | 结束日期 (YYYY-MM-DD) |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |

**响应**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 100,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "date": "2024-11-29",
        "joined": 10,
        "left": 3,
        "net_change": 7
      }
    ]
  }
}
```

### 3. 发言统计

**接口地址**: `GET /api/stats/messages`

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| group_id | int64 | 是 | 群号 |
| period | string | 否 | 统计周期: "hour" / "day" / "week"，默认 "day" |
| start_date | string | 否 | 开始日期 |
| end_date | string | 否 | 结束日期 |

**响应**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "day",
    "stats": [
      {
        "time": "2024-11-29",
        "total_messages": 1500,
        "active_users": 80,
        "by_type": {
          "text": 1200,
          "image": 200,
          "file": 50,
          "other": 50
        }
      }
    ],
    "hourly_distribution": [
      { "hour": 0, "count": 50 },
      { "hour": 1, "count": 30 }
    ]
  }
}
```

### 4. 活跃用户排行

**接口地址**: `GET /api/stats/active-users`

**查询参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| group_id | int64 | 是 | 群号 |
| period | string | 否 | 统计周期: "day" / "week" / "month" |
| limit | int | 否 | 返回数量，默认 10 |

**响应**:

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "period": "day",
    "users": [
      {
        "rank": 1,
        "user_id": 987654321,
        "nickname": "活跃用户",
        "message_count": 200,
        "last_active": 1700000000
      }
    ]
  }
}
```

---

## 错误码

| 错误码 | 说明         |
| ------ | ------------ |
| 0      | 成功         |
| 1001   | 参数错误     |
| 1002   | 认证失败     |
| 1003   | 权限不足     |
| 2001   | 数据库错误   |
| 2002   | 数据不存在   |
| 3001   | 请求频率超限 |
