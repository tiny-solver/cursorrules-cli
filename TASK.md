# Cursor Rules CLI 개발 작업 계획

## Phase 1: 기본 CLI 구조 (3일)

### 1.1 프로젝트 초기화 ✅
- [x] Go 프로젝트 설정
  ```bash
  mkdir rules-cli
  cd rules-cli
  go mod init github.com/tinysolver/rules-cli
  ```

- [x] 의존성 설치
  ```bash
  go get github.com/spf13/cobra
  go get github.com/spf13/viper
  go get github.com/google/go-github/v58
  ```

### 1.2 CLI 명령어 구현 ✅
- [x] 기본 명령어 구조
  ```go
  package main

  import (
      "fmt"
      "os"
      "github.com/spf13/cobra"
  )

  var rootCmd = &cobra.Command{
      Use:   "cursorrules",
      Short: "Cursor Rules CLI",
      Long:  "Cursor Rules CLI는 터미널에서 Cursor rules 파일을 관리하는 도구입니다.",
  }

  var authCmd = &cobra.Command{
      Use:   "auth",
      Short: "GitHub Personal Access Token 설정",
      Run: func(cmd *cobra.Command, args []string) {
          fmt.Println("GitHub 토큰 설정...")
          // TODO: 토큰 설정 로직 구현
      },
  }

  var listCmd = &cobra.Command{
      Use:   "list",
      Short: "템플릿 목록 출력",
      Run: func(cmd *cobra.Command, args []string) {
          fmt.Println("템플릿 목록 출력...")
          // TODO: 목록 출력 로직 구현
      },
  }

  var syncCmd = &cobra.Command{
      Use:   "sync [name]",
      Short: "템플릿 다운로드",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {
          fmt.Printf("템플릿 %s 다운로드 중...\n", args[0])
          // TODO: 다운로드 로직 구현
      },
  }

  var pushCmd = &cobra.Command{
      Use:   "push [name]",
      Short: "로컬 템플릿 업로드",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {
          fmt.Printf("템플릿 %s 업로드 중...\n", args[0])
          // TODO: 업로드 로직 구현
      },
  }

  func init() {
      rootCmd.AddCommand(authCmd)
      rootCmd.AddCommand(listCmd)
      rootCmd.AddCommand(syncCmd)
      rootCmd.AddCommand(pushCmd)
  }

  func main() {
      if err := rootCmd.Execute(); err != nil {
          fmt.Println(err)
          os.Exit(1)
      }
  }
  ```

- [x] 명령어 테스트
  - [x] `cursorrules --help`
  - [x] `cursorrules auth`
  - [x] `cursorrules list`
  - [x] `cursorrules sync <name>`
  - [x] `cursorrules push <name>`

## Phase 2: 인증 시스템 (2일)

### 2.1 Personal Access Token 관리
- [ ] 토큰 설정 구현
  ```go
  package config

  import (
      "fmt"
      "os"
      "path/filepath"
      "github.com/spf13/viper"
  )

  func InitConfig() error {
      home, err := os.UserHomeDir()
      if err != nil {
          return err
      }

      configPath := filepath.Join(home, ".cursorrules")
      if err := os.MkdirAll(configPath, 0755); err != nil {
          return err
      }

      viper.SetConfigName("config-cli")
      viper.SetConfigType("json")
      viper.AddConfigPath(configPath)

      return viper.ReadInConfig()
  }

  func SaveToken(token string) error {
      viper.Set("github_token", token)
      return viper.WriteConfig()
  }

  func GetToken() string {
      return viper.GetString("github_token")
  }
  ```

### 2.2 토큰 검증
- [ ] GitHub API로 토큰 유효성 검사
- [ ] 토큰 권한 확인 (gist 접근 권한)
- [ ] 에러 처리

## Phase 3: Gist 연동 (3일)

### 3.1 Gist API 클라이언트
- [ ] GitHub API 클라이언트 구현
  ```go
  package github

  import (
      "context"
      "github.com/google/go-github/v58/github"
      "golang.org/x/oauth2"
  )

  type Client struct {
      client *github.Client
  }

  func NewClient(token string) *Client {
      ctx := context.Background()
      ts := oauth2.StaticTokenSource(
          &oauth2.Token{AccessToken: token},
      )
      tc := oauth2.NewClient(ctx, ts)
      return &Client{
          client: github.NewClient(tc),
      }
  }
  ```

### 3.2 데이터 구조
- [ ] 템플릿 구조 정의
  ```go
  package models

  type Rule struct {
      Name    string `json:"name"`
      Content string `json:"content"`
      Path    string `json:"path"`
  }

  type Template struct {
      Files map[string]Rule `json:"files"`
  }
  ```

## Phase 4: 파일 시스템 (2일)

### 4.1 로컬 파일 관리
- [ ] 파일 읽기/쓰기
- [ ] JSON 변환
- [ ] 경로 처리

### 4.2 동기화 로직
- [ ] 다운로드 구현
- [ ] 업로드 구현
- [ ] 충돌 처리

## 당장 해야할 작업 (오늘)

1. **인증 시스템 구현**
   - Personal Access Token 설정 구현
   - 토큰 저장 및 검증 로직 구현
   - 설정 파일 구조 구현

2. **설정 파일 구조 구현**
   - `~/.cursorrules/config-cli.json` 생성
   - 토큰 저장 로직 구현
   - 설정 읽기/쓰기 구현

## 리스크 관리

### 기술적 리스크
1. GitHub API 제한
   - API 호출 최소화
   - 캐싱 전략 수립

2. 파일 시스템
   - 권한 문제
   - 경로 처리
   - 충돌 해결

## 성공 지표

### 기술적 지표
- 명령어 실행 성공률
- API 응답 시간
- 파일 처리 안정성

### 사용자 경험 지표
- 명령어 직관성
- 에러 메시지 명확성
- 작업 완료 시간 