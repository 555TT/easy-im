#!/usr/bin/env bash
set -euo pipefail

# 项目根目录与常用子目录
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
GO_IM_DIR="$ROOT_DIR/go-im"
WEB_DIR="$ROOT_DIR/easy-im-web"
RUN_DIR="/opt/easy-im-run"
LOG_DIR="$RUN_DIR/logs"
PID_DIR="$RUN_DIR/pids"

# 运行时输出目录：日志与 pid 文件，统一放到仓库外
mkdir -p "$LOG_DIR" "$PID_DIR"

# 读取 VM_IP，供启动完成后打印宿主机访问地址
# 例如前端页面、Navicat、Jaeger 等
# shellcheck disable=SC1091
source "$GO_IM_DIR/infra.host.env"

# 后台启动单个服务，并记录日志与 pid
start_service() {
  local name="$1"
  local command="$2"
  local log_file="$LOG_DIR/${name}.log"
  local pid_file="$PID_DIR/${name}.pid"

  echo "Starting $name ..."
  nohup bash -lc "$command" >"$log_file" 2>&1 &
  local pid=$!
  echo "$pid" >"$pid_file"
  echo "$name started, log: $log_file, pid: $pid"
}

# 1. 先启动基础设施容器，确保 etcd/mysql/redis/mongo/kafka 等依赖可用
echo "[1/4] Starting infrastructure ..."
cd "$GO_IM_DIR"
docker compose --env-file infra.host.env -f docker-compose.infra.yaml up -d
docker compose --env-file infra.host.env -f docker-compose.infra.yaml ps

# 等待基础设施完成初始启动
sleep 3

# 2. 启动底层服务：rpc 与 ws，供上层 api 和 mq 依赖
echo "[2/4] Starting rpc and ws services ..."
start_service "user-rpc" "cd '$GO_IM_DIR' && go run apps/user/rpc/user.go -f apps/user/rpc/etc/user.yaml"
start_service "social-rpc" "cd '$GO_IM_DIR' && go run apps/social/rpc/social.go -f apps/social/rpc/etc/social.yaml"
start_service "im-rpc" "cd '$GO_IM_DIR' && go run apps/im/rpc/im.go -f apps/im/rpc/etc/im.yaml"
start_service "im-ws" "cd '$GO_IM_DIR' && go run apps/im/ws/im.go -f apps/im/ws/etc/im.yaml"

# 给 rpc/ws 一点注册与建连时间
sleep 3

# 3. 启动 api 服务；这些服务依赖前面的 rpc 注册到 etcd
echo "[3/4] Starting api services ..."
start_service "user-api" "cd '$GO_IM_DIR' && go run apps/user/api/user.go -f apps/user/api/etc/user.yaml"
start_service "social-api" "cd '$GO_IM_DIR' && go run apps/social/api/social.go -f apps/social/api/etc/social.yaml"
start_service "im-api" "cd '$GO_IM_DIR' && go run apps/im/api/im.go -f apps/im/api/etc/im.yaml"

# 给 api 一点启动时间
sleep 3

# 4. 最后启动消息消费者与前端
echo "[4/4] Starting task service and frontend ..."
start_service "task-mq" "cd '$GO_IM_DIR' && go run apps/task/mq/task.go -f apps/task/mq/etc/task.yaml"
start_service "easy-im-web" "cd '$WEB_DIR' && if [ ! -d node_modules ]; then npm install; fi && npm run dev -- --host"

# 输出常用信息
echo "All services have been started in background."
echo "Logs directory: $LOG_DIR"
echo "Frontend in VM: http://127.0.0.1:3000"
echo "Frontend from Windows host: http://${VM_IP}:3000"
echo "MySQL for Navicat: ${VM_IP}:13306"
echo "MongoDB for Navicat: ${VM_IP}:47017"
