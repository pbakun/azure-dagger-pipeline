
$ResourceGroupName = $env:ResourceGroupName
$AppServiceName = $env:WebAppName

Write-Host "Stop staging slot with production webapp"

Stop-AzWebAppSlot -ResourceGroupName $ResourceGroupName `
                  -Name $AppServiceName `
                  -Slot staging