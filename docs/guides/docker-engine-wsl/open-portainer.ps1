<#
.SYNOPSIS
    Opens Portainer in Google Chrome, starting the WSL distro and the Portainer container first when they are not running.
.DESCRIPTION
    Idempotent and per-user: no administrator rights, no UAC prompt. If Portainer already answers on the URL, Chrome opens at once.
    Otherwise the script starts the WSL distro, starts the container, waits for the page to answer, then opens Chrome.
    Chrome is started by its full path, never through the App Paths registry key, because that key can point to a different Chromium browser.
    With -CreateShortcut the script copies itself to %LOCALAPPDATA%\Programs\docker-wsl and creates Portainer.lnk on the Desktop and in the Start Menu.
    The shortcut icon is the Portainer logo, fetched from the running Portainer page and saved as portainer.ico beside the script; without a running Portainer the icon falls back to Chrome's.
.PARAMETER Distro
    Name of the WSL distro that runs Docker Engine.
.PARAMETER Container
    Name of the Portainer container.
.PARAMETER Url
    Address Portainer listens on.
.PARAMETER BrowserPath
    Full path of the browser to use instead of the default search for Google Chrome.
.PARAMETER TimeoutSeconds
    How long to wait for the engine and the page before giving up.
.PARAMETER CreateShortcut
    Install the script under %LOCALAPPDATA%\Programs\docker-wsl and create the Desktop and Start Menu shortcuts, then exit.
.PARAMETER NoBrowser
    Prepare the environment but do not open a browser.
.PARAMETER NoDialog
    Print errors to the console only, without a message box.
#>
[CmdletBinding()]
param(
    [string]$Distro = "Ubuntu-24.04",
    [string]$Container = "portainer",
    [string]$Url = "https://localhost:9443/",
    [string]$BrowserPath,
    [int]$TimeoutSeconds = 90,
    [switch]$CreateShortcut,
    [switch]$NoBrowser,
    [switch]$NoDialog
)

$ErrorActionPreference = "Stop"

# Native commands write progress to stderr; under "Stop" Windows PowerShell 5.1 turns that into a terminating error.
function Invoke-Native {
    param([string]$File, [string[]]$Arguments)
    $previous = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    try {
        $output = & $File @Arguments 2>&1 | Out-String
        return [pscustomobject]@{ ExitCode = $LASTEXITCODE; Output = $output.Trim() }
    } finally {
        $ErrorActionPreference = $previous
    }
}

function Get-HttpCode {
    $result = Invoke-Native "curl.exe" @("-k", "-s", "-o", "NUL", "-w", "%{http_code}", "-m", "3", $Url)
    return $result.Output
}

function Test-PageUp {
    return ((Get-HttpCode) -match '^[23]\d\d$')
}

function Find-Chrome {
    $roots = @($env:ProgramFiles, ${env:ProgramFiles(x86)}, $env:LOCALAPPDATA) | Where-Object { $_ }
    foreach ($root in $roots) {
        $candidate = Join-Path $root "Google\Chrome\Application\chrome.exe"
        if (Test-Path $candidate) { return $candidate }
    }
    return $null
}

# Portainer lists its logo as PNG files with hashed names in <link> tags. The PNG is wrapped in a single-image .ico, which a shortcut can use. The logo comes from the running instance and is not stored in this repository.
function Save-PortainerIcon {
    param([string]$Directory)
    $html = (Invoke-Native "curl.exe" @("-k", "-s", "-m", "10", $Url)).Output
    $match = [regex]::Match($html, '<link rel="apple-touch-icon"[^>]*href="([^"]+\.png)"')
    if (-not $match.Success) { return $null }
    $pngUrl = (New-Object System.Uri ([System.Uri]$Url), $match.Groups[1].Value).AbsoluteUri
    $png = Join-Path $Directory "portainer.png"
    $download = Invoke-Native "curl.exe" @("-k", "-s", "-m", "20", "-o", $png, $pngUrl)
    if ($download.ExitCode -ne 0 -or -not (Test-Path $png)) { return $null }
    $bytes = [System.IO.File]::ReadAllBytes($png)
    Remove-Item $png -Force
    if ($bytes.Length -lt 24 -or $bytes[0] -ne 0x89 -or $bytes[1] -ne 0x50 -or $bytes[2] -ne 0x4E -or $bytes[3] -ne 0x47) { return $null }
    $width = $bytes[18] * 256 + $bytes[19]
    $height = $bytes[22] * 256 + $bytes[23]
    $stream = New-Object System.IO.MemoryStream
    $writer = New-Object System.IO.BinaryWriter($stream)
    $writer.Write([uint16]0)
    $writer.Write([uint16]1)
    $writer.Write([uint16]1)
    $writer.Write([byte]$(if ($width -ge 256) { 0 } else { $width }))
    $writer.Write([byte]$(if ($height -ge 256) { 0 } else { $height }))
    $writer.Write([byte]0)
    $writer.Write([byte]0)
    $writer.Write([uint16]1)
    $writer.Write([uint16]32)
    $writer.Write([uint32]$bytes.Length)
    $writer.Write([uint32]22)
    $writer.Write($bytes)
    $writer.Flush()
    $ico = Join-Path $Directory "portainer.ico"
    [System.IO.File]::WriteAllBytes($ico, $stream.ToArray())
    $writer.Close()
    return $ico
}

function Install-Shortcuts {
    $installDir = Join-Path $env:LOCALAPPDATA "Programs\docker-wsl"
    New-Item -ItemType Directory -Force $installDir | Out-Null
    $installed = Join-Path $installDir "open-portainer.ps1"
    if ($PSCommandPath -and ($PSCommandPath -ne $installed)) {
        Copy-Item -Path $PSCommandPath -Destination $installed -Force
    }

    $arguments = '-NoProfile -File "' + $installed + '"'
    foreach ($name in @("Distro", "Container", "Url", "BrowserPath", "TimeoutSeconds")) {
        if ($PSBoundParameters.ContainsKey($name)) {
            $arguments += " -$name `"$($PSBoundParameters[$name])`""
        }
    }

    $powershell = Join-Path $env:SystemRoot "System32\WindowsPowerShell\v1.0\powershell.exe"
    $chrome = Find-Chrome
    $icon = if ($chrome) { "$chrome,0" } else { "$powershell,0" }
    $portainerIcon = if (Test-PageUp) { Save-PortainerIcon -Directory $installDir } else { $null }
    if ($portainerIcon) {
        $icon = "$portainerIcon,0"
    } else {
        Write-Host "Portainer did not answer, so the shortcut uses a fallback icon. Run -CreateShortcut again while Portainer runs to get the Portainer icon." -ForegroundColor Yellow
    }

    $shell = New-Object -ComObject WScript.Shell
    foreach ($folder in @([Environment]::GetFolderPath("Desktop"), [Environment]::GetFolderPath("Programs"))) {
        $path = Join-Path $folder "Portainer.lnk"
        $link = $shell.CreateShortcut($path)
        $link.TargetPath = $powershell
        $link.Arguments = $arguments
        $link.WorkingDirectory = $installDir
        $link.WindowStyle = 7
        $link.IconLocation = $icon
        $link.Description = "Start WSL Docker and Portainer if needed, then open Portainer in Chrome"
        $link.Save()
        Write-Host "Created $path" -ForegroundColor Green
    }
    Write-Host "Script installed at $installed" -ForegroundColor Green
}

function Start-Environment {
    Write-Host "[1/4] Checking Portainer at $Url ..." -ForegroundColor Cyan
    if (Test-PageUp) {
        Write-Host "Portainer already answers." -ForegroundColor Green
        return
    }

    Write-Host "[2/4] Starting WSL distro $Distro ..." -ForegroundColor Cyan
    $installed = (Invoke-Native "wsl.exe" @("-l", "-q")).Output -replace "`0", ""
    if (($installed -split "\r?\n" | ForEach-Object { $_.Trim() }) -notcontains $Distro) {
        throw "WSL distro '$Distro' is not installed. Install it with the steps in the guide."
    }
    $started = Invoke-Native "wsl.exe" @("-d", $Distro, "--exec", "true")
    if ($started.ExitCode -ne 0) {
        throw "WSL distro '$Distro' did not start: $($started.Output -replace "`0", '')"
    }

    Write-Host "[3/4] Starting container $Container ..." -ForegroundColor Cyan
    $deadline = (Get-Date).AddSeconds($TimeoutSeconds)
    $running = $null
    while ($true) {
        $inspect = Invoke-Native "wsl.exe" @("-d", $Distro, "--", "docker", "inspect", "-f", "{{.State.Running}}", $Container)
        if ($inspect.ExitCode -eq 0) { $running = $inspect.Output; break }
        if ($inspect.Output -match "No such object") {
            throw "Container '$Container' does not exist in distro '$Distro'. Create it with the Portainer step in the guide."
        }
        if ((Get-Date) -gt $deadline) {
            throw "Docker did not answer within $TimeoutSeconds seconds: $($inspect.Output)"
        }
        Start-Sleep -Seconds 2
    }
    if ($running -ne "true") {
        $start = Invoke-Native "wsl.exe" @("-d", $Distro, "--", "docker", "start", $Container)
        if ($start.ExitCode -ne 0) {
            throw "Could not start container '$Container': $($start.Output)"
        }
    }

    Write-Host "[4/4] Waiting for $Url ..." -ForegroundColor Cyan
    while (-not (Test-PageUp)) {
        if ((Get-Date) -gt $deadline) {
            throw "Portainer did not answer on $Url within $TimeoutSeconds seconds."
        }
        Start-Sleep -Seconds 2
    }
    Write-Host "Portainer is ready." -ForegroundColor Green
}

function Open-Browser {
    $browser = if ($BrowserPath) { $BrowserPath } else { Find-Chrome }
    if ($browser -and (Test-Path $browser)) {
        Start-Process -FilePath $browser -ArgumentList $Url
    } else {
        Write-Host "Google Chrome was not found; opening the default browser." -ForegroundColor Yellow
        Start-Process $Url
    }
}

try {
    if ($CreateShortcut) {
        Install-Shortcuts
        return
    }
    Start-Environment
    if (-not $NoBrowser) { Open-Browser }
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
    if (-not $NoDialog) {
        Add-Type -AssemblyName System.Windows.Forms
        [void][System.Windows.Forms.MessageBox]::Show($_.Exception.Message, "Portainer launcher", "OK", "Error")
    }
    exit 1
}
