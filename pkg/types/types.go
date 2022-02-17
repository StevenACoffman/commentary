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
	ID     string `json:"id"`
	URL    string `json:"url"`
	Author Author `json:"author"`
	Body   string `json:"body"`
}

type PageInfo struct {
	HasNextPage bool `json:"hasNextPage"`
}

type Comments struct {
	Nodes    []CommentNodes `json:"nodes"`
	PageInfo PageInfo       `json:"pageInfo"`
}

type PRNodes struct {
	Number   int      `json:"number"`
	ID       string   `json:"id"`
	Comments Comments `json:"comments"`
}

type AssociatedPullRequests struct {
	PRNodes []PRNodes `json:"nodes"`
}

type PullRequest struct {
	Number int    `json:"number"`
	ID     string `json:"id"`
}
