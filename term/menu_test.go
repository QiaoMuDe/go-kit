package term

import (
	"strings"
	"testing"
)

// TestValidateMenu 测试菜单验证
func TestValidateMenu(t *testing.T) {
	tests := []struct {
		name    string
		menu    *Menu
		wantErr bool
	}{
		{
			name:    "空菜单",
			menu:    nil,
			wantErr: true,
		},
		{
			name: "空菜单项",
			menu: &Menu{
				Title: "测试菜单",
				Items: []MenuItem{},
			},
			wantErr: true,
		},
		{
			name: "空键",
			menu: &Menu{
				Title: "测试菜单",
				Items: []MenuItem{
					{Key: "", Value: "选项1"},
				},
			},
			wantErr: true,
		},
		{
			name: "重复键",
			menu: &Menu{
				Title: "测试菜单",
				Items: []MenuItem{
					{Key: "1", Value: "选项1"},
					{Key: "1", Value: "选项2"},
				},
			},
			wantErr: true,
		},
		{
			name: "无效默认值",
			menu: &Menu{
				Title:   "测试菜单",
				Items:   []MenuItem{{Key: "1", Value: "选项1"}},
				Default: "99",
			},
			wantErr: true,
		},
		{
			name: "允许退出但无退出键",
			menu: &Menu{
				Title:     "测试菜单",
				Items:     []MenuItem{{Key: "1", Value: "选项1"}},
				AllowExit: true,
				ExitKey:   "",
			},
			wantErr: true,
		},
		{
			name: "有效菜单",
			menu: &Menu{
				Title: "测试菜单",
				Items: []MenuItem{
					{Key: "1", Value: "选项1"},
					{Key: "2", Value: "选项2"},
				},
				Default:   "1",
				AllowExit: true,
				ExitKey:   "q",
				ExitText:  "退出",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMenu(tt.menu)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMenu() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRenderMenu 测试菜单渲染
func TestRenderMenu(t *testing.T) {
	tests := []struct {
		name  string
		menu  *Menu
		style *MenuStyle
	}{
		{
			name: "基本渲染",
			menu: &Menu{
				Title: "主菜单",
				Items: []MenuItem{
					{Key: "1", Value: "查看列表"},
					{Key: "2", Value: "添加项目"},
				},
			},
		},
		{
			name: "带标题",
			menu: &Menu{
				Title: "测试菜单",
				Items: []MenuItem{{Key: "1", Value: "选项1"}},
			},
		},
		{
			name: "不带标题",
			menu: &Menu{
				Title: "",
				Items: []MenuItem{{Key: "1", Value: "选项1"}},
			},
		},
		{
			name: "带退出选项",
			menu: &Menu{
				Title:     "主菜单",
				Items:     []MenuItem{{Key: "1", Value: "选项1"}},
				AllowExit: true,
				ExitKey:   "q",
				ExitText:  "退出",
			},
		},
		{
			name: "自定义样式",
			menu: &Menu{
				Title: "主菜单",
				Items: []MenuItem{{Key: "1", Value: "选项1"}},
			},
			style: &MenuStyle{
				Prefix:    ">> ",
				Separator: " | ",
				Indent:    2,
				ShowTitle: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderMenu(tt.menu, tt.style)
			if result == "" {
				t.Error("RenderMenu() 返回空字符串")
			}
		})
	}
}

// TestShowMenu 测试菜单显示
func TestShowMenu(t *testing.T) {
	tests := []struct {
		name    string
		menu    *Menu
		input   string
		wantKey string
		wantErr bool
	}{
		{
			name: "正常选择",
			menu: &Menu{
				Title: "主菜单",
				Items: []MenuItem{
					{Key: "1", Value: "选项1"},
					{Key: "2", Value: "选项2"},
				},
				Prompt: "请选择: ",
			},
			input:   "1\n",
			wantKey: "1",
			wantErr: false,
		},
		{
			name: "字母选择",
			menu: &Menu{
				Title: "主菜单",
				Items: []MenuItem{
					{Key: "a", Value: "选项1"},
					{Key: "b", Value: "选项2"},
				},
				Prompt: "请选择: ",
			},
			input:   "a\n",
			wantKey: "a",
			wantErr: false,
		},
		{
			name: "空输入有默认值",
			menu: &Menu{
				Title:   "主菜单",
				Items:   []MenuItem{{Key: "1", Value: "选项1"}},
				Prompt:  "请选择: ",
				Default: "1",
			},
			input:   "\n",
			wantKey: "1",
			wantErr: false,
		},
		{
			name: "空输入无默认值",
			menu: &Menu{
				Title:  "主菜单",
				Items:  []MenuItem{{Key: "1", Value: "选项1"}},
				Prompt: "请选择: ",
			},
			input:   "\n",
			wantKey: "",
			wantErr: true,
		},
		{
			name: "无效选择",
			menu: &Menu{
				Title:  "主菜单",
				Items:  []MenuItem{{Key: "1", Value: "选项1"}},
				Prompt: "请选择: ",
			},
			input:   "99\n",
			wantKey: "",
			wantErr: true,
		},
		{
			name: "选择退出键",
			menu: &Menu{
				Title:     "主菜单",
				Items:     []MenuItem{{Key: "1", Value: "选项1"}},
				Prompt:    "请选择: ",
				AllowExit: true,
				ExitKey:   "q",
				ExitText:  "退出",
			},
			input:   "q\n",
			wantKey: "q",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			key, err := ShowMenu(tt.menu, reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShowMenu() error = %v, wantErr %v", err, tt.wantErr)
			}
			if key != tt.wantKey {
				t.Errorf("ShowMenu() = %v, want %v", key, tt.wantKey)
			}
		})
	}
}

// TestMenuError 测试菜单错误
func TestMenuError(t *testing.T) {
	tests := []struct {
		name    string
		err     *MenuError
		wantMsg string
	}{
		{
			name: "带详细信息的错误",
			err: &MenuError{
				Code:    "TEST",
				Message: "测试错误",
				Detail:  "详细信息",
			},
			wantMsg: "TEST: 测试错误 (详细信息)",
		},
		{
			name: "不带详细信息的错误",
			err: &MenuError{
				Code:    "TEST",
				Message: "测试错误",
				Detail:  "",
			},
			wantMsg: "TEST: 测试错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			if msg != tt.wantMsg {
				t.Errorf("Error() = %v, want %v", msg, tt.wantMsg)
			}
		})
	}
}

// TestGetDefaultMenuStyle 测试默认样式
func TestGetDefaultMenuStyle(t *testing.T) {
	style := GetDefaultMenuStyle()
	if style.Separator != ". " {
		t.Errorf("Separator = %v, want \". \"", style.Separator)
	}
	if style.Indent != 0 {
		t.Errorf("Indent = %v, want 0", style.Indent)
	}
	if !style.ShowTitle {
		t.Error("ShowTitle = false, want true")
	}
}
