package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thanhpk/go-favicon"
)

func ScrapeLogo(url string) (string, error) {
	icons, err := favicon.Find(url)
	if err != nil {
		return "", err
	}
	remaining := make([]*favicon.Icon, 0, len(icons))
	for _, icon := range icons {
		if icon.FileExt == "svg" {
			logo, err := fetchLogo(icon)
			if err == nil {
				return logo, nil
			}
		} else {
			remaining = append(remaining, icon)
		}
	}
	for _, icon := range remaining {
		logo, err := fetchLogo(icon)
		if err == nil {
			return logo, nil
		} else {
			// ignore error
		}
	}
	return "", nil
}

func fetchLogo(icon *favicon.Icon) (string, error) {
	mimeType := icon.MimeType
	if mimeType == "image/vnd.microsoft.icon" {
		mimeType = "image/x-icon"
	}
	res, err := http.Get(icon.URL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %v", res.StatusCode)
	}
	if contentType := res.Header.Get("Content-Type"); contentType != "" && !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf(`unexpected Content-Type %q`, contentType)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return base64DataURL(mimeType, data), nil
}

func base64DataURL(mime string, data []byte) string {
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(data))
}

// func urlEncodedDataURL(mime string, data []byte) string {
// 	encoded := strings.ReplaceAll(url.QueryEscape(string(data)), "+", "%20")
// 	return fmt.Sprintf("data:%s;charset=utf-8,%s", mime, encoded)
// }
