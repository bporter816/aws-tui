package view

type SecretsManager struct {
}

func (s SecretsManager) GetService() string {
	return "Secrets Manager"
}
