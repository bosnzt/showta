package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"showta.cc/app/internal/memcache"
	"showta.cc/app/internal/sign"
	"showta.cc/app/lib/util"
	"showta.cc/app/storage"
	"showta.cc/app/system/conf"
	"showta.cc/app/system/log"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
	"strings"
)

func ListFile(ctx context.Context, rpath string) (list []msg.Finfo, err error) {
	rpath = util.StandardPath(rpath)
	//Virtual mounting directory
	if rpath == "/" {
		storageMap.Range(func(key, value interface{}) bool {
			v := value.(storage.Storage)
			mountPath := v.GetData().MountPath
			list = append(list, &msg.FileInfo{
				Path:     mountPath,
				Name:     util.SimplePath(mountPath),
				Size:     0,
				Modified: v.GetData().UpdatedAt,
				IsFolder: true,
			})

			return true
		})
	} else {
		var store storage.Storage
		storageMap.Range(func(key, value interface{}) bool {
			v := value.(storage.Storage)
			mountPath := v.GetData().MountPath
			if strings.Contains(rpath, mountPath) {
				store = v
				return false
			}

			return true
		})

		if store != nil {
			if store.AllowCache() {
				list, err = cacheListFile(rpath, store)
			} else {
				list, err = store.List(&msg.FileInfo{Path: rpath})
			}
		}
	}

	return
}

func ViewListFile(ctx context.Context, rpath string) (resp msg.ListFileResp, err error) {
	list, err := ListFile(ctx, rpath)
	if err != nil {
		return
	}

	setting, err := model.GetFolderSettingByFolder(rpath)
	if err != nil {
		return
	}

	var dataList []msg.FileInfo
	for _, v := range list {
		dataList = append(dataList, msg.FileInfo{
			Path:     v.GetPath(),
			Name:     v.GetName(),
			Size:     v.GetSize(),
			Modified: v.ModTime(),
			IsFolder: v.IsDir(),
		})
	}
	resp.List = dataList
	if setting.ID > 0 {
		resp.Topmd = setting.Topmd
		resp.Readme = setting.Readme
	}

	return
}

func GetStorageFile(c *gin.Context, rpath string) (resp msg.GetFileResp, err error) {
	info, err := GetFile(c, rpath)
	if err != nil {
		return
	}

	isDir := info.IsDir()
	name := info.GetName()
	rawUrl := info.GetRaw()
	var ptype int
	if !isDir {
		ptype = getPreviewType(name)
		if rawUrl == "" {
			var param string
			if conf.GlobalSign {
				param = "?sig=" + sign.Gen(rpath, "")
			}

			rawUrl = fmt.Sprintf("%s/fd%s%s", getHost(c.Request), rpath, param)
		}
	}

	resp = msg.GetFileResp{
		FileInfo: msg.FileInfo{
			Path:     info.GetPath(),
			Name:     name,
			Size:     info.GetSize(),
			Modified: info.ModTime(),
			IsFolder: isDir,
			Ptype:    ptype,
		},
		RawUrl: rawUrl,
	}

	return
}

func ProxyFile(r *http.Request, w http.ResponseWriter, rpath string) {
	rpath = util.StandardPath(rpath)
	var store storage.Storage
	storageMap.Range(func(key, value interface{}) bool {
		v := value.(storage.Storage)
		mountPath := v.GetData().MountPath
		if strings.Contains(rpath, mountPath) {
			store = v
			return false
		}

		return true
	})

	if store == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "no such file:", rpath)
		return
	}

	var linkInfo *msg.LinkInfo
	var err error
	if store.AllowCache() {
		item, err := findCacheFile(rpath, store)
		if err != nil || item.IsDir() {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "link file error:", err)
			return
		}

		linkInfo, err = cacheFileLink(item, store)
	} else {
		linkInfo, err = store.Link(&msg.FileInfo{Path: rpath})
	}

	link := linkInfo.Url
	if !store.IsDirect() {
		http.Redirect(w, r, link, http.StatusFound)
		return
	}

	f, err := os.Open(link)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "open file error:", err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "call stat error:", err)
		return
	}

	if fi.IsDir() {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "no such file:", link)
		return
	}

	setAttach(w, fi)
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}

func setAttach(w http.ResponseWriter, fi os.FileInfo) {
	name := fi.Name()
	ext := path.Ext(name)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, name, url.QueryEscape(name)))
	w.Header().Set("Content-Type", mimeType)
}

func GetFile(ctx context.Context, rpath string) (info msg.Finfo, err error) {
	rpath = util.StandardPath(rpath)
	//Virtual mounting directory
	if rpath == "/" {
		info = &msg.FileInfo{
			Name:     "/",
			IsFolder: true,
		}
		return
	}

	var store storage.Storage
	storageMap.Range(func(key, value interface{}) bool {
		v := value.(storage.Storage)
		mountPath := v.GetData().MountPath
		if strings.Contains(rpath, mountPath) {
			store = v
			return false
		}

		return true
	})

	if store == nil {
		err = errors.New("dir not exist")
		return
	}

	getter, ok := store.(storage.Getter)
	if ok {
		info, err = getter.Get(rpath)
		return
	}

	if store.AllowCache() {
		item, err := findCacheFile(rpath, store)
		if err != nil {
			return nil, err
		}

		if item.IsDir() {
			return item, nil
		}

		linkInfo, err := cacheFileLink(item, store)
		if err != nil {
			return nil, err
		}

		info = &msg.FileInfo{
			FileId:   item.GetFileId(),
			Path:     item.GetPath(),
			Name:     item.GetName(),
			Size:     item.GetSize(),
			Modified: item.ModTime(),
			IsFolder: item.IsDir(),
			RawUrl:   linkInfo.Url,
		}

		return info, nil
	}

	err = errors.New("dir err")
	return
}

func Subdir(ctx context.Context, rpath string) (list []msg.Finfo, err error) {
	fileList, err := ListFile(ctx, rpath)
	if err != nil {
		return
	}

	for _, v := range fileList {
		if v.IsDir() {
			list = append(list, v)
		}
	}

	return
}

func cacheListFile(rpath string, store storage.Storage) (list []msg.Finfo, err error) {
	data, found := memcache.Get(memcache.List, rpath)
	if found {
		log.Debugf("[CACHE] list path: %+s", rpath)
		list = data.([]msg.Finfo)
		return
	}

	info := msg.FileInfo{Path: rpath}
	parentPath := util.GetParentDir(rpath)
	if parentPath != "/" {
		pdata, pfound := memcache.Get(memcache.List, parentPath)
		if pfound {
			plist := pdata.([]msg.Finfo)
			for _, item := range plist {
				if item.GetPath() == rpath {
					info.FileId = item.GetFileId()
					break
				}
			}
		}
	}

	list, err = store.List(&info)
	if err != nil {
		return
	}

	memcache.Set(memcache.List, rpath, list)
	return
}

func findCacheFile(rpath string, store storage.Storage) (info msg.Finfo, err error) {
	dpath, fname := util.SplitPath(rpath)
	if dpath == "/" {
		err = errors.New("dir err")
		return
	}

	list, err := cacheListFile(dpath, store)
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		if item.GetName() == fname {
			return item, nil
		}
	}

	return
}

func cacheFileLink(info msg.Finfo, store storage.Storage) (linkInfo *msg.LinkInfo, err error) {
	rpath := info.GetPath()
	data, found := memcache.Get(memcache.Link, rpath)
	if found {
		log.Debugf("[CACHE] link path: %+s", rpath)
		linkInfo = data.(*msg.LinkInfo)
		return
	}

	linkInfo, err = store.Link(info)
	if err != nil {
		return nil, err
	}

	memcache.Expire(memcache.Link, rpath, linkInfo, linkInfo.Expire)
	return
}
