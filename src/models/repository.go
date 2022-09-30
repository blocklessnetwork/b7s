package models

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	CID     string `json:"cid"`
}

type Repo struct {
	version      string   `json:"version"`
	id           string   `json:"id"`
	pageType     string   `json:"type"`
	pages        []string `json:"pages"`
	lastModified string   `json:"lastModified"`
}

type RepoPage struct {
	version      string    `json:"version"`
	id           string    `json:"id"`
	pageType     string    `json:"type"`
	packages     []Package `json:"packages"`
	lastModified string    `json:"lastModified"`
}
