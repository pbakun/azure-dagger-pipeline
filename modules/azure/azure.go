package azure

import (
	"reflect"

	"dagger.io/dagger"
)

type AzureServicePrincipal struct {
	TenantId     string
	ClientId     string
	ClientSecret string
}

type DeploymentVariables struct {
	ResourceGroupName string
	WebAppName        string
	Location          string
	ArtifactName      string
}

func DeployWebApp(servicePrincipal AzureServicePrincipal, resourceGroupName string) error {

	// cred, err := azidentity.NewClientSecretCredential(servicePrincipal.TenantId,
	// 	servicePrincipal.ClientId,
	// 	servicePrincipal.ClientSecret,
	// 	&azidentity.ClientSecretCredentialOptions{})
	// if err != nil {
	// 	return err
	// }

	// appServiceClientFactory, err := armappservice.NewClientFactory(servicePrincipal.SubscriptionId, cred, nil)
	// if err != nil {
	// 	return nil
	// }

	// plansClient := appServiceClientFactory.NewPlansClient()
	// webAppsClient := appServiceClientFactory.NewWebAppsClient()
	// appServiceClientFactory.NewStaticSitesClient().
	// 	armappservice.StaticSiteZipDeploymentARMResource{}
	// webAppsClient.Begin

	// fmt.Print("Deploy to Azure")
	return nil
}

func GetAzPwsh(client dagger.Client, servicePrincipal AzureServicePrincipal) *dagger.Container {

	scriptsDir := client.Host().Directory("./scripts/")

	container := client.Container().
		From("mcr.microsoft.com/azure-powershell").
		WithMountedDirectory("/scripts", scriptsDir).
		WithWorkdir("/scripts")

	// use reflection to loop through properties of service principal and set them as secrets in container
	servicePrincipalReflect := reflect.ValueOf(&servicePrincipal).Elem()
	for i := 0; i < servicePrincipalReflect.NumField(); i++ {
		field := servicePrincipalReflect.Field(i)
		name := servicePrincipalReflect.Type().Field(i).Name
		secret := client.SetSecret(name, field.String())
		container = container.WithSecretVariable(name, secret)
	}

	return container.
		WithExec([]string{"pwsh", "utilities/AzLogin.ps1"})
}
