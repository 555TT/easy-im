#!/usr/bin/env bash
set -euo pipefail

# 项目根目录与 pid 目录
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
GO_IM_DIR="$ROOT_DIR/go-im"
PID_DIR="$ROOT_DIR/pids"

# 按 pid 文件停止单个服务；如果进程已不存在则清理残留 pid 文件
stop_service() {
  local name="$1"
  local pid_file="$PID_DIR/${name}.pid"

  if [ ! -f "$pid_file" ]; then
    echo "$name pid file not found, skip"
    return
  fi

  local pid
  pid="$(cat "$pid_file")"

  if kill -0 "$pid" >/dev/null 2>&1; then
    echo "Stopping $name (pid: $pid) ..."
    kill "$pid"
  else
    echo "$name process $pid not running, cleaning pid file"
  fi

  rm -f "$pid_file"
}

# 先停上层服务，再停底层 rpc/ws，避免依赖中的服务先被杀掉
echo "Stopping app services ..."
stop_service "easy-im-web"
stop_service "task-mq"
stop_service "im-api"
stop_service "social-api"
stop_service "user-api"
stop_service "im-ws"
stop_service "im-rpc"
stop_service "social-rpc"
stop_service "user-rpc"

# 最后停止基础设施容器
echo "Stopping infrastructure ..."
cd "$GO_IM_DIR"
docker compose --env-file infra.host.env -f docker-compose.infra.yaml down

echo "All services stopped."
