#!/usr/bin/env bash
set -euo pipefail

# 项目根目录与常用目录
ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
GO_IM_DIR="$ROOT_DIR/go-im"
PID_DIR="$ROOT_DIR/pids"
LOG_DIR="$ROOT_DIR/logs"

# 读取 VM_IP，用于打印宿主机/外部访问地址
# shellcheck disable=SC1091
source "$GO_IM_DIR/infra.host.env"

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
echo "=== Service Endpoints ==="
echo "Frontend:"
printf "  easy-im-web   http://127.0.0.1:3000   (VM 内访问)\n"
printf "  easy-im-web   http://%s:3000   (Windows/宿主机访问，前端页面)\n" "$VM_IP"
echo
printf "  user-api      http://127.0.0.1:8880\n"
printf "  social-api    http://127.0.0.1:8881\n"
printf "  im-api        http://127.0.0.1:8882\n"
printf "  im-ws         ws://127.0.0.1:10090/ws?userId={userId}\n"
echo
printf "  user-rpc      127.0.0.1:10000   (内部 RPC)\n"
printf "  social-rpc    127.0.0.1:10001   (内部 RPC)\n"
printf "  im-rpc        127.0.0.1:10002   (内部 RPC)\n"
printf "  task-mq       127.0.0.1:10091   (消息消费服务)\n"
echo
printf "  MySQL         %s:13306   (Navicat)\n" "$VM_IP"
printf "  MongoDB       %s:47017   (Navicat/外部工具)\n" "$VM_IP"
printf "  Redis         127.0.0.1:16379   (VM 内部)\n"
printf "  Etcd          %s:3379   (advertise/external)\n" "$VM_IP"
printf "  Kafka         %s:9092   (advertise/external)\n" "$VM_IP"
printf "  Jaeger        http://127.0.0.1:16686   (VM 内访问 UI)\n"
printf "  Jaeger        http://%s:16686   (Windows/宿主机访问 UI)\n" "$VM_IP"

echo
echo "=== Docker Compose Status ==="
cd "$GO_IM_DIR"
docker compose --env-file infra.host.env -f docker-compose.infra.yaml ps || true
