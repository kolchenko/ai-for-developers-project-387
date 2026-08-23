### Hexlet tests and linter status:
[![Actions Status](https://github.com/kolchenko/ai-for-developers-project-386/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/kolchenko/ai-for-developers-project-386/actions)

# Календарь звонков

Frontend + backend для приложения бронирования звонков. Пользователь выбирает тип события, видит свободные слоты на ближайшие 14 дней и бронирует удобное время. Администратор управляет типами событий и просматривает бронирования.

## Стек

- **Frontend:** React 19, Mantine 9, React Query 5, Vite
- **Backend:** Go 1.22+ (`net/http`), SQLite (`modernc.org/sqlite`, без cgo)
- **API-контракт:** TypeSpec (`main.tsp`) → OpenAPI (`tsp-output/openapi.yaml`)
- **Деплой:** Docker, Render

## Возможности

- Список типов событий для гостей
- Выбор слота и бронирование по типу события (сетка из 14 дней, шаг `durationMinutes`)
- Админ-панель: CRUD типов событий, просмотр бронирований
- Аутентификация отсутствует

## Быстрый старт

```sh
npm install

# Frontend dev-сервер на :5173 (проксирует /api -> http://127.0.0.1:4010)
npm run dev

# Backend на :4010
cd backend
go run ./cmd/server
```

## Команды

```sh
npm run dev       # Vite dev-сервер
npm run openapi   # скомпилировать main.tsp -> tsp-output/openapi.yaml
npm run build     # проверка: tsc --noEmit && vite build

cd backend
go test ./...     # unit + handler тесты
go vet ./...      # статические проверки
```

`npm run openapi` / mock требуют Node >= 22 (см. `.nvmrc`). `dev`/`build` работают на Node 20.19+.

## API-контракт

`main.tsp` — источник правды для API. После изменения запусти `npm run openapi` и проверь `tsp-output/openapi.yaml` перед изменением UI-кода. `src/api/endpoints.ts` должен совпадать с путями OpenAPI.

Известная особенность: операция `upcoming()` под `@route("/admin/bookings")` эмитится как `GET /admin/bookings`, а не `/admin/bookings/upcoming`.

## Структура проекта

```
main.tsp              # API-контракт (TypeSpec)
src/api/              # типы, axios-клиент, endpoints, React Query хуки
src/pages/            # компоненты маршрутов
src/components/       # общий UI (layout, модалки)
src/utils/time.ts     # хелперы для времени (24-часовой формат)
backend/              # Go + SQLite backend
e2e/                  # Playwright e2e-тесты
```

## Маршруты

- `/` — список типов событий (гости)
- `/event-types/:eventTypeId` — слоты + бронирование
- `/admin` — бронирования
- `/admin/event-types` — CRUD типов событий

## Деплой

Публичная ссылка: <https://call-calendar-kid4.onrender.com>

Приложение упаковано в Docker-образ (см. `Dockerfile`) и запускается в контейнере на порту из переменной окружения `PORT`. Деплой на Render выполнен через GitHub-репозиторий: Render собирает образ из `Dockerfile` в `master` и запускает сервис на публичном URL выше.
