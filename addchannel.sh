#!/bin/bash

# Проверяем, что передан один аргумент
if [ "$#" -ne 1 ]; then
    echo "Использование: $0 URL"
    exit 1
fi

URL=$1
FILE="channels.txt"

# Проверяем, существует ли файл
if [ ! -f "$FILE" ]; then
    # Если файл не существует, создаем его
    touch "$FILE"
fi

# Проверяем, содержится ли URL в файле
if grep -Fxq "$URL" "$FILE"; then
    echo "Ссылка уже существует в файле."
else
    echo "$URL" >> "$FILE"
    echo "Ссылка добавлена в файл."
fi