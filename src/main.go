package main

import (
	"bakson/dagger/azure/modules/secrets"
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"dagger.io/dagger"
)

func main() {
	client, err := dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stdout))
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	// secretsArr := secrets.GetSecrets(os.Args)

	// err = deployEnv(context.Background(), *client, secretsArr)
	os.RemoveAll("./build")
	err = build(context.Background(), *client)
	if err != nil {
		fmt.Println(err)
	}
}

func deployAzure(ctx context.Context) error {
	_, err := azidentity.NewClientSecretCredential("", "", "", &azidentity.ClientSecretCredentialOptions{})
	if err != nil {
		return err
	}
	return err
}

func deployEnv(ctx context.Context, client dagger.Client, secrets []secrets.Secret) error {

	dir := client.Host().Directory("./scripts")

	pwsh := GetAzPwsh(client, dir, secrets).
		WithExec([]string{"pwsh", "deployment/step1.ps1"}).
		WithExec([]string{"pwsh", "deployment/step2.ps1"})

	_, err := pwsh.ExitCode(ctx)

	if err != nil {
		return err
	}
	return nil
}

func GetAzPwsh(c dagger.Client, dir *dagger.Directory, secrets []secrets.Secret) *dagger.Container {

	container := c.Container().
		From("mcr.microsoft.com/azure-powershell").
		WithMountedDirectory("/deployScripts", dir).
		WithWorkdir("/deployScripts")

	for _, secret := range secrets {
		sec := c.SetSecret(secret.Name, secret.Secret)
		container = container.WithSecretVariable(secret.Name, sec)
	}

	return container.
		// WithExec([]string{"pwsh", "utilities/InstallAzModules.ps1"}).
		WithExec([]string{"pwsh", "utilities/AzLogin.ps1"})
}

func build(ctx context.Context, client dagger.Client) error {

	dir := client.Host().Directory("./app")

	output := client.Directory()

	container := client.Container().
		From("mcr.microsoft.com/dotnet/sdk:7.0").
		WithMountedDirectory("/app", dir).
		WithWorkdir("/app/WebApp")

	container = container.
		WithExec([]string{"dotnet", "build", "-c", "Release"}).
		WithExec([]string{"dotnet", "publish", "-c", "Release", "-o", "/build"})

	output = output.WithDirectory(".", container.Directory("/build"))

	_, err := output.Export(ctx, "./build")
	if err != nil {
		return err
	}
	return nil
}
