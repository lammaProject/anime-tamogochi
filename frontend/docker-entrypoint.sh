#!/bin/sh
set -e

# Получаем DNS из /etc/resolv.conf и подставляем в nginx.conf
DNS=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}' | head -1)
echo "Using DNS resolver: $DNS"

# Подставляем DNS и BACKEND_URL в шаблон
export DOCKER_DNS_PLACEHOLDER="$DNS"
envsubst '${BACKEND_URL} ${DOCKER_DNS_PLACEHOLDER}' \
  < /etc/nginx/templates/default.conf.template \
  > /etc/nginx/conf.d/default.conf

echo "Generated nginx config:"
cat /etc/nginx/conf.d/default.conf

# Запускаем nginx
exec nginx -g "daemon off;"
