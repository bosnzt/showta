package logic

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"showta.cc/app/internal/jwt"
	"showta.cc/app/lib/util"
	"showta.cc/app/system/conf"
	"showta.cc/app/system/log"
	"showta.cc/app/system/model"
	"showta.cc/app/system/msg"
	"time"
)

func checkDefaultUser() {
	total, err := model.CountUser()
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}

	if total > 0 {
		return
	}

	salt, encryptPwd := encryptPassword("123456")
	var addList = []*model.User{
		{
			Username:   "admin",
			EncryptPwd: encryptPwd,
			PwdStamp:   time.Now().UnixNano(),
			Salt:       salt,
			Role:       conf.SuperAdmin,
			Enable:     true,
		},
		{
			Username: "guest",
			Role:     conf.Guest,
		},
	}
	err = model.BatchCreateUser(addList)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
}

func Auth(c *gin.Context, username string, password string) (resp msg.LoginResp, err error) {
	user, err := model.GetUserByName(username)
	if err != nil {
		return
	}

	if user.ID == 0 {
		err = msg.ErrAuthAccount
		return
	}

	mdStr := util.ToMD5(password + user.Salt)
	if mdStr != user.EncryptPwd {
		err = msg.ErrAuthAccount
		return
	}

	token, err := jwt.GenToken(user.Username, user.PwdStamp)
	if err != nil {
		return
	}

	user.LoginIp = util.ClientIP(c.Request)
	user.LoginTime = time.Now()
	model.UpdateUser(user)

	resp.Token = token
	resp.Username = user.Username
	return
}

func ResetPwd(c *gin.Context, password string) (err error) {
	user := c.MustGet("identity").(*model.User)
	user.Salt, user.EncryptPwd = encryptPassword(password)
	err = model.UpdateUser(user)
	return
}

func ListUser(ctx context.Context, req msg.ListUserReq) (resp msg.ListUserResp, err error) {
	total, err := model.CountUserBySearch(req.Query)
	if err != nil {
		return
	}

	offset := (req.Pagenum - 1) * req.Pagesize
	userList, err := model.GeUserListBySearch(req.Query, req.Pagesize, offset)
	if err != nil {
		return
	}

	resp.Total = total
	resp.UserList = userList

	return
}

func EnableUser(ctx context.Context, req model.User) (err error) {
	user, err := model.GetUser(req.ID)
	if err != nil {
		return
	}

	if user.Role == conf.SuperAdmin {
		return fmt.Errorf("super user cannot be disabled")
	}

	if user.Enable == req.Enable {
		return
	}

	user.Enable = req.Enable
	err = model.UpdateUser(user)
	return
}

func AddUser(ctx context.Context, req model.User) (err error) {
	user, err := model.GetUserByName(req.Username)
	if err != nil {
		return
	}

	if user.ID > 0 {
		return fmt.Errorf("username already exist")
	}

	salt, encryptPwd := encryptPassword(req.Password)
	data := model.User{
		Username:   req.Username,
		EncryptPwd: encryptPwd,
		PwdStamp:   time.Now().UnixNano(),
		Salt:       salt,
		Role:       conf.Viewer,
		Enable:     req.Enable,
		Perm:       req.Perm,
	}
	err = model.CreateUser(&data)
	return
}

func UpdateUser(ctx context.Context, req model.User) (err error) {
	user, err := model.GetUser(req.ID)
	if err != nil {
		return
	}

	isSame := true
	if req.Username != user.Username {
		data, err := model.GetUserByName(req.Username)
		if err != nil {
			return err
		}

		if data.ID > 0 {
			return fmt.Errorf("username already exist")
		}

		user.Username = req.Username
		isSame = false
	}

	if req.Password != "" && util.ToMD5(req.Password+user.Salt) != user.EncryptPwd {
		user.Salt, user.EncryptPwd = encryptPassword(req.Password)
		user.PwdStamp = time.Now().UnixNano()
		isSame = false
	}

	if req.Enable != user.Enable {
		user.Enable = req.Enable
		isSame = false
	}

	if req.Perm != user.Perm {
		user.Perm = req.Perm
		isSame = false
	}

	if isSame {
		return
	}

	err = model.UpdateUser(user)
	return
}

func DeleteUser(ctx context.Context, req model.User) (err error) {
	user, err := model.GetUser(req.ID)
	if err != nil {
		return
	}

	if user.Role != conf.Viewer {
		return fmt.Errorf("default user cannot be deleted")
	}

	err = model.DeleteUser(user.ID)
	return
}

func encryptPassword(password string) (string, string) {
	salt := util.GenRandStr(16)
	encryptPwd := util.ToMD5(password + salt)
	return salt, encryptPwd
}
