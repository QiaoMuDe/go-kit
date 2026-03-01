package term

import (
	"strings"
	"testing"
)

// TestRenderBasicMenu 测试基础菜单渲染
func TestRenderBasicMenu(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		options []string
		style   *MenuStyle
	}{
		{
			name:    "基本渲染",
			title:   "主菜单",
			options: []string{"查看列表", "添加项目"},
		},
		{
			name:    "带标题",
			title:   "测试菜单",
			options: []string{"选项1"},
		},
		{
			name:    "不带标题",
			title:   "",
			options: []string{"选项1"},
		},
		{
			name:    "空选项列表",
			title:   "主菜单",
			options: []string{},
		},
		{
			name:    "自定义样式",
			title:   "主菜单",
			options: []string{"选项1"},
			style: &MenuStyle{
				Prefix:    ">> ",
				Separator: " | ",
				Indent:    2,
				ShowTitle: true,
			},
		},
		{
			name:    "多选项",
			title:   "主菜单",
			options: []string{"选项1", "选项2", "选项3", "选项4", "选项5"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderBasicMenu(tt.title, tt.options, tt.style)
			if result == "" && len(tt.options) > 0 {
				t.Error("RenderBasicMenu() 返回空字符串")
			}
		})
	}
}

// TestShowBasicMenu 测试基础菜单显示
func TestShowBasicMenu(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		options []string
		input   string
		wantIdx int
		wantErr bool
	}{
		{
			name:    "正常选择",
			title:   "主菜单",
			options: []string{"查看列表", "添加项目", "删除项目"},
			input:   "1\n",
			wantIdx: 0,
			wantErr: false,
		},
		{
			name:    "选择最后一个",
			title:   "主菜单",
			options: []string{"选项1", "选项2", "选项3"},
			input:   "3\n",
			wantIdx: 2,
			wantErr: false,
		},
		{
			name:    "空输入",
			title:   "主菜单",
			options: []string{"选项1"},
			input:   "\n",
			wantIdx: 0,
			wantErr: true,
		},
		{
			name:    "非数字输入",
			title:   "主菜单",
			options: []string{"选项1"},
			input:   "abc\n",
			wantIdx: 0,
			wantErr: true,
		},
		{
			name:    "部分数字",
			title:   "主菜单",
			options: []string{"选项1"},
			input:   "1abc\n",
			wantIdx: 0,
			wantErr: true,
		},
		{
			name:    "超出范围",
			title:   "主菜单",
			options: []string{"选项1", "选项2", "选项3"},
			input:   "99\n",
			wantIdx: 0,
			wantErr: true,
		},
		{
			name:    "小于1",
			title:   "主菜单",
			options: []string{"选项1", "选项2"},
			input:   "0\n",
			wantIdx: 0,
			wantErr: true,
		},
		{
			name:    "负数",
			title:   "主菜单",
			options: []string{"选项1"},
			input:   "-1\n",
			wantIdx: 0,
			wantErr: true,
		},
		{
			name:    "浮点数",
			title:   "主菜单",
			options: []string{"选项1"},
			input:   "1.5\n",
			wantIdx: 0,
			wantErr: true,
		},
		{
			name:    "空选项列表",
			title:   "主菜单",
			options: []string{},
			input:   "1\n",
			wantIdx: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			idx, err := ShowBasicMenu(tt.title, tt.options, reader, "请选择: ")
			if (err != nil) != tt.wantErr {
				t.Errorf("ShowBasicMenu() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && idx != tt.wantIdx {
				t.Errorf("ShowBasicMenu() = %v, want %v", idx, tt.wantIdx)
			}
		})
	}
}

// TestShowBasicMenuLoop 测试基础菜单循环
func TestShowBasicMenuLoop(t *testing.T) {
	t.Run("handler返回false", func(t *testing.T) {
		reader := strings.NewReader("1\n")

		err := ShowBasicMenuLoop("主菜单", []string{"选项1"}, reader, "请选择: ", func(index int) bool {
			return false
		})

		if err != nil {
			t.Errorf("ShowBasicMenuLoop() error = %v", err)
		}
	})

	t.Run("正常循环一次", func(t *testing.T) {
		reader := strings.NewReader("1\n")
		count := 0

		err := ShowBasicMenuLoop("主菜单", []string{"选项1", "选项2"}, reader, "请选择: ", func(index int) bool {
			count++
			return false // 只循环一次
		})

		if err != nil {
			t.Errorf("ShowBasicMenuLoop() error = %v", err)
		}
		if count != 1 {
			t.Errorf("循环次数 = %v, want 1", count)
		}
	})

	t.Run("空选项列表", func(t *testing.T) {
		reader := strings.NewReader("1\n")

		err := ShowBasicMenuLoop("主菜单", []string{}, reader, "请选择: ", func(index int) bool {
			return true
		})

		if err != ErrEmptyMenu {
			t.Errorf("ShowBasicMenuLoop() error = %v, want ErrEmptyMenu", err)
		}
	})
}

// TestShowBasicMenuLine 测试便捷函数
func TestShowBasicMenuLine(t *testing.T) {
	t.Skip("跳过 TestShowBasicMenuLine: 便捷函数使用 os.Stdin，难以在单元测试中模拟输入")
}

// TestShowBasicMenuLoopLine 测试循环便捷函数
func TestShowBasicMenuLoopLine(t *testing.T) {
	t.Skip("跳过 TestShowBasicMenuLoopLine: 便捷函数使用 os.Stdin，难以在单元测试中模拟输入")
}
