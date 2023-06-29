$ResourceGroupName = $env:ResourceGroupName
$AppServiceName = $env:WebAppName

Write-Host "Warm up staging slot"

$appServiceSlot = Get-AzWebAppSlot -Name $AppServiceName `
                               -ResourceGroupName $ResourceGroupName `
                               -Slot staging

Invoke-WebRequest -Uri $appServiceSlot.DefaultHostName `
                  -TimeoutSec 20 `
                  -MaximumRetryCount 100
