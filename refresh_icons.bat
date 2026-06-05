@echo off
taskkill /f /im explorer.exe > nul
CD /d %userprofile%\AppData\Local
DEL IconCache.db /a > nul 2>&1
CD /d %userprofile%\AppData\Local\Microsoft\Windows\Explorer
DEL iconcache* /a > nul 2>&1
DEL thumbcache* /a > nul 2>&1
start explorer.exe
echo Иконки успешно перестроены!