package webdav

import (
	"net/http"
	"path"
	"showta.cc/app/system/logic"
	"showta.cc/app/system/msg"
)

func (h *Handler) ServeHTTPOverride(w http.ResponseWriter, r *http.Request) {
	status, err := http.StatusBadRequest, errUnsupportedMethod
	if h.FileSystem == nil {
		status, err = http.StatusInternalServerError, errNoFileSystem
	} else if h.LockSystem == nil {
		status, err = http.StatusInternalServerError, errNoLockSystem
	} else {
		switch r.Method {
		case "OPTIONS":
			status, err = h.handleOptions(w, r)
		case "GET", "HEAD", "POST":
			status, err = h.handleGetHeadPostOverride(w, r)
		case "DELETE":
			status, err = h.handleDelete(w, r)
		case "PUT":
			status, err = h.handlePut(w, r)
		case "MKCOL":
			status, err = h.handleMkcol(w, r)
		case "COPY", "MOVE":
			status, err = h.handleCopyMove(w, r)
		case "LOCK":
			status, err = h.handleLock(w, r)
		case "UNLOCK":
			status, err = h.handleUnlock(w, r)
		case "PROPFIND":
			status, err = h.handlePropfindOverride(w, r)
		case "PROPPATCH":
			status, err = h.handleProppatch(w, r)
		}
	}

	if status != 0 {
		w.WriteHeader(status)
		if status != http.StatusNoContent {
			w.Write([]byte(StatusText(status)))
		}
	}
	if h.Logger != nil {
		h.Logger(r, err)
	}
}

func (h *Handler) handleGetHeadPostOverride(w http.ResponseWriter, r *http.Request) (status int, err error) {
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil {
		return status, err
	}
	// TODO: check locks for read-only access??
	ctx := r.Context()
	// f, err := h.FileSystem.OpenFile(ctx, reqPath, os.O_RDONLY, 0)
	// if err != nil {
	// 	return http.StatusNotFound, err
	// }
	// defer f.Close()
	// fi, err := f.Stat()
	fi, err := logic.GetFile(ctx, reqPath)
	if err != nil {
		return http.StatusNotFound, err
	}
	if fi.IsDir() {
		return http.StatusMethodNotAllowed, nil
	}
	etag, err := findETagOverride(ctx, h.FileSystem, h.LockSystem, reqPath, fi)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Header().Set("ETag", etag)
	// Let ServeContent determine the Content-Type header.
	// http.ServeContent(w, r, reqPath, fi.ModTime(), f)
	logic.ProxyFile(r, w, reqPath)

	return 0, nil
}

func (h *Handler) handlePropfindOverride(w http.ResponseWriter, r *http.Request) (status int, err error) {
	reqPath, status, err := h.stripPrefix(r.URL.Path)
	if err != nil {
		return status, err
	}
	ctx := r.Context()
	// fi, err := h.FileSystem.Stat(ctx, reqPath)
	fi, err := logic.GetFile(ctx, reqPath)
	if err != nil {
		// if os.IsNotExist(err) {
		return http.StatusNotFound, err
		// }
		// return http.StatusMethodNotAllowed, err
	}
	depth := infiniteDepth
	if hdr := r.Header.Get("Depth"); hdr != "" {
		depth = parseDepth(hdr)
		if depth == invalidDepth {
			return http.StatusBadRequest, errInvalidDepth
		}
	}
	pf, status, err := readPropfind(r.Body)
	if err != nil {
		return status, err
	}

	mw := multistatusWriter{w: w}

	walkFn := func(reqPath string, info msg.Finfo, err error) error {
		if err != nil {
			return err
			// return handlePropfindError(err, info)
		}

		var pstats []Propstat
		if pf.Propname != nil {
			pnames, err := propnames(ctx, h.FileSystem, h.LockSystem, reqPath)
			if err != nil {
				return err
				// return handlePropfindError(err, info)
			}
			pstat := Propstat{Status: http.StatusOK}
			for _, xmlname := range pnames {
				pstat.Props = append(pstat.Props, Property{XMLName: xmlname})
			}
			pstats = append(pstats, pstat)
		} else if pf.Allprop != nil {
			pstats, err = allpropOverride(ctx, info, h.LockSystem, reqPath, pf.Prop)
		} else {
			pstats, err = propsOverride(ctx, info, h.LockSystem, reqPath, pf.Prop)
		}
		if err != nil {
			return err
			// return handlePropfindError(err, info)
		}
		href := path.Join(h.Prefix, reqPath)
		if href != "/" && info.IsDir() {
			href += "/"
		}
		return mw.write(makePropstatResponse(href, pstats))
	}

	walkErr := walkFSOverride(ctx, h.FileSystem, depth, reqPath, fi, walkFn)
	closeErr := mw.close()
	if walkErr != nil {
		return http.StatusInternalServerError, walkErr
	}
	if closeErr != nil {
		return http.StatusInternalServerError, closeErr
	}
	return 0, nil
}
