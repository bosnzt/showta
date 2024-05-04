package conf

var (
	PreviewVideo   = []string{"avi", "flv", "m3u8", "mkv", "mov", "mp4", "swf", "rmvb", "webm"}
	PreviewPicture = []string{"bmp", "gif", "ico", "jpeg", "jpg", "png", "svg", "tiff", "webp"}
	PreviewText    = []string{"ass", "bat", "c", "conf", "cpp", "go", "h", "hpp", "htm", "html", "ini", "java", "js", "log", "lrc", "md", "pdf", "php", "py", "rs", "sh", "sql", "srt", "tsx", "txt", "vtt", "vue", "xml", "yaml", "yml"}
	PreviewAudio   = []string{"flac", "m4a", "mp3", "ogg", "opus", "wav", "wma"}

	PreviewOffice  string
	DocLocal       = "local"
	DocMS          = "ms"
	OfficeLocalMap = map[string]int{
		"docx": 1,
		"xlsx": 1,
	}

	OfficeMSMap = map[string]int{
		"doc":  1,
		"docm": 1,
		"docx": 1,
		"dot":  1,
		"dotm": 1,
		"dotx": 1,
		"csv":  1,
		"xlam": 1,
		"xls":  1,
		"xlsb": 1,
		"xlsm": 1,
		"xlsx": 1,
		"xlt":  1,
		"xltm": 1,
		"xltx": 1,
		"ppt":  1,
		"pptx": 1,
	}

	SiteDefaultTitle = "ShowTa"
	SiteTitle        string
	SiteLogo         string
	SiteFavicon      string
	SiteDomain       string
	SiteNotice       string
	GlobalSign       bool
	SignExpiration   int64
)
