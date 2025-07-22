# ClinicHub – Мобильное приложение для врачей

## 🛠️ Используемые технологии  
 
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Git](https://img.shields.io/badge/Git-F05032?style=for-the-badge&logo=git&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)
![GORM](https://img.shields.io/badge/GORM-000000?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZpZXdCb3g9IjAgMCAyNTAgMjUwIj48cGF0aCBmaWxsPSIjMDAwMDAwIiBkPSJNMTI1IDI1MEM1NiAyNTAgMCAxOTQgMCAxMjVTNTYgMCAxMjUgMHMyMjUgNTYgMjI1IDEyNS01NiAxMjUtMTI1IDEyNXoiLz48cGF0aCBmaWxsPSIjRkZGRkZGIiBkPSJNMTI1IDQwYzQ3IDAgODUgMzggODUgODVzLTM4IDg1LTg1IDg1LTg1LTM4LTg1LTg1IDM4LTg1IDg1LTg1bTAtMTBjLTUyIDAtOTUgNDMtOTUgOTVzNDMgOTUgOTUgOTUgOTUtNDMgOTUtOTUtNDMtOTUtOTUtOTV6Ii8+PC9zdmc+&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white)


- **🗃️ PostgreSQL** – Реляционная база данных  
- **🐳 Docker** – Контейнеризация приложения  
- **🦾 Golang** – Язык backend-разработки  
- **🔧 Git** – Система контроля версий  
- **📜 Swagger** – Документирование API  
- **💾 GORM** – ORM для работы с БД  
- **🔐 JWT** – Аутентификация и авторизация 


## 📌 Основные функции

### 📅 Управление расписанием
- Интерактивный календарь с почасовым планом
- Быстрая навигация по дням/неделям
- Цветовая маркировка типов приёмов по статусам (завершился, не явился и др.)
- Возможность просматривать информацию о пациенте и редактировать её
- Можно начать приём, а также отметить то, что он состоялся

### 🚑 Экстренные вызовы
- Список вызовов с приоритетами
- На каждом вызове можно посмотреть всех пострадавших пациентов
- Если пациенты не были указаны системой - врач может указать их сам
- Возможность взять вызов в работу
- Редактирование статуса в реальном времени для информирования других врачей

### 👨‍⚕️ Список пациентов
- Краткая информация о пациенте: номер комнаты, диагноз и ФИО
- Полный доступ к медицинским картам
- Возможность просматривать информацию о пациенте и редактировать её
- История приёмов и заключений
---
### 🔒 Авторизация
**Вход в систему осуществляется через:**
- Номер телефона
- Пароль (выдаётся администрацией)
---
### 💻 Запуск для разработки
**Требования:**
- Go 1.20+
- PostgreSQL 14+
- Git

**Файл окружения:**
```.env
# Database settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=medapp

# App settings
APP_PORT=8080
JWT_SECRET=your_strong_jwt_secret_here
GIN_MODE=debug
```

**Запуск на Windows:**
``` PowerShell
git clone https://github.com/AlexanderMorozov1919/mobileapp.git
cd mobileapp
cp .env.example .env
.\run.bat
```