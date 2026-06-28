package Globals

type Config struct {
	PKGType      string `toml:"PKGType"`
	DownloadPath string `toml:"DownloadPath"`
	ExtractPath  string `toml:"ExtractPath"`
	AppImagePath string `toml:"AppImagePath"`
	IconPath     string `toml:"IconPath"`
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

var SupportedIconSuffix = []string{
	".png",
	".jpg",
	".jpeg",
	".ico",
	".svg",
	".xpm",
}

// Must be setup in the main func
var ConfigPath string = ""
