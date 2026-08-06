# Validate a draft map issue body before gh issue edit --body-file.
# Usage: .\validate-map-body.ps1 path\to\map-body.md
param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$Path
)

$ErrorActionPreference = "Stop"

if (-not (Test-Path $Path)) {
    Write-Error "FAIL: file not found: $Path"
    exit 1
}

$content = [System.IO.File]::ReadAllText($Path, [System.Text.UTF8Encoding]::new($false))
$lines = ($content -split "`r?`n").Count

if ($lines -lt 40) {
    Write-Error "FAIL: body has $lines lines (expected >= 40; collapsed maps are often ~10 lines)"
    exit 1
}

$requiredSections = @(
    "## To Do",
    "## Completed",
    "## Decision coverage"
)

foreach ($section in $requiredSections) {
    $pattern = "(?m)^$([regex]::Escape($section))\s*$"
    if ($content -notmatch $pattern) {
        Write-Error "FAIL: missing section header on its own line: $section"
        exit 1
    }
}

if ($content -match "[\u2013\u2014\u00b7]") {
    Write-Error "FAIL: unicode em dash, en dash, or middle dot found - use ASCII hyphen only"
    exit 1
}

# Common UTF-8 misread as Latin-1/Windows-1252 (mojibake)
if ($content -match '\u00c3|\u0393|\u256c|\u251c|\u2524|\u252c|\u2510') {
    Write-Error "FAIL: mojibake sequences detected (encoding corruption)"
    exit 1
}

Write-Host ('OK: map body valid ({0} lines, required sections present, ASCII punctuation)' -f $lines)
exit 0
