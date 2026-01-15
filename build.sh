#!/bin/bash

echo "Сборка всех Go утилит..."

success=0
failed=0

# Автоматически добавляем PROMPT_COMMAND в ~/.bashrc
BASHRC="$HOME/.bashrc"
PROMPT_LINE='export PROMPT_COMMAND="history -a"'

if ! grep -qF "$PROMPT_LINE" "$BASHRC" 2>/dev/null; then
    echo "$PROMPT_LINE" >> "$BASHRC"
    NEED_SOURCE=1
else
    NEED_SOURCE=0
fi

for go_file in *.go; do
    if [[ -f "$go_file" && "$go_file" != "build.sh" ]]; then
        echo "Сборка $go_file → ${go_file%.go}"
        go build -o "${go_file%.go}" "$go_file"
        if [ $? -eq 0 ]; then
            echo "Сборка успешна"
            ((success++))
        else
            echo "Ошибка сборки"
            ((failed++))
        fi
    fi
done

echo ""
echo "Готово! Исполняемые файлы созданы."

if [ $failed -ne 0 ]; then
    exit 1
fi

if [ "$NEED_SOURCE" -eq 1 ]; then
    echo ""
    echo "Выполните: source ~/.bashrc"
fi

