$ErrorActionPreference = 'Stop'; # stop on all errors
$toolsDir   = "$(Split-Path -parent $MyInvocation.MyCommand.Definition)"
$url64      = "https://github.com/kubeshop/tracetest/releases/download/v${env:chocolateyPackageVersion}/tracetest_${env:chocolateyPackageVersion}_windows_amd64.tar.gz" # 64bit URL here (HTTPS preferred) or remove - if installer contains both (very rare), use $url

$packageArgs = @{
  packageName   = $env:ChocolateyPackageName
  unzipLocation = $toolsDir
  fileType      = 'exe'
  url64bit      = $url64

  softwareName  = 'tracetest*'

  # %placeholder% will be replaced by the correct value in the CI pipeline
  checksum      = '%checksum%'
  checksumType  = 'sha256'
  checksum64    = '%checksum%'
  checksumType64= 'sha256'
}

Install-ChocolateyZipPackage @packageArgs
Get-ChocolateyUnzip -FileFullPath "$toolsDir/tracetest_${env:chocolateyPackageVersion}_windows_amd64.tar" -Destination $toolsDir