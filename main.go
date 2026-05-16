package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"file_tool/provider"
	"file_tool/domain/entity"
)

func setup(container provider.Container) {
	container.SetConfigUseCase.SetConfigValue(entity.CurrDir, getCurrDir())
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
	var container *provider.Container
	// Khởi tạo Container và các thành phần cần thiết
	container = provider.NewContainer()
	setup(*container)

	// 2. Tạo Command gốc (Root Command)
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "App là một công cụ CLI siêu nhỏ",
		Long:  `Một ứng dụng CLI mẫu được xây dựng bằng Cobra trong Golang.`,
	}

	initGetConfigCmd(rootCmd, container)
	initWriteCoreLibCmd(rootCmd, container)

	// 6. Thực thi ứng dụng
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
			// 3. Kiểm tra nếu người dùng truyền cờ --all
			if getAll {
				allConfigs, err := container.GetAllConfigUseCase.GetAllConfigValues()
				if err != nil {
					fmt.Println("Error fetching all configs:", err)
					return
				}
				
				// Duyệt và in ra toàn bộ map cấu hình (Key-Value)
				fmt.Println("--- All Configurations ---")
				for key, val := range allConfigs {
					fmt.Printf("%s: %s\n", key, val)
				}
				return
			}

			// 4. Logic cũ khi không có cờ --all (Lấy 1 config cụ thể)
			val, err := container.GetConfigUseCase.GetConfigValue(entity.CurrDir)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(val)
		},
	}

	// 2. Gắn cờ (Flag) vào getConfigCmd
	// "all" là tên cờ dài (--all)
	// "a" là tên viết tắt (-a)
	// false là giá trị mặc định nếu người dùng không gõ cờ này
	getConfigCmd.Flags().BoolVarP(&getAll, "all", "a", false, "Get all configuration values instead of a single one")

	rootCmd.AddCommand(getConfigCmd)
}