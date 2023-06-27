package main

import (
	"bakson/pipelines/dagger/modules/azure"
	"bakson/pipelines/dagger/modules/secrets"
	"context"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
)

func main() {
	client, err := dagger.Connect(context.Background(), dagger.WithLogOutput(os.Stdout))
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	pipelineMode := secrets.GetArgByCode(os.Args, "-m")
	if strings.EqualFold("ci", pipelineMode) {
		outputDir := secrets.GetArgByCode(os.Args, "-o")
		if outputDir == "" {
			fmt.Println("Missing build output directory. Pass it using -o command argument like: go run main.go -o <path>")
			return
		}

		os.RemoveAll(outputDir)
		err = build(context.Background(), *client, outputDir)
	} else if strings.EqualFold("cd", pipelineMode) {
		azure.DeployWebApp()
	}
	// secretsArr := secrets.GetSecrets(os.Args)

	// err = deployEnv(context.Background(), *client, secretsArr)

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

func build(ctx context.Context, client dagger.Client, outputDirectory string) error {

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

	_, err := output.Export(ctx, outputDirectory)
	if err != nil {
		return err
	}
	return nil
}
