@echo off
REM Configures the environment variables required to build [neon-kubernetes] projects.
REM 
REM 	buildenv [ <source folder> ]
REM
REM Note that <source folder> defaults to the folder holding this
REM batch file.
REM
REM This must be [RUN AS ADMINISTRATOR].

echo ============================================
echo * NEONCLOUD Build Environment Configurator *
echo ============================================

REM Default NKK_ROOT to the folder holding this batch file after stripping
REM off the trailing backslash.

set NKK_ROOT=%~dp0 
set NKK_ROOT=%NKK_ROOT:~0,-2%

if not [%1]==[] set NKK_ROOT=%1

if exist %NKK_ROOT%\code-of-conduct.md goto goodPath
echo The [%NKK_ROOT%\code-of-conduct.md] file does not exist.  Please pass the path
echo to the [neon-kubernetes] solution folder.
goto done

:goodPath 

echo.
echo Configuring...
echo.

REM Persist the environment variables.

setx NKK_ROOT "%NKK_ROOT%" /M > nul

:done
echo ============================================================================================
echo * Be sure to close and reopen Visual Studio and any command windows to pick up the changes *
echo ============================================================================================
pause
