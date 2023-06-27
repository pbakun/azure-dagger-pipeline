package azure

import (
	"fmt"
)

type AzureServicePrincipal struct {
	TenantId       string
	ClientId       string
	ClientSecret   string
	SubscriptionId string
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

	fmt.Print("Deploy to Azure")
	return nil
}
