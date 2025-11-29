# 机器人端 API 变更通知

## 变更概述

**版本**: v2.0  
**日期**: 2024-11-29

API 已简化，**不再需要上报具体成员信息**，只需要上报数量变化。

---

## 移除的接口

### ❌ `POST /api/report/members` - 成员列表同步

此接口已移除，不再需要同步完整成员列表。

---

## 变更的接口

### `POST /api/report/member` - 成员变动上报

**旧格式** (已废弃):

```json
{
  "group_id": 123456789,
  "event_type": "join",
  "member": {
    "user_id": 987654321,
    "nickname": "用户昵称",
    "card": "群名片"
  },
  "operator_id": 111111111,
  "timestamp": 1700000000
}
```

**新格式**:

```json
{
  "group_id": 123456789,
  "event_type": "join",
  "count": 1,
  "current_total": 500,
  "timestamp": 1700000000
}
```

**字段说明**:

| 字段          | 类型   | 必填 | 说明                                         |
| ------------- | ------ | ---- | -------------------------------------------- |
| group_id      | int64  | ✅   | 群号                                         |
| event_type    | string | ✅   | 事件类型: `join` / `leave` / `kick` / `sync` |
| count         | int    | ✅   | 本次变动的人数（正数）                       |
| current_total | int    | ✅   | 变动后群的当前总人数                         |
| timestamp     | int64  | ✅   | 事件时间戳                                   |

**事件类型说明**:

| 类型    | 说明         | 示例场景                      |
| ------- | ------------ | ----------------------------- |
| `join`  | 新成员加入   | 有人入群                      |
| `leave` | 成员主动离开 | 有人退群                      |
| `kick`  | 成员被踢出   | 管理员踢人                    |
| `sync`  | 仅同步人数   | 启动时/定时同步（不记录事件） |

---

### `POST /api/report/heartbeat` - 心跳上报

**新增可选字段**: `member_count`

```json
{
  "group_id": 123456789,
  "bot_id": 111111111,
  "status": "online",
  "member_count": 500,
  "timestamp": 1700000000
}
```

心跳时可以顺便同步当前群人数，服务端会自动更新。

---

## 机器人端实现指南

### 1. 成员加入事件

```python
def on_member_join(event):
    # 获取当前群人数
    current_total = bot.get_group_info(event.group_id).member_count

    data = {
        "group_id": event.group_id,
        "event_type": "join",
        "count": 1,  # 单人加入
        "current_total": current_total,
        "timestamp": int(time.time())
    }
    post("/api/report/member", data)
```

### 2. 成员离开事件

```python
def on_member_leave(event):
    current_total = bot.get_group_info(event.group_id).member_count

    data = {
        "group_id": event.group_id,
        "event_type": "leave",  # 或 "kick"
        "count": 1,
        "current_total": current_total,
        "timestamp": int(time.time())
    }
    post("/api/report/member", data)
```

### 3. 启动时同步人数

```python
def on_bot_ready():
    for group_id in watch_groups:
        info = bot.get_group_info(group_id)

        data = {
            "group_id": group_id,
            "event_type": "sync",  # 使用 sync 类型，仅同步人数
            "count": 0,
            "current_total": info.member_count,
            "timestamp": int(time.time())
        }
        post("/api/report/member", data)
```

### 4. 心跳（推荐带上人数）

```python
def heartbeat():
    for group_id in watch_groups:
        info = bot.get_group_info(group_id)

        data = {
            "group_id": group_id,
            "bot_id": bot.self_id,
            "status": "online",
            "member_count": info.member_count,  # 推荐带上
            "timestamp": int(time.time())
        }
        post("/api/report/heartbeat", data)
```

---

## 消息上报（无变化）

消息上报接口保持不变：

```json
{
  "group_id": 123456789,
  "user_id": 987654321,
  "message_id": "msg_123456",
  "message_type": "text",
  "content_length": 50,
  "timestamp": 1700000000
}
```

---

## 完整事件处理示例

```python
import time
import requests

class GuildMonitor:
    def __init__(self, base_url, api_key, watch_groups):
        self.base_url = base_url
        self.headers = {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {api_key}"
        }
        self.watch_groups = watch_groups

    def report(self, endpoint, data):
        try:
            resp = requests.post(
                f"{self.base_url}/api/report/{endpoint}",
                headers=self.headers,
                json=data,
                timeout=5
            )
            return resp.json()
        except Exception as e:
            print(f"Report failed: {e}")
            return None

    # 成员变动
    def on_member_change(self, group_id, event_type, count=1):
        if group_id not in self.watch_groups:
            return

        current_total = bot.get_group_info(group_id).member_count

        self.report("member", {
            "group_id": group_id,
            "event_type": event_type,
            "count": count,
            "current_total": current_total,
            "timestamp": int(time.time())
        })

    # 消息
    def on_message(self, group_id, user_id, msg_id, msg_type, length):
        if group_id not in self.watch_groups:
            return

        self.report("message", {
            "group_id": group_id,
            "user_id": user_id,
            "message_id": msg_id,
            "message_type": msg_type,
            "content_length": length,
            "timestamp": int(time.time())
        })

    # 心跳
    def heartbeat(self):
        for group_id in self.watch_groups:
            info = bot.get_group_info(group_id)
            self.report("heartbeat", {
                "group_id": group_id,
                "bot_id": bot.self_id,
                "status": "online",
                "member_count": info.member_count,
                "timestamp": int(time.time())
            })

    # 启动同步
    def sync_all(self):
        for group_id in self.watch_groups:
            info = bot.get_group_info(group_id)
            self.report("member", {
                "group_id": group_id,
                "event_type": "sync",
                "count": 0,
                "current_total": info.member_count,
                "timestamp": int(time.time())
            })


# 使用示例
monitor = GuildMonitor(
    base_url="http://your-server:8080",
    api_key="your-secret-api-key",
    watch_groups=[123456789]
)

# 启动时同步
monitor.sync_all()

# 事件处理
@bot.on_member_join
def handle_join(event):
    monitor.on_member_change(event.group_id, "join")

@bot.on_member_leave
def handle_leave(event):
    monitor.on_member_change(event.group_id, "leave")

@bot.on_member_kick
def handle_kick(event):
    monitor.on_member_change(event.group_id, "kick")

@bot.on_group_message
def handle_message(event):
    msg_type = "text"
    if event.has_image:
        msg_type = "image"
    elif event.has_file:
        msg_type = "file"

    monitor.on_message(
        event.group_id,
        event.user_id,
        event.message_id,
        msg_type,
        len(event.raw_message)
    )

# 定时心跳
import threading
def heartbeat_loop():
    while True:
        monitor.heartbeat()
        time.sleep(60)

threading.Thread(target=heartbeat_loop, daemon=True).start()
```

---

## 核心变化总结

| 项目         | 旧版                         | 新版            |
| ------------ | ---------------------------- | --------------- |
| 成员列表同步 | 需要                         | ❌ 移除         |
| 单个成员信息 | 需要 user_id, nickname, card | ❌ 不需要       |
| 变动上报     | 每个成员单独上报             | ✅ 只上报数量   |
| 当前总人数   | 服务端计算                   | ✅ 机器人端提供 |
| 心跳带人数   | 不支持                       | ✅ 可选支持     |

**关键**: 服务端现在只关心 **数字**，不关心具体是谁。
