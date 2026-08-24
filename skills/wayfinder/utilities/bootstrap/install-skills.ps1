# Legacy global install - copies wayfinder pack to ~/.cursor/skills/
# Prefer agent-config-wizard for desktop project apply (see repo README).
# Cloud AFK: wrapper around install-skills.sh for local smoke tests.
param(
    [string]$Tag = "v1.0.0",
    [string]$Repo = "KroniK907/agent-config"
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ShScript = Join-Path $ScriptDir "install-skills.sh"

$env:WAYFINDER_SKILLS_TAG = $Tag
$env:WAYFINDER_SKILLS_REPO = $Repo

if (Get-Command bash -ErrorAction SilentlyContinue) {
    Push-Location $ScriptDir
    $prevEap = $ErrorActionPreference
    $ErrorActionPreference = 'Continue'
    try {
        bash ./install-skills.sh 2>&1 | ForEach-Object { Write-Host $_ }
        exit $LASTEXITCODE
    } finally {
        $ErrorActionPreference = $prevEap
        Pop-Location
    }
}

# Fallback when bash unavailable (Git for Windows usually provides it)
$Dest = Join-Path $env:USERPROFILE ".cursor\skills"
$CacheName = ($Repo -replace '/', '-')
$Cache = Join-Path $env:USERPROFILE ".cache\wayfinder-skills\$CacheName"

New-Item -ItemType Directory -Force -Path $Dest, $Cache | Out-Null

$cloneUrl = "https://github.com/$Repo.git"
if ($env:GH_TOKEN) {
    $cloneUrl = "https://x-access-token:$($env:GH_TOKEN)@github.com/$Repo.git"
}

if (-not (Test-Path (Join-Path $Cache ".git"))) {
    git clone --depth 1 --branch $Tag $cloneUrl $Cache
} else {
    git -C $Cache fetch origin tag $Tag --force 2>$null
    if ($LASTEXITCODE -ne 0) { git -C $Cache fetch origin --tags --force }
    git -C $Cache checkout $Tag
}

function Sync-Dir {
    param([string]$Src, [string]$Name)
    if (Test-Path $Src) {
        $target = Join-Path $Dest $Name
        if (Test-Path $target) { Remove-Item $target -Recurse -Force }
        Copy-Item $Src $target -Recurse -Force
        Write-Host "  synced $Name"
    }
}

Write-Host "Installing skills from ${Repo}@${Tag} -> $Dest"
Sync-Dir (Join-Path $Cache "wayfinder") "wayfinder"
@(
    "tdd", "commit", "writing-for-agents", "write-a-prd", "prd-to-plan",
    "prd-to-issues", "request-refactor-plan", "triage-issue",
    "improve-codebase-architecture", "ubiquitous-language", "write-a-skill"
) | ForEach-Object { Sync-Dir (Join-Path $Cache $_) $_ }

Write-Host "Skills install complete."
