param(
    [switch]$build_darwin,
    [switch]$build_linux_arm,
    [switch]$build_linux_x64,
    [switch]$build,
    [switch]$bump_version,
    [switch]$start,
    [switch]$test,
    [switch]$test_integration,
    [switch]$clean,
    [switch]$help,
    [switch]$vet,
    [switch]$fmt,
    [switch]$lint,
    [switch]$tools
)

function showHelp{
    Write-Host "Usage: 
    .\make.ps1 -build_darwin            build an executable for MacOS
    .\make.ps1 -build_linux_arm         build an executable for linux on arm platform
    .\make.ps1 -build_linux_x64         build an executable for linux on x64 platform
    .\make.ps1 -build                   build the windows executable
    .\make.ps1 -bump_version            bump the app version
    .\make.ps1 -start                   start the ergo proxy
    .\make.ps1 -fmt                     run gofmt on the source code    
    .\make.ps1 -vet                     run go vet on the source code    
    .\make.ps1 -lint                    run golint on the source code    
    .\make.ps1 -tools                   obtain the tools needed to run different make targets
    .\make.ps1 -test                    run the tests
    .\make.ps1 -test-integration        run the integration tests
    .\make.ps1 -clean                   remove all the created executables
    "
}

function fmt() {
    Write-Host "Running fmt ..." -ForegroundColor Green
    $rez = $(gofmt -l .)
    if($rez.Length -eq 0){
        return 
    }
    $noOfLines = $(Measure-Object $rez | Select-Object -ExpandProperty Count)
    if($noOfLines -ne 0){
        Write-Warning -Message "The following files are not correctly formated:"
        foreach($r in $rez){
            Write-Host $r -ForegroundColor Red
        }
        exit $noOfLines
    }
}
function vet(){
    Write-Host "Running vet ..." -ForegroundColor Green
    go vet ./...
    if($LASTEXITCODE -ne 0){
        exit $LASTEXITCODE
    }
}

function tools(){
    go get -u github.com/golang/lint/golint
}

function lint(){  
    tools  
    Write-Host "Running lint ..." -ForegroundColor Green
    $Env:Path+=";"+"$Env:GOPATH\bin"
    $rez = $(golint ./...)
    if($rez.Length -eq 0){
        return 
    }
    $noOfLines = $(Measure-Object $rez | Select-Object -ExpandProperty Count)
    if($noOfLines -ne 0){
        Write-Warning -Message "Lint says the following files are not ok:"
        foreach($r in $rez){
            Write-Host $r -ForegroundColor Red
        }
        exit $noOfLines
    }
}
function bumpVersion {
    git tag --sort=committerdate | Select-Object -Last 1 > .version
    Get-Content .version
}
function build(){
    bumpVersion
    go build -ldflags "-w -s -X main.VERSION=$(Get-Content .\.version)" -o bin/ergo.exe
}

if($build_darwin_arm) {
    Write-Host "Building darwin executable ..." -ForegroundColor Green
    &{$CGO_ENABLED=0;$GOOS="darwin";$GOARCH="amd64"; go build -o bin/darwin/ergo}    
}
if($build_linux_arm) {
    Write-Host "Building linux executable for the arm platform ..." -ForegroundColor Green
    &{$CGO_ENABLED=0;$GOOS="linux";$GOARCH="arm64"; go build -o bin/darwin/ergo}    
}
if($build_linux_x64) {
    Write-Host "Building linux executable for the 64 bit platform ..." -ForegroundColor Green
    &{$CGO_ENABLED=0;$GOOS="linux";$GOARCH="amd64"; go build -o bin/ergo}
}
if($build){
    Write-Host "Building windows executable ..." -ForegroundColor Green
    build
}
if($clean){
    Write-Host "Cleaning ..." -ForegroundColor Green
    Remove-Item -Recurse -Force .\bin\*
}
if($test){    
    fmt
    vet
    lint
    Write-Host "Running tests ..." -ForegroundColor Green
    go test -v ./... 
}

if($test_integration) {
    Write-Host "Running integration tests ..." -ForegroundColor Green
    build
    go test -tags=integration -v ./... 
    
}

if($bump_version){
    bumpVersion
}

if($help){
    showHelp
}

if($fmt){
    fmt
}

if($vet){
    vet
}

if($lint){
    lint
}

if($tools){
    tools
}

if(!$build -and !$build_darwin -and 
    !$build_linux_arm -and !$build_linux_x64 -and
    !$bump_version -and !$start -and
    !$test -and !$clean -and 
    !$test_integration -and !$help -and
    !$vet -and !$fmt -and !$lint){
        showHelp
}


