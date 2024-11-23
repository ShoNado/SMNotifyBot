@echo off
chcp 65001 >nul

REM Установите имя службы
set "ServiceName=TelegramBot"
set "ExePath=%~dp0SMBot_v3_go.exe"

REM Удаляем службу
cd %~dp0
sc stop "%ServiceName%"
sc delete "%ServiceName%"

REM Выводим сообщение
echo Служба "%ServiceName%" была удалена.

REM Установка службы
nssm install "%ServiceName%" "%ExePath%"
nssm set "%ServiceName%" Start SERVICE_AUTO_START
nssm set "%ServiceName%" AppExit Default Restart

REM Запуск службы
nssm start "%ServiceName%"

echo Служба "%ServiceName%" успешно создана и запущена.
pause