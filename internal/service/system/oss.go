package system

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gone-io/gone/goner/gin"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"io"
	"time"
)

var supportFileTypes = []string{
	"pdf",
}

func (s *svc) GenOssToken(userId int64, fileType string) (*entity.OssToken, error) {
	if fileType != "" {
		support := false
		for _, ext := range supportFileTypes {
			if ext == fileType {
				support = true
				break
			}
		}
		if !support {
			return nil, gin.NewParameterError(fmt.Sprintf("not support %s upload", fileType))
		}
	}

	now := time.Now()
	filename := fmt.Sprintf("%s/%d/%s/%d", s.OssUseDir, userId, now.Format("2006-01-02"), now.UnixMicro())
	policy, signature := s.imageOssPolicy(filename, time.Now().Add(s.OssTokenExpiresIn), fileType)

	return &entity.OssToken{
		AccessId:  s.OssAccessId,
		Host:      s.OssHost,
		Policy:    policy,
		Signature: signature,
		Key:       filename,
	}, nil
}

type jsonNode map[string]interface{}

const iso8601 = "2006-01-02T15:04:05Z"

func (s *svc) imageOssPolicy(filename string, expiredAt time.Time, fileType string) (policy string, signature string) {
	//默认支持图片上传
	var fileLimit []string
	if fileType != "" {
		fileLimit = []string{
			"eq", "$content-type", "application/" + fileType,
		}
	} else {
		fileLimit = []string{
			"starts-with", "$content-type", "image/",
		}
	}

	policyContent := jsonNode{
		"expiration": expiredAt.Format(iso8601),
		"conditions": [][]string{
			{"eq", "$key", filename},
			fileLimit,
		},
	}

	jsonBytes, _ := json.Marshal(policyContent)
	policy = base64.StdEncoding.EncodeToString(jsonBytes)

	h := hmac.New(sha1.New, []byte(s.OssAccessSecret))
	_, _ = io.WriteString(h, policy)
	signature = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return
}
