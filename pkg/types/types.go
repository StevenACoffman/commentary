package types

type GitObject struct {
	Typename               string                 `json:"__typename"`
	AssociatedPullRequests AssociatedPullRequests `json:"associatedPullRequests"`
}

type Author struct {
	Typename string `json:"__typename"`
	Login    string `json:"login"`
}

type CommentNodes struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
	URL      string `json:"url"`
	Author   Author `json:"author"`
	Body     string `json:"body"`
}

type PageInfo struct {
	HasNextPage bool `json:"hasNextPage"`
}

type Comments struct {
	Typename string         `json:"__typename"`
	Nodes    []CommentNodes `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

type PullRequest struct {
	Typename     string   `json:"__typename"`
	Number       int      `json:"number"`
	ID           string   `json:"id"`
	Comments     Comments `json:"comments"`
	ResourcePath string   `json:"resourcePath"`
	URL          string   `json:"url"`
}

type AssociatedPullRequests struct {
	PRNodes []PullRequest `json:"nodes"`
}
