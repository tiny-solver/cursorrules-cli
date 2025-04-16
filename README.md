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

### 기본 명령어

```bash
# GitHub OAuth 인증
cursorrules auth

# 템플릿 목록 보기
cursorrules list

# 템플릿 다운로드
cursorrules sync <template-name>

# 로컬 템플릿 업로드
cursorrules push <template-name>
```

### 상세 설명

1. **인증**
   ```bash
   cursorrules auth
   ```
   - GitHub OAuth를 통한 인증
   - 토큰은 `~/.cursorrules/config-cli.json`에 저장

2. **템플릿 목록**
   ```bash
   cursorrules list
   ```
   - GitHub Gist에 저장된 템플릿 목록 출력
   - 각 템플릿의 이름과 설명 표시

3. **템플릿 다운로드**
   ```bash
   cursorrules sync <template-name>
   ```
   - 지정한 템플릿을 현재 디렉토리에 다운로드
   - 기존 파일이 있으면 덮어쓰기 전 확인

4. **템플릿 업로드**
   ```bash
   cursorrules push <template-name>
   ```
   - `~/.cursor/rules/*`의 모든 파일을 JSON으로 변환
   - GitHub Gist에 업로드
   - 파일 경로와 내용을 포함

## 파일 구조

### 로컬 저장소
```
~/.cursor/rules/
  ├── rule1.json
  ├── rule2.json
  └── ...
```

### 설정 파일
```
~/.cursorrules/config-cli.json
  ├── github_token
  └── gist_id
```

## 기여하기

프로젝트에 기여하고 싶으시다면 다음 단계를 따라주세요:

1. 이슈 생성
2. 브랜치 생성
3. 변경사항 커밋
4. PR 생성

## 라이선스

MIT License
