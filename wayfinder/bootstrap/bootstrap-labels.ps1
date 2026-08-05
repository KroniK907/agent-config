# Bootstrap wayfinder:* GitHub labels from labels-manifest.json.
# Idempotent: gh label create --force updates existing labels.
# Usage: .\bootstrap-labels.ps1 [-Repo owner/repo] [-Manifest path\to\labels-manifest.json]
param(
    [string]$Repo = "",
    [string]$Manifest = ""
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
if (-not $Manifest) { $Manifest = Join-Path $ScriptDir "labels-manifest.json" }
if (-not $Repo) {
    $Repo = (gh repo view --json nameWithOwner -q .nameWithOwner)
}

$labels = Get-Content $Manifest -Raw | ConvertFrom-Json
Write-Host "Bootstrapping $($labels.Count) wayfinder labels on $Repo ..."

foreach ($label in $labels) {
    $desc = if ($label.description) { $label.description } else { "" }
    gh label create $label.name `
        --repo $Repo `
        --color $label.color `
        --description $desc `
        --force 2>$null
    if ($LASTEXITCODE -ne 0) {
        Write-Warning "Could not create or update label: $($label.name)"
    } else {
        Write-Host "  OK $($label.name)"
    }
}

Write-Host "Done."
