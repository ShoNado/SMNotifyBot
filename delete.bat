@echo off
chcp 65001 >nul

REM Укажите имя службы и путь к файлу logs.txt
set "ServiceName=TelegramBot"
set "LogFilePath=%~dp0logs.txt"

REM Удаляем службу
cd %~dp0
sc stop "%ServiceName%"
sc delete "%ServiceName%"

REM Выводим сообщение
echo Служба "%ServiceName%" была удалена.

REM Удаление файла logs.txt
if exist "%LogFilePath%" (
    del /f "%LogFilePath%"
    echo "%LogFilePath%" successfully deleted
)
pause
