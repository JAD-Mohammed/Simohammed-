package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/github/github-mcp-server/pkg/translations"
	"github.com/google/go-github/v69/github"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// getPullRequest creates a tool to get details of a specific pull request.
func getPullRequest(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pull_request",
			mcp.WithDescription(t("TOOL_GET_PULL_REQUEST_DESCRIPTION", "Get details of a specific pull request")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("pull_number",
				mcp.Required(),
				mcp.Description("Pull request number"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			pullNumber := int(request.Params.Arguments["pull_number"].(float64))

			pr, resp, err := client.PullRequests.Get(ctx, owner, repo, pullNumber)
			if err != nil {
				return nil, fmt.Errorf("failed to get pull request: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get pull request: %s", string(body))), nil
			}

			r, err := json.Marshal(pr)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// listPullRequests creates a tool to list and filter repository pull requests.
func listPullRequests(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_pull_requests",
			mcp.WithDescription(t("TOOL_LIST_PULL_REQUESTS_DESCRIPTION", "List and filter repository pull requests")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithString("state",
				mcp.Description("Filter by state ('open', 'closed', 'all')"),
			),
			mcp.WithString("head",
				mcp.Description("Filter by head user/org and branch"),
			),
			mcp.WithString("base",
				mcp.Description("Filter by base branch"),
			),
			mcp.WithString("sort",
				mcp.Description("Sort by ('created', 'updated', 'popularity', 'long-running')"),
			),
			mcp.WithString("direction",
				mcp.Description("Sort direction ('asc', 'desc')"),
			),
			mcp.WithNumber("per_page",
				mcp.Description("Results per page (max 100)"),
			),
			mcp.WithNumber("page",
				mcp.Description("Page number"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			state := ""
			if s, ok := request.Params.Arguments["state"].(string); ok {
				state = s
			}
			head := ""
			if h, ok := request.Params.Arguments["head"].(string); ok {
				head = h
			}
			base := ""
			if b, ok := request.Params.Arguments["base"].(string); ok {
				base = b
			}
			sort := ""
			if s, ok := request.Params.Arguments["sort"].(string); ok {
				sort = s
			}
			direction := ""
			if d, ok := request.Params.Arguments["direction"].(string); ok {
				direction = d
			}
			perPage := 30
			if pp, ok := request.Params.Arguments["per_page"].(float64); ok {
				perPage = int(pp)
			}
			page := 1
			if p, ok := request.Params.Arguments["page"].(float64); ok {
				page = int(p)
			}

			opts := &github.PullRequestListOptions{
				State:     state,
				Head:      head,
				Base:      base,
				Sort:      sort,
				Direction: direction,
				ListOptions: github.ListOptions{
					PerPage: perPage,
					Page:    page,
				},
			}

			prs, resp, err := client.PullRequests.List(ctx, owner, repo, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list pull requests: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list pull requests: %s", string(body))), nil
			}

			r, err := json.Marshal(prs)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// mergePullRequest creates a tool to merge a pull request.
func mergePullRequest(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("merge_pull_request",
			mcp.WithDescription(t("TOOL_MERGE_PULL_REQUEST_DESCRIPTION", "Merge a pull request")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("pull_number",
				mcp.Required(),
				mcp.Description("Pull request number"),
			),
			mcp.WithString("commit_title",
				mcp.Description("Title for merge commit"),
			),
			mcp.WithString("commit_message",
				mcp.Description("Extra detail for merge commit"),
			),
			mcp.WithString("merge_method",
				mcp.Description("Merge method ('merge', 'squash', 'rebase')"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			pullNumber := int(request.Params.Arguments["pull_number"].(float64))
			commitTitle := ""
			if ct, ok := request.Params.Arguments["commit_title"].(string); ok {
				commitTitle = ct
			}
			commitMessage := ""
			if cm, ok := request.Params.Arguments["commit_message"].(string); ok {
				commitMessage = cm
			}
			mergeMethod := ""
			if mm, ok := request.Params.Arguments["merge_method"].(string); ok {
				mergeMethod = mm
			}

			options := &github.PullRequestOptions{
				CommitTitle: commitTitle,
				MergeMethod: mergeMethod,
			}

			result, resp, err := client.PullRequests.Merge(ctx, owner, repo, pullNumber, commitMessage, options)
			if err != nil {
				return nil, fmt.Errorf("failed to merge pull request: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to merge pull request: %s", string(body))), nil
			}

			r, err := json.Marshal(result)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// getPullRequestFiles creates a tool to get the list of files changed in a pull request.
func getPullRequestFiles(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pull_request_files",
			mcp.WithDescription(t("TOOL_GET_PULL_REQUEST_FILES_DESCRIPTION", "Get the list of files changed in a pull request")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("pull_number",
				mcp.Required(),
				mcp.Description("Pull request number"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			pullNumber := int(request.Params.Arguments["pull_number"].(float64))

			opts := &github.ListOptions{}
			files, resp, err := client.PullRequests.ListFiles(ctx, owner, repo, pullNumber, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to get pull request files: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get pull request files: %s", string(body))), nil
			}

			r, err := json.Marshal(files)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// getPullRequestStatus creates a tool to get the combined status of all status checks for a pull request.
func getPullRequestStatus(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pull_request_status",
			mcp.WithDescription(t("TOOL_GET_PULL_REQUEST_STATUS_DESCRIPTION", "Get the combined status of all status checks for a pull request")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("pull_number",
				mcp.Required(),
				mcp.Description("Pull request number"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			pullNumber := int(request.Params.Arguments["pull_number"].(float64))

			// First get the PR to find the head SHA
			pr, resp, err := client.PullRequests.Get(ctx, owner, repo, pullNumber)
			if err != nil {
				return nil, fmt.Errorf("failed to get pull request: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get pull request: %s", string(body))), nil
			}

			// Get combined status for the head SHA
			status, resp, err := client.Repositories.GetCombinedStatus(ctx, owner, repo, *pr.Head.SHA, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to get combined status: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get combined status: %s", string(body))), nil
			}

			r, err := json.Marshal(status)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// updatePullRequestBranch creates a tool to update a pull request branch with the latest changes from the base branch.
func updatePullRequestBranch(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("update_pull_request_branch",
			mcp.WithDescription(t("TOOL_UPDATE_PULL_REQUEST_BRANCH_DESCRIPTION", "Update a pull request branch with the latest changes from the base branch")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("pull_number",
				mcp.Required(),
				mcp.Description("Pull request number"),
			),
			mcp.WithString("expected_head_sha",
				mcp.Description("The expected SHA of the pull request's HEAD ref"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			pullNumber := int(request.Params.Arguments["pull_number"].(float64))
			expectedHeadSHA := ""
			if sha, ok := request.Params.Arguments["expected_head_sha"].(string); ok {
				expectedHeadSHA = sha
			}

			opts := &github.PullRequestBranchUpdateOptions{}
			if expectedHeadSHA != "" {
				opts.ExpectedHeadSHA = github.Ptr(expectedHeadSHA)
			}

			result, resp, err := client.PullRequests.UpdateBranch(ctx, owner, repo, pullNumber, opts)
			if err != nil {
				// Check if it's an acceptedError. An acceptedError indicates that the update is in progress,
				// and it's not a real error.
				if resp != nil && resp.StatusCode == http.StatusAccepted && isAcceptedError(err) {
					return mcp.NewToolResultText("Pull request branch update is in progress"), nil
				}
				return nil, fmt.Errorf("failed to update pull request branch: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusAccepted {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to update pull request branch: %s", string(body))), nil
			}

			r, err := json.Marshal(result)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// getPullRequestComments creates a tool to get the review comments on a pull request.
func getPullRequestComments(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pull_request_comments",
			mcp.WithDescription(t("TOOL_GET_PULL_REQUEST_COMMENTS_DESCRIPTION", "Get the review comments on a pull request")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("pull_number",
				mcp.Required(),
				mcp.Description("Pull request number"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			pullNumber := int(request.Params.Arguments["pull_number"].(float64))

			opts := &github.PullRequestListCommentsOptions{
				ListOptions: github.ListOptions{
					PerPage: 100,
				},
			}

			comments, resp, err := client.PullRequests.ListComments(ctx, owner, repo, pullNumber, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to get pull request comments: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get pull request comments: %s", string(body))), nil
			}

			r, err := json.Marshal(comments)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// getPullRequestReviews creates a tool to get the reviews on a pull request.
func getPullRequestReviews(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_pull_request_reviews",
			mcp.WithDescription(t("TOOL_GET_PULL_REQUEST_REVIEWS_DESCRIPTION", "Get the reviews on a pull request")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("pull_number",
				mcp.Required(),
				mcp.Description("Pull request number"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			pullNumber := int(request.Params.Arguments["pull_number"].(float64))

			reviews, resp, err := client.PullRequests.ListReviews(ctx, owner, repo, pullNumber, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to get pull request reviews: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get pull request reviews: %s", string(body))), nil
			}

			r, err := json.Marshal(reviews)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// createPullRequestReview creates a tool to submit a review on a pull request.
func createPullRequestReview(client *github.Client, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_pull_request_review",
			mcp.WithDescription(t("TOOL_CREATE_PULL_REQUEST_REVIEW_DESCRIPTION", "Create a review on a pull request")),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("Repository owner"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("Repository name"),
			),
			mcp.WithNumber("pull_number",
				mcp.Required(),
				mcp.Description("Pull request number"),
			),
			mcp.WithString("body",
				mcp.Description("Review comment text"),
			),
			mcp.WithString("event",
				mcp.Required(),
				mcp.Description("Review action ('APPROVE', 'REQUEST_CHANGES', 'COMMENT')"),
			),
			mcp.WithString("commit_id",
				mcp.Description("SHA of commit to review"),
			),
			mcp.WithArray("comments",
				mcp.Description("Line-specific comments array of objects, each object with path (string), position (number), and body (string)"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner := request.Params.Arguments["owner"].(string)
			repo := request.Params.Arguments["repo"].(string)
			pullNumber := int(request.Params.Arguments["pull_number"].(float64))
			event := request.Params.Arguments["event"].(string)

			// Create review request
			reviewRequest := &github.PullRequestReviewRequest{
				Event: github.Ptr(event),
			}

			// Add body if provided
			if body, ok := request.Params.Arguments["body"].(string); ok && body != "" {
				reviewRequest.Body = github.Ptr(body)
			}

			// Add commit ID if provided
			if commitID, ok := request.Params.Arguments["commit_id"].(string); ok && commitID != "" {
				reviewRequest.CommitID = github.Ptr(commitID)
			}

			// Add comments if provided
			if commentsObj, ok := request.Params.Arguments["comments"].([]interface{}); ok && len(commentsObj) > 0 {
				comments := []*github.DraftReviewComment{}

				for _, c := range commentsObj {
					commentMap, ok := c.(map[string]interface{})
					if !ok {
						return mcp.NewToolResultError("each comment must be an object with path, position, and body"), nil
					}

					path, ok := commentMap["path"].(string)
					if !ok || path == "" {
						return mcp.NewToolResultError("each comment must have a path"), nil
					}

					positionFloat, ok := commentMap["position"].(float64)
					if !ok {
						return mcp.NewToolResultError("each comment must have a position"), nil
					}
					position := int(positionFloat)

					body, ok := commentMap["body"].(string)
					if !ok || body == "" {
						return mcp.NewToolResultError("each comment must have a body"), nil
					}

					comments = append(comments, &github.DraftReviewComment{
						Path:     github.Ptr(path),
						Position: github.Ptr(position),
						Body:     github.Ptr(body),
					})
				}

				reviewRequest.Comments = comments
			}

			review, resp, err := client.PullRequests.CreateReview(ctx, owner, repo, pullNumber, reviewRequest)
			if err != nil {
				return nil, fmt.Errorf("failed to create pull request review: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to create pull request review: %s", string(body))), nil
			}

			r, err := json.Marshal(review)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
