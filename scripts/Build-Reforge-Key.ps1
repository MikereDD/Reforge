[CmdletBinding()]
param(
    [ValidateSet("amd64", "arm64")]
    [string]$Architecture = "amd64",

    [switch]$Clean
)

$ErrorActionPreference = "Stop"

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$bundleName = "reforge-key-linux-$Architecture"
$bundleRoot = Join-Path $repoRoot "dist\$bundleName"

Push-Location $repoRoot

$previousGoos = $env:GOOS
$previousGoarch = $env:GOARCH
$previousCgo = $env:CGO_ENABLED

try {
    if ($Clean -and (Test-Path $bundleRoot)) {
        Remove-Item $bundleRoot -Recurse -Force
    }

    New-Item -ItemType Directory -Force `
        (Join-Path $bundleRoot "bin") |
        Out-Null

    New-Item -ItemType Directory -Force `
        (Join-Path $bundleRoot "manifests") |
        Out-Null

    New-Item -ItemType Directory -Force `
        (Join-Path $bundleRoot "scripts") |
        Out-Null

    $env:GOOS = "linux"
    $env:GOARCH = $Architecture
    $env:CGO_ENABLED = "0"

    $binaryPath = Join-Path $bundleRoot "bin\reforge"

    go build `
        -trimpath `
        -ldflags "-s -w" `
        -o $binaryPath `
        ./cmd/reforge

    if ($LASTEXITCODE -ne 0) {
        throw "Go build failed with exit code $LASTEXITCODE"
    }

    Copy-Item `
        "manifests\catalog.example.json" `
        (Join-Path $bundleRoot "manifests\catalog.example.json") `
        -Force

    Copy-Item `
        "scripts\Prepare-Reforge-Key.sh" `
        (Join-Path $bundleRoot "scripts\prepare.sh") `
        -Force

    Copy-Item `
        "scripts\Run-Reforge-Key.sh" `
        (Join-Path $bundleRoot "run.sh") `
        -Force

    $hash = (Get-FileHash -Algorithm SHA256 $binaryPath).Hash.ToLowerInvariant()
    "$hash  bin/reforge" |
        Set-Content `
            (Join-Path $bundleRoot "SHA256SUMS") `
            -Encoding ascii

    Write-Host
    Write-Host "Reforge Key bundle created:"
    Write-Host "  $bundleRoot"
    Write-Host
    Write-Host "Target:"
    Write-Host "  linux/$Architecture"
    Write-Host
    Write-Host "Binary SHA-256:"
    Write-Host "  $hash"
}
finally {
    $env:GOOS = $previousGoos
    $env:GOARCH = $previousGoarch
    $env:CGO_ENABLED = $previousCgo

    Pop-Location
}
