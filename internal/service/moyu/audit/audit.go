package audit

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gone-io/gone"
	"github.com/google/uuid"
	"github.com/imroc/req"
	"github.com/sirupsen/logrus"
	"gitlab.openviewtech.com/moyu-chat/moyu-server/internal/interface/entity"
	"strings"
	"time"
)

type Response struct {
	Code int                  `json:"code"`
	Data []entity.AuditResult `json:"data"`
	Msg  string               `json:"msg"`
}

type Audit struct {
	gone.Flag
	AccessKey    string `gone:"config,aliyun-scan.access-key-id"`
	AccessSecret string `gone:"config,aliyun-scan.access-key-secret"`
	Endpoint     string `gone:"config,aliyun-scan.endpoint"`
}

//go:gone
func NewAudit() gone.Goner {
	return &Audit{}
}

func (a *Audit) ScanImage(imgUrl string) (*entity.AuditResult, error) {
	body := createScanImageBody(imgUrl)

	data, err := a.request("/green/image/scan", body)
	if err != nil {
		return nil, err
	}

	return &data[0], nil
}

func (a *Audit) ScanText(text string) (*entity.AuditResult, error) {
	body := createScanTextBody(text)

	data, err := a.request("/green/text/scan", body)
	if err != nil {
		return nil, err
	}

	return &data[0], nil
}

func md5Base64(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	sum := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

func createScanTextBody(text string) string {
	m := map[string]any{
		"bizType": "default",
		"scenes":  []string{"antispam"},
		"tasks":   []map[string]any{{"dataId": uuid.NewString(), "content": text}},
	}

	b, _ := json.Marshal(m)
	return string(b)
}

func createScanImageBody(url string) string {
	return fmt.Sprintf(`{"scenes":["porn","terrorism","ad","live","qrcode","logo"],"tasks":[{"dataId":"%s","url":"%s"}]}`, uuid.New(), url)
}

func sign(header req.Header, accessKey, accessSecret, uri string) {
	var keys []string

	for k := range header {
		if strings.HasPrefix(k, "x-acs") {
			keys = append(keys, k)
		}
	}

	toSign := bytes.Buffer{}
	toSign.WriteString("POST\n")
	toSign.WriteString(header["Content-Type"])
	toSign.WriteString("\n")
	toSign.WriteString(header["Content-MD5"])
	toSign.WriteString("\n")
	toSign.WriteString(header["Accept"])
	toSign.WriteString("\n")
	toSign.WriteString(header["Date"])
	toSign.WriteString("\n")
	toSign.WriteString("x-acs-signature-method:")
	toSign.WriteString(header["x-acs-signature-method"])
	toSign.WriteString("\n")
	toSign.WriteString("x-acs-signature-nonce:")
	toSign.WriteString(header["x-acs-signature-nonce"])
	toSign.WriteString("\n")
	toSign.WriteString("x-acs-signature-version:")
	toSign.WriteString(header["x-acs-signature-version"])
	toSign.WriteString("\n")
	toSign.WriteString("x-acs-version:")
	toSign.WriteString(header["x-acs-version"])
	toSign.WriteString("\n")
	toSign.WriteString(uri)
	fmt.Println(toSign.String())

	h := hmac.New(sha1.New, []byte(accessSecret))
	h.Write([]byte(toSign.String()))
	sum := h.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(sum)
	header["Authorization"] = fmt.Sprintf("acs %s:%s", accessKey, signature)
}

func FormatGMT(t time.Time) string {
	week := t.Weekday().String()[:3]
	day := t.Day()
	month := t.Month().String()[:3]
	year := t.Year()
	time := t.Format("15:04:05")
	return fmt.Sprintf(fmt.Sprintf("%s, %d %s %d %s GMT", week, day, month, year, time))
}

func (a *Audit) createHeader(uri, body string) req.Header {
	header := req.Header{
		"Content-Type":            "application/json",
		"Content-MD5":             md5Base64(body),
		"Accept":                  "application/json",
		"Date":                    FormatGMT(time.Now().UTC()),
		"x-acs-signature-method":  "HMAC-SHA1",
		"x-acs-signature-nonce":   uuid.NewString(),
		"x-acs-signature-version": "1",
		"x-acs-version":           "2018-05-09",
	}
	sign(header, a.AccessKey, a.AccessSecret, uri)
	return header
}

func (a *Audit) request(uri, body string) ([]entity.AuditResult, error) {
	resp, err := req.Post(fmt.Sprintf("https://%s/%s", a.Endpoint, uri), a.createHeader(uri, body), body)
	if err != nil {
		return nil, err
	}

	var res Response
	logrus.Infof("ScanRequest: uri=%s resp=%s", uri, resp.String())
	err = json.Unmarshal(resp.Bytes(), &res)
	if err != nil {
		return nil, err
	}

	if res.Code != 200 {
		return nil, errors.New(res.Msg)
	}

	return res.Data, nil
}
