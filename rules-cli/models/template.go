package models

// Rule 개별 규칙 파일
type Rule struct {
	Name    string `json:"name"`    // 파일 이름
	Content string `json:"content"` // 파일 내용
	Path    string `json:"path"`    // 파일 경로
}

// Template 프로젝트 단위의 규칙 모음
type Template struct {
	Name        string          `json:"name"`        // 프로젝트 이름
	Description string          `json:"description"` // 프로젝트 설명
	Files       map[string]Rule `json:"files"`       // 규칙 파일 목록
}

// NewTemplate 새로운 템플릿 생성
func NewTemplate(name, description string) *Template {
	return &Template{
		Name:        name,
		Description: description,
		Files:       make(map[string]Rule),
	}
}

// AddFile 템플릿에 파일 추가
func (t *Template) AddFile(name, content, path string) {
	t.Files[name] = Rule{
		Name:    name,
		Content: content,
		Path:    path,
	}
}

// GetFile 템플릿에서 파일 조회
func (t *Template) GetFile(name string) (Rule, bool) {
	rule, exists := t.Files[name]
	return rule, exists
}

// RemoveFile 템플릿에서 파일 제거
func (t *Template) RemoveFile(name string) {
	delete(t.Files, name)
} 