package logic

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"showta.cc/app/lib/util"
	"showta.cc/app/system/conf"
	"showta.cc/app/system/log"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
	"strconv"
	"strings"
	"sync"
)

var (
	ptypeMap       sync.Map
	previewVideo   []string
	previewPicture []string
	previewText    []string
	previewAudio   []string
)

const (
	SiteTitleKey      = "site_title"
	SiteLogoKey       = "site_logo"
	SiteFaviconKey    = "site_favicon"
	SiteDomainKey     = "site_domain"
	SiteNoticeKey     = "site_notice"
	GlobalSignKey     = "global_sign"
	SignExpirationKey = "sign_expiration"
	PreviewVideoKey   = "preview_video"
	PreviewPictureKey = "preview_picture"
	PreviewTextKey    = "preview_text"
	PreviewAudioKey   = "preview_audio"
	TremSite          = 1
	TremDisplay       = 2
)

func checkDefaultPreference() {
	dataList, err := model.GetAllPreference()
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}

	if len(dataList) > 0 {
		for _, v := range dataList {
			if v.Key == SiteTitleKey {
				setSiteTitle(v.Value)
			} else if v.Key == SiteLogoKey {
				conf.SiteLogo = v.Value
			} else if v.Key == SiteFaviconKey {
				conf.SiteFavicon = v.Value
			} else if v.Key == SiteDomainKey {
				conf.SiteDomain = v.Value
			} else if v.Key == SiteNoticeKey {
				conf.SiteNotice = v.Value
			} else if v.Key == GlobalSignKey {
				conf.GlobalSign = util.StrToBool(v.Value)
			} else if v.Key == SignExpirationKey {
				conf.SignExpiration, _ = strconv.ParseInt(v.Value, 10, 64)
			} else if v.Key == PreviewVideoKey {
				previewVideo = strings.Split(v.Value, ",")
			} else if v.Key == PreviewPictureKey {
				previewPicture = strings.Split(v.Value, ",")
			} else if v.Key == PreviewTextKey {
				previewText = strings.Split(v.Value, ",")
			} else if v.Key == PreviewAudioKey {
				previewAudio = strings.Split(v.Value, ",")
			}
		}

		loadPreviewConf()
	} else {
		var addList = []*model.Preference{
			{
				Key:   SiteTitleKey,
				Value: conf.SiteDefaultTitle,
				Term:  TremSite,
			},
			{
				Key:   SiteLogoKey,
				Value: "",
				Term:  TremSite,
			},
			{
				Key:   SiteFaviconKey,
				Value: "",
				Term:  TremSite,
			},
			{
				Key:   SiteDomainKey,
				Value: "",
				Term:  TremSite,
			},
			{
				Key:   SiteNoticeKey,
				Value: "",
				Term:  TremSite,
			},
			{
				Key:   GlobalSignKey,
				Value: "true",
				Term:  TremSite,
			},
			{
				Key:   SignExpirationKey,
				Value: "0",
				Term:  TremSite,
			},
			{
				Key:   PreviewVideoKey,
				Value: toString(conf.PreviewVideo),
				Term:  TremDisplay,
			},
			{
				Key:   PreviewPictureKey,
				Value: toString(conf.PreviewPicture),
				Term:  TremDisplay,
			},
			{
				Key:   PreviewTextKey,
				Value: toString(conf.PreviewText),
				Term:  TremDisplay,
			},
			{
				Key:   PreviewAudioKey,
				Value: toString(conf.PreviewAudio),
				Term:  TremDisplay,
			},
		}
		err = model.BatchCreatePreference(addList)
		if err != nil {
			log.Error(err)
			os.Exit(0)
		}

		previewVideo = conf.PreviewVideo
		previewPicture = conf.PreviewPicture
		previewText = conf.PreviewText
		previewAudio = conf.PreviewAudio
		loadPreviewConf()
	}
}

func loadPreviewConf() {
	for _, v := range previewVideo {
		ptypeMap.Store(v, conf.Video)
	}

	for _, v := range previewPicture {
		ptypeMap.Store(v, conf.Picture)
	}

	for _, v := range previewText {
		ptypeMap.Store(v, conf.Text)
	}

	for _, v := range previewAudio {
		ptypeMap.Store(v, conf.Audio)
	}
}

func GetPreference() (resp map[string]string, err error) {
	dataList, err := model.GetPreferenceByTerm(TremSite)
	if err != nil {
		return
	}

	resp = map[string]string{}
	for _, v := range dataList {
		resp[v.Key] = v.Value
	}

	return
}

func GetDisplay() (resp msg.GetDisplayResp, err error) {
	dataList, err := model.GetPreferenceByTerm(TremDisplay)
	if err != nil {
		return
	}

	resp = msg.GetDisplayResp{
		Template: msg.DisplayTemplate{
			Video:   conf.PreviewVideo,
			Picture: conf.PreviewPicture,
			Text:    conf.PreviewText,
			Audio:   conf.PreviewAudio,
		},
		Data: dataList,
	}
	return
}

func UpdateDisplay(ctx context.Context, data msg.UpdateDisplayReq) (err error) {
	video := data.Video
	if len(video) > 0 {
		video = findIntersection(conf.PreviewVideo, video)
	}

	err = model.UpdatePreference(PreviewVideoKey, toString(video))
	if err != nil {
		return
	}

	picture := data.Picture
	if len(picture) > 0 {
		picture = findIntersection(conf.PreviewPicture, picture)
	}

	err = model.UpdatePreference(PreviewPictureKey, toString(picture))
	if err != nil {
		return
	}

	text := data.Text
	if len(text) > 0 {
		text = findIntersection(conf.PreviewText, text)
	}

	err = model.UpdatePreference(PreviewTextKey, toString(text))
	if err != nil {
		return
	}

	audio := data.Audio
	if len(audio) > 0 {
		audio = findIntersection(conf.PreviewAudio, audio)
	}

	err = model.UpdatePreference(PreviewAudioKey, toString(audio))
	if err != nil {
		return
	}

	previewVideo = video
	previewPicture = picture
	previewText = text
	previewAudio = audio
	loadPreviewConf()

	return
}

func GetSite() (resp msg.GetSiteResp, err error) {
	dataList, err := model.GetPreferenceByTerm(TremSite)
	if err != nil {
		return
	}

	resp = msg.GetSiteResp{
		Data: dataList,
	}
	return
}

func UpdateSite(ctx context.Context, data msg.UpdateSiteReq) (err error) {
	err = model.UpdatePreference(SiteTitleKey, data.Title)
	if err != nil {
		return
	}

	err = model.UpdatePreference(SiteLogoKey, data.Logo)
	if err != nil {
		return
	}

	err = model.UpdatePreference(SiteFaviconKey, data.Favicon)
	if err != nil {
		return
	}

	err = model.UpdatePreference(SiteDomainKey, data.Domain)
	if err != nil {
		return
	}

	err = model.UpdatePreference(SiteNoticeKey, data.Notice)
	if err != nil {
		return
	}

	err = model.UpdatePreference(GlobalSignKey, fmt.Sprintf("%t", data.GlobalSign))
	if err != nil {
		return
	}

	err = model.UpdatePreference(SignExpirationKey, strconv.FormatInt(data.SignExpiration, 10))
	if err != nil {
		return
	}

	setSiteTitle(data.Title)
	conf.SiteLogo = data.Logo
	conf.SiteFavicon = data.Favicon
	conf.SiteDomain = data.Domain
	conf.SiteNotice = data.Notice
	conf.GlobalSign = data.GlobalSign
	conf.SignExpiration = data.SignExpiration

	return
}

func setSiteTitle(title string) {
	if title == "" {
		conf.SiteTitle = conf.SiteDefaultTitle
	} else {
		conf.SiteTitle = title
	}
}

func getPreviewType(name string) int {
	ext := strings.TrimPrefix(filepath.Ext(name), ".")
	if ext == "" {
		return conf.Other
	}

	ext = strings.ToLower(ext)
	if ptype, ok := ptypeMap.Load(ext); ok {
		return ptype.(int)
	}

	return conf.Other
}

func findIntersection(sliceA, sliceB []string) []string {
	elements := make(map[string]bool)
	intersection := []string{}

	for _, element := range sliceA {
		elements[element] = true
	}

	for _, element := range sliceB {
		if _, ok := elements[element]; ok {
			intersection = append(intersection, element)
		}
	}

	return intersection
}

func toString(data []string) string {
	return strings.Join(data, ",")
}

func getHost(r *http.Request) string {
	if conf.SiteDomain != "" {
		return conf.SiteDomain
	}

	scheme := "http://"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https://"
	}

	scheme += r.Host
	return scheme
}
