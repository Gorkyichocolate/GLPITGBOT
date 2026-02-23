Создание тикета в glpi
Шаг 1: Зарегистрировать OAuth-клиент (однократно)

    В GLPI: Setup > OAuth Clients > Add
    После создания сохраните client_id и client_secret api_doc.MD:22-26 .

Шаг 2: Получить access-токен (Password Grant)

    Method: POST
    URL: {{GLPI_URL}}/api.php/token
    Headers: Content-Type: application/json
    Body (raw JSON):

{  
  "grant_type": "password",  
  "client_id": "ВАШ_CLIENT_ID",  
  "client_secret": "ВАШ_CLIENT_SECRET",  
  "username": "ВАШ_ЛОГИН",  
  "password": "ВАШ_ПАРОЛЬ",  
  "scope": "api"  
}

    Ответ (200 OK):

{  
  "token_type": "Bearer",  
  "expires_in": 3600,  
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1..."  
}

api_doc.MD:40-55
Шаг 3: Создать тикет

    Method: POST
    URL: {{GLPI_URL}}/api.php/Ticket (или /api.php/v2/Ticket для high‑level API) api_doc.MD:191-193
    Headers:
        Authorization: Bearer {{access_token}}
        Content-Type: application/json
    Body (raw JSON):

{  
  "input": {  
    "name": "Проблема с принтером",  
    "content": "Принтер не печатает. Проверьте кабель и драйверы.",  
    "entities_id": 0,  
    "type": 1,  
    "urgency": 3,  
    "requesttypes_id": 1,  
    "_users_id_requester": 2  
  }  
}

    Ответ (201 Created):

{  
  "id": 123,  
  "message": ""  
}

apirest.md:1022-1062
Шаг 4: Проверить созданный тикет (опционально)

    Method: GET
    URL: {{GLPI_URL}}/api.php/Ticket/{{id}}
    Headers: Authorization: Bearer {{access_token}}
    Ответ: JSON с полями тикета.

Подробности и вариации
Альтернативный URL с версией

    Указание версии: /api.php/v2/Ticket api_doc.MD:191-200 .
    Без версии используется последняя.

Дополнительные поля (опционально)

    type: 1 = Incident, 2 = Request ITILController.php:362-371 .
    status: ID статуса из Ticket::getAllStatusArray() ITILController.php:383-401 .
    entity: объект сущности (dropdown) ITILController.php:402-404 .
    _actors для назначения групп/пользователей через JSON (как в web) ticket.form.php:69-73 .

Пример из тестов (Cypress)

    В тестах используется cy.createWithAPI('Ticket', { name: 'Test ticket', content: 'Test ticket' }) .
    Внутри это делает POST на /Ticket с input.

Обработка ошибок

    400 Bad Request: неверные параметры.
    401 Unauthorized: неверный/просроченный токен.
    403 Forbidden: нет прав на создание в указанной сущности.

Важные замечания

    Убедитесь, что у пользователя есть право CREATE на тикеты в целевой сущности.
    Токен по умолчанию живёт 1 час api_doc.MD:52 .
    Для загрузки документов используйте multipart/form-data и uploadManifest apirest.md:1035-1038 .
