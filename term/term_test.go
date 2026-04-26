package term

import (
	"os"
	"testing"
)

// TestIsStdinPipeWithError 测试 IsStdinPipeWithError 函数
//
// 注意: 在终端直接运行测试时，stdin 是终端，返回 false
// 在管道中运行测试时，stdin 是管道，返回 true
func TestIsStdinPipeWithError(t *testing.T) {
	// 测试函数能正常执行而不 panic
	isPipe, err := IsStdinPipeWithError()

	// 错误应该为 nil（因为 stdin 总是可以 stat）
	if err != nil {
		t.Errorf("IsStdinPipeWithError() 返回错误: %v", err)
	}

	// 记录当前环境
	t.Logf("IsStdinPipeWithError() = %v, err = %v", isPipe, err)
	t.Logf("当前环境: stdin 是终端=%v, 是管道=%v", !isPipe, isPipe)
}

// TestIsStdinPipe 测试 IsStdinPipe 函数
//
// 注意: 在终端直接运行测试时返回 false，在管道中运行返回 true
func TestIsStdinPipe(t *testing.T) {
	// 测试函数能正常执行而不 panic
	isPipe := IsStdinPipe()

	// 记录当前环境
	t.Logf("IsStdinPipe() = %v", isPipe)
	t.Logf("当前环境: stdin 是终端=%v, 是管道=%v", !isPipe, isPipe)
}

// TestGetSafeTerminalWidth 测试 GetSafeTerminalWidth 函数
func TestGetSafeTerminalWidth(t *testing.T) {
	width := GetSafeTerminalWidth()

	// 检查返回值在合理范围内 [40, 1200]
	if width < 40 {
		t.Errorf("GetSafeTerminalWidth() = %d, 小于最小值 40", width)
	}
	if width > 1200 {
		t.Errorf("GetSafeTerminalWidth() = %d, 大于最大值 1200", width)
	}

	t.Logf("GetSafeTerminalWidth() = %d", width)
}

// TestGetSafeTerminalWidth_Default 测试默认宽度
//
// 在非终端环境下应该返回默认值 80
func TestGetSafeTerminalWidth_Default(t *testing.T) {
	width := GetSafeTerminalWidth()

	// 在非终端环境（如 CI）中，应该返回默认值
	// 但我们不能确定测试环境，所以只检查范围
	if width < 40 || width > 1200 {
		t.Errorf("GetSafeTerminalWidth() = %d, 超出有效范围 [40, 1200]", width)
	}

	t.Logf("终端宽度: %d", width)
}

// TestIsStdinPipe_Consistency 测试两个函数的一致性
//
// IsStdinPipe 应该与 IsStdinPipeWithError 返回相同的布尔值
func TestIsStdinPipe_Consistency(t *testing.T) {
	isPipe1 := IsStdinPipe()
	isPipe2, err := IsStdinPipeWithError()

	if err != nil {
		t.Errorf("IsStdinPipeWithError() 返回错误: %v", err)
	}

	if isPipe1 != isPipe2 {
		t.Errorf("两个函数返回不一致: IsStdinPipe()=%v, IsStdinPipeWithError()=%v", isPipe1, isPipe2)
	}
}

// TestGetSafeTerminalWidth_Environment 测试 COLUMNS 环境变量
//
// 设置 COLUMNS 环境变量后，函数应该返回该值
func TestGetSafeTerminalWidth_Environment(t *testing.T) {
	// 保存原始值
	originalCols := os.Getenv("COLUMNS")
	defer func() {
		_ = os.Setenv("COLUMNS", originalCols)
	}()

	tests := []struct {
		name      string
		columns   string
		wantMin   int
		wantMax   int
		wantExact int
	}{
		{
			name:    "有效值 100",
			columns: "100",
			wantMin: 100,
			wantMax: 100,
		},
		{
			name:    "边界值 40",
			columns: "40",
			wantMin: 40,
			wantMax: 40,
		},
		{
			name:    "边界值 1200",
			columns: "1200",
			wantMin: 1200,
			wantMax: 1200,
		},
		{
			name:    "无效值 太小",
			columns: "30",
			wantMin: 40, // 应该被限制到最小值
			wantMax: 1200,
		},
		{
			name:    "无效值 太大",
			columns: "2000",
			wantMin: 40,
			wantMax: 1200, // 应该被限制到最大值
		},
		{
			name:    "无效值 非数字",
			columns: "abc",
			wantMin: 40,
			wantMax: 1200, // 应该忽略，使用其他方式获取
		},
		{
			name:    "空值",
			columns: "",
			wantMin: 40,
			wantMax: 1200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Setenv("COLUMNS", tt.columns); err != nil {
				t.Fatalf("设置环境变量失败: %v", err)
			}
			got := GetSafeTerminalWidth()

			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("GetSafeTerminalWidth() = %d, want between %d and %d", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

// BenchmarkIsStdinPipe 基准测试 IsStdinPipe
func BenchmarkIsStdinPipe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsStdinPipe()
	}
}

// BenchmarkIsStdinPipeWithError 基准测试 IsStdinPipeWithError
func BenchmarkIsStdinPipeWithError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = IsStdinPipeWithError()
	}
}

// BenchmarkGetSafeTerminalWidth 基准测试 GetSafeTerminalWidth
func BenchmarkGetSafeTerminalWidth(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetSafeTerminalWidth()
	}
}
