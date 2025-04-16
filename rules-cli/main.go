package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tinysolver/rules-cli/config"
	"github.com/tinysolver/rules-cli/gist"
	"github.com/tinysolver/rules-cli/filesystem"
	"github.com/tinysolver/rules-cli/models"
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
		client, err := gist.NewGistClient()
		if err != nil {
			fmt.Printf("Gist 클라이언트 생성 실패: %v\n", err)
			return
		}

		gists, err := client.ListGists()
		if err != nil {
			fmt.Printf("Gist 목록 조회 실패: %v\n", err)
			return
		}

		if len(gists) == 0 {
			fmt.Println("저장된 템플릿이 없습니다.")
			return
		}

		fmt.Println("저장된 템플릿 목록:")
		for _, g := range gists {
			name := gist.GetProjectName(g.GetDescription())
			if name == "" {
				name = "(이름 없음)"
			}
			fmt.Printf("- %s\n", name)
		}
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download [name]",
	Short: "템플릿 다운로드",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		force, _ := cmd.Flags().GetBool("force")

		client, err := gist.NewGistClient()
		if err != nil {
			fmt.Printf("Gist 클라이언트 생성 실패: %v\n", err)
			return
		}

		gist, err := client.FindGistByDescription(projectName)
		if err != nil {
			fmt.Printf("프로젝트 '%s'를 찾을 수 없습니다: %v\n", projectName, err)
			return
		}

		contents, err := client.GetGistContent(gist.GetID())
		if err != nil {
			fmt.Printf("템플릿 내용 조회 실패: %v\n", err)
			return
		}

		// JSON 데이터를 템플릿으로 변환
		template, err := models.FromJSON([]byte(contents["template.json"]))
		if err != nil {
			fmt.Printf("템플릿 변환 실패: %v\n", err)
			return
		}

		// 충돌 확인
		conflicts, err := filesystem.CheckConflicts(template)
		if err != nil {
			fmt.Printf("충돌 확인 실패: %v\n", err)
			return
		}

		if len(conflicts) > 0 {
			if force {
				fmt.Println("강제로 덮어쓰기 모드로 진행합니다...")
			} else {
				fmt.Println("다음 파일들이 충돌합니다:")
				for _, conflict := range conflicts {
					fmt.Printf("  - %s\n", conflict)
				}
				fmt.Println("\n다음 옵션 중 하나를 선택하세요:")
				fmt.Println("  --force: 강제로 덮어쓰기")
				return
			}
		}

		// 템플릿 저장
		version := models.NewTemplateVersion(projectName, "v1.0.0")
		if err := filesystem.MergeTemplate(template, version); err != nil {
			fmt.Printf("템플릿 저장 실패: %v\n", err)
			return
		}

		fmt.Printf("프로젝트 '%s'의 템플릿이 다운로드되었습니다.\n", projectName)
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload [name]",
	Short: "로컬 템플릿 업로드",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		client, err := gist.NewGistClient()
		if err != nil {
			fmt.Printf("Gist 클라이언트 생성 실패: %v\n", err)
			return
		}

		template, version, err := filesystem.LoadLocalTemplate()
		if err != nil {
			fmt.Printf("로컬 템플릿 로드 실패: %v\n", err)
			return
		}

		// 템플릿을 JSON으로 변환
		jsonData, err := template.ToJSON()
		if err != nil {
			fmt.Printf("JSON 변환 실패: %v\n", err)
			return
		}

		// 버전 정보를 JSON으로 변환
		versionData, err := version.ToJSONString()
		if err != nil {
			fmt.Printf("버전 정보 변환 실패: %v\n", err)
			return
		}

		// Gist에 업로드
		files := map[string]string{
			"template.json": string(jsonData),
			"version.json": versionData,
		}

		_, err = client.CreateGist(projectName, files)
		if err != nil {
			fmt.Printf("템플릿 업로드 실패: %v\n", err)
			return
		}

		fmt.Printf("프로젝트 '%s'의 템플릿이 업로드되었습니다.\n", projectName)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "템플릿 삭제",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		force, _ := cmd.Flags().GetBool("force")

		client, err := gist.NewGistClient()
		if err != nil {
			fmt.Printf("Gist 클라이언트 생성 실패: %v\n", err)
			return
		}

		gist, err := client.FindGistByDescription(projectName)
		if err != nil {
			fmt.Printf("프로젝트 '%s'를 찾을 수 없습니다: %v\n", projectName, err)
			return
		}

		if !force {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("정말로 '%s' 템플릿을 삭제하시겠습니까? (y/N): ", projectName)
			answer, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(answer)) != "y" {
				fmt.Println("삭제가 취소되었습니다.")
				return
			}
		}

		if err := client.DeleteGist(gist.GetID()); err != nil {
			fmt.Printf("템플릿 삭제 실패: %v\n", err)
			return
		}

		fmt.Printf("프로젝트 '%s'의 템플릿이 삭제되었습니다.\n", projectName)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(deleteCmd)

	downloadCmd.Flags().BoolP("force", "f", false, "강제로 덮어쓰기")
	deleteCmd.Flags().BoolP("force", "f", false, "확인 없이 강제 삭제")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 