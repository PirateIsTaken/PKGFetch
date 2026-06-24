package Globals

type Config struct {
	PKGType      string `toml:"PKGType"`
	DownloadPath string `toml:"DownloadPath"`
	AppImagePath string `toml:"AppImagePath"`
	SymlnkPath   string `toml:"SymlnkPath"`
}

var AppConfig Config

var SupportedAssets []string = []string{
	".AppImage",
	".appimage",
	".tar.gz",
	".tar.xz",
	".zip",
}
var UnsupportedKeywords = []string{
	"windows",
	"win64",
	"win32",
	"macos",
	"osx",
	"darwin",
}

// Must be setup in the main func
var ConfigPath string = ""
