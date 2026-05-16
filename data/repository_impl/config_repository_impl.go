package repository_impl

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"file_tool/domain/entity"
	"file_tool/domain/repository"
)

type configRepositoryImpl struct{}

var cache = make(map[entity.ConfigKey]string)

func NewConfigRepository() repository.ConfigRepository {
	return &configRepositoryImpl{}
}

func (r *configRepositoryImpl) GetConfigValue(key entity.ConfigKey) (string, error) {
	// Implementation for retrieving configuration value
	if len(cache) == 0 {
		cachedData, err := LoadCacheFromFile("config_cache.json")
		if err != nil {
			return "", fmt.Errorf("lỗi tải cache từ file: %w", err)
		}
		cache = cachedData
	}
	return cache[key], nil
}

func (r *configRepositoryImpl) SetConfigValue(key entity.ConfigKey, value string) error {
	// Implementation for setting configuration value
	cache[key] = value
	SaveCacheToFile(cache[entity.CurrDir] + "/config/config_cache.json", cache)
	return nil
}

// GetAllConfigValues implements [repository.ConfigRepository].
func (r *configRepositoryImpl) GetAllConfigValues() (map[entity.ConfigKey]string, error) {
	return cache, nil
}

func SaveCacheToFile(filePath string, cache map[entity.ConfigKey]string) error {
	// 1. Đảm bảo thư mục cha của file tồn tại (Tránh lỗi nếu folder chưa được tạo)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("không thể tạo thư mục %s: %w", dir, err)
	}

	// 2. Chuyển đổi map thành dữ liệu JSON dạng byte[]
	// MarshalIndent giúp file JSON được xuống dòng và thụt lề đẹp mắt, dễ đọc bằn mắt thường
	jsonData, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("lỗi tuần tự hóa JSON: %w", err)
	}

	// 3. Ghi dữ liệu byte vào file
	// Quy định quyền 0644 (Chủ sở hữu có quyền Đọc/Ghi, người khác chỉ có quyền Đọc)
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("lỗi ghi file: %w", err)
	}

	return nil
}

func LoadCacheFromFile(filePath string) (map[entity.ConfigKey]string, error) {
	cache := make(map[entity.ConfigKey]string)

	// 1. Đọc toàn bộ file byte
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 2. Giải mã JSON ngược lại vào cấu trúc map
	err = json.Unmarshal(jsonData, &cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}
