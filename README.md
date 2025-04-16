# Cursor Rules CLI

Cursor Rules CLI는 터미널에서 Cursor rules(프롬프트/지침 템플릿) 파일을 쉽게 관리할 수 있는 도구입니다.

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

### 1. GitHub 토큰 설정

먼저 GitHub Personal Access Token을 발급받아야 합니다:

1. GitHub 계정으로 로그인
2. Settings > Developer settings > Personal access tokens > Tokens (classic)로 이동
3. "Generate new token (classic)" 클릭
4. 토큰 설정:
   - Note: `cursorrules-cli` (또는 원하는 이름)
   - Expiration: 원하는 만료 기간 선택
   - Select scopes: `gist` 체크
5. "Generate token" 클릭하여 토큰 발급
6. 발급된 토큰을 복사하여 다음 명령어로 설정:
   ```bash
   cursorrules auth
   ```

### 2. 기본 명령어

```bash
# 템플릿 목록 보기
cursorrules list

# 템플릿 다운로드
cursorrules pull 'template-name'

# 로컬 템플릿 업로드
cursorrules push 'template-name'
```

### 상세 설명

1. **인증**
   ```bash
   cursorrules auth
   ```
   - GitHub Personal Access Token을 통한 인증
   - 토큰은 `~/.cursorrules/config-cli.json`에 저장

2. **템플릿 목록**
   ```bash
   cursorrules list
   ```
   - GitHub Gist에 저장된 템플릿 목록 출력
   - 각 템플릿의 이름과 설명 표시

3. **템플릿 다운로드**
   ```bash
   cursorrules pull <template-name>
   ```
   - 지정한 템플릿을 현재 디렉토리에 다운로드
   - 기존 파일이 있으면 덮어쓰기 전 확인

4. **템플릿 업로드**
   ```bash
   cursorrules push <template-name>
   ```
   - 현재 디렉토리의 `.cursor/rules/**/*.mdc` 파일들을 찾아서 JSON으로 변환
   - GitHub Gist에 업로드
   - 파일 경로와 내용을 포함

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

기여는 언제나 환영합니다! 이슈를 열거나 PR을 보내주세요.

## 라이선스

MIT License
