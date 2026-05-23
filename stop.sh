#!/usr/bin/env bash
set -euo pipefail

# 项目根目录与运行时目录
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
GO_IM_DIR="$ROOT_DIR/go-im"
RUN_DIR="/opt/easy-im-run"
PID_DIR="$RUN_DIR/pids"

# 按 pid 文件停止单个服务，pid 文件中实际保存的是进程组 ID (PGID)
# 通过 kill -- -PGID 一次性结束 bash 包装、go run 以及真正的服务二进制
stop_service() {
  local name="$1"
  local pid_file="$PID_DIR/${name}.pid"

  if [ ! -f "$pid_file" ]; then
    echo "$name pid file not found, skip"
    return
  fi

  local pid
  pid="$(cat "$pid_file")"

  if [ -z "$pid" ]; then
    echo "$name pid file empty, cleaning"
    rm -f "$pid_file"
    return
  fi

  if kill -0 "-$pid" >/dev/null 2>&1; then
    echo "Stopping $name (pgid: $pid) ..."
    kill -TERM "-$pid" 2>/dev/null || true

    # 最多等待 5 秒让进程优雅退出
    local i
    for i in 1 2 3 4 5; do
      if kill -0 "-$pid" >/dev/null 2>&1; then
        sleep 1
      else
        break
      fi
    done

    # 仍未退出则强制结束整个进程组
    if kill -0 "-$pid" >/dev/null 2>&1; then
      echo "$name still alive, sending SIGKILL"
      kill -KILL "-$pid" 2>/dev/null || true
    fi
  else
    echo "$name process group $pid not running, cleaning pid file"
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
