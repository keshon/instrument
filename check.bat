@echo off
rem ============================================================================
rem  instrument -- run the kit's checks without typing the commands.
rem
rem  Same calls as .github/workflows/ci.yml, same order: first what counts
rem  numbers, then what reads the documentation, and the site build last. This
rem  is NOT a second list of checks -- a line added here must also appear in
rem  CI, otherwise local is green while push is red.
rem
rem     check              every check
rem     check contrast     one of: contrast targets proportion docscheck registry
rem     check dist         rebuild dist/ from src/
rem     check site         build the documentation site
rem     check serve        build and serve on :4322
rem     check pixels       rendered-pixel audit; needs Chrome and a live server
rem
rem  WHY THIS FILE IS IN ENGLISH, unlike everything else in the repository.
rem  Not a language decision -- a cmd.exe one. A .bat is parsed byte by byte in
rem  whatever code page is active when the line is READ, so multi-byte text in
rem  the file itself desynchronises the parser: it lands mid-line, runs the
rem  fragment as a command and falls through into unrelated labels. Writing the
rem  file in cp866 would work and would rot the moment an editor saved it back
rem  as UTF-8. ASCII has neither failure mode.
rem
rem  The Go tools print UTF-8, so the code page is switched to 65001 for their
rem  output and restored afterwards: leaving it changed in someone else's
rem  console is not ours to do.
rem ============================================================================

setlocal enabledelayedexpansion
for /f "tokens=2 delims=:." %%a in ('chcp') do set "_cp=%%a"
chcp 65001 >nul
cd /d "%~dp0"

set "LOG=%TEMP%\instrument-check.log"
set "fails=0"

if "%~1"==""          goto all
if /i "%~1"=="all"    goto all
if /i "%~1"=="dist"   goto build
if /i "%~1"=="site"   goto site
if /i "%~1"=="serve"  goto serve
if /i "%~1"=="pixels" goto pixels

call :gate %~1
goto done

rem -- everything ------------------------------------------------------------
:all
echo.
echo   instrument
echo   ------------------------------------------------------------
call :gate contrast
call :gate targets
call :gate proportion
call :gate docscheck
call :gate registry
call :gate dist -check
call :sitegate
call :vetgate
echo   ------------------------------------------------------------
if "%fails%"=="0" (
  echo   all green
) else (
  echo   FAILURES: %fails%
)
goto done

rem -- one tool from tools/cmd -----------------------------------------------
:gate
set "n=%~1"
call :label "%n%"
go -C tools run ./cmd/%n% %2 %3 >"%LOG%" 2>&1
if errorlevel 1 (
  echo FAILED
  type "%LOG%"
  set /a fails+=1
) else (
  call :tail
)
exit /b

rem -- the site build is its own check: broken link, unexpanded fence,
rem -- missing sprite symbol, empty example
:sitegate
call :label "site"
go -C site run ./cmd/site >"%LOG%" 2>&1
if errorlevel 1 (
  echo FAILED
  type "%LOG%"
  set /a fails+=1
) else (
  call :tail
)
exit /b

rem -- go vet over both modules ----------------------------------------------
rem
rem gofmt is deliberately NOT run here. On a working copy checked out before
rem .gitattributes pinned *.go to LF, gofmt flags every file for line endings
rem and a real violation drowns in the noise. CI runs it on Linux where the
rem checkout is LF, which is the only place the answer is meaningful.
:vetgate
call :label "vet"
(go -C tools vet ./... && go -C site vet ./...) >"%LOG%" 2>&1
if errorlevel 1 (
  echo FAILED
  type "%LOG%"
  set /a fails+=1
) else (
  echo clean
)
exit /b

rem -- fixed-width name so the summaries line up in a column ------------------
:label
set "pad=%~1            "
<nul set /p "=  !pad:~0,12!"
exit /b

rem -- last non-empty line of the log: every tool ends with its summary -------
:tail
set "last="
for /f "usebackq delims=" %%l in ("%LOG%") do set "last=%%l"
echo !last!
exit /b

rem -- rebuild the distribution ----------------------------------------------
:build
echo.
go -C tools run ./cmd/dist
if errorlevel 1 set /a fails+=1
goto done

:site
echo.
go -C site run ./cmd/site
if errorlevel 1 set /a fails+=1
goto done

:serve
echo.
echo   http://localhost:4322    Ctrl+C to stop
echo   No watch mode: restart after editing sources.
echo.
go -C site run ./cmd/site -serve :4322
goto done

rem -- rendered-pixel audit ---------------------------------------------------
rem
rem Kept out of "all" on purpose: it needs a live server and Chrome, runs for
rem minutes, and measures what actually landed on screen -- contrast in six
rem themes, tap targets in three densities. The text gates cannot see that by
rem construction.
rem
rem   check pixels                     every page
rem   check pixels /components/        one section
:pixels
echo.
echo   Needs "check serve" running in another window.
echo.
node tools\audit-run.mjs --base http://localhost:4322 %2
if errorlevel 1 set /a fails+=1
goto done

:done
chcp %_cp% >nul
endlocal & exit /b %fails%
