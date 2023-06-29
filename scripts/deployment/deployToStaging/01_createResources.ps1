
$ResourceGroupName = $env:ResourceGroupName
$Location = $env:Location
$AppServiceName = $env:WebAppName
$ArtifactName = $env:ArtifactName

$resourceGroup = Get-AzResourceGroup -Name $ResourceGroupName
if ($resourceGroup -eq $null) {
    Write-Host "Resource group doesn't exist. Creating..."
    New-AzResourceGroup -Name $ResourceGroupName -Location $Location
}

$appServicePlan = Get-AzAppServicePlan -Name $AppServiceName -ResourceGroupName $ResourceGroupName
if($appServicePlan -eq $null) {
    Write-Host "App Service Plan doesn't exist. Creating..."
    New-AzAppServicePlan -Name $AppServiceName -Location $Location -ResourceGroupName $ResourceGroupName -Tier S1
}

$appService = Get-AzWebApp -Name $AppServiceName -ResourceGroupName $ResourceGroupName
if($appService -eq $null) {
    Write-Host "App Service doesn't exist. Creating..."
    New-AzWebApp -Name $AppServiceName -Location $Location -AppServicePlan $AppServiceName -ResourceGroupName $ResourceGroupName
}

