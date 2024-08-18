@echo off

setLocal

for /F %%i in ('go env GOROOT') do ( set target_file=%%i/src/runtime/ellyn_goid.go)

echo target_file=%target_file%

if exist %target_file% exit
::set target_file=%target_file%_test
echo package runtime > %target_file%
echo. >>  %target_file%
echo func EllynGetGoid() uint64 { >> %target_file%
echo    return uint64(getg().goid) >> %target_file%
echo } >> %target_file%
echo write success