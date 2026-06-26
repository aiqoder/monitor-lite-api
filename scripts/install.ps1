# wtv (monitor-lite-api) Windows 一键安装脚本
# 用法（管理员 PowerShell）:
#   irm https://github.com/<owner>/<repo>/releases/latest/download/install.ps1 | iex
#   irm ... | iex -Install
#   irm ... | iex -Update
#   irm ... | iex -Uninstall

param(
    [switch]$Install,
    [switch]$Update,
    [switch]$Uninstall,
    [string]$InstallPath = "$env:ProgramFiles\wtv",
    [string]$GitHubRepo = "aiqoder/monitor-lite-api",
    [string]$Version = "latest",
    [string]$BinaryPrefix = "monitor-lite-api",
    [string]$ServiceName = "wtv"
)

$ErrorActionPreference = "Stop"

function Write-Info($Message) { Write-Host $Message -ForegroundColor Green }
function Write-Warn($Message) { Write-Host $Message -ForegroundColor Yellow }
function Write-Err($Message) { Write-Host $Message -ForegroundColor Red }

function Get-Arch {
    switch ((Get-CimInstance Win32_Processor).Architecture) {
        9 { return "amd64" }   # x64
        12 { return "arm64" }  # ARM64
        default {
            if ($env:PROCESSOR_ARCHITECTURE -match "ARM64") { return "arm64" }
            if ($env:PROCESSOR_ARCHITECTURE -eq "AMD64") { return "amd64" }
            throw "不支持的 CPU 架构"
        }
    }
}

function Get-DownloadUrl {
    param([string]$Arch)
    $asset = "$BinaryPrefix-windows-$Arch.exe"
    if ($Version -eq "latest") {
        return "https://github.com/$GitHubRepo/releases/latest/download/$asset"
    }
    return "https://github.com/$GitHubRepo/releases/download/$Version/$asset"
}

function Ensure-Admin {
    $current = New-Object Security.Principal.WindowsPrincipal(
        [Security.Principal.WindowsIdentity]::GetCurrent()
    )
    if (-not $current.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        throw "请以管理员身份运行 PowerShell"
    }
}

function Ensure-Ffmpeg {
    if (Get-Command ffmpeg -ErrorAction SilentlyContinue) { return }
    Write-Warn "未检测到 ffmpeg，请先安装并加入 PATH"
    Write-Warn "推荐: winget install Gyan.FFmpeg  或  choco install ffmpeg"
    throw "缺少 ffmpeg 依赖"
}

function Install-Service {
    $binary = Join-Path $InstallPath "wtv.exe"
    $existing = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
    if ($existing) {
        Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
        sc.exe delete $ServiceName | Out-Null
        Start-Sleep -Seconds 2
    }

    New-Service `
        -Name $ServiceName `
        -BinaryPathName "`"$binary`"" `
        -DisplayName "wtv IPTV Monitor" `
        -Description "monitor-lite-api IPTV 直播源监控服务" `
        -StartupType Automatic | Out-Null

    Start-Service -Name $ServiceName
}

function Install-Wtv {
    Ensure-Admin
    Ensure-Ffmpeg

    $arch = Get-Arch
    $dest = Join-Path $InstallPath "wtv.exe"
    if (Test-Path $dest) {
        Write-Warn "检测到已安装在 $InstallPath，请使用 -Update 更新"
        return
    }

    New-Item -ItemType Directory -Force -Path $InstallPath | Out-Null
    $url = Get-DownloadUrl -Arch $arch
    Write-Info "下载 $BinaryPrefix-windows-$arch.exe ..."
    Invoke-WebRequest -Uri $url -OutFile $dest -UseBasicParsing
    Install-Service

    Write-Info "安装完成"
    Write-Host ""
    Write-Host "管理地址: http://localhost:9876/admin"
    Write-Host "配置文件: $InstallPath\etc\tv.yaml"
    Write-Host "默认账号: admin / admin123（请尽快修改）"
    Write-Host ""
    Write-Host "常用命令:"
    Write-Host "  Get-Service $ServiceName"
    Write-Host "  Restart-Service $ServiceName"
}

function Update-Wtv {
    Ensure-Admin

    $dest = Join-Path $InstallPath "wtv.exe"
    if (-not (Test-Path $dest)) {
        throw "未检测到已安装的 wtv，请先执行 -Install"
    }

    $arch = Get-Arch
    $url = Get-DownloadUrl -Arch $arch
    $backup = Join-Path $env:TEMP "wtv.bak.exe"
    $newFile = Join-Path $env:TEMP "wtv.new.exe"

    Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
    Copy-Item $dest $backup -Force

    try {
        Write-Info "下载更新包..."
        Invoke-WebRequest -Uri $url -OutFile $newFile -UseBasicParsing
        Move-Item $newFile $dest -Force
    }
    catch {
        Move-Item $backup $dest -Force
        Start-Service -Name $ServiceName -ErrorAction SilentlyContinue
        throw
    }

    Remove-Item $backup -Force -ErrorAction SilentlyContinue
    Start-Service -Name $ServiceName
    Write-Info "已更新到最新版本"
}

function Uninstall-Wtv {
    Ensure-Admin

    Stop-Service -Name $ServiceName -Force -ErrorAction SilentlyContinue
    sc.exe delete $ServiceName | Out-Null
    Remove-Item -Recurse -Force $InstallPath -ErrorAction SilentlyContinue
    Write-Info "wtv 已卸载"
}

if ($Install) { Install-Wtv; return }
if ($Update) { Update-Wtv; return }
if ($Uninstall) { Uninstall-Wtv; return }

# 通过 irm | iex 调用时默认安装
Install-Wtv
