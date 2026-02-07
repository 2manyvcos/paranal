package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/2manyvcos/paranal/utils"
	"github.com/thanhpk/go-favicon"
)

func PostFindWebsiteLogo(res http.ResponseWriter, req *http.Request) {
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
		}
	}
	return ""
}

func fetchLogo(icon *favicon.Icon) (string, error) {
	mimeType := icon.MimeType
	if mimeType == "image/vnd.microsoft.icon" {
		mimeType = "image/x-icon"
	}
	resp, err := http.Get(icon.URL)
	if err != nil {
		log.Printf("Error fetching logo - %s\n", err)
		return "", err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(data)), nil
}
