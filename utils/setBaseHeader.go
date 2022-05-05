package utils

import (
	"bilibiliOpen/content"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewBiliBiReq(path string, jsonBody string) *http.Request {
	body := []byte(jsonBody)
	md5Body := MD5V(body)
	req, _ := http.NewRequest("POST", content.SERVER+path, strings.NewReader(string(body)))
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	u, _ := uuid.NewUUID()
	nonce := u.String()
	token := makeToken(md5Body, nonce, timestamp)
	h := req.Header
	h.Set("Accept", "application/json")
	h.Set("Content-Type", "application/json")
	h.Set("x-bili-content-md5", md5Body)
	h.Set("x-bili-timestamp", timestamp)
	h.Set("x-bili-signature-method", "HMAC-SHA256")
	h.Set("x-bili-signature-nonce", nonce)
	h.Set("x-bili-accesskeyid", content.ACCESS_KEY_ID)
	h.Set("x-bili-signature-version", "1.0")
	h.Set("Authorization", token)
	return req
}

func makeToken(md5Body, nonce, timestamp string) string {
	dataStr := "x-bili-accesskeyid:" + content.ACCESS_KEY_ID + "\n" +
		"x-bili-content-md5:" + md5Body + "\n" +
		"x-bili-signature-method:HMAC-SHA256" + "\n" +
		"x-bili-signature-nonce:" + nonce + "\n" +
		"x-bili-signature-version:1.0" + "\n" +
		"x-bili-timestamp:" + timestamp
	token := hmacSHA256(content.ACCESS_KEY_SECRET, dataStr)
	return token
}

func hmacSHA256(key string, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}
