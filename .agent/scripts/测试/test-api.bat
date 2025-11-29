@echo off
chcp 65001 >nul
echo ====================================
echo SeekYourGuild - API 测试脚本
echo ====================================
echo.

set API_KEY=your-secret-api-key
set BASE_URL=http://localhost:8080/api

echo [测试1] 上报成员加入事件...
curl -X POST %BASE_URL%/report/member ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %API_KEY%" ^
  -d "{\"group_id\": 123456789, \"event_type\": \"join\", \"member\": {\"user_id\": 987654321, \"nickname\": \"测试用户\", \"card\": \"测试名片\"}, \"timestamp\": %time:~0,2%%time:~3,2%%time:~6,2%}"

echo.
echo.

echo [测试2] 上报消息...
curl -X POST %BASE_URL%/report/message ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %API_KEY%" ^
  -d "{\"group_id\": 123456789, \"user_id\": 987654321, \"message_id\": \"msg_%random%\", \"message_type\": \"text\", \"content_length\": 50, \"timestamp\": %time:~0,2%%time:~3,2%%time:~6,2%}"

echo.
echo.

echo [测试3] 上报心跳...
curl -X POST %BASE_URL%/report/heartbeat ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %API_KEY%" ^
  -d "{\"group_id\": 123456789, \"bot_id\": 111111111, \"status\": \"online\", \"timestamp\": %time:~0,2%%time:~3,2%%time:~6,2%}"

echo.
echo.

echo [测试4] 获取群概览...
curl -X GET "%BASE_URL%/stats/overview?group_id=123456789"

echo.
echo.

echo [测试5] 获取消息统计...
curl -X GET "%BASE_URL%/stats/messages?group_id=123456789"

echo.
echo.

echo [测试6] 获取活跃用户...
curl -X GET "%BASE_URL%/stats/active-users?group_id=123456789"

echo.
echo ====================================
echo 测试完成!
echo ====================================
pause
