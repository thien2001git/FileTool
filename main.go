package main

import (
	"file_tool/domain/entity"
	"file_tool/provider"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func setup(container provider.Container) {
	// Sửa errcheck: Kiểm tra lỗi trả về từ SetConfigValue
	if err := container.SetConfigUseCase.SetConfigValue(entity.CurrDir, getCurrDir()); err != nil {
		log.Printf("Warning: failed to set config value: %v", err)
	}
}

func getCurrDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return ""
	}
	return dir
}

func main() {
	// Sửa gosimple (S1021): Gộp khai báo biến và gán giá trị làm một dòng
	container := provider.NewContainer()
	setup(*container)

	// Tạo Command gốc (Root Command)
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "App là một công cụ CLI siêu nhỏ",
		Long:  `Một ứng dụng CLI mẫu được xây dựng bằng Cobra trong Golang.`,
	}

	initGetConfigCmd(rootCmd, container)
	initWriteCoreLibCmd(rootCmd, container)
	initGenCleanDirsCmd(rootCmd, container)
	initGenCleanFilesCmd(rootCmd, container)

	// Thực thi ứng dụng
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initGenCleanFilesCmd(rootCmd *cobra.Command, container *provider.Container) {
	var path string
	var name string
	var fwStr string

	cmd := &cobra.Command{
		Use:   "gen-files",
		Short: "Generate clean files for a specific framework",
		Long:  `This command generates clean structure files based on the provided framework, path, and file name.`,
		Run: func(cmd *cobra.Command, args []string) {
			useCase := container.GenCleanFilesUseCase
			frameworkEntity := entity.Framework(fwStr)

			err := useCase.GenCleanFiles(frameworkEntity, path, name)
			if err != nil {
				log.Fatalf("Error generating clean files: %v", err)
			}

			fmt.Println("Clean files generated successfully!")
		},
	}

	// Định nghĩa các Flags cho Command
	cmd.Flags().StringVarP(&path, "path", "p", "", "Target directory path (required)")
	cmd.Flags().StringVarP(&name, "name", "n", "", "File name to generate (required)")
	cmd.Flags().StringVarP(&fwStr, "framework", "f", "", "Framework type (e.g., gin, fiber) (required)")

	// Sửa errcheck: Kiểm tra lỗi trả về từ MarkFlagRequired
	if err := cmd.MarkFlagRequired("path"); err != nil {
		log.Fatalf("Fatal error setting flag 'path' as required: %v", err)
	}
	if err := cmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("Fatal error setting flag 'name' as required: %v", err)
	}
	if err := cmd.MarkFlagRequired("framework"); err != nil {
		log.Fatalf("Fatal error setting flag 'framework' as required: %v", err)
	}

	// Add command vào rootCmd
	rootCmd.AddCommand(cmd)
}

func initGenCleanDirsCmd(rootCmd *cobra.Command, container *provider.Container) {
	var path string
	var fwStr string

	cmd := &cobra.Command{
		Use:   "gen-dirs",
		Short: "Generate clean directories for a specific framework",
		Long:  `This command generates a clean directory structure based on the provided framework and target path.`,
		Run: func(cmd *cobra.Command, args []string) {
			useCase := container.GenCleanDirsUseCase
			frameworkEntity := entity.Framework(fwStr)

			resultPath, err := useCase.GenCleanDirs(frameworkEntity, path)
			if err != nil {
				log.Fatalf("Error generating clean directories: %v", err)
			}

			fmt.Printf("Clean directories created successfully at: %s\n", resultPath)
		},
	}

	// Định nghĩa các Flags cho Command
	cmd.Flags().StringVarP(&path, "path", "p", "", "Target directory path (required)")
	cmd.Flags().StringVarP(&fwStr, "framework", "f", "", "Framework type (e.g., gin, fiber) (required)")

	// Sửa errcheck: Kiểm tra lỗi trả về từ MarkFlagRequired
	if err := cmd.MarkFlagRequired("path"); err != nil {
		log.Fatalf("Fatal error setting flag 'path' as required: %v", err)
	}
	if err := cmd.MarkFlagRequired("framework"); err != nil {
		log.Fatalf("Fatal error setting flag 'framework' as required: %v", err)
	}

	// Add command vào rootCmd
	rootCmd.AddCommand(cmd)
}

func initWriteCoreLibCmd(rootCmd *cobra.Command, container *provider.Container) {
	writeCoreLibCmd := &cobra.Command{
		Use:   "writecorelibs",
		Short: "Write core libraries",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			success := container.WriteCoreLibsUseCase.WriteCoreLibs()
			if !success {
				fmt.Println("Error writing core libraries")
				return
			}
			fmt.Println("Core libraries written successfully")
		},
	}
	rootCmd.AddCommand(writeCoreLibCmd)
}

func initGetConfigCmd(rootCmd *cobra.Command, container *provider.Container) {
	var getAll bool

	getConfigCmd := &cobra.Command{
		Use:   "getconfig",
		Short: "Get configuration value",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if getAll {
				allConfigs, err := container.GetAllConfigUseCase.GetAllConfigValues()
				if err != nil {
					fmt.Println("Error fetching all configs:", err)
					return
				}

				fmt.Println("--- All Configurations ---")
				for key, val := range allConfigs {
					fmt.Printf("%s: %s\n", key, val)
				}
				return
			}

			val, err := container.GetConfigUseCase.GetConfigValue(entity.CurrDir)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(val)
		},
	}

	getConfigCmd.Flags().BoolVarP(&getAll, "all", "a", false, "Get all configuration values instead of a single one")
	rootCmd.AddCommand(getConfigCmd)
}