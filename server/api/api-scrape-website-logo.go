package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/2manyvcos/paranal/utils"
	"github.com/thanhpk/go-favicon"
)

func PostScrapeWebsiteLogo(res http.ResponseWriter, req *http.Request) {
	if !utils.JsonRegex.MatchString(req.Header.Get("Content-Type")) {
		http.Error(res, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
		return
	}

	var requestPayload struct {
		URL string `json:"url"`
	}
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&requestPayload)
	if err != nil || requestPayload.URL == "" {
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	icons, err := favicon.Find(requestPayload.URL)
	if err != nil {
		log.Printf("Error fetching logo - %s\n", err)
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	logo := selectLogo(icons)
	if logo == "" {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		Logo string `json:"logo"`
	}{
		Logo: logo,
	})
	if err != nil {
		log.Printf("Error encoding response payload - %s\n", err)
	}
}

func selectLogo(icons []*favicon.Icon) string {
	remaining := make([]*favicon.Icon, 0, len(icons))
	for _, icon := range icons {
		if icon.FileExt == "svg" {
			logo, err := fetchLogo(icon)
			if err == nil {
				return logo
			}
		} else {
			remaining = append(remaining, icon)
		}
	}
	for _, icon := range remaining {
		logo, err := fetchLogo(icon)
		if err == nil {
			return logo
		} else {
			log.Printf("Error fetching logo - %s\n", err)
		}
	}
	return ""
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
		return "", fmt.Errorf(`unexpected Content-Type "%s"`, contentType)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data)), nil
}
