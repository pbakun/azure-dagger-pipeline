package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

type AzureServicePrincipal struct {
	TenantId     string
	AppId        string
	ClientSecret string
}

func DeployWebApp() {
	// _, err := azidentity.NewClientSecretCredential("", "", "", &azidentity.ClientSecretCredentialOptions{})
	fmt.Println("deploy webapp")
}

func test() error {

	_, err := azidentity.NewClientSecretCredential("", "", "", &azidentity.ClientSecretCredentialOptions{})
	if err != nil {
		return err
	}
	return nil
}
