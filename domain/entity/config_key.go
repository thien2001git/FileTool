package entity

type ConfigKey int

const (
	CurrDir ConfigKey = iota
	// Add more configuration keys as needed
)

func (k ConfigKey) String() string {
	// Khai báo một mảng các chuỗi ký tự tương ứng chính xác với thứ tự hằng số ở trên
	keys := [...]string{
		"CurrDir",
	}

	// Kiểm tra tránh lỗi vượt quá chỉ mục mảng (Index out of range) 
	// Nếu k là một số âm hoặc lớn hơn số lượng phần tử có trong mảng
	if k < 0 || int(k) >= len(keys) {
		return "UnknownConfigKey"
	}

	return keys[k]
}