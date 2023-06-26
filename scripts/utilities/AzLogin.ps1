Write-Host "Logging to Azure with Service Principal"

$appId = $env:appId
$principalPass = $env:principalPass
$tenant = $env:tenant

$SecureStringPwd = $principalPass | ConvertTo-SecureString -AsPlainText -Force
$pscredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $appId, $SecureStringPwd
Connect-AzAccount -ServicePrincipal -Credential $pscredential -Tenant $tenant