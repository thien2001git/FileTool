package repository_impl

import (
	"bytes"
	"file_tool/domain/entity"
	"file_tool/domain/repository"
	"file_tool/domain/usecase/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type genCleanFilesRepositoryImpl struct {
	getConfigUseCase config.GetConfigUseCase
}

func NewGenCleanFilesRepository(getConfigUc config.GetConfigUseCase) repository.GenCleanFilesRepository {
	return &genCleanFilesRepositoryImpl{
		getConfigUseCase: getConfigUc,
	}
}

// Hàm bổ trợ lấy danh sách các thư mục con cần tạo theo từng framework
func getDirsForFramework(fw entity.Framework) []string {
	switch fw {
	case entity.FrameworkGo:
		return []string{
			"cmd/api",
			"internal/domain",
			"internal/usecase",
			"internal/repository",
			"internal/delivery/http",
			"config",
		}
	case entity.FrameworkFlutter:
		return []string{
			"lib/domain/entity",
			"lib/domain/usecase",
			"lib/domain/repository",
			"lib/data/repository",
			"lib/data/datasources",
			"lib/presentation/screen",
			"lib/presentation/cubit",
			"lib/shared",
		}
	case entity.FrameworkReactNative:
		return []string{
			"src/domain/entities",
			"src/domain/usecases",
			"src/data/repositories",
			"src/presentation/screens",
			"src/presentation/components",
			"src/presentation/hooks",
		}
	case entity.FrameworkNest:
		return []string{
			"src/domain",
			"src/presentation/controllers",
			"src/presentation/services",
			"src/domain/repositories",
			"src/domain/entities",
			"src/data/repositories",
			"src/domain/usecases",
		}
	case entity.FrameworkNext:
		return []string{
			"src/domain/entity",
			"src/domain/usecase",
			"src/data/repository",
			"src/domain/repository",
			"src/presentation/context",
			"src/presentation/services",
			"src/presentation/hooks",
			"src/shared",
		}
	case entity.FrameworkSpringBoot:
		return []string{
			"domain/entity",
			"domain/repository",
			"service",
			"controller",
		}
	case entity.FrameworkAndroid:
		return []string{
			"domain/entity",
			"domain/usecase",
			"data/repository",
			"domain/repository",
			"presentation/screen",
			"presentation/viewmodel",
		}
	default:
		return []string{"domain", "data", "presentation"}
	}
}

func (r *genCleanFilesRepositoryImpl) GenCleanDirs(frame_work entity.Framework, path string) (string, error) {
	subDirs := getDirsForFramework(frame_work)
	if len(subDirs) == 0 {
		return "", fmt.Errorf("không hỗ trợ framework: %s", frame_work)
	}

	for _, subDir := range subDirs {
		targetPath := filepath.Join(path, subDir)
		err := os.MkdirAll(targetPath, 0755)
		if err != nil {
			return "", fmt.Errorf("lỗi tạo thư mục %s: %w", targetPath, err)
		}
	}

	return fmt.Sprintf("Đã khởi tạo thành công %d thư mục cho %s", len(subDirs), frame_work), nil
}

func (r *genCleanFilesRepositoryImpl) getTemplateFolder(frame_work entity.Framework) (string, error) {
	rootPath, err := r.getConfigUseCase.GetConfigValue(entity.CurrDir)
	if err != nil {
		return "", err
	}
	return filepath.Join(rootPath, "lang_template", string(frame_work)), nil
}

// Cấu trúc map ánh xạ file template sang sub-folder tương ứng của từng framework
func getTemplateMapping(fw entity.Framework) map[string]string {
	switch fw {
	case entity.FrameworkGo:
		return map[string]string{
			"entity.template":          "internal/domain",
			"repository.template":      "internal/domain",
			"repository.impl.template": "internal/repository",
			"usecase.template":         "internal/usecase",
			"viewmodel.template":       "internal/delivery/http",
		}
	case entity.FrameworkFlutter:
		return map[string]string{
			"entity.template":          "lib/domain/entity",
			"repository.template":      "lib/domain/repository",
			"repository.impl.template": "lib/data/repository",
			"usecase.template":         "lib/domain/usecase",
			"viewmodel.template":       "lib/presentation/cubit",
			"cubit.template":           "lib/presentation/cubit",
		}
	case entity.FrameworkReactNative:
		return map[string]string{
			"entity.template":          "src/domain/entities",
			"repository.template":      "src/domain/repositories",
			"repository.impl.template": "src/data/repositories",
			"usecase.template":         "src/domain/usecases",
			"viewmodel.template":       "src/presentation/hooks",
		}
	case entity.FrameworkNest:
		return map[string]string{
			"entity.template":          "src/domain/entities",
			"repository.template":      "src/domain/repositories",
			"repository.impl.template": "src/data/repositories",
			"usecase.template":         "src/domain/usecases",
			"viewmodel.template":       "src/presentation/services",
		}
	case entity.FrameworkNext:
		return map[string]string{
			"entity.template":          "src/domain/entity",
			"repository.template":      "src/domain/repository",
			"repository.impl.template": "src/data/repository",
			"usecase.template":         "src/domain/usecase",
			"viewmodel.template":       "src/presentation/hooks",
		}
	case entity.FrameworkSpringBoot:
		return map[string]string{
			"entity.template":          "domain/entity",
			"repository.template":      "domain/repository",
			"repository.impl.template": "infrastructure/repository",
			"usecase.template":         "domain/usecase",
			"viewmodel.template":       "presentation/model",
		}
	case entity.FrameworkAndroid:
		return map[string]string{
			"entity.template":          "domain/entity",
			"repository.template":      "domain/repository",
			"repository.impl.template": "data/repository",
			"usecase.template":         "domain/usecase",
			"viewmodel.template":       "presentation/viewmodel",
		}
	default:
		return map[string]string{
			"entity.template":     "domain",
			"repository.template": "domain",
			"usecase.template":    "domain",
		}
	}
}

// Helper sinh tên file đúng quy chuẩn đặt tên (Naming Convention) từng ngôn ngữ
func resolveTargetFileName(fw entity.Framework, templateName string, entityName string) string {
	lowerEntity := strings.ToLower(entityName)

	switch fw {
	case entity.FrameworkGo:
		switch templateName {
		case "entity.template":
			return fmt.Sprintf("%s.go", lowerEntity)
		case "repository.template":
			return fmt.Sprintf("%s_repository.go", lowerEntity)
		case "repository.impl.template":
			return fmt.Sprintf("%s_repository_impl.go", lowerEntity)
		case "usecase.template":
			return fmt.Sprintf("%s_usecase.go", lowerEntity)
		default:
			return fmt.Sprintf("%s.go", lowerEntity)
		} // <-- ĐÃ SỬA: Bổ sung dấu đóng ngoặc switch templateName cho Go ở đây

	case entity.FrameworkFlutter:
		switch templateName {
		case "entity.template":
			return fmt.Sprintf("%s_entity.dart", lowerEntity)
		case "repository.template":
			return fmt.Sprintf("i_%s_repository.dart", lowerEntity)
		case "repository.impl.template":
			return fmt.Sprintf("%s_repository_impl.dart", lowerEntity)
		case "usecase.template":
			return fmt.Sprintf("%s_usecase.dart", lowerEntity)
		case "viewmodel.template", "cubit.template":
			return fmt.Sprintf("%s_cubit.dart", lowerEntity)
		}
	case entity.FrameworkNest, entity.FrameworkNext, entity.FrameworkReactNative:
		ext := "ts"
		if fw == entity.FrameworkNext || fw == entity.FrameworkReactNative {
			if templateName == "viewmodel.template" {
				ext = "tsx"
			}
		}
		switch templateName {
		case "entity.template":
			return fmt.Sprintf("%s.entity.%s", lowerEntity, ext)
		case "repository.template":
			return fmt.Sprintf("i-%s.repository.%s", lowerEntity, ext)
		case "repository.impl.template":
			return fmt.Sprintf("%s.repository.impl.%s", lowerEntity, ext)
		case "usecase.template":
			return fmt.Sprintf("%s.usecase.%s", lowerEntity, ext)
		case "viewmodel.template", "cubit.template":
			return fmt.Sprintf("use-%s.viewmodel.%s", lowerEntity, ext)
		}
	case entity.FrameworkSpringBoot, entity.FrameworkAndroid:
		ext := "java"
		if fw == entity.FrameworkAndroid {
			ext = "kt"
		}
		switch templateName {
		case "entity.template":
			return fmt.Sprintf("%sEntity.%s", entityName, ext)
		case "repository.template":
			return fmt.Sprintf("%sRepository.%s", entityName, ext)
		case "repository.impl.template":
			return fmt.Sprintf("%sRepositoryImpl.%s", entityName, ext)
		case "usecase.template":
			return fmt.Sprintf("Get%sUseCase.%s", entityName, ext)
		case "viewmodel.template", "cubit.template":
			return fmt.Sprintf("%sViewModel.%s", entityName, ext)
		}
	}
	return fmt.Sprintf("%s_%s", entityName, strings.Replace(templateName, ".template", "", 1))
}

// Triển khai chi tiết cấu trúc tạo file từ template
func (r *genCleanFilesRepositoryImpl) GenCleanFiles(frame_work entity.Framework, path string, name string) error {
	// 1. Lấy đường dẫn tuyệt đối đến thư mục chứa template nguồn
	templateFolder, err := r.getTemplateFolder(frame_work)
	if err != nil {
		return fmt.Errorf("không thể lấy thư mục chứa template: %w", err)
	}

	// 2. Lấy danh sách ánh xạ các file template cần sinh ra
	mappings := getTemplateMapping(frame_work)

	// 3. Chuẩn bị hàm bổ trợ (FuncMap) cho Go text/template
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	// 4. Định nghĩa Data DTO sẽ truyền vào Go template token
	templateData := struct {
		Name        string
		ModulePath  string
		PackageName string
	}{
		Name:        name,
		ModulePath:  strings.ToLower(name),
		PackageName: "com.system.generated",
	}

	// 5. Duyệt qua từng template được định nghĩa để render
	for templateName, subDir := range mappings {
		templateFilePath := filepath.Join(templateFolder, templateName)

		// Kiểm tra xem file template nguồn có thực sự tồn tại hay không
		if _, err := os.Stat(templateFilePath); os.IsNotExist(err) {
			continue 
		}

		// Đọc nội dung file template mẫu
		tmplContent, err := os.ReadFile(templateFilePath)
		if err != nil {
			return fmt.Errorf("lỗi đọc file template %s: %w", templateName, err)
		}

		// Khởi tạo engine Go template parser
		tmpl, err := template.New(templateName).Funcs(funcMap).Parse(string(tmplContent))
		if err != nil {
			return fmt.Errorf("lỗi parse cấu trúc cú pháp template %s: %w", templateName, err)
		}

		// Thực thi ghi dữ liệu động vào buffer vùng nhớ tạm
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, templateData); err != nil {
			return fmt.Errorf("lỗi render template %s với dữ liệu entity: %w", templateName, err)
		}

		// Xác định tên file đầu ra chuẩn hóa
		targetFileName := resolveTargetFileName(frame_work, templateName, name)

		// Xây dựng đường dẫn file đích chính xác trong cấu trúc Clean Architecture
		targetFileDirectory := filepath.Join(path, subDir)
		targetFullFilePath := filepath.Join(targetFileDirectory, targetFileName)

		// Đảm bảo thư mục cha chứa file tồn tại
		if err := os.MkdirAll(targetFileDirectory, 0755); err != nil {
			return fmt.Errorf("lỗi tạo thư mục đích %s: %w", targetFileDirectory, err)
		}

		// Ghi đè file đầu ra với quyền thực thi chuẩn
		if err := os.WriteFile(targetFullFilePath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("lỗi tạo file code đầu ra tại vị trí %s: %w", targetFullFilePath, err)
		}
	}

	return nil
}