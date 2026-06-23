package Globals

type Config struct {
	Distro       string
	DownloadPath string
	SymlnkPath   string
}

var AppConfig Config

var SupportedAssets []string = []string{
	".AppImage",
	".appimage",
	".tar.gz",
	".tar.xz",
	".zip",
}
