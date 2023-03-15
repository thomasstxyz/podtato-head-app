package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/podtato-head/podtato-head-app/pkg/assets"
	"github.com/podtato-head/podtato-head-app/pkg/version"
	"io/fs"
	"log"
	"net/http"
	"os"
)

type PartResponse struct {
	PartNumber string `json:"partNumber"`
	Image      string `json:"image"`
	ServedBy   string `json:"servedBy"`
	Version    string `json:"version"`
}

func PartHandler(w http.ResponseWriter, r *http.Request) {
	partName, found := mux.Vars(r)["partName"]
	if !found {
		// shouldn't happen...
		http.Error(w, fmt.Sprintf("part name %s not found in URL", partName), http.StatusNotFound)
		return
	}

	desiredPartNumber := version.PartNumber(partName)
	imagePath := fmt.Sprintf("images/%s/%s-%s.svg", partName, partName, desiredPartNumber)

	log.Printf("returning file %s", imagePath)
	image, err := fs.ReadFile(assets.Assets, imagePath)
	if err != nil {
		log.Printf("failed to read file %s: %v", imagePath, err)
		http.Error(w, fmt.Sprintf("failed to read file %s", imagePath), http.StatusNotFound)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("failed to get hostname: %v", err)
	}

	data := PartResponse{
		PartNumber: desiredPartNumber,
		Image:      toBase64(image),
		ServedBy:   hostname,
		Version:    version.ServiceVersion(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
