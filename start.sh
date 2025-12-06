#!/bin/sh

set -e # 遇到错误立即退出

echo "run db migration"
# 执行迁移命令（前提是你的 Dockerfile 把 migrate 工具和 migration 文件都拷进去了）
/app/migrate -path /app/migration -database "$DBSOURCE" -verbose up

echo "start the app"
# 执行 Dockerfile CMD 中传进来的命令（通常是 /app/main）
exec "$@"