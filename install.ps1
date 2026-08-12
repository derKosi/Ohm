# Ohm installer for Windows (PowerShell).
# Usage:
#   irm https://github.com/derKosi/Ohm/releases/latest/download/install.ps1 | iex
#
# Downloads the latest Windows release, extracts ohm.exe to
# $env:USERPROFILE\.ohm\bin, and adds that directory to the user PATH.
# No admin privileges required.
[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'

$Owner = 'derKosi'
$Repo  = 'Ohm'

function Write-Step($msg) { Write-Host "==> $msg" -ForegroundColor Cyan }
function Write-Ok($msg)   { Write-Host "    $msg" -ForegroundColor Green }
function Write-Warn($msg) { Write-Host "    $msg" -ForegroundColor Yellow }
function Die($msg)        { Write-Host "Error: $msg" -ForegroundColor Red; exit 1 }

# --- Detect architecture ---
$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { 'x86_64' }
    'ARM64' { 'arm64' }
    default { Die "Unsupported architecture: $env:PROCESSOR_ARCHITECTURE. Download a binary from https://github.com/$Owner/$Repo/releases" }
}

# --- Find the latest release tag via redirect ---
Write-Step 'Detecting latest release...'
$latestUrl = "https://github.com/$Owner/$Repo/releases/latest"
try {
    $resp = Invoke-WebRequest -Uri $latestUrl -MaximumRedirection 0 -ErrorAction SilentlyContinue
} catch {
    # A 302 redirect surfaces as a non-terminating error here; grab the Location.
    $resp = $_.Exception.Response
}
if (-not $resp) {
    Die "Could not reach $latestUrl. Check your internet connection."
}
$finalUrl = if ($resp.Headers.Location) { $resp.Headers.Location.ToString() } else { $latestUrl }
$tag = $finalUrl -replace '.*/tag/', ''
if (-not $tag) { Die "Could not determine latest release tag." }

$version = $tag -replace '^v', ''
$asset = "ohm_${version}_Windows_${arch}.zip"
$url = "https://github.com/$Owner/$Repo/releases/download/$tag/$asset"

Write-Step "Downloading Ohm $tag for Windows/$arch..."

# --- Download to temp ---
$tempDir = Join-Path $env:TEMP "ohm-install-$(Get-Random)"
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
$zipPath = Join-Path $tempDir $asset

try {
    Invoke-WebRequest -Uri $url -OutFile $zipPath -UseBasicParsing
} catch {
    Die "Download failed: $($_.Exception.Message)"
}

# --- Extract ---
Write-Step 'Extracting...'
Expand-Archive -Path $zipPath -DestinationPath $tempDir -Force

# --- Install to $USERPROFILE\.ohm\bin ---
$installDir = Join-Path $env:USERPROFILE '.ohm\bin'
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}
$dest = Join-Path $installDir 'ohm.exe'
Move-Item -Path (Join-Path $tempDir 'ohm.exe') -Destination $dest -Force

Write-Ok "Installed ohm.exe to $installDir"

# --- Add to user PATH (no admin needed) ---
$userPath = [Environment]::GetEnvironmentVariable('PATH', 'User')
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable('PATH', "$userPath;$installDir", 'User')
    Write-Ok "Added $installDir to your user PATH."
    Write-Warn "Restart your terminal (or open a new one) for PATH changes to take effect."
} else {
    Write-Ok "$installDir is already in your PATH."
}

# --- Cleanup ---
Remove-Item -Recurse -Force $tempDir -ErrorAction SilentlyContinue

Write-Host ''
Write-Host '    Run "ohm scan" to scan your system for AI software.' -ForegroundColor White
Write-Host "    https://github.com/$Owner/$Repo" -ForegroundColor DarkGray
