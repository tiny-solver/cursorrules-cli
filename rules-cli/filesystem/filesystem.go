package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tinysolver/rules-cli/models"
)

const (
	rulesDir = ".cursor/rules"
)

// GetRulesDir 규칙 파일 디렉토리 경로 조회
func GetRulesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("홈 디렉토리를 찾을 수 없습니다: %v", err)
	}

	rulesPath := filepath.Join(home, rulesDir)
	if err := os.MkdirAll(rulesPath, 0755); err != nil {
		return "", fmt.Errorf("규칙 디렉토리를 생성할 수 없습니다: %v", err)
	}

	return rulesPath, nil
}

// LoadLocalTemplate 로컬 템플릿 로드
func LoadLocalTemplate() (*models.Template, error) {
	rulesPath, err := GetRulesDir()
	if err != nil {
		return nil, err
	}

	template := models.NewTemplate("local", "로컬 규칙")
	err = filepath.Walk(rulesPath, func(path string, info os.FileInfo, err error) error {
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

		relPath, err := filepath.Rel(rulesPath, path)
		if err != nil {
			return fmt.Errorf("상대 경로 변환 실패: %v", err)
		}

		template.AddFile(info.Name(), string(content), relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("템플릿 로드 실패: %v", err)
	}

	return template, nil
}

// SaveLocalTemplate 로컬 템플릿 저장
func SaveLocalTemplate(template *models.Template) error {
	rulesPath, err := GetRulesDir()
	if err != nil {
		return err
	}

	// 기존 파일 백업
	backupPath := rulesPath + ".backup"
	if err := os.Rename(rulesPath, backupPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("백업 실패: %v", err)
	}

	// 새 디렉토리 생성
	if err := os.MkdirAll(rulesPath, 0755); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %v", err)
	}

	// 파일 저장
	for _, rule := range template.Files {
		filePath := filepath.Join(rulesPath, rule.Path)
		dirPath := filepath.Dir(filePath)

		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("디렉토리 생성 실패: %v", err)
		}

		if err := os.WriteFile(filePath, []byte(rule.Content), 0644); err != nil {
			return fmt.Errorf("파일 저장 실패: %v", err)
		}
	}

	// 백업 삭제
	if err := os.RemoveAll(backupPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("백업 삭제 실패: %v", err)
	}

	return nil
}

// FindLocalRules 현재 디렉토리의 .mdc 파일 찾기
func FindLocalRules() (*models.Template, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("현재 디렉토리를 찾을 수 없습니다: %v", err)
	}

	rulesPath := filepath.Join(cwd, rulesDir)
	if _, err := os.Stat(rulesPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("규칙 디렉토리를 찾을 수 없습니다: %s", rulesPath)
	}

	template := models.NewTemplate(filepath.Base(cwd), "로컬 규칙")
	err = filepath.Walk(rulesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".mdc" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("파일 읽기 실패: %v", err)
		}

		relPath, err := filepath.Rel(rulesPath, path)
		if err != nil {
			return fmt.Errorf("상대 경로 변환 실패: %v", err)
		}

		template.AddFile(info.Name(), string(content), relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("규칙 파일 검색 실패: %v", err)
	}

	if len(template.Files) == 0 {
		return nil, fmt.Errorf("규칙 파일을 찾을 수 없습니다: %s/*.mdc", rulesPath)
	}

	return template, nil
} 