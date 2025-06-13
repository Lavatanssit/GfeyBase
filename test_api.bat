@echo off
echo 测试首页 GET /
curl http://localhost:9999/
echo.

echo 测试 GET /news 带参数
curl "http://localhost:9999/news?title=TimesNews&content=ConferencesStart"
echo.

echo 测试 POST /login
curl -X POST -d "username=admin&password=123456" http://localhost:9999/login
echo.

echo 测试 POST /login 返回 JSON
curl -X POST -d "username=admin&password=123456" -H "Content-Type: application/x-www-form-urlencoded" http://localhost:9999/login
echo.

pause