@echo off

set TARGET=%1
set BASE=%~dp0

if "%TARGET%" == "" (
	echo setup.cmd targetfolder\
	exit /B 1
)

cd %BASE%
git pull

set EXEID=updater
set BUILDEXE=%BASE%%EXEID%.exe
set TARGETEXE=%TARGET%%EXEID%.exe

go build -o %BUILDEXE% main.go

if exist %TARGETEXE%.old del /F %TARGETEXE%.old
if exist %TARGETEXE% ren %TARGETEXE% %TARGETEXE%.old
copy %BUILDEXE% %TARGETEXE%

exit /B 0
