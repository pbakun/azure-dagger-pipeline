package secrets

import (
	"strings"
)

type Secret struct {
	Name   string
	Secret string
}

func NewSecret(name string, secret string) *Secret {

	return &Secret{
		Name:   name,
		Secret: secret}
}

func NewSecretFromVariable(variable string) *Secret {

	firstEqualSign := strings.Index(variable, "=")
	if firstEqualSign == -1 {
		return nil
	}
	first := variable[0:firstEqualSign]
	second := variable[firstEqualSign+1:]

	return &Secret{
		Name:   first,
		Secret: second}
}

func GetSecrets(args []string) []Secret {

	secrets := []Secret{}
	candidates := getCandidateIndexes(args)

	for _, candidate := range candidates {
		newSecret := NewSecretFromVariable(args[candidate])
		if newSecret != nil {
			secrets = append(secrets, *newSecret)
		}
	}
	return secrets
}

func getCandidateIndexes(args []string) []int {

	indexes := []int{}
	argCode := "-s"
	for i, arg := range args {
		if len(args) == i+1 {
			continue
		}
		if strings.EqualFold(arg, argCode) {
			indexes = append(indexes, i+1)
		}
	}
	return indexes
}
