package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gorilla/mux"
)

var (
	errUnknownContentType error = errors.New("Unknown content type")
	errInvalidUser        error = errors.New("Invalid user permissions for operation")

	contentTypeToFileExtension map[string]string = map[string]string{
		"image/bmp":                "bmp",
		"image/png":                "png",
		"image/gif":                "gif",
		"image/jpeg":               "jpeg",
		"image/svg+xml":            "svg",
		"image/webp":               "webp",
		"image/tiff":               "tiff",
		"image/vnd.microsoft.icon": "ico",
	}
)

type (
	// CreatePictureResponse - struct representing response after creating an image in azure blob store
	CreatePictureResponse struct {
		URL string
	}

	imageBlobStoreClient struct {
		containerURL azblob.ContainerURL
		urlPrefix    string
	}
)

func newImageBlobStoreClient(cfg *imageBlobStoreConfig) *imageBlobStoreClient {
	credential, err := azblob.NewSharedKeyCredential(cfg.acctName, cfg.acctKey)
	if err != nil {
		log.Fatal("Problem with azure blob credentials: ", err)
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	urlString := fmt.Sprintf("https://%s.blob.core.windows.net/%s", cfg.acctName, cfg.containerName)
	url, _ := url.Parse(urlString)
	containerURL := azblob.NewContainerURL(*url, p)
	log.WithFields(log.Fields{
		"Url": urlString,
	}).Info("Connected to Azure Blob Store")

	return &imageBlobStoreClient{
		containerURL: containerURL,
		urlPrefix:    urlString + "/",
	}
}

func (i *imageBlobStoreClient) SendCreateBlobRequest(ctx context.Context, data []byte, contentType string) (string, error) {
	if ext, ok := contentTypeToFileExtension[contentType]; ok {
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

func (i *imageBlobStoreClient) SendDeleteBlobRequest(ctx context.Context, url string) (*azblob.BlobDeleteResponse, error) {
	//TODO there is a better way to strip the file from the absolute path
	split := strings.Split(url, "/")
	file := split[len(split)-1]
	blob := i.containerURL.NewBlockBlobURL(file)
	return blob.Delete(ctx, "", azblob.BlobAccessConditions{})
}

func (s *server) CreateEventImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	claims := AccessTokenClaimsFromContext(r.Context())

	if claims.UserID == 0 && !claims.IsAdmin {
		http.Error(w, errInvalidUser.Error(), http.StatusUnauthorized)
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

	resp := CreatePictureResponse{URL: s.imageClient.urlPrefix + newImageURL}

	s.db.updateTimelineEventWithImageURL(id, resp.URL)
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
func (s *server) UpdateEventImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ctx := r.Context()

	claims := AccessTokenClaimsFromContext(ctx)

	if claims.UserID == 0 && !claims.IsAdmin {
		http.Error(w, errInvalidUser.Error(), http.StatusUnauthorized)
		return
	}
	imageURL, err := s.db.readImageURLFromEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := s.imageClient.SendDeleteBlobRequest(ctx, imageURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.Response().StatusCode == http.StatusAccepted {
		s.db.deleteImageURLFromEvent(id)
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

	cpr := CreatePictureResponse{URL: s.imageClient.urlPrefix + newImageURL}

	s.db.updateTimelineEventWithImageURL(id, cpr.URL)
	respJSON, err := json.Marshal(cpr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func (s *server) DeleteEventImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ctx := r.Context()

	claims := AccessTokenClaimsFromContext(ctx)

	if claims.UserID == 0 && !claims.IsAdmin {
		http.Error(w, errInvalidUser.Error(), http.StatusUnauthorized)
		return
	}
	imageURL, err := s.db.readImageURLFromEvent(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := s.imageClient.SendDeleteBlobRequest(ctx, imageURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.Response().StatusCode == http.StatusAccepted {
		s.db.deleteImageURLFromEvent(id)
	}

	w.WriteHeader(http.StatusAccepted)

	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

func randomString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.Itoa(r.Int())
}
