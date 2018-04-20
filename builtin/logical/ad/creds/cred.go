package creds

type credential struct {
	Username        string `json:"username"`
	CurrentPassword string `json:"current_password"`
	LastPassword    string `json:"last_password,omitempty"`
}

func (c *credential) Map() map[string]interface{} {
	m := map[string]interface{}{
		"username":         c.Username,
		"current_password": c.CurrentPassword,
	}
	if c.LastPassword != "" {
		m["last_password"] = c.LastPassword
	}
	return m
}
