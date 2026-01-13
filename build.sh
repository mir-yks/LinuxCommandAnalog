#!/bin/bash

echo "Сборка всех Go утилит..."

for go_file in *.go; do
    if [[ -f "$go_file" && "$go_file" != "build.sh" ]]; then
        echo "Сборка $go_file → ${go_file%.go}"
        go build -o "${go_file%.go}" "$go_file"
        if [ $? -eq 0 ]; then
            echo "Сборка успешна"
        else
            echo "Ошибка сборки"
        fi
    fi
done

echo ""
echo "Готово! Исполняемые файлы созданы."
