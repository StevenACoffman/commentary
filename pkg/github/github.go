package github

import (
	"context"
	"github.com/Khan/genqlient/graphql"
	"github.com/StevenACoffman/commentary/pkg/generated/genqlient"
	"github.com/StevenACoffman/commentary/pkg/types"
)

func GetPullRequestAndCommentsForCommit(ctx context.Context, graphqlClient graphql.Client, sha, repo, org string) (types.PullRequest, []types.CommentNodes, error) {
	resp, err := genqlient.GetPullRequestAndCommentsForCommit(ctx, graphqlClient, org, repo, sha)
	if err != nil {
		return types.PullRequest{}, nil, err
	}
	for _, node := range resp.Repository.Commit.AssociatedPullRequests.PRNodes {
		pr := types.PullRequest{
			Number: node.Number,
			ID:     node.ID,
		}
		return pr, node.Comments.Nodes, nil
	}

	return types.PullRequest{}, nil, nil
}

func UpdateComment(ctx context.Context, graphqlClient graphql.Client,
	commentId string,
	body string) (string, error) {
	resp, err := genqlient.UpdatePullRequestComment(ctx, graphqlClient, commentId, body)

	return resp.UpdateIssueComment.IssueComment.Id, err
}

func CreateNewPullRequestComment(ctx context.Context, graphqlClient graphql.Client,
	prId string,
	body string) (string, error) {
	resp, err := genqlient.CreateNewPullRequestComment(ctx, graphqlClient, prId, body)

	return resp.AddComment.CommentEdge.Node.Id, err
}
