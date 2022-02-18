package github

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"

	"github.com/StevenACoffman/commentary/pkg/generated/genqlient"
	"github.com/StevenACoffman/commentary/pkg/types"
)

func GetPullRequestAndCommentsForCommit(ctx context.Context, graphqlClient graphql.Client, sha, repo, org string) (types.PullRequest, []types.CommentNodes, error) {
	resp, err := genqlient.GetPullRequestAndCommentsForCommit(ctx, graphqlClient, org, repo, sha)
	fmt.Printf("%+v", resp)
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

func GetPullRequestByBranch(ctx context.Context, graphqlClient graphql.Client, owner, repo, headref, baseref string) (types.PullRequest, []types.CommentNodes, error) {
	resp, err := genqlient.GetPullRequestForBranch(ctx, graphqlClient, owner, repo, headref, baseref)

	if err != nil {
		return types.PullRequest{}, nil, err
	}
	for _, node := range resp.Repository.PullRequests.Nodes {
		pr := types.PullRequest{
			Number: node.Number,
			ID:     node.Id,
		}
		var comments []types.CommentNodes
		for _, comment := range node.Comments.Nodes {
			comments = append(comments, types.CommentNodes{
				ID:  comment.Id,
				URL: comment.Url,
				Author: types.Author{
					Login: comment.Author.GetLogin(),
				},
				Body: comment.Body,
			})
		}
		return pr, comments, nil
	}
	return types.PullRequest{}, nil, nil
}

func GetPullRequestByURI(ctx context.Context, graphqlClient graphql.Client, uri string) (types.PullRequest, []types.CommentNodes, error) {
	resp, err := genqlient.GetCommentsForPullRequest(ctx, graphqlClient, uri)
	fmt.Println("TYPENAME", resp.Resource.GetTypename())
	if err != nil {
		return types.PullRequest{}, nil, err
	}

	b, err := resp.MarshalJSON()
	if err != nil {
		return types.PullRequest{}, nil, err
	}
	fmt.Printf("STARTING: %+v\n\n", string(b))
	//if err != nil {
	//	return types.PullRequest{}, nil, err
	//}
	//for _, node := range resp.Resource.GetTypename() resp.Repository.Commit.AssociatedPullRequests.PRNodes {
	//	pr := types.PullRequest{
	//		Number: node.Number,
	//		ID:     node.ID,
	//	}
	//	return pr, node.Comments.Nodes, nil
	//}

	return types.PullRequest{}, nil, nil
}
