package main

import (
	"bakson/dagger/azure/modules/secrets"
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	client, err := dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stdout))
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	secretsArr := secrets.GetSecrets(os.Args)

	err = deployEnv(context.Background(), *client, secretsArr)

	if err != nil {
		fmt.Println(err)
	}
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
