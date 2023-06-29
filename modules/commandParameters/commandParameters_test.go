package commandParameters

import (
	"testing"
)

func TestFindSecrets(t *testing.T) {
	arr := []string{"apple", "-s", "cherry"}
	secrets := GetSecrets(arr)
	if len(secrets) != 0 {
		t.Errorf("The array should be empty but has %d elements", len(secrets))
	}

	arr = []string{"test", "-s", "name1=pass1"}
	secrets = GetSecrets(arr)
	if len(secrets) != 1 {
		t.Errorf("It should find only one secret and has %d", len(secrets))
	}

	arr = []string{"test", "-s", "name1=pass1", "-s"}
	secrets = GetSecrets(arr)
	if len(secrets) != 1 {
		t.Errorf("It should find only one secret and has %d", len(secrets))
	}

	arr = []string{"test", "-s", "name1=pass1="}
	secret := GetSecrets(arr)[0]
	if !(secret.Name == "name1" && secret.Secret == "pass1=") {
		t.Logf("\n Secret name %s \n Secret value %s", secret.Name, secret.Secret)
		t.Fail()
		// t.Errorf("Secret doesn't have correct values")
	}
}

func TestGetSecretByName(t *testing.T) {
	secretsArr := []Secret{}
	expectedSecretValue := "123456"
	secretsArr = append(secretsArr,
		Secret{Name: "appId", Secret: "test123456"},
		Secret{Name: "tenantId", Secret: expectedSecretValue})

	secret := GetSecretByName(secretsArr, "tenantId")
	if secret != expectedSecretValue {
		t.Errorf("The expected secret should be %s, but was %s", expectedSecretValue, secret)
	}

	secret = GetSecretByName(secretsArr, "TenantId")
	if secret != expectedSecretValue {
		t.Errorf("The expected secret should be %s, but was %s", expectedSecretValue, secret)
	}
}
