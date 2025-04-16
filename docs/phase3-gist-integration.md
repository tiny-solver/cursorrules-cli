# Phase 3: Gist 연동 상세 구현

## 3.1 Gist API 클라이언트 구현

### GitHub API 클라이언트
```go
type GistClient struct {
    client *github.Client
}

func NewGistClient() (*GistClient, error) {
    token := config.GetToken()
    if token == "" {
        return nil, fmt.Errorf("GitHub 토큰이 설정되지 않았습니다")
    }

    ts := github.BasicAuthTransport{
        Username: "token",
        Password: token,
    }
    client := github.NewClient(ts.Client())

    return &GistClient{client: client}, nil
}
```

### Gist 목록 조회
- 1년 이내의 Gist만 조회하도록 제한
- 프로젝트 이름(description)으로 Gist 식별
- 에러 처리 및 사용자 친화적 메시지

```go
func (g *GistClient) ListGists() ([]*github.Gist, error) {
    ctx := context.Background()
    gists, _, err := g.client.Gists.List(ctx, "", &github.GistListOptions{
        Since: time.Now().Add(-24 * 365 * time.Hour),
    })
    if err != nil {
        return nil, fmt.Errorf("Gist 목록 조회 실패: %v", err)
    }
    return gists, nil
}
```

### 프로젝트 이름 기반 Gist 관리
- Gist의 description을 프로젝트 이름으로 사용
- 프로젝트 이름으로 Gist 검색 기능

```go
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
```

## 3.2 데이터 구조 (진행 중)

### 템플릿 구조
- Rule: 개별 규칙 파일
- Template: 프로젝트 단위의 규칙 모음

```go
type Rule struct {
    Name    string `json:"name"`
    Content string `json:"content"`
    Path    string `json:"path"`
}

type Template struct {
    Files map[string]Rule `json:"files"`
}
```

### JSON 변환
- 로컬 파일 → JSON 변환
- JSON → 로컬 파일 변환
- 에러 처리 및 유효성 검사

### 파일 시스템 연동
- 로컬 파일 읽기/쓰기
- 경로 처리 및 권한 관리
- 충돌 처리 및 백업

## 다음 단계
1. 템플릿 구조 구현 완료
2. JSON 변환 로직 구현
3. 파일 시스템 연동 구현 