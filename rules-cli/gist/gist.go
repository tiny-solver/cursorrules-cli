package gist

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/tinysolver/rules-cli/config"
)

// GistClient GitHub Gist API 클라이언트
type GistClient struct {
	client *github.Client
}

// NewGistClient 새로운 Gist 클라이언트 생성
func NewGistClient() (*GistClient, error) {
	token := config.GetToken()
	if token == "" {
		return nil, fmt.Errorf("GitHub 토큰이 설정되지 않았습니다. 'cursorrules auth' 명령어로 토큰을 설정하세요")
	}

	ctx := context.Background()
	ts := github.BasicAuthTransport{
		Username: "token",
		Password: token,
	}
	client := github.NewClient(ts.Client())

	return &GistClient{client: client}, nil
}

// ListGists 저장된 Gist 목록 조회
func (g *GistClient) ListGists() ([]*github.Gist, error) {
	ctx := context.Background()
	gists, _, err := g.client.Gists.List(ctx, "", &github.GistListOptions{
		Since: time.Now().Add(-24 * 365 * time.Hour), // 1년 이내의 Gist만 조회
	})
	if err != nil {
		return nil, fmt.Errorf("Gist 목록 조회 실패: %v", err)
	}

	return gists, nil
}

// GetGistContent Gist 내용 조회
func (g *GistClient) GetGistContent(gistID string) (map[string]string, error) {
	ctx := context.Background()
	gist, _, err := g.client.Gists.Get(ctx, gistID)
	if err != nil {
		return nil, fmt.Errorf("Gist 내용 조회 실패: %v", err)
	}

	contents := make(map[string]string)
	for _, file := range gist.Files {
		if file.Content != nil {
			contents[file.GetFilename()] = *file.Content
		}
	}

	return contents, nil
}

// CreateGist 새로운 Gist 생성
func (g *GistClient) CreateGist(description string, files map[string]string) (*github.Gist, error) {
	ctx := context.Background()
	gist := &github.Gist{
		Description: &description,
		Public:      github.Bool(false),
		Files:       make(map[github.GistFilename]github.GistFile),
	}

	for filename, content := range files {
		gist.Files[github.GistFilename(filename)] = github.GistFile{
			Content: &content,
		}
	}

	newGist, _, err := g.client.Gists.Create(ctx, gist)
	if err != nil {
		return nil, fmt.Errorf("Gist 생성 실패: %v", err)
	}

	return newGist, nil
} 