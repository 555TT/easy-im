#!/usr/bin/env bash
set -euo pipefail

# 项目根目录与常用目录
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
GO_IM_DIR="$ROOT_DIR/go-im"
PID_DIR="$ROOT_DIR/pids"
LOG_DIR="$ROOT_DIR/logs"

SERVICES=(
  "user-rpc"
  "social-rpc"
  "im-rpc"
  "im-ws"
  "user-api"
  "social-api"
  "im-api"
  "task-mq"
  "easy-im-web"
)

# 展示单个服务状态：运行中 / pid 文件残留 / 未启动
show_service_status() {
  local name="$1"
  local pid_file="$PID_DIR/${name}.pid"

  if [ ! -f "$pid_file" ]; then
    printf "%-12s %s\n" "$name" "not started"
    return
  fi

  local pid
  pid="$(cat "$pid_file")"

  if kill -0 "$pid" >/dev/null 2>&1; then
    printf "%-12s %s (pid: %s, log: %s)\n" "$name" "running" "$pid" "$LOG_DIR/${name}.log"
  else
    printf "%-12s %s (stale pid: %s)\n" "$name" "stopped" "$pid"
  fi
}

echo "=== App Process Status ==="
for service in "${SERVICES[@]}"; do
  show_service_status "$service"
done

echo
echo "=== Docker Compose Status ==="
cd "$GO_IM_DIR"
docker compose --env-file infra.host.env -f docker-compose.infra.yaml ps || true
