$installDir = "$Env:UserProfile\Ergo"
$downloadTo = "$Env:TEMP\ergo.exe"

function Test-RegistryValue {
    param (
     [parameter(Mandatory=$true)]
     [ValidateNotNullOrEmpty()]$Path,
    [parameter(Mandatory=$true)]
     [ValidateNotNullOrEmpty()]$Value
    )
    
    try {    
        Get-ItemProperty -Path $Path | Select-Object -ExpandProperty $Value -ErrorAction Stop | Out-Null
        return $true
    }
    catch {
        return $false
    }
}
    

Invoke-WebRequest -Uri "https://raw.githubusercontent.com/cristianoliveira/ergo/master/.version" -OutFile .version

$v = Get-Content .version | Out-String
$v = $v.Trim()
Write-Host "Ergo $v will be installed to: $installDir"

Invoke-WebRequest -Uri "https://github.com/cristianoliveira/ergo/releases/download/$v/ergo.exe" -OutFile $downloadTo


If(!(Test-Path $installDir))
{
    New-Item -ItemType Directory -Force -Path $installDir
}

Copy-Item -Path "$downloadTo" -Dest "$installDir\ergo.exe"

if(!(Test-Path "HKCU:\Environment"))
{
    New-Item -Path "HKCU:\Environment" -Force
    #reg add "HKCU\Environment" /v Path /f /t REG_SZ /d "$Env:path;$installDir"
}

$pathNow = $installDir

if(Test-RegistryValue -Path "HKCU:\Environment" -Value "Path"){
    $pathNow = Get-ItemPropertyValue -Path HKCU:\Environment  -Name "Path"
    $pathNow = "$pathNow;$installDir"
}

Write-Output $pathNow

New-ItemProperty -Path "HKCU:\Environment" -Name "Path" -Value "$pathNow" -PropertyType String -Force

setx ERGO_PATH "$installDir"