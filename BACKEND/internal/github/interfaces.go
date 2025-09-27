package github

type Client interface {
	GetRepos(owner string) ([]map[string]interface{}, error)
	CheckCollaborator(owner, repo, user string) (bool, error)
}
