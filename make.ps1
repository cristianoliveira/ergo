param(
    [switch]$build_darwin,
    [switch]$build_linux_arm,
    [switch]$build_linux_x64,
    [switch]$build,
    [switch]$bump_version,
    [switch]$start,
    [switch]$test,
    [switch]$clean
)

function build(){
    go build -o bin/ergo.exe
}

if($build_darwin_arm) {
    Write-Host "Building darwin executable ..."
    &{$CGO_ENABLED=0;$GOOS="darwin";$GOARCH="amd64"; go build -o bin/darwin/ergo}    
}
if($build_linux_arm) {
    Write-Host "Building linux executable for the arm platform ..."
    &{$CGO_ENABLED=0;$GOOS="linux";$GOARCH="arm64"; go build -o bin/darwin/ergo}    
}
if($build_linux_x64) {
    Write-Host "Building linux executable for the 64 bit platform ..."
    &{$CGO_ENABLED=0;$GOOS="linux";$GOARCH="amd64"; go build -o bin/ergo}
}
if($build){
    Write-Host "Building windows executable ..."
    build
}
if($clean){
    Write-Host "Cleaning ..."
    Remove-Item -Recurse -Force .\bin\*
}
if($test){
    Write-Host "Starting tests ..."
    build
    go test -v -tags=integration ./... 
}
if($bump_version){
    git tag --sort=committerdate | tail -n 1 > .version
    Get-Content .version
}

function showHelp{
    Write-Host "Usage: 
    .\make.ps1 -build_darwin            build an executable for MacOS
    .\make.ps1 -build_linux_arm         build an executable for linux on arm platform
    .\make.ps1 -build_linux_x64         build an executable for linux on x64 platform
    .\make.ps1 -build                   build an executable for windows
    .\make.ps1 -bump_version            bump the app version
    .\make.ps1 -start                   start the ergo proxy
    .\make.ps1 -test                    run the tests
    .\make.ps1 -clean                   remove all the created executables
    "
}

if(!$build -and !$build_darwin -and 
    !$build_linux_arm -and !$build_linux_x64 -and
    !$bump_version -and !$start -and
    !$test -and !$clean){
        showHelp
}


