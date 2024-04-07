package sign

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"showta.cc/app/system/conf"
	"strconv"
	"time"
)

var (
	ErrInvalidSign = errors.New("Invalid Signature")
	ErrExpiredSign = errors.New("Expired signature")
)

func Gen(rpath string, stamp string) string {
	if conf.SignExpiration > 0 && stamp == "" {
		stamp = fmt.Sprintf("%d", time.Now().Unix()+conf.SignExpiration*3600)
	}

	key := []byte(conf.AppConf.Secure.SignKey)
	hasher := hmac.New(md5.New, key)
	hasher.Write([]byte(rpath + stamp))
	enstr := hex.EncodeToString(hasher.Sum([]byte("")))
	return enstr + stamp
}

func Verify(rpath string, data string) error {
	if conf.SignExpiration == 0 {
		return nil
	}

	if len(data) != 42 {
		return ErrInvalidSign
	}

	stamp := data[32:]
	unixTime, err := strconv.ParseInt(stamp, 10, 64)
	if err != nil {
		return ErrInvalidSign
	}

	timeDiff := unixTime - time.Now().Unix()
	if timeDiff < 0 {
		return ErrExpiredSign
	}

	if timeDiff > (conf.SignExpiration * 3600) {
		return ErrInvalidSign
	}

	if Gen(rpath, stamp) != data {
		return ErrInvalidSign
	}

	return nil
}
