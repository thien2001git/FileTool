package repository_impl

import (
	"bytes"
	"encoding/json"
	"file_tool/domain/entity"
	"file_tool/domain/repository"
	"file_tool/domain/usecase/config"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type coreLibRepositoryImpl struct {
	getConfigUseCase config.GetConfigUseCase
}

func NewCoreLibRepository(getConfigUc config.GetConfigUseCase) repository.CoreLibRepository {
	return &coreLibRepositoryImpl{
		getConfigUseCase: getConfigUc,
	}
}

func (r *coreLibRepositoryImpl) GetCurrentCoreLibs() ([]entity.CoreLib, error) {
	var coreLibs []entity.CoreLib

	// 2. Lấy đường dẫn file json cấu hình từ hàm helper
	jsonPath, err := r.getCoreLibPath()
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy đường dẫn file cấu hình: %w", err)
	}

	// 3. Đọc dữ liệu từ file config/core_lib_keys.json
	jsonFile, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("không thể đọc file %s: %w", jsonPath, err)
	}

	// Giả định file JSON có cấu trúc dạng mảng: ["flutter", "go", "firebase-cli"]
	var targetKeys []string
	if err := json.Unmarshal(jsonFile, &targetKeys); err != nil {
		return nil, fmt.Errorf("lỗi phân tích file JSON: %w", err)
	}

	// 4. Ánh xạ các key từ file JSON sang câu lệnh Terminal tương ứng
	// Việc này giúp bạn quản lý tập trung, nếu JSON thiếu phần tử nào thì vòng lặp chỉ chạy phần tử có sẵn
	cmdMapping := map[string]struct {
		techName string
		args     []string
	}{
		"flutter":      {techName: "Flutter", args: []string{"flutter", "--version"}},
		"go":           {techName: "Go", args: []string{"go", "version"}},
		"firebase-cli": {techName: "Firebase-CLI", args: []string{"firebase", "--version"}},
	}

	// 5. Duyệt qua danh sách lấy được từ file JSON để thực thi lệnh
	for _, key := range targetKeys {
		// Chuẩn hóa chữ thường để so sánh chính xác
		loweredKey := strings.ToLower(strings.TrimSpace(key))

		cmdInfo, exists := cmdMapping[loweredKey]
		if !exists {
			// Nếu trong file JSON có một công cụ lạ chưa được định nghĩa lệnh, bỏ qua hoặc log ra
			continue
		}

		// Khởi tạo lệnh dựa trên mapping
		cmd := exec.Command(cmdInfo.args[0], cmdInfo.args[1:]...)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		err := cmd.Run()

		lib := entity.CoreLib{
			Name:    cmdInfo.techName,
			Version: "Not Installed",
		}

		if err == nil {
			outputStr := strings.TrimSpace(out.String())
			lib.Version = parseVersionOutput(cmdInfo.techName, outputStr)

		}

		coreLibs = append(coreLibs, lib)
	}

	return coreLibs, nil
}

func (r *coreLibRepositoryImpl) WriteCoreLibs() bool {
	// 1. Lấy đường dẫn file cần ghi từ cấu hình
	path, err := r.GetCoreLibWritePath()
	if err != nil {
		return false
	}

	// 2. Lấy danh sách core libraries hiện tại qua Terminal
	coreLibs, err := r.GetCurrentCoreLibs()
	if err != nil {
		return false
	}

	// 3. Tuần tự hóa danh sách struct thành dữ liệu JSON (định dạng đẹp mắt với thụt lề)
	jsonData, err := json.MarshalIndent(coreLibs, "", "  ")
	if err != nil {
		return false
	}

	// 4. Đảm bảo thư mục đích tồn tại trước khi ghi file
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return false
	}

	// 5. Ghi trực tiếp mảng JSON thu được xuống file vật lý
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return false
	}

	return true
}

func (r *coreLibRepositoryImpl) getCoreLibPath() (string, error) {
	rootPath, err := r.getConfigUseCase.GetConfigValue(entity.CurrDir)
	if err != nil {
		return "", err
	}
	// ✅ Thay vì cộng chuỗi "/", dùng filepath.Join tự tương thích Windows/Mac/Linux
	return filepath.Join(rootPath, "config", "core_lib_keys.json"), nil
}

func (r *coreLibRepositoryImpl) GetCoreLibWritePath() (string, error) {
	rootPath, err := r.getConfigUseCase.GetConfigValue(entity.CurrDir)
	if err != nil {
		return "", err
	}
	// ✅ Thay vì cộng chuỗi "/", dùng filepath.Join tự tương thích Windows/Mac/Linux
	return filepath.Join(rootPath, "..", "core_lib.json"), nil
}

func parseVersionOutput(_, output string) string {
	if output == "" {
		return "Installed"
	}
	// Lấy dòng đầu tiên của kết quả trả về
	lines := strings.Split(output, "\n")
	firstLine := strings.TrimSpace(lines[0])

	// Bạn có thể viết thêm logic regex tại đây nếu chỉ muốn lấy đúng cụm số (VD: 3.41.9)
	return firstLine
}
