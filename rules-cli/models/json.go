package models

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// ToJSON 템플릿을 JSON으로 변환
func (t *Template) ToJSON() ([]byte, error) {
	return json.MarshalIndent(t, "", "  ")
}

// FromJSON JSON을 템플릿으로 변환
func FromJSON(data []byte) (*Template, error) {
	var template Template
	if err := json.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("JSON 변환 실패: %v", err)
	}
	return &template, nil
}

// SaveToFile 템플릿을 파일로 저장
func (t *Template) SaveToFile(path string) error {
	data, err := t.ToJSON()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("파일 저장 실패: %v", err)
	}

	return nil
}

// LoadFromFile 파일에서 템플릿 로드
func LoadFromFile(path string) (*Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("파일 읽기 실패: %v", err)
	}

	return FromJSON(data)
}

// LoadFromReader Reader에서 템플릿 로드
func LoadFromReader(r io.Reader) (*Template, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("데이터 읽기 실패: %v", err)
	}

	return FromJSON(data)
}

// FromJSON JSON 데이터를 Template으로 변환
func (t *Template) FromJSON(data []byte) error {
	var temp Template
	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("JSON 파싱 실패: %v", err)
	}
	*t = temp
	return nil
} 