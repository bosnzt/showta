package alipan

import (
	"time"
)

type SimpleResp struct {
	Code string `json:"code"`
	Msg  string `json:"message"`
}

type OauthErrResp struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}

type AccessTokenResp struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type GetDriveInfoResp struct {
	DefaultDriveId  string `json:"default_drive_id"`
	ResourceDriveId string `json:"resource_drive_id"`
	BackupDriveId   string `json:"backup_drive_id"`
}

type FileItem struct {
	DriveId       string    `json:"drive_id"`
	FileId        string    `json:"file_id"`
	ParentFileId  string    `json:"parent_file_id"`
	Name          string    `json:"name"`
	Size          int64     `json:"size"`
	FileExtension string    `json:"file_extension"`
	ContentHash   string    `json:"content_hash"`
	Category      string    `json:"category"`
	Type          string    `json:"type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ListResp struct {
	ItemList []FileItem `json:"items"`
}

type getDownloadUrlResp struct {
	Url        string `json:"url"`
	Expiration string `json:"expiration"`
	Method     string `json:"method"`
}
