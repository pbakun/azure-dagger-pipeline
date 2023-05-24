
## About

This app is a template for Azure deployment with dagger.

It pulls Azure Powershell docker image, sets Azure Context by given Service Principal, executes scripts in that environment

#### How to start

```pwsh
go run main.go -s appId=<application-id> -s tenant=<azure-tenant> -s principalPass=<service-principal-password>
```

