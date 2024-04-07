package conf

var (
	PreviewVideo     = []string{"avi", "flv", "mkv", "mov", "mp4", "rmvb", "webm"}
	PreviewPicture   = []string{"bmp", "gif", "ico", "jpeg", "jpg", "png", "svg", "swf", "tiff", "webp"}
	PreviewText      = []string{"ass", "bat", "c", "conf", "cpp", "go", "h", "hpp", "htm", "html", "ini", "java", "js", "log", "lrc", "md", "php", "py", "rs", "sh", "sql", "srt", "tsx", "txt", "vtt", "vue", "xml", "yaml", "yml"}
	PreviewAudio     = []string{"flac", "m4a", "mp3", "ogg", "opus", "wav", "wma"}
	SiteDefaultTitle = "ShowTa"
	SiteTitle        string
	SiteLogo         string
	SiteFavicon      string
	SiteDomain       string
	SiteNotice       string
	GlobalSign       bool
	SignExpiration   int64
)
