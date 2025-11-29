@echo off
chcp 65001 >nul
echo ====================================
echo SeekYourGuild - 模拟数据生成脚本
echo ====================================
echo.

set API_KEY=your-secret-api-key
set BASE_URL=http://localhost:8080/api
set GROUP_ID=123456789

echo 正在生成模拟数据...
echo.

REM 生成10个模拟成员
for /L %%i in (1,1,10) do (
    echo 添加成员 %%i...
    curl -s -X POST %BASE_URL%/report/member ^
      -H "Content-Type: application/json" ^
      -H "Authorization: Bearer %API_KEY%" ^
      -d "{\"group_id\": %GROUP_ID%, \"event_type\": \"join\", \"member\": {\"user_id\": 10000%%i, \"nickname\": \"用户%%i\", \"card\": \"群名片%%i\"}, \"timestamp\": 1732896000}"
)

echo.

REM 生成模拟消息
for /L %%h in (0,1,23) do (
    for /L %%m in (1,1,5) do (
        set /a USER_ID=10000 + %%m
        curl -s -X POST %BASE_URL%/report/message ^
          -H "Content-Type: application/json" ^
          -H "Authorization: Bearer %API_KEY%" ^
          -d "{\"group_id\": %GROUP_ID%, \"user_id\": !USER_ID!, \"message_id\": \"msg_%%h_%%m_%random%\", \"message_type\": \"text\", \"content_length\": 50, \"timestamp\": 1732896000}"
    )
)

echo.

REM 发送心跳
curl -s -X POST %BASE_URL%/report/heartbeat ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %API_KEY%" ^
  -d "{\"group_id\": %GROUP_ID%, \"bot_id\": 111111111, \"status\": \"online\", \"timestamp\": 1732896000}"

echo.
echo ====================================
echo 模拟数据生成完成!
echo ====================================
pause
