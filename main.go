package main

import (
	"bakson/pipelines/dagger/modules/azure"
	"bakson/pipelines/dagger/modules/commandParameters"
	"context"
	"fmt"
	"io/ioutil"
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

	pipelineMode := commandParameters.GetArgByCode(os.Args, "-m")
	if strings.EqualFold("ci", pipelineMode) {
		outputDir := commandParameters.GetArgByCode(os.Args, "-o")
		if outputDir == "" {
			fmt.Println("Missing build output directory. Pass it using -o command argument like: go run main.go -o <path>")
			return
		}

		os.RemoveAll(outputDir)
		err = build(context.Background(), *client, outputDir)

	} else if strings.EqualFold("cd", pipelineMode) {
		secretsArr := commandParameters.GetSecrets(os.Args)
		azureServicePrincipal := azure.AzureServicePrincipal{
			TenantId:     commandParameters.GetSecretByName(secretsArr, "tenantId"),
			ClientId:     commandParameters.GetSecretByName(secretsArr, "appId"),
			ClientSecret: commandParameters.GetSecretByName(secretsArr, "password"),
		}
		variablesArr := commandParameters.GetVariables(os.Args)
		deploymentVariables := azure.DeploymentVariables{
			ResourceGroupName: commandParameters.GetVariableByName(variablesArr, "resourceGroupName"),
			WebAppName:        commandParameters.GetVariableByName(variablesArr, "webAppName"),
			Location:          commandParameters.GetVariableByName(variablesArr, "location"),
			ArtifactName:      commandParameters.GetVariableByName(variablesArr, "ArtifactName"),
		}
		artifactPath := commandParameters.GetArgByCode(os.Args, "-a")
		stepsFolder := commandParameters.GetArgByCode(os.Args, "-f")
		// azure.DeployWebApp(azureServicePrincipal, "dagger-webapp")
		err = deploy(context.Background(), *client, artifactPath, stepsFolder, azureServicePrincipal, deploymentVariables)
	}

	if err != nil {
		fmt.Println(err)
	}
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

func deploy(ctx context.Context,
	client dagger.Client,
	artifactPath string,
	stepsFolder string,
	azureServicePrincipal azure.AzureServicePrincipal,
	deploymentVariable azure.DeploymentVariables) error {

	container := azure.GetAzPwsh(client, azureServicePrincipal)

	publishDir := client.Host().Directory(artifactPath)
	container = container.WithMountedDirectory("/publish", publishDir).
		WithEnvVariable("ResourceGroupName", deploymentVariable.ResourceGroupName).
		WithEnvVariable("WebAppName", deploymentVariable.WebAppName).
		WithEnvVariable("Location", deploymentVariable.Location).
		WithEnvVariable("ArtifactName", deploymentVariable.ArtifactName)

	stepsPath := fmt.Sprintf("./scripts/deployment/%s", stepsFolder)
	steps, err := ioutil.ReadDir(stepsPath)
	if err != nil {
		return err
	}

	for _, file := range steps {
		stepPath := fmt.Sprintf("deployment/%s/%s", stepsFolder, file.Name())
		container = container.WithExec([]string{"pwsh", stepPath})
	}
	_, err = container.ExitCode(ctx)

	if err != nil {
		return err
	}
	return nil
}
