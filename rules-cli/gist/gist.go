package gist

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v58/github"
	"github.com/tinysolver/rules-cli/config"
)

const (
	// GistTag Cursor Rules CLI에서 사용하는 Gist임을 식별하는 태그
	GistTag = "[cursor-rules-cli]"
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

	ts := github.BasicAuthTransport{
		Username: "token",
		Password: token,
	}
	client := github.NewClient(ts.Client())

	return &GistClient{client: client}, nil
}

// IsCursorRulesGist Gist가 Cursor Rules CLI에서 사용하는 Gist인지 확인
func IsCursorRulesGist(gist *github.Gist) bool {
	return strings.Contains(gist.GetDescription(), GistTag)
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

	// Cursor Rules CLI에서 사용하는 Gist만 필터링
	var cursorGists []*github.Gist
	for _, gist := range gists {
		if IsCursorRulesGist(gist) {
			cursorGists = append(cursorGists, gist)
		}
	}

	return cursorGists, nil
}

// GistInfo Gist의 상세 정보
type GistInfo struct {
	ID          string            `json:"id"`
	Description string            `json:"description"`
	Public      bool              `json:"public"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Files       map[string]FileInfo `json:"files"`
	Owner       string            `json:"owner"`
	Version     string            `json:"version"`
}

// FileInfo Gist 파일의 상세 정보
type FileInfo struct {
	Filename  string `json:"filename"`
	Type      string `json:"type"`
	Language  string `json:"language"`
	Size      int    `json:"size"`
	RawURL    string `json:"raw_url"`
	Content   string `json:"content"`
}

// GetGistInfo Gist의 상세 정보 조회
func (g *GistClient) GetGistInfo(gistID string) (*GistInfo, error) {
	ctx := context.Background()
	gist, _, err := g.client.Gists.Get(ctx, gistID)
	if err != nil {
		return nil, fmt.Errorf("Gist 정보 조회 실패: %v", err)
	}

	info := &GistInfo{
		ID:          gist.GetID(),
		Description: gist.GetDescription(),
		Public:      gist.GetPublic(),
		CreatedAt:   gist.GetCreatedAt(),
		UpdatedAt:   gist.GetUpdatedAt(),
		Files:       make(map[string]FileInfo),
		Owner:       gist.GetOwner().GetLogin(),
		Version:     gist.GetHistory()[0].GetVersion(),
	}

	for filename, file := range gist.Files {
		info.Files[filename] = FileInfo{
			Filename:  file.GetFilename(),
			Type:      file.GetType(),
			Language:  file.GetLanguage(),
			Size:      file.GetSize(),
			RawURL:    file.GetRawURL(),
			Content:   file.GetContent(),
		}
	}

	return info, nil
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
	// 설명에 태그 추가
	description = fmt.Sprintf("%s %s", GistTag, description)

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

// FindGistByDescription 프로젝트 이름으로 Gist 찾기
func (g *GistClient) FindGistByDescription(description string) (*github.Gist, error) {
	gists, err := g.ListGists()
	if err != nil {
		return nil, err
	}

	for _, gist := range gists {
		if gist.GetDescription() == description {
			return gist, nil
		}
	}

	return nil, fmt.Errorf("프로젝트 '%s'를 찾을 수 없습니다", description)
}

// DeleteGist Gist 삭제
func (g *GistClient) DeleteGist(gistID string) error {
	ctx := context.Background()
	_, err := g.client.Gists.Delete(ctx, gistID)
	if err != nil {
		return fmt.Errorf("Gist 삭제 실패: %v", err)
	}
	return nil
} 