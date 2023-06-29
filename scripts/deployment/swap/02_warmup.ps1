$ResourceGroupName = $env:ResourceGroupName
$AppServiceName = $env:WebAppName

Write-Host "Warm up production slot"

$appServiceSlot = Get-AzWebAppSlot -Name $AppServiceName `
                               -ResourceGroupName $ResourceGroupName `
                               -Slot production

Invoke-WebRequest -Uri $appServiceSlot.DefaultHostName `
                  -TimeoutSec 20 `
                  -MaximumRetryCount 100
