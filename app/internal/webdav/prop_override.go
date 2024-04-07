package webdav

import (
	"context"
	"encoding/xml"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"showta.cc/app/system/msg"
	"strconv"
)

var livePropsOverride = map[xml.Name]struct {
	// findFn implements the propfind function of this property. If nil,
	// it indicates a hidden property.
	findFn func(context.Context, FileSystem, LockSystem, string, msg.Finfo) (string, error)
	// dir is true if the property applies to directories.
	dir bool
}{
	{Space: "DAV:", Local: "resourcetype"}: {
		findFn: findResourceTypeOverride,
		dir:    true,
	},
	{Space: "DAV:", Local: "displayname"}: {
		findFn: findDisplayNameOverride,
		dir:    true,
	},
	{Space: "DAV:", Local: "getcontentlength"}: {
		findFn: findContentLengthOverride,
		dir:    false,
	},
	{Space: "DAV:", Local: "getlastmodified"}: {
		findFn: findLastModifiedOverride,
		// http://webdav.org/specs/rfc4918.html#PROPERTY_getlastmodified
		// suggests that getlastmodified should only apply to GETable
		// resources, and this package does not support GET on directories.
		//
		// Nonetheless, some WebDAV clients expect child directories to be
		// sortable by getlastmodified date, so this value is true, not false.
		// See golang.org/issue/15334.
		dir: true,
	},
	{Space: "DAV:", Local: "creationdate"}: {
		findFn: nil,
		dir:    false,
	},
	{Space: "DAV:", Local: "getcontentlanguage"}: {
		findFn: nil,
		dir:    false,
	},
	{Space: "DAV:", Local: "getcontenttype"}: {
		findFn: findContentTypeOverride,
		dir:    false,
	},
	{Space: "DAV:", Local: "getetag"}: {
		findFn: findETagOverride,
		// findETag implements ETag as the concatenated hex values of a file's
		// modification time and size. This is not a reliable synchronization
		// mechanism for directories, so we do not advertise getetag for DAV
		// collections.
		dir: false,
	},

	// TODO: The lockdiscovery property requires LockSystem to list the
	// active locks on a resource.
	{Space: "DAV:", Local: "lockdiscovery"}: {},
	{Space: "DAV:", Local: "supportedlock"}: {
		findFn: findSupportedLockOverride,
		dir:    true,
	},
}

func propsOverride(ctx context.Context, fi msg.Finfo, ls LockSystem, name string, pnames []xml.Name) ([]Propstat, error) {
	// f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	// if err != nil {
	// 	return nil, err
	// }
	// defer f.Close()
	// fi, err := f.Stat()
	// if err != nil {
	// 	return nil, err
	// }
	isDir := fi.IsDir()

	var deadProps map[xml.Name]Property
	// if dph, ok := f.(DeadPropsHolder); ok {
	// 	deadProps, err = dph.DeadProps()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	pstatOK := Propstat{Status: http.StatusOK}
	pstatNotFound := Propstat{Status: http.StatusNotFound}
	for _, pn := range pnames {
		// If this file has dead properties, check if they contain pn.
		if dp, ok := deadProps[pn]; ok {
			pstatOK.Props = append(pstatOK.Props, dp)
			continue
		}
		// Otherwise, it must either be a live property or we don't know it.
		if prop := livePropsOverride[pn]; prop.findFn != nil && (prop.dir || !isDir) {
			innerXML, err := prop.findFn(ctx, nil, ls, name, fi)
			if err != nil {
				return nil, err
			}
			pstatOK.Props = append(pstatOK.Props, Property{
				XMLName:  pn,
				InnerXML: []byte(innerXML),
			})
		} else {
			pstatNotFound.Props = append(pstatNotFound.Props, Property{
				XMLName: pn,
			})
		}
	}
	return makePropstats(pstatOK, pstatNotFound), nil
}

func propnamesOverride(ctx context.Context, fi msg.Finfo, ls LockSystem, name string) ([]xml.Name, error) {
	// f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	// if err != nil {
	// 	return nil, err
	// }
	// defer f.Close()
	// fi, err := f.Stat()
	// if err != nil {
	// 	return nil, err
	// }
	isDir := fi.IsDir()

	var deadProps map[xml.Name]Property
	// if dph, ok := f.(DeadPropsHolder); ok {
	// 	deadProps, err = dph.DeadProps()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	pnames := make([]xml.Name, 0, len(liveProps)+len(deadProps))
	for pn, prop := range liveProps {
		if prop.findFn != nil && (prop.dir || !isDir) {
			pnames = append(pnames, pn)
		}
	}
	for pn := range deadProps {
		pnames = append(pnames, pn)
	}
	return pnames, nil
}

func allpropOverride(ctx context.Context, fs msg.Finfo, ls LockSystem, name string, include []xml.Name) ([]Propstat, error) {
	pnames, err := propnamesOverride(ctx, fs, ls, name)
	if err != nil {
		return nil, err
	}
	// Add names from include if they are not already covered in pnames.
	nameset := make(map[xml.Name]bool)
	for _, pn := range pnames {
		nameset[pn] = true
	}
	for _, pn := range include {
		if !nameset[pn] {
			pnames = append(pnames, pn)
		}
	}
	return propsOverride(ctx, fs, ls, fs.GetName(), pnames)
}

func findResourceTypeOverride(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi msg.Finfo) (string, error) {
	if fi.IsDir() {
		return `<D:collection xmlns:D="DAV:"/>`, nil
	}
	return "", nil
}

func findDisplayNameOverride(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi msg.Finfo) (string, error) {
	if slashClean(name) == "/" {
		// Hide the real name of a possibly prefixed root directory.
		return "", nil
	}
	return escapeXML(fi.GetName()), nil
}

func findContentLengthOverride(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi msg.Finfo) (string, error) {
	return strconv.FormatInt(fi.GetSize(), 10), nil
}

func findLastModifiedOverride(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi msg.Finfo) (string, error) {
	return fi.ModTime().UTC().Format(http.TimeFormat), nil
}

func findContentTypeOverride(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi msg.Finfo) (string, error) {
	// if do, ok := fi.(ContentTyper); ok {
	// 	ctype, err := do.ContentType(ctx)
	// 	if err != ErrNotImplemented {
	// 		return ctype, err
	// 	}
	// }
	// f, err := fs.OpenFile(ctx, name, os.O_RDONLY, 0)
	// if err != nil {
	// 	return "", err
	// }
	// defer f.Close()
	// This implementation is based on serveContent's code in the standard net/http package.
	ctype := mime.TypeByExtension(filepath.Ext(name))
	// if ctype != "" {
	return ctype, nil
	// }
	// // Read a chunk to decide between utf-8 text and binary.
	// var buf [512]byte
	// n, err := io.ReadFull(f, buf[:])
	// if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
	// 	return "", err
	// }
	// ctype = http.DetectContentType(buf[:n])
	// // Rewind file.
	// _, err = f.Seek(0, io.SeekStart)
	// return ctype, err
}

func findETagOverride(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi msg.Finfo) (string, error) {
	if do, ok := fi.(ETager); ok {
		etag, err := do.ETag(ctx)
		if err != ErrNotImplemented {
			return etag, err
		}
	}
	// The Apache http 2.4 web server by default concatenates the
	// modification time and size of a file. We replicate the heuristic
	// with nanosecond granularity.
	return fmt.Sprintf(`"%x%x"`, fi.ModTime().UnixNano(), fi.GetSize()), nil
}

func findSupportedLockOverride(ctx context.Context, fs FileSystem, ls LockSystem, name string, fi msg.Finfo) (string, error) {
	return `` +
		`<D:lockentry xmlns:D="DAV:">` +
		`<D:lockscope><D:exclusive/></D:lockscope>` +
		`<D:locktype><D:write/></D:locktype>` +
		`</D:lockentry>`, nil
}
