package models

type Node struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type RepoPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	CID     string `json:"cid"`
}
