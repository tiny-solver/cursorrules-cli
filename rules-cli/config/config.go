package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configDir  = ".cursorrules"
	configFile = "config-cli.json"
)

// InitConfig 설정 파일 초기화
func InitConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("홈 디렉토리를 찾을 수 없습니다: %v", err)
	}

	configPath := filepath.Join(home, configDir)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("설정 디렉토리를 생성할 수 없습니다: %v", err)
	}

	viper.SetConfigName(configFile[:len(configFile)-5]) // .json 확장자 제거
	viper.SetConfigType("json")
	viper.AddConfigPath(configPath)

	// 설정 파일이 없으면 생성
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return viper.WriteConfigAs(filepath.Join(configPath, configFile))
		}
		return fmt.Errorf("설정 파일을 읽을 수 없습니다: %v", err)
	}

	return nil
}

// SaveToken GitHub 토큰 저장
func SaveToken(token string) error {
	viper.Set("github_token", token)
	return viper.WriteConfig()
}

// GetToken 저장된 GitHub 토큰 조회
func GetToken() string {
	return viper.GetString("github_token")
}

// GetConfigPath 설정 파일 경로 조회
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
} 