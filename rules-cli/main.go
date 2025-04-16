package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cursorrules",
	Short: "Cursor Rules CLI",
	Long:  "Cursor Rules CLI는 터미널에서 Cursor rules 파일을 관리하는 도구입니다.",
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "GitHub OAuth 인증",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("인증 시작...")
		// TODO: 인증 로직 구현
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