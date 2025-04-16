package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// VersionInfo 파일의 버전 정보
type VersionInfo struct {
	LastModified time.Time `json:"last_modified"` // 파일의 마지막 수정 시간
	LastSynced   time.Time `json:"last_synced"`   // 마지막 동기화 시간
	Hash         string    `json:"hash"`          // 파일 내용의 해시
}

// TemplateVersion 템플릿의 버전 정보
type TemplateVersion struct {
	Name        string                `json:"name"`         // 템플릿 이름
	Version     string                `json:"version"`      // 버전 (예: v1.0.0)
	Files       map[string]VersionInfo `json:"files"`        // 파일별 버전 정보
	Description string                `json:"description"`  // 설명
	CreatedAt   time.Time             `json:"created_at"`   // 생성 시간
	UpdatedAt   time.Time             `json:"updated_at"`   // 마지막 업데이트 시간
}

// NewTemplateVersion 새로운 템플릿 버전 생성
func NewTemplateVersion(name, version string) *TemplateVersion {
	return &TemplateVersion{
		Name:        name,
		Version:     version,
		Files:       make(map[string]VersionInfo),
		Description: "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// AddFile 파일 추가
func (tv *TemplateVersion) AddFile(path string, lastModified time.Time, hash string) {
	tv.Files[path] = VersionInfo{
		LastModified: lastModified,
		LastSynced:   time.Now(),
		Hash:         hash,
	}
	tv.UpdatedAt = time.Now()
}

// GetFile 파일 정보 조회
func (tv *TemplateVersion) GetFile(path string) (VersionInfo, bool) {
	info, exists := tv.Files[path]
	return info, exists
}

// UpdateFile 파일 정보 업데이트
func (tv *TemplateVersion) UpdateFile(path string, lastModified time.Time, hash string) {
	if info, exists := tv.Files[path]; exists {
		info.LastModified = lastModified
		info.Hash = hash
		info.LastSynced = time.Now()
		tv.Files[path] = info
		tv.UpdatedAt = time.Now()
	}
}

// ToJSONString JSON 문자열로 변환
func (v *TemplateVersion) ToJSONString() (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON 변환 실패: %v", err)
	}
	return string(data), nil
}

// FromJSONString JSON 문자열에서 템플릿 버전 생성
func FromJSONString(jsonStr string) (*TemplateVersion, error) {
	var version TemplateVersion
	if err := json.Unmarshal([]byte(jsonStr), &version); err != nil {
		return nil, fmt.Errorf("JSON 파싱 실패: %v", err)
	}
	return &version, nil
} 