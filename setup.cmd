rem @echo off

set ZH_HOME=%1
set SCRIPTDIR=%~dp0

if "%ZH_HOME%" == "" (
	echo setup.cmd targetfolder\
	exit /B 1
)

cd %SCRIPTDIR%
git pull

set EXEID=updater
set SOURCE=%SCRIPTDIR%%EXEID%.exe
set TARGET=%ZH_HOME%%EXEID%.exe

go build -o %SOURCE% main.go

if exist %TARGET%.old del /Y %TARGET%.old
if exist %TARGET% ren %TARGET% %TARGET%.old
copy %SOURCE% %TARGET%

exit /B 0
