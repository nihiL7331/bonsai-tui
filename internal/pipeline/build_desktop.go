package pipeline

func GetDesktopBinName(targetOS string) string {
	if targetOS == "windows" {
		return "game_desktop.exe"
	}
	return "game_desktop.bin"
}
