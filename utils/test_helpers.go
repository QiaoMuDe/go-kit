package utils

import "runtime"

// getEchoCommand 获取跨平台的echo命令
func getEchoCommand(text string) []string {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/c", "echo", text}
	}
	return []string{"echo", text}
}

// getSleepCommand 获取跨平台的sleep命令
func getSleepCommand(seconds string) []string {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/c", "timeout", "/t", seconds, "/nobreak"}
	}
	return []string{"sleep", seconds}
}

// getWriteFileCommand 获取跨平台的写文件命令
func getWriteFileCommand(content, filename string) []string {
	if runtime.GOOS == "windows" {
		return []string{"cmd", "/c", "echo", content, ">", filename}
	}
	return []string{"sh", "-c", "echo '" + content + "' > " + filename}
}
