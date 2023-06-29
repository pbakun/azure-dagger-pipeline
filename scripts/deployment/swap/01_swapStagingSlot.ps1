
$ResourceGroupName = $env:ResourceGroupName
$AppServiceName = $env:WebAppName

Write-Host "Swapping staging slot with production webapp"
Switch-AzWebAppSlot -Name $AppServiceName `
                    -ResourceGroupName $ResourceGroupName `
                    -SourceSlotName staging `
                    -DestinationSlotName production