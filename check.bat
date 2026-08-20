@echo off
rem ============================================================================
rem  instrument -- run the kit's checks without typing the commands.
rem
rem  Same calls as .github/workflows/ci.yml, same order: first what counts
rem  numbers, then what reads the documentation, and the site build last. This
rem  is NOT a second list of checks -- a line added here must also appear in
rem  CI, otherwise local is green while push is red.
rem
rem     check                    every check
rem     check <name>             one of: contrast targets proportion docscheck
rem                              registry mutate dist lang
rem     check <name> -v          the flag is forwarded to the tool
rem     check dist               rebuild dist/ from src/
rem     check site               build the documentation site
rem     check serve [port]       build and serve; default :4322
rem     check pixels [/section/] rendered-pixel audit; needs Chrome and a server
rem     check pixels --mutate    the pixel audit, checked
rem     check behavior [/sect/]  keyboard contract of instrument.js; same needs
rem     check behavior --mutate  the behaviour checks, checked
rem     check fmt                gofmt over both modules
rem     check help               this list
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

rem The port lives in ONE place. It used to be typed into both `serve` and
rem `pixels` separately, and CONTRIBUTING showed a third value for the manual
rem invocation -- so serving on one port and auditing another looked like a
rem clean run over a server that was never asked anything.
rem
rem   check serve 4321
rem   set INSTRUMENT_PORT=4321
if not defined INSTRUMENT_PORT set "INSTRUMENT_PORT=4322"

if "%~1"==""          goto all
if /i "%~1"=="all"    goto all
if /i "%~1"=="dist"   goto build
if /i "%~1"=="site"   goto site
if /i "%~1"=="serve"  goto serve
if /i "%~1"=="pixels" goto pixels
if /i "%~1"=="behavior" goto behavior
if /i "%~1"=="fmt"    goto fmt
if /i "%~1"=="help"   goto usage
if /i "%~1"=="-h"     goto usage

rem -- One gate by name. The name is matched against the list instead of being
rem -- handed to `go run` as it comes: a typo would otherwise surface as a Go
rem -- build error about a missing package, which says nothing about what was
rem -- actually mistyped. Trailing arguments are forwarded, so `-v` works --
rem -- and `-v` is exactly what is wanted the moment a check goes red.
for %%g in (contrast targets proportion docscheck registry mutate lang) do (
  if /i "%~1"=="%%g" (
    call :gate %~1 %2 %3
    goto done
  )
)
echo.
echo   unknown check: %~1
goto usage

rem -- what this file offers -------------------------------------------------
:usage
echo.
echo   check                    every check, same as CI
echo   check ^<name^>             one of: contrast targets proportion docscheck registry mutate lang
echo   check ^<name^> -v          verbose; the flag is forwarded to the tool
echo.
echo   check dist               rebuild dist/ from src/
echo   check site               build the documentation site
echo   check serve [port]       build and serve; default :4322
echo   check pixels [/section/] rendered-pixel audit; needs "check serve" running
echo   check pixels --mutate    break one token, demand the audit goes red
echo   check behavior [/sect/]  keyboard contract of instrument.js; same needs
echo   check behavior --mutate  break one promise, demand the check goes red
echo.
echo   check fmt                gofmt over both modules. CI runs it; "all" does not
echo.
set "fails=0"
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
call :gate mutate
call :gate dist -check
call :gate lang
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
  rem A normal run prints one summary line per check: a list nobody reads is
  rem noise, not a report. `-v` asked for detail, so swallowing it after that
  rem would be perverse.
  rem
  rem The test is `-v` specifically, not "any argument": the full run passes
  rem `-check` to dist, which is plumbing, not a request to be shown everything.
  if /i "%~2"=="-v" (
    echo.
    type "%LOG%"
  ) else (
    call :tail
  )
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
if not "%~2"=="" set "INSTRUMENT_PORT=%~2"
echo.
echo   http://localhost:%INSTRUMENT_PORT%    Ctrl+C to stop
echo   No watch mode: restart after editing sources.
echo.
go -C site run ./cmd/site -serve :%INSTRUMENT_PORT%
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
rem   check pixels --mutate            the audit, checked
rem   check pixels --jobs 8            tabs at once; 4 by default, 8 is not faster
:pixels
echo.
echo   Needs "check serve" running in another window on :%INSTRUMENT_PORT%.
echo.
node tools\audit-run.mjs --base http://localhost:%INSTRUMENT_PORT% %2
if errorlevel 1 set /a fails+=1
goto done

rem -- what instrument.js promises the KEYBOARD, on the rendered page.
rem
rem Out of "all" for the same reason as `pixels`: a live server and Chrome.
rem Every other gate reads text and therefore cannot see markup that is right
rem on load and wrong a keystroke later -- a collapsed tree node taking the
rem only Tab stop with it was published in the reference itself.
rem
rem   check behavior                   every page
rem   check behavior /agent/           one section
rem   check behavior --mutate          the checks, checked
:behavior
echo.
echo   Needs "check serve" running in another window on :%INSTRUMENT_PORT%.
echo.
node tools\behavior-run.mjs --base http://localhost:%INSTRUMENT_PORT% %2 %3
if errorlevel 1 set /a fails+=1
goto done

rem -- formatting -------------------------------------------------------------
rem
rem Out of "all" on purpose, and the vet section above says why: on a working
rem copy checked out before .gitattributes pinned *.go to LF, gofmt flags every
rem file for line endings and a real violation drowns in the noise. That is a
rem reason to keep it off the default run, not a reason to make it unreachable
rem -- CI runs it, so it has to be runnable here too.
:fmt
echo.
gofmt -l tools site
if errorlevel 1 set /a fails+=1
echo   (empty above means formatted; names listed are not)
goto done

:done
chcp %_cp% >nul
endlocal & exit /b %fails%
