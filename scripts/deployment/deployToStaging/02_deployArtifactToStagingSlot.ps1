$ResourceGroupName = $env:ResourceGroupName
$AppServiceName = $env:WebAppName
$ArtifactName = $env:ArtifactName

New-AzWebAppSlot -Name $AppServiceName `
                 -ResourceGroupName $ResourceGroupName `
                 -Slot staging

Write-Host "Publishing $ArtifactName"

Publish-AzWebApp -ResourceGroupName $ResourceGroupName `
                 -Name $AppServiceName `
                 -ArchivePath "/publish/$ArtifactName" `
                 -Slot staging `
                 -Clean `
                 -Force

Start-AzWebAppSlot -ResourceGroupName $ResourceGroupName `
                   -Name $AppServiceName `
                   -Slot staging
