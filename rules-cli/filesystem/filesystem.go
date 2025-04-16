package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"time"

	"github.com/tinysolver/rules-cli/models"
)

const (
	rulesDir = ".cursor/rules"
	versionFile = "version.json"
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
func LoadLocalTemplate() (*models.Template, *models.TemplateVersion, error) {
	dir, err := GetRulesDir()
	if err != nil {
		return nil, nil, err
	}

	template := models.NewTemplate()
	version := models.NewTemplateVersion("local", "v1.0.0")

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// version.json 파일은 건너뛰기
		if filepath.Base(path) == versionFile {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("파일 읽기 실패: %v", err)
		}

		// 파일 내용의 해시 계산
		hash := sha256.Sum256(content)
		hashStr := hex.EncodeToString(hash[:])

		// 템플릿에 파일 추가
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return fmt.Errorf("상대 경로 변환 실패: %v", err)
		}
		template.AddFile(relPath, string(content))

		// 버전 정보 추가
		version.AddFile(relPath, info.ModTime(), hashStr)

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return template, version, nil
}

// SaveLocalTemplate 로컬 템플릿 저장
func SaveLocalTemplate(template *models.Template, version *models.TemplateVersion) error {
	dir, err := GetRulesDir()
	if err != nil {
		return err
	}

	// 기존 파일 백업
	backupDir := filepath.Join(dir, "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %v", err)
	}

	// 각 파일 저장
	for path, content := range template.Files {
		fullPath := filepath.Join(dir, path)

		// 디렉토리 생성
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("디렉토리 생성 실패: %v", err)
		}

		// 기존 파일이 있으면 백업
		if _, err := os.Stat(fullPath); err == nil {
			backupPath := filepath.Join(backupDir, path)
			if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
				return fmt.Errorf("백업 디렉토리 생성 실패: %v", err)
			}
			if err := os.Rename(fullPath, backupPath); err != nil {
				return fmt.Errorf("파일 백업 실패: %v", err)
			}
		}

		// 새 파일 저장
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("파일 저장 실패: %v", err)
		}

		// 버전 정보 업데이트
		hash := sha256.Sum256([]byte(content))
		hashStr := hex.EncodeToString(hash[:])
		version.UpdateFile(path, time.Now(), hashStr)
	}

	// 버전 정보 저장
	versionPath := filepath.Join(dir, versionFile)
	versionData, err := version.ToJSON()
	if err != nil {
		return fmt.Errorf("버전 정보 JSON 변환 실패: %v", err)
	}
	if err := os.WriteFile(versionPath, versionData, 0644); err != nil {
		return fmt.Errorf("버전 정보 저장 실패: %v", err)
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

// CheckConflicts 로컬 파일과 템플릿의 충돌 확인
func CheckConflicts(template *models.Template) ([]string, error) {
	dir, err := GetRulesDir()
	if err != nil {
		return nil, err
	}

	var conflicts []string
	for path := range template.Files {
		fullPath := filepath.Join(dir, path)
		if _, err := os.Stat(fullPath); err == nil {
			conflicts = append(conflicts, path)
		}
	}

	return conflicts, nil
}

// MergeTemplate 로컬 파일과 템플릿 병합
func MergeTemplate(template *models.Template, version *models.TemplateVersion) error {
	dir, err := GetRulesDir()
	if err != nil {
		return err
	}

	// 기존 파일 백업
	backupDir := filepath.Join(dir, "backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("백업 디렉토리 생성 실패: %v", err)
	}

	// 각 파일 병합
	for path, content := range template.Files {
		fullPath := filepath.Join(dir, path)

		// 디렉토리 생성
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("디렉토리 생성 실패: %v", err)
		}

		// 기존 파일이 있으면 백업
		if _, err := os.Stat(fullPath); err == nil {
			backupPath := filepath.Join(backupDir, path)
			if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
				return fmt.Errorf("백업 디렉토리 생성 실패: %v", err)
			}
			if err := os.Rename(fullPath, backupPath); err != nil {
				return fmt.Errorf("파일 백업 실패: %v", err)
			}
		}

		// 새 파일 저장
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("파일 저장 실패: %v", err)
		}

		// 버전 정보 업데이트
		hash := sha256.Sum256([]byte(content))
		hashStr := hex.EncodeToString(hash[:])
		version.UpdateFile(path, time.Now(), hashStr)
	}

	// 버전 정보 저장
	versionPath := filepath.Join(dir, versionFile)
	versionData, err := version.ToJSON()
	if err != nil {
		return fmt.Errorf("버전 정보 JSON 변환 실패: %v", err)
	}
	if err := os.WriteFile(versionPath, versionData, 0644); err != nil {
		return fmt.Errorf("버전 정보 저장 실패: %v", err)
	}

	return nil
}

// generateUnifiedDiff 두 파일의 차이를 unified diff 포맷으로 생성
func generateUnifiedDiff(localContent, gistContent string) string {
	localLines := strings.Split(localContent, "\n")
	gistLines := strings.Split(gistContent, "\n")

	var diff strings.Builder
	diff.WriteString("@@ -1," + strconv.Itoa(len(localLines)) + " +1," + strconv.Itoa(len(gistLines)) + " @@\n")

	// 로컬 버전
	for _, line := range localLines {
		diff.WriteString("-" + line + "\n")
	}

	// Gist 버전
	for _, line := range gistLines {
		diff.WriteString("+" + line + "\n")
	}

	return diff.String()
} 