# TodoApp

Веб-сервис "список дел"

## Доменная модель

### Users
- id: UUID
- version: int (оптимистичная блокировка)
- full_name: string (обязательно 3-100)
- phone_number: string (необязателен, формат +цифры, длина 10-15)

### Tasks
- id: UUID
- version: int (оптимистичная блокировка)
- title: string (1-100 символов) обязательное поле
- description: string (необязательно(длинна 1-1000)) 
- completed: bool (флаг выполнения)
- created_at: timestamp с часовым поясом
- completed_at: timestamp с часовым поясом (необязательно)
- author_user_id: UUID (связь с users) 

### Statistics
Вычисляемый агрегат по задачам. Не хранится в БД
Создано, выполнено, процент выполнения, среднее время до выполнения

## REST API

### Users (/api/v1/users)
- POST /users создать пользователя
- GET /users получить список пользователей(с пагинацией)
- GET /users/{id} получить одного пользователя по id
- PATCH /users/{id} частично обновить пользователя
- DELETE /users/{id} удалить пользователя 

### Tasks (/api/v1/tasks)
- POST /tasks создать задачу
- GET /tasks получить список задач(с пагинацией + фильтр по user_id)
- GET /tasks/{id} получить одну задачу по id
- PATCH /tasks/{id} частично обновить задачу
- DELETE /tasks/{id} удалить задачу 

### Statistics(/api/v1/statistics)
- GET /statistics - агрегаты (фильтры user_id, from, to)
