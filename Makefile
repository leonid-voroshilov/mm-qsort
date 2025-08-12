# Makefile для проекта компаратора

# Переменные
GO = go
BINARY_NAME = comparator
BUILD_DIR = build
CMD_DIR = cmd
PKG_DIR = .

# Цели по умолчанию
.PHONY: all build test clean run help benchmark coverage lint format deps

# Сборка исполняемого файла
build:
	@echo "Сборка проекта..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Исполняемый файл создан: $(BUILD_DIR)/$(BINARY_NAME)"

# Запуск программы
run: build
	@echo "Запуск программы..."
	./$(BUILD_DIR)/$(BINARY_NAME)


# Запуск всех тестов
test:
	@echo "Запуск тестов..."
	$(GO) test ./... -v

# Запуск тестов с кратким выводом
test-short:
	@echo "Запуск тестов (краткий режим)..."
	$(GO) test ./...

# Запуск тестов с покрытием
coverage:
	@echo "Запуск тестов с анализом покрытия..."
	$(GO) test ./... -cover
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Отчет о покрытии создан: coverage.html"

# Бенчмарки
benchmark:
	@echo "Запуск бенчмарков..."
	$(GO) test ./... -bench=. -benchmem

# Линтинг кода
lint:
	@echo "Проверка кода линтерами..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint не установлен. Используем go vet..."; \
		$(GO) vet ./...; \
	fi

# Форматирование кода
format:
	@echo "Форматирование кода..."
	$(GO) fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	fi

# Установка зависимостей
deps:
	@echo "Установка зависимостей..."
	$(GO) mod tidy
	$(GO) mod download

# Обновление зависимостей
deps-update:
	@echo "Обновление зависимостей..."
	$(GO) get -u ./...
	$(GO) mod tidy

# Полная проверка перед коммитом
pre-commit: format-check lint test
	@echo "Все проверки пройдены успешно!"

# Справка
help:
	@echo "Доступные команды:"
	@echo "  build          - Сборка исполняемого файла"
	@echo "  run            - Запуск программы"
	@echo "  test           - Запуск всех тестов"
	@echo "  test-short     - Запуск тестов (краткий вывод)"
	@echo "  coverage       - Запуск тестов с анализом покрытия"
	@echo "  benchmark      - Запуск бенчмарков"
	@echo "  profile        - Профилирование производительности"
	@echo "  lint           - Проверка кода линтерами"
	@echo "  format         - Форматирование кода"
	@echo "  deps           - Установка зависимостей"
	@echo "  deps-update    - Обновление зависимостей"
	@echo "  pre-commit     - Полная проверка перед коммитом"
	@echo ""
	@echo "Примеры:"
	@echo "  make run"
	@echo "  make test"
	@echo "  make coverage"
	@echo "  make release"