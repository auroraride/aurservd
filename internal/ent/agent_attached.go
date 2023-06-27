package ent

import "github.com/auroraride/adapter"

func (a *Agent) Meta() string {
	return a.Name + " - " + a.Phone
}

func (a *Agent) AdapterUser() *adapter.User {
	return &adapter.User{
		Type: adapter.UserTypeAgent,
		ID:   a.Phone,
	}
}
