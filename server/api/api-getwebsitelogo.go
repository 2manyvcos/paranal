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

func PostGetWebsiteLogo(res http.ResponseWriter, req *http.Request) {
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
	if logo == nil {
		http.Error(res, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	if logo.MimeType == "image/vnd.microsoft.icon" {
		logo.MimeType = "image/x-icon"
	}
	resp, err := http.Get(logo.URL)
	if err != nil {
		log.Printf("Error fetching logo - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error fetching logo - %s\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(struct {
		Logo string `json:"logo"`
	}{
		Logo: fmt.Sprintf("data:%s;base64,%s", logo.MimeType, base64.StdEncoding.EncodeToString(data)),
	})
}

func selectLogo(icons []*favicon.Icon) *favicon.Icon {
	for _, icon := range icons {
		if icon.FileExt == "svg" {
			return icon
		}
	}
	if len(icons) > 0 {
		return icons[0]
	}
	return nil
}
