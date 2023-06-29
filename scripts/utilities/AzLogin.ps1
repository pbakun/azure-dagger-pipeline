Write-Host "Logging to Azure with Service Principal"

$appId = $env:ClientId
$principalPass = $env:ClientSecret
$tenant = $env:TenantId

$SecureStringPwd = $principalPass | ConvertTo-SecureString -AsPlainText -Force
$pscredential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList $appId, $SecureStringPwd
Connect-AzAccount -ServicePrincipal -Credential $pscredential -Tenant $tenant