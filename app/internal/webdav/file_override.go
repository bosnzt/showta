package webdav

import (
	"context"
	"path"
	"path/filepath"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/msg"
)

type WalkFunc func(pathStr string, info msg.Finfo, err error) error

func walkFSOverride(ctx context.Context, fs FileSystem, depth int, name string, info msg.Finfo, walkFn WalkFunc) error {
	// This implementation is based on Walk's code in the standard path/filepath package.
	err := walkFn(name, info, nil)
	if err != nil {
		if info.IsDir() && err == filepath.SkipDir {
			return nil
		}
		return err
	}
	if !info.IsDir() || depth == 0 {
		return nil
	}
	if depth == 1 {
		depth = 0
	}

	// Read directory names.
	fileInfos, err := logic.ListFile(ctx, name)
	// f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	// if err != nil {
	// 	return walkFn(name, info, err)
	// }
	// fileInfos, err := f.Readdir(0)
	// f.Close()
	if err != nil {
		return walkFn(name, info, err)
	}

	for _, fileInfo := range fileInfos {
		filename := path.Join(name, fileInfo.GetName())
		// fileInfo, err := fs.Stat(ctx, filename)
		if err != nil {
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir {
				return err
			}
		} else {
			err = walkFSOverride(ctx, fs, depth, filename, fileInfo, walkFn)
			if err != nil {
				if !fileInfo.IsDir() || err != filepath.SkipDir {
					return err
				}
			}
		}
	}
	return nil
}
