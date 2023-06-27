package secrets

import (
	"strings"

	"golang.org/x/exp/slices"
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
	candidates := getCandidateIndexes(args, "-s")

	for _, candidate := range candidates {
		newSecret := NewSecretFromVariable(args[candidate])
		if newSecret != nil {
			secrets = append(secrets, *newSecret)
		}
	}
	return secrets
}

func GetSecretByName(secrets []Secret, name string) string {
	index := slices.IndexFunc(secrets, func(s Secret) bool { return strings.EqualFold(s.Name, name) })
	if index != -1 {
		return secrets[index].Secret
	}
	return ""
}

func GetArgByCode(args []string, argCode string) string {
	candidates := getCandidateIndexes(args, argCode)
	if len(candidates) > 0 {
		return args[candidates[0]]
	}
	return ""
}

func getCandidateIndexes(args []string, argCode string) []int {

	indexes := []int{}
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
