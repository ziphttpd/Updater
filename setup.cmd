rem @echo off

set ZH_HOME=%1
set SCRIPTDIR=%~dp0

if "%ZH_HOME%" == "" (
	echo setup.cmd targetfolder\
	exit /B 1
)

cd %SCRIPTDIR%
git pull

set FILE=updater.exe
set SOURCE=%SCRIPTDIR%%FILE%
set TARGET=%ZH_HOME%%FILE%

go build -o %SOURCE% main.go

if exist %TARGET%.old del /Y %TARGET%.old
if exist %TARGET% ren %TARGET% %FILE%.old
copy %SOURCE% %TARGET%

exit /B 0
