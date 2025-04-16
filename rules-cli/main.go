package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tinysolver/rules-cli/config"
)

var rootCmd = &cobra.Command{
	Use:   "cursorrules",
	Short: "Cursor Rules CLI",
	Long:  "Cursor Rules CLI는 터미널에서 Cursor rules 파일을 관리하는 도구입니다.",
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "GitHub Personal Access Token 설정",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.InitConfig(); err != nil {
			fmt.Printf("설정 초기화 실패: %v\n", err)
			return
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("GitHub Personal Access Token을 입력하세요: ")
		token, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("토큰 입력 실패: %v\n", err)
			return
		}

		token = strings.TrimSpace(token)
		if err := config.SaveToken(token); err != nil {
			fmt.Printf("토큰 저장 실패: %v\n", err)
			return
		}

		configPath, err := config.GetConfigPath()
		if err != nil {
			fmt.Printf("설정 파일 경로 조회 실패: %v\n", err)
			return
		}

		fmt.Printf("토큰이 %s에 저장되었습니다.\n", configPath)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "템플릿 목록 출력",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("템플릿 목록 출력...")
		// TODO: 목록 출력 로직 구현
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync [name]",
	Short: "템플릿 다운로드",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("템플릿 %s 다운로드 중...\n", args[0])
		// TODO: 다운로드 로직 구현
	},
}

var pushCmd = &cobra.Command{
	Use:   "push [name]",
	Short: "로컬 템플릿 업로드",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("템플릿 %s 업로드 중...\n", args[0])
		// TODO: 업로드 로직 구현
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(pushCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 