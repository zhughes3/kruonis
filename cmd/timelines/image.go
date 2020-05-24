package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gorilla/mux"
)

var errUnknownContentType error = errors.New("Unknown content type.")

type ImageClient struct {
	containerURL azblob.ContainerURL
	urlPrefix    string
}

var ContentTypeToFileExtension map[string]string = map[string]string{
	"image/bmp":                "bmp",
	"image/png":                "png",
	"image/gif":                "gif",
	"image/jpeg":               "jpeg",
	"image/svg+xml":            "svg",
	"image/webp":               "webp",
	"image/tiff":               "tiff",
	"image/vnd.microsoft.icon": "ico",
}

func (i *ImageClient) SendCreateBlobRequest(ctx context.Context, data []byte, contentType string) (string, error) {
	if ext, ok := ContentTypeToFileExtension[contentType]; ok {
		filename := randomString() + "." + ext

		url := i.containerURL.NewBlockBlobURL(filename)

		_, err := azblob.UploadBufferToBlockBlob(ctx, data, url, azblob.UploadToBlockBlobOptions{
			BlockSize:   4 * 1024 * 1024,
			Parallelism: 16,
			BlobHTTPHeaders: azblob.BlobHTTPHeaders{
				ContentType: contentType,
			},
		})

		if err != nil {
			return "", err
		}

		return filename, nil
	}

	return "", errUnknownContentType
}

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}
func NewImageClient(cfg *imageServerConfig) *ImageClient {
	credential, err := azblob.NewSharedKeyCredential(cfg.acctName, cfg.acctKey)
	if err != nil {
		log.Fatal("Problem with azure blob credentials: ", err)
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	urlString := fmt.Sprintf("https://%s.blob.core.windows.net/%s", cfg.acctName, cfg.containerName)
	url, _ := url.Parse(urlString)
	return &ImageClient{
		containerURL: azblob.NewContainerURL(*url, p),
		urlPrefix:    urlString + "/",
	}
}

type CreatePictureResponse struct {
	Url string
}

func (s *server) CreatePictureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("content-type")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newImageURL, err := s.imageClient.SendCreateBlobRequest(r.Context(), body, contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := CreatePictureResponse{Url: s.imageClient.urlPrefix + newImageURL}

	s.db.updateTimelineEventWithImageURL(eventID, resp.Url)
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
