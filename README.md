# Cursor Rules CLI

Cursor Rules CLI는 터미널에서 Cursor rules 파일을 관리하는 도구입니다.

## 주요 기능

1. GitHub Personal Access Token을 통한 간편한 인증
2. GitHub Gist를 이용한 규칙 파일 동기화
3. 안전한 파일 관리 (덮어쓰기 확인 및 백업)
4. 템플릿 삭제 기능

## 설치 방법

### 바이너리 다운로드 (추후 업데이트 예정)
```bash
# Linux
curl -L https://github.com/tinysolver/rules-cli/releases/latest/download/rules-cli-linux-amd64 -o /usr/local/bin/cursorrules
chmod +x /usr/local/bin/cursorrules

# macOS
brew install tinysolver/tap/cursorrules

# Windows
# 추후 업데이트 예정
```

### 소스코드로 설치
```bash
go install github.com/tinysolver/rules-cli@latest
```

## 사용 방법

### 1. 인증 설정

```bash
cursorrules auth
```

GitHub Personal Access Token을 입력하여 인증을 설정합니다.

### 2. 템플릿 목록 보기

```bash
cursorrules list
```

GitHub Gist에 저장된 템플릿 목록을 보여줍니다.

### 3. 템플릿 다운로드

```bash
cursorrules download <템플릿이름>
```

지정한 템플릿을 다운로드합니다. 로컬 파일과 충돌이 있는 경우:
- 파일의 수정 시간을 비교하여 알려줍니다
- `--force` 옵션을 사용하면 강제로 덮어쓸 수 있습니다

### 4. 템플릿 업로드

```bash
cursorrules upload <템플릿이름>
```

현재 디렉토리의 규칙 파일을 GitHub Gist에 업로드합니다. 기존 템플릿이 있으면 덮어씁니다.

### 5. 템플릿 삭제

```bash
cursorrules delete <템플릿이름>
```

지정한 템플릿을 GitHub Gist에서 삭제합니다. 삭제 전 확인 메시지가 표시됩니다.

## 파일 구조

### 로컬 저장소
```
.cursor/rules/
  ├── rule1.mdc
  ├── rule2.mdc
  └── ...
```

### 설정 파일
```
~/.cursorrules/config-cli.json
  └── github_token
```

## 기여하기

기여는 언제나 환영합니다! 버그 리포트, 기능 제안, 풀 리퀘스트 등을 통해 참여해주세요.

## 라이선스

MIT License
