package models

type Officer struct {
	User
}

func (o *Officer) getRole() string {
	return "officer"
}
