# GLPITGBOT

Telegram-бот + HTTP API для создания заявок в GLPI, хранения локальной истории в PostgreSQL и отправки уведомлений по вебхукам GLPI.

---

## 1) Что делает проект

- авторизует пользователя по его `user_token` GLPI;
- сохраняет `api_token` и `session_token` пользователя локально;
- создаёт заявку в GLPI;
- хранит локальную копию заявки в PostgreSQL;
- связывает локальную заявку с внешним `external_ticket_id` из GLPI;
- принимает вебхуки (изменение статуса/комментарии) и отправляет уведомления в Telegram.

---

## 2) Переменные окружения

Минимально нужны:

- `PORT` — порт HTTP-сервера (по умолчанию `7000`);
- `TG_BOT_TOKEN` — токен Telegram-бота;
- `GLPI_URL` (или `GLPI_API`) — базовый URL GLPI;
- `APP_TOKEN` — app token GLPI;
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` — PostgreSQL;
- `WEBHOOK_KEY` / `WEBHOOK_SECRET` (опционально, но рекомендуется) — защита вебхуков;
- `WEBHOOK_TICKET`, `WEBHOOK_COMMENT` — дополнительные ключи для разных типов вебхуков;
- `WEBHOOK_DEBUG` — `1/true/yes`, чтобы логировать подробности входящего вебхука.

---

## 3) Запуск

1. Поднять PostgreSQL.
2. Создать `.env` с переменными выше.
3. Применить SQL-миграции из папки migrations.
4. Запустить:

```bash
go run .
```

Сервис поднимет:

- Telegram long polling;
- HTTP API с маршрутами:
    - `GET /api/auth/check?token=...`
    - `POST /api/ticket/create?session_token=...`
    - `POST /api/webhook`

---

## 4) Краткая архитектура

- [main.go](main.go) — точка входа;
- [db/postgres.go](db/postgres.go) — подключение и базовая синхронизация схемы;
- [https/auth.go](https/auth.go) — авторизация в GLPI по user token;
- [https/create_ticket.go](https/create_ticket.go) — создание заявки в GLPI;
- [https/webhook.go](https/webhook.go) — приём вебхуков и рассылка уведомлений;
- [repository/users_repository.go](repository/users_repository.go), [repository/tickets_repository.go](repository/tickets_repository.go) — работа с БД;
- [telegram](telegram) — бот, диалоги, клавиатуры, локальные сессии;
- [notifications](notifications) — форматирование текста и отправка сообщений в Telegram;
- [models](models) — модели данных и JSON payload-структуры.

---

## 5) Построчное объяснение кода (по файлам)

Ниже — подробный разбор каждой рабочей строки (пустые строки и одиночные `}` как структурные разделители не дублируются комментариями отдельно).

### 5.1 [go.mod](go.mod)

- `module GLPITGBOT` — имя модуля.
- `go 1.24.0` — версия языка.
- Блок `require` (основной):
    - `gin` — HTTP-фреймворк.
    - `telegram-bot-api/v5` — Telegram API.
    - `godotenv` — загрузка `.env`.
    - `lib/pq` — драйвер PostgreSQL.
- Второй `require` — косвенные (`// indirect`) зависимости, подтянутые транзитивно.

### 5.2 [main.go](main.go)

- `package main` — главный исполняемый пакет.
- Импорты `db`, `https`, `telegram` — внутренние пакеты приложения.
- Импорты `log`, `os` — логирование и чтение окружения.
- Импорт `gin` — HTTP-сервер.
- `func main()` — точка входа процесса.
- `db.Connect()` — подключение к PostgreSQL и проверка доступности.
- `go telegram.Bot()` — запуск Telegram-бота в отдельной goroutine.
- `r := gin.Default()` — создание роутера с logger/recovery middleware.
- `r.SetTrustedProxies(...)` — доверенные прокси для корректного `ClientIP`.
- `api := r.Group("/api")` — общий префикс API.
- `api.GET("/auth/check", https.CheckAuth)` — проверка user token и выдача session token.
- `api.POST("/ticket/create", https.CreateTicket)` — создание заявки в GLPI.
- `api.POST("/webhook", https.WebhookGLPI)` — endpoint вебхуков GLPI.
- Чтение `PORT` из окружения.
- Если порт пустой — используется `7000`.
- `r.Run(":" + port)` — запуск HTTP-сервера.
- `log.Fatal(err)` — аварийное завершение, если сервер не стартовал.

### 5.3 [db/postgres.go](db/postgres.go)

- `var DB *sql.DB` — глобальный пул подключений.
- `godotenv.Load()` — загрузка `.env`.
- Если `.env` не найден — `log.Fatal` (процесс завершается).
- `os.Getenv(...)` — чтение параметров БД.
- Проверка `DB_PORT` на пустоту — ранняя диагностика конфигурации.
- `dsn := fmt.Sprintf(...)` — сборка строки подключения PostgreSQL.
- `sql.Open("postgres", dsn)` — создание клиента БД.
- `DB.Ping()` — проверка фактического соединения.
- `ALTER TABLE users ADD COLUMN ... session_token` — мягкая синхронизация столбца.
- `ALTER TABLE tickets ADD COLUMN ... external_ticket_id` — добавление внешнего ID.
- `CREATE INDEX ... idx_tickets_external_ticket_id` — индекс по внешнему ID.
- `log.Println("✅ PostgreSQL connected")` — подтверждение успешного старта.

### 5.4 [https/auth.go](https/auth.go)

- `type authResponse` — модель JSON-ответа GLPI с `session_token`.
- `AuthByUserToken(userToken string)`:
    - читает `GLPI_URL` (fallback: `GLPI_API`);
    - читает `APP_TOKEN`;
    - триммит входной `userToken`;
    - валидирует обязательные значения;
    - формирует endpoint `/apirest.php/initSession`;
    - создаёт `GET`-запрос;
    - выставляет заголовки `Authorization: user_token ...` и `App-Token`;
    - отправляет через `http.Client` с timeout 10s;
    - читает body;
    - проверяет `StatusOK`;
    - парсит JSON;
    - проверяет, что `session_token` не пуст;
    - возвращает `session_token`.
- `CheckAuth(c *gin.Context)`:
    - берёт `token` из query;
    - вызывает `AuthByUserToken`;
    - при ошибке — `401`;
    - при успехе — `200` и JSON `{ok:true, session_token:...}`.

### 5.5 [https/create_ticket.go](https/create_ticket.go)

- `createTicketResponse` — модель ответа с полем `id` (как `RawMessage`, т.к. ID может быть строкой/числом).
- `ExtractCreatedTicketID(body []byte)`:
    - парсит JSON;
    - проверяет наличие `id`;
    - пытается распарсить `id` как строку;
    - затем как число;
    - приводит к строке и возвращает.
- `createTicketRequest` — обёртка `{ "input": ... }` под формат GLPI.
- `CreateTicketWithSession(sessionToken, input)`:
    - читает `GLPI_URL`/`GLPI_API` и `APP_TOKEN`;
    - валидирует обязательные поля;
    - сериализует payload;
    - создаёт `POST /apirest.php/Ticket`;
    - ставит `Content-Type`, `App-Token`, `Session-Token`;
    - отправляет запрос с timeout 10s;
    - читает body;
    - если статус не `200/201` — возвращает тело + ошибку;
    - иначе возвращает body и статус.
- `CreateTicket(c *gin.Context)`:
    - читает `session_token` из query;
    - проверяет JSON body;
    - вызывает `CreateTicketWithSession`;
    - проксирует статус/тело ответа GLPI клиенту.

### 5.6 [https/webhook.go](https/webhook.go)

- `webhookDebugBodyPreviewLimit` — ограничение лога тела запроса.
- `webhookDebugEnabled()` — включает debug-режим по `WEBHOOK_DEBUG`.
- `debugWebhookRequest(...)` — логирует метаданные запроса и урезанное тело.
- `readIncomingWebhookKey(...)`:
    - читает ключ из набора заголовков/параметров;
    - поддерживает `Bearer ...`.
- `uniqueNonEmpty(...)` — убирает пустые и дубликаты.
- `verifyHMACSHA256(...)` — проверка подписи HMAC SHA-256.
- `verifyGLPISignature(...)`:
    - нормализует `sha256=...`;
    - проверяет подпись по `body`;
    - при наличии `timestamp` проверяет варианты `body+timestamp` и `timestamp+body`.
- `verifyWebhookAuth(...)`:
    - сначала прямое сравнение ключей;
    - если не подошло — проверка HMAC-подписи через секреты из env.
- `WebhookGLPI(c *gin.Context)`:
    - читает тело;
    - логирует в debug;
    - проверяет авторизацию вебхука (`401` при отказе);
    - распознаёт тип payload (ticket/followup) через probe-модели;
    - направляет в `handleFollowupWebhook` или `handleTicketWebhook`.
- `handleFollowupWebhook(...)`:
    - извлекает ID тикета;
    - формирует текст уведомления о комментарии;
    - находит получателей по `external_ticket_id`;
    - отправляет сообщения в Telegram;
    - возвращает HTTP-статус результата.
- `handleTicketWebhook(...)`:
    - извлекает ID тикета;
    - вычисляет статус (предпочтительно `status_name`);
    - формирует уведомление об изменении статуса;
    - находит получателей и рассылает сообщения.

### 5.7 [models/headers_model.go](models/headers_model.go)

- `type Headers` — структура заголовков HTTP.
- Поля с JSON-тегами:
    - `ContentType` ↔ `Content-Type`;
    - `Authorization` ↔ `Authorization`;
    - `AppToken` ↔ `App-Token`;
    - `SessionToken` ↔ `Session-Token`.

### 5.8 [models/ticket_model.go](models/ticket_model.go)

- `type Ticket` — локальная доменная модель заявки.
- Поля `Id/Title/Description/UserID/...` — данные заявки в приложении.
- `type CreateTicketInput` — формат `input` для API GLPI.
- JSON-теги (`name`, `content`, `entities_id`, ...) определяют точные ключи payload.

### 5.9 [models/user_model.go](models/user_model.go)

- `type User` — модель пользователя.
- `TelegramID`, `Username` — идентификация из Telegram.
- `ApiToken`, `SessionToken` — учётные данные GLPI.
- `Lang` — язык интерфейса.
- `ActiveTicket *Ticket` — активная заявка в диалоге (сессионный контекст).

### 5.10 [models/webhooks_model.go](models/webhooks_model.go)

- Импорт `encoding/json` нужен для `json.RawMessage`.
- `BaseWebhookPayload` — минимальная модель: `event` + сырой `item`.
- `WebhookItemTypeProbe` — probe для `itemtype`.
- `UpdateEvent`, `UpdateItem` — структура ticket update webhook.
- `AddFollowupEvent` — структура followup webhook с `parent_item`.
- `FullFollowupWebhookPayload`, `FullTicketWebhookPayload` — полные модели полезной нагрузки.
- `FollowupItem` — данные комментария/фоллоуапа.
- `ParentItemRef` — ссылка на родительскую заявку.
- Использование `json.RawMessage` позволяет гибко принимать строки/числа/объекты.

### 5.11 [notifications/unmarshal.go](notifications/unmarshal.go)

- `unmarshalStringOrNumber(data)`:
    - нормализует пустое/null;
    - пытается распарсить как строку;
    - затем как `json.Number`;
    - затем как bool;
    - возвращает строковое представление.
- `UnmarshalStringOrNumber` — публичная обёртка с trim.
- `UnmarshalStatusValue`:
    - поддерживает статус в формате primitive или объекта `{id,name}`;
    - возвращает `name`, если есть;
    - иначе возвращает `id`.

### 5.12 [notifications/telegram_notifier.go](notifications/telegram_notifier.go)

- `htmlTagRe` — regex для удаления HTML-тегов.
- `FormatCommentText(raw)`:
    - trim;
    - заменяет `<br>`, `</p>`, `</div>` на переносы;
    - удаляет остальные HTML-теги;
    - декодирует HTML entities;
    - удаляет `\u00a0`;
    - чистит пустые строки.
- `normalizeTicketStatusRU(status)`:
    - приводит много вариантов статуса к русским каноническим значениям.
- `BuildTicketCommentNotification(...)` — шаблон сообщения о комментарии.
- `BuildTicketStatusChangedNotification(...)` — шаблон сообщения о новом статусе.
- `SendTelegramMessageToMany(chatIDs, text)`:
    - проверяет `TG_BOT_TOKEN` и непустой текст;
    - создаёт Telegram API клиент;
    - отправляет сообщение каждому `chat_id`.

### 5.13 [repository/users_repository.go](repository/users_repository.go)

- `ensureDB()` — проверка, что подключение к БД инициализировано.
- `GetUserByTelegramID(...)`:
    - валидирует входной ID;
    - делает `SELECT` пользователя;
    - сканирует в `models.User`.
- `SaveUser(user)`:
    - валидации (`nil`, `telegram_id`);
    - `INSERT ... ON CONFLICT (telegram_id) DO UPDATE`;
    - возвращает сохранённую запись через `RETURNING`.
- `EnsureUser(...)`:
    - пытается получить пользователя;
    - если нет (`sql.ErrNoRows`) — создаёт и сохраняет нового.
- `IsUserAuthorized(user)` — проверяет, что есть и `api_token`, и `session_token`.
- `ListAuthorizedTelegramIDs()` — выбирает всех пользователей с непустыми токенами.

### 5.14 [repository/tickets_repository.go](repository/tickets_repository.go)

- Константы `TicketStatus*` — внутренние русские статусы.
- `normalizeTicketStatus(status)` — приведение разных статусов к канону.
- `GetLastUserTicketsText(telegramID, limit)`:
    - джойн `tickets` + `users`;
    - сортировка по `created_at DESC`;
    - сборка человекочитаемого текста.
- `CreateTicketByTelegramID(...)` — обёртка с дефолтным статусом ожидания.
- `CreateTicketByTelegramIDWithStatus(...)`:
    - валидации входных данных;
    - вставка новой заявки для пользователя по `telegram_id`;
    - `RETURNING id` локальной заявки.
- `UpdateTicketStatus(ticketID, status)`:
    - обновляет статус и `updated_at`;
    - проверяет, что строка действительно обновилась.
- `BindExternalTicketID(localTicketID, externalTicketID)`:
    - записывает `external_ticket_id` локальной заявке.
- `ListTelegramIDsByExternalTicketID(externalTicketID)`:
    - ищет всех уникальных пользователей, связанных с внешним тикетом.

### 5.15 [telegram/commands.go](telegram/commands.go)

- Набор `const` callback-идентификаторов inline-кнопок.
- Значения (`cb_create_ticket`, `lang_ru`, ...) используются в switch-обработчиках callback.

### 5.16 [telegram/session.go](telegram/session.go)

- `sessionData` — данные шага диалога + активная заявка.
- `userSessions sync.Map` — потокобезопасное хранилище сессий по `telegramID`.
- Константы `Step*` — finite-state-machine шагов диалога.
- `getSession` — безопасно читает сессию, при отсутствии возвращает `StepIdle`.
- `setSession` — записывает/обновляет сессию.
- `resetSession` — удаляет сессию пользователя.

### 5.17 [telegram/keyboards.go](telegram/keyboards.go)

- `MainMenuKeyboard(lang)` — reply-клавиатура главного меню.
- `NotificationsKeyboard(lang)` — inline-клавиатура уведомлений.
- `PreferencesKeyboard(lang)` — reply-клавиатура настроек.
- `LanguageKeyboard()` — inline-кнопки выбора языка + кнопка назад.
- `CancelTicketKeyboard(lang)` — inline-кнопка отмены создания заявки.

### 5.18 [telegram/callback.go](telegram/callback.go)

- `handleCallback(bot, update)`:
    - проверяет наличие callback и сообщения;
    - `EnsureUser` для автора callback;
    - выставляет язык, если он ещё не задан;
    - создаёт сообщение-ответ;
    - получает текущую сессию.
- Блок неавторизованного пользователя:
    - ставит шаг `StepWaitApiToken`;
    - просит ввести API token;
    - подтверждает callback (`NewCallback`);
    - отправляет сообщение и сохраняет пользователя.
- `switch q.Data`:
    - `CbCreateTicket` — старт диалога создания заявки;
    - `CbCancelTicket` — отмена и возврат в меню;
    - `CbLastTickets` — вывод последних заявок;
    - `CbPreferences` / `CbNotifications` / `CbOpenLanguage` — экраны настроек;
    - `CbLangRU/EN/KK` — смена языка;
    - `CbStart/CbExit/CbMainMenu` — возврат в стартовое меню.
- После switch:
    - при необходимости сохраняет сессию;
    - сохраняет пользователя;
    - отправляет callback-ack и сообщение.

### 5.19 [telegram/handlers.go](telegram/handlers.go)

- `dueDateLayout` — формат даты `YYYY-MM-DD HH:MM:SS`.
- `HandleUpdate(bot, update)`:
    - если пришёл callback — делегирует в `handleCallback`;
    - если нет `Message` — игнорирует update;
    - извлекает `chatID`, `telegramID`, `text`.
- `EnsureUser(...)` — гарантирует существование пользователя в БД.
- Инициализация языка по `DetectLang`.
- `session := getSession(telegramID)` — загрузка состояния FSM.

**Ветка авторизации (`StepWaitApiToken`)**

- если токен пуст — повторно просит токен;
- вызывает `AuthByUserToken(text)`;
- при ошибке — сообщает об ошибке авторизации;
- при успехе сохраняет `ApiToken` и `SessionToken` в БД;
- сбрасывает сессию и возвращает успех.

**Глобальная проверка авторизации**

- если пользователь не авторизован — переводит в `StepWaitApiToken`.

**Команда `/start`**

- сбрасывает сессию;
- показывает главное меню.

**FSM создания заявки**

- `StepWaitTicketTitle` — сохраняет тему, просит описание;
- `StepWaitTicketDesc` — сохраняет описание, просит приоритет;
- `StepWaitTicketPriority` — валидирует число `0..5`, просит дедлайн;
- `StepWaitTicketDueDate`:
    - валидирует формат даты;
    - формирует `models.CreateTicketInput`;
    - создаёт локальную заявку в PostgreSQL;
    - создаёт заявку в GLPI через `CreateTicketWithSession`;
    - при `ERROR_SESSION_TOKEN_INVALID` пробует переавторизацию и повтор;
    - при успехе извлекает `external_ticket_id`;
    - связывает внешний ID с локальной заявкой;
    - переводит локальный статус в `В рассмотрении`;
    - отправляет пользователю сообщение об успехе.

**Default-ветка обычных текстовых сообщений**

- обрабатывает кнопочные тексты настроек/выхода/языка/создания/последних заявок;
- при `logout` очищает токены и снова переводит в `StepWaitApiToken`;
- для неизвестного текста показывает «выберите пункт меню».

**Финал функции**

- `repository.SaveUser(user)` — сохраняет изменения пользователя;
- `bot.Send(msg)` — отправляет итоговый ответ.

### 5.20 [telegram/bot.go](telegram/bot.go)

- загружает `.env`;
- читает `TG_BOT_TOKEN`;
- создаёт Telegram API client;
- включает debug (`bot.Debug = true`);
- запускает long polling `GetUpdatesChan`;
- в цикле передаёт каждый update в `HandleUpdate`.

### 5.21 [telegram/i18n/i18n.go](telegram/i18n/i18n.go)

- `ENG`, `RUS`, `KAZ` — словари переводов ключей интерфейса.
- `translations` — реестр словарей по коду языка.
- `T(lang, key)`:
    - ищет ключ в текущем языке;
    - fallback на русский;
    - если нет нигде — возвращает сам `key`.
- `DetectLang(code)`:
    - берёт первые 2 символа языка Telegram;
    - поддерживает `ru`, `kk`, `en`;
    - по умолчанию возвращает `ru`.

### 5.22 SQL-миграции

#### [migrations/001_init.up.sql](migrations/001_init.up.sql)

- создаёт таблицу `users`:
    - `id`, `telegram_id`, `username`, `api_token`, `session_token`, `lang`, `created_at`, `updated_at`;
- создаёт таблицу `tickets`:
    - `id`, `user_id` (FK на `users`), `title`, `description`, `status`, `created_at`, `updated_at`;
- создаёт индексы:
    - `idx_users_telegram_id`;
    - `idx_tickets_user_id`;
    - `idx_tickets_created_at`.

#### [migrations/001_init.down.sql](migrations/001_init.down.sql)

- удаляет таблицы `tickets`, затем `users`.

#### [migrations/002_add_session_token.up.sql](migrations/002_add_session_token.up.sql)

- добавляет колонку `session_token` в `users`, если отсутствует.

#### [migrations/002_add_session_token.down.sql](migrations/002_add_session_token.down.sql)

- удаляет колонку `session_token`.

#### [migrations/003_add_external_ticket_id.up.sql](migrations/003_add_external_ticket_id.up.sql)

- добавляет `external_ticket_id` в `tickets`;
- создаёт индекс `idx_tickets_external_ticket_id`.

#### [migrations/003_add_external_ticket_id.down.sql](migrations/003_add_external_ticket_id.down.sql)

- удаляет индекс `idx_tickets_external_ticket_id`;
- удаляет колонку `external_ticket_id`.

---

## 6) HTTP API (практическое использование)

### Проверка авторизации пользователя

- `GET /api/auth/check?token=<GLPI_USER_TOKEN>`
- Успех: `200` + `session_token`.
- Ошибка: `401`.

### Создание заявки

- `POST /api/ticket/create?session_token=<SESSION_TOKEN>`
- Body:

```json
{
    "input": {
        "name": "Проблема с принтером",
        "content": "Не печатает",
        "entities_id": 0,
        "status": 1,
        "priority": 3,
        "due_date": "2026-02-28 18:00:00"
    }
}
```

### Вебхук

- `POST /api/webhook`
- Авторизация:
    - через `X-Webhook-Key`/`Authorization: Bearer ...`;
    - или через `X-GLPI-Signature` + секрет.

---

## 7) Важные замечания

- `session_token` GLPI может протухать — в коде есть повторная авторизация.
- Если вебхук приходит без подходящего ключа/подписи — возвращается `401`.
- Локальный статус заявки и внешний статус GLPI синхронизируются через вебхуки.

---

## 8) Что можно улучшить дальше

- вынести миграции в отдельный мигратор (golang-migrate);
- добавить unit/integration тесты на repository и webhook-парсер;
- добавить retry/backoff для отправки Telegram-уведомлений;
- вынести конфиг в явную структуру с централизованной валидацией;
- добавить structured logging (например, `zerolog`).
