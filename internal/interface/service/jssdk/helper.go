package jssdk

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func createSignature(nonceStr, ticket, timestamp, url string) string {
	b := bytes.Buffer{}
	b.WriteString("jsapi_ticket=")
	b.WriteString(ticket)
	b.WriteString("&noncestr=")
	b.WriteString(nonceStr)
	b.WriteString("&timestamp=")
	b.WriteString(timestamp)
	b.WriteString("&url=")
	b.WriteString(extractMainUrl(url))

	hash := sha1.New()
	hash.Write(b.Bytes())
	return hex.EncodeToString(hash.Sum(nil))
}

// extractMainUrl 去掉 URL 的锚点
func extractMainUrl(url string) string {
	index := strings.Index(url, "#")
	if index != -1 {
		return url[:index]
	}

	return url
}
