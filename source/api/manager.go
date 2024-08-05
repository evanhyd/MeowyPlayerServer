package api

import (
	"encoding/json"
	"fmt"
	"meowyplayerserver/api/account"
	"meowyplayerserver/api/album"
	"meowyplayerserver/api/logger"
	"net/http"
)

type apiManager struct {
	loggerComponent  logger.Component
	accountComponent account.Component
	albumComponent   album.Component
}

func MakeAPIManager() apiManager {
	return apiManager{
		loggerComponent:  logger.MakeComponent(),
		accountComponent: account.MakeComponent(),
		albumComponent:   album.MakeComponent(),
	}
}

func (m *apiManager) Initialize() error {
	if err := m.loggerComponent.Initialize(); err != nil {
		return err
	}

	if err := m.accountComponent.Initialize(); err != nil {
		return err
	}

	if err := m.albumComponent.Initialize(); err != nil {
		return err
	}

	m.loggerComponent.RegisterAPI("/stats", m.statsHandler)
	m.loggerComponent.RegisterAPI("/logs", m.logsHandler)
	m.loggerComponent.RegisterAPI("/register", m.registerHandler)
	m.loggerComponent.RegisterAPI("/upload", m.uploadHandler)
	m.loggerComponent.RegisterAPI("/download", m.downloadHandler)
	return nil
}

func (m *apiManager) statsHandler(w http.ResponseWriter, _ *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m.loggerComponent.GetRecords()); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}
}

func (m *apiManager) logsHandler(w http.ResponseWriter, _ *http.Request) {
	n, err := m.loggerComponent.DumpLog(w)
	if err != nil {
		http.Error(w, fmt.Sprintf("[logsHandler] %v bytes sent, %v\n", n, err), http.StatusForbidden)
	}
}

func (m *apiManager) registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.User.Username()
	password, ok := r.URL.User.Password()
	if !ok {
		http.Error(w, "[registerHandler] missing password field\n", http.StatusForbidden)
		return
	}
	err := m.accountComponent.Register(username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("[registerHandler] %v\n", err), http.StatusForbidden)
	}
}

func (m *apiManager) uploadHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.User.Username()
	password, _ := r.URL.User.Password()
	if !m.accountComponent.Authorize(username, password) {
		http.Error(w, "[uploadHandler] failed to authorize\n", http.StatusUnauthorized)
	}

	var album album.Album
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, fmt.Sprintf("[uploadHandler] failed to decode album data: %v\n", err), http.StatusInternalServerError)
	}
	if err := m.albumComponent.Upload(w, album); err != nil {
		http.Error(w, fmt.Sprintf("[uploadHandler] failed to upload album data: %v\n", err), http.StatusInternalServerError)
	}
}

func (m *apiManager) downloadHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.User.Username()
	password, _ := r.URL.User.Password()
	if !m.accountComponent.Authorize(username, password) {
		http.Error(w, "[downloadHandler] failed to authorize\n", http.StatusUnauthorized)
	}

	const kAlbumParameter = "albumKey"
	key := r.URL.Query().Get(kAlbumParameter)
	if err := m.albumComponent.Download(w, album.AlbumKey(key)); err != nil {
		http.Error(w, fmt.Sprintf("[downloadHandler] failed to download album: %v\n", err), http.StatusInternalServerError)
	}
}
