package main

type Entry struct {
	Site     string `json:"site"`
	Username string `json:"username"`
	Password string `json:"password"`
	Note     string `json:"note,omitempty"`
}

type Vault struct {
	Entries []Entry `json:"entries"`
}
