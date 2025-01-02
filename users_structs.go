package main

type UserCreateStruct struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Authenticated struct {
	Authenticated string `json:"authenticated"`
}
