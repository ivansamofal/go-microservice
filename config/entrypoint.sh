#!/bin/bash
# Запускаем cron в фоне
/usr/sbin/cron -f &

# Если требуется установить крон-задачи из файла, например, /app/crontab, можно выполнить:
# crontab /app/crontab

# Запускаем основное приложение
exec /app/cmd/main
