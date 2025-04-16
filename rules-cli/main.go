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
		if len(args) != 1 {
			fmt.Println("템플릿 이름을 지정해주세요")
			return
		}
		templateName := args[0]

		client, err := gist.NewGistClient()
		if err != nil {
			fmt.Printf("Gist 클라이언트 생성 실패: %v\n", err)
			return
		}

		// Gist에서 템플릿 다운로드
		gistObj, err := client.FindGistByDescription(templateName)
		if err != nil {
			fmt.Printf("템플릿 다운로드 실패: %v\n", err)
			return
		}

		contents, err := client.GetGistContent(gistObj.GetID())
		if err != nil {
			fmt.Printf("템플릿 내용 조회 실패: %v\n", err)
			return
		}

		// 템플릿 생성
		template := &models.Template{
			Files: make(map[string]models.Rule),
		}

		// 각 파일을 템플릿에 추가
		for filename, content := range contents {
			// 파일 구조 보존을 위해 경로를 키로 사용
			rule := models.Rule{
				Name:    filename,
				Content: content,
				Path:    filename,
			}
			template.Files[filename] = rule
		}

		// 로컬에 저장
		if err := filesystem.SaveLocalTemplate(template, nil); err != nil {
			fmt.Printf("템플릿 저장 실패: %v\n", err)
			return
		}

		fmt.Printf("템플릿 '%s'이(가) 성공적으로 다운로드되었습니다.\n", templateName)
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload [name]",
	Short: "로컬 템플릿 업로드",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("템플릿 이름을 지정해주세요")
			return
		}
		templateName := args[0]

		client, err := gist.NewGistClient()
		if err != nil {
			fmt.Printf("Gist 클라이언트 생성 실패: %v\n", err)
			return
		}

		// 로컬 템플릿 로드
		localTemplate, _, err := filesystem.LoadLocalTemplate()
		if err != nil {
			fmt.Printf("로컬 템플릿 로드 실패: %v\n", err)
			return
		}

		// Gist에서 기존 템플릿 확인
		gist, err := client.FindGistByDescription(templateName)
		if err == nil {
			// 기존 템플릿이 있는 경우 버전 비교
			contents, err := client.GetGistContent(gist.GetID())
			if err != nil {
				fmt.Printf("템플릿 내용 조회 실패: %v\n", err)
				return
			}

			needsUpdate := false
			for filename, content := range contents {
				localRule, exists := localTemplate.Files[filename]
				if !exists || localRule.Content != content {
					needsUpdate = true
					fmt.Printf("업데이트 필요: %s\n", filename)
				}
			}

			if !needsUpdate {
				fmt.Println("모든 파일이 최신 상태입니다.")
				return
			}
		}

		// 템플릿 업로드
		files := make(map[string]string)
		for _, rule := range localTemplate.Files {
			files[rule.Name] = rule.Content
		}

		if gist == nil {
			// 새 Gist 생성
			_, err = client.CreateGist(templateName, files)
		} else {
			// 기존 Gist 업데이트
			if err := client.DeleteGist(gist.GetID()); err != nil {
				fmt.Printf("기존 템플릿 삭제 실패: %v\n", err)
				return
			}
			_, err = client.CreateGist(templateName, files)
		}

		if err != nil {
			fmt.Printf("템플릿 업로드 실패: %v\n", err)
			return
		}

		fmt.Printf("템플릿 '%s'이(가) 성공적으로 업로드되었습니다.\n", templateName)
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
	downloadCmd.Flags().BoolP("merge", "m", false, "로컬 파일과 병합")
	deleteCmd.Flags().BoolP("force", "f", false, "확인 없이 강제 삭제")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
} 