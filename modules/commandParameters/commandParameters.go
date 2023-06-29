package commandParameters

import (
	"strings"

	"golang.org/x/exp/slices"
)

type Secret struct {
	Name   string
	Secret string
}

type Variable struct {
	Name  string
	Value string
}

func NewSecret(name string, secret string) *Secret {
	return &Secret{
		Name:   name,
		Secret: secret}
}

func NewVariable(name string, value string) *Variable {
	return &Variable{
		Name:  name,
		Value: value}
}

func GetSecrets(args []string) []Secret {

	secrets := []Secret{}
	candidates := getCandidateIndexes(args, "-s")

	for _, candidate := range candidates {
		secretName, secretValue := newNameValuePair(args[candidate])
		if secretName != "" && secretValue != "" {
			secrets = append(secrets, *NewSecret(secretName, secretValue))
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

func GetVariables(args []string) []Variable {
	variables := []Variable{}
	candidates := getCandidateIndexes(args, "-v")

	for _, candidate := range candidates {
		name, value := newNameValuePair(args[candidate])
		if name != "" && value != "" {
			variables = append(variables, *NewVariable(name, value))
		}
	}
	return variables
}

func GetVariableByName(variables []Variable, name string) string {
	index := slices.IndexFunc(variables, func(v Variable) bool { return strings.EqualFold(v.Name, name) })
	if index != -1 {
		return variables[index].Value
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

func newNameValuePair(variable string) (string, string) {

	firstEqualSign := strings.Index(variable, "=")
	if firstEqualSign == -1 {
		return "", ""
	}
	first := variable[0:firstEqualSign]
	second := variable[firstEqualSign+1:]

	return first, second
}
