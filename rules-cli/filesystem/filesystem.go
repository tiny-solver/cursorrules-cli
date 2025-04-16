package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"crypto/md5"

	"github.com/tinysolver/rules-cli/models"
)

const (
	rulesDir = ".cursor/rules"
)

// GetRulesDir 규칙 디렉토리 경로 조회
func GetRulesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("홈 디렉토리를 찾을 수 없습니다: %v", err)
	}

	dir := filepath.Join(home, rulesDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("규칙 디렉토리를 생성할 수 없습니다: %v", err)
	}

	return dir, nil
}

// LoadLocalTemplate 로컬 템플릿 로드
func LoadLocalTemplate() (*models.Template, *models.TemplateVersion, error) {
	rulesDir, err := GetRulesDir()
	if err != nil {
		return nil, nil, err
	}

	template := &models.Template{
		Files: make(map[string]models.Rule),
	}

	version := models.NewTemplateVersion("local", "v1.0.0")

	err = filepath.Walk(rulesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("파일 읽기 실패: %v", err)
		}

		relPath, err := filepath.Rel(rulesDir, path)
		if err != nil {
			return fmt.Errorf("상대 경로 변환 실패: %v", err)
		}

		rule := models.Rule{
			Name:    relPath,
			Content: string(content),
		}

		template.Files[relPath] = rule

		// 버전 정보 추가
		hash := fmt.Sprintf("%x", md5.Sum(content))
		version.AddFile(relPath, info.ModTime(), hash)

		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("템플릿 로드 실패: %v", err)
	}

	return template, version, nil
}

// SaveLocalTemplate 로컬 템플릿 저장
func SaveLocalTemplate(template *models.Template, version *models.TemplateVersion) error {
	dir, err := GetRulesDir()
	if err != nil {
		return err
	}

	// 버전 정보 저장
	versionData, err := version.ToJSONString()
	if err != nil {
		return fmt.Errorf("버전 정보 변환 실패: %v", err)
	}

	versionPath := filepath.Join(dir, "version.json")
	if err := os.WriteFile(versionPath, []byte(versionData), 0644); err != nil {
		return fmt.Errorf("버전 정보 저장 실패: %v", err)
	}

	// 파일 저장
	for _, file := range template.Files {
		filePath := filepath.Join(dir, file.Name)
		content := []byte(file.Content)
		
		// 기존 파일 백업
		if _, err := os.Stat(filePath); err == nil {
			backupPath := filePath + ".bak"
			if err := os.Rename(filePath, backupPath); err != nil {
				return fmt.Errorf("백업 생성 실패: %v", err)
			}
		}

		// 새 파일 저장
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return fmt.Errorf("파일 저장 실패: %v", err)
		}
	}

	return nil
}

// CheckConflicts 충돌 확인
func CheckConflicts(template *models.Template) ([]string, error) {
	dir, err := GetRulesDir()
	if err != nil {
		return nil, err
	}

	var conflicts []string
	for _, file := range template.Files {
		filePath := filepath.Join(dir, file.Name)
		if _, err := os.Stat(filePath); err == nil {
			conflicts = append(conflicts, file.Name)
		}
	}

	return conflicts, nil
}

// MergeTemplate 템플릿 병합
func MergeTemplate(template *models.Template, version *models.TemplateVersion) error {
	dir, err := GetRulesDir()
	if err != nil {
		return err
	}

	// 버전 정보 저장
	versionData, err := version.ToJSONString()
	if err != nil {
		return fmt.Errorf("버전 정보 변환 실패: %v", err)
	}

	versionPath := filepath.Join(dir, "version.json")
	if err := os.WriteFile(versionPath, []byte(versionData), 0644); err != nil {
		return fmt.Errorf("버전 정보 저장 실패: %v", err)
	}

	// 파일 저장
	for _, file := range template.Files {
		filePath := filepath.Join(dir, file.Name)
		content := []byte(file.Content)
		
		// 기존 파일 백업
		if _, err := os.Stat(filePath); err == nil {
			backupPath := filePath + ".bak"
			if err := os.Rename(filePath, backupPath); err != nil {
				return fmt.Errorf("백업 생성 실패: %v", err)
			}
		}

		// 새 파일 저장
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return fmt.Errorf("파일 저장 실패: %v", err)
		}
	}

	return nil
} 