#!/bin/sh
set -e

# Получаем первый IPv4 DNS из /etc/resolv.conf
DNS=$(grep nameserver /etc/resolv.conf | grep -v ':' | awk '{print $2}' | head -1)

# Если IPv4 не найден — берём любой (включая IPv6)
if [ -z "$DNS" ]; then
  DNS=$(grep nameserver /etc/resolv.conf | awk '{print $2}' | head -1)
fi

echo "Using DNS resolver: $DNS"

# Сначала подставляем BACKEND_URL через envsubst
# Потом заменяем DNS_PLACEHOLDER через sed
envsubst '${BACKEND_URL}' \
  < /etc/nginx/templates/default.conf.template \
  | sed "s|DNS_PLACEHOLDER|$DNS|g" \
  > /etc/nginx/conf.d/default.conf

echo "Generated nginx config:"
cat /etc/nginx/conf.d/default.conf

exec nginx -g "daemon off;"
