@echo off
chcp 65001 >nul 
echo === Generating Swagger documentation ===
swag init -g .\cmd\app\main.go

IF %ERRORLEVEL% NEQ 0 (
    echo Swagger generation failed. Stopping execution.
    exit /b %ERRORLEVEL%
)

echo === Starting the application ===
go run .\cmd\app\main.go
