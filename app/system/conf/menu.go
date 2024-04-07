package conf

// Mtype:{1: 'Normal menu', 2: 'Jump menu inside', 3: 'Jump menu outside'}
type Menu struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsAdmin  bool   `json:"-"`
	Icon     string `json:"icon"`
	Mtype    int    `json:"mtype"`
	Children []Menu `json:"children"`
}

var (
	CommonMenuList = []Menu{}
	AdminMenuList  = []Menu{
		{
			Name:     "home",
			Path:     "/@admin/home",
			Icon:     "Postcard",
			Children: []Menu{},
		},
		{
			Name:     "storage",
			Path:     "/@admin/storage",
			IsAdmin:  true,
			Icon:     "MostlyCloudy",
			Children: []Menu{},
		},
		{
			Name:     "folder",
			Path:     "/@admin/folder",
			IsAdmin:  true,
			Icon:     "FolderOpened",
			Children: []Menu{},
		},
		{
			Name:     "user",
			Path:     "/@admin/user",
			IsAdmin:  true,
			Icon:     "User",
			Children: []Menu{},
		},
		{
			Name:     "site",
			Path:     "/@admin/preference/site",
			IsAdmin:  true,
			Icon:     "Monitor",
			Children: []Menu{},
		},
		{
			Name:     "display",
			Path:     "/@admin/preference/display",
			IsAdmin:  true,
			Icon:     "Camera",
			Children: []Menu{},
		},
		{
			Name:     "backup",
			Path:     "/@admin/backup",
			IsAdmin:  true,
			Icon:     "Sort",
			Children: []Menu{},
		},
		{
			Name:     "documentation",
			Path:     "https://www.showta.cc/intro/",
			Icon:     "Star",
			Mtype:    2,
			Children: []Menu{},
		},
		{
			Name:     "cdhomepage",
			Path:     "/",
			Icon:     "House",
			Mtype:    1,
			Children: []Menu{},
		},
	}
)

func init() {
	for _, v := range AdminMenuList {
		if !v.IsAdmin {
			CommonMenuList = append(CommonMenuList, v)
		}
	}
}
