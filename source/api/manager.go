package api

import (
	"encoding/json"
	"fmt"
	"log"
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
	m.loggerComponent.RegisterAPI("/login", m.loginHandler)
	m.loggerComponent.RegisterAPI("/upload", m.uploadHandler)
	m.loggerComponent.RegisterAPI("/download", m.downloadHandler)
	m.loggerComponent.RegisterAPI("/downloadAll", m.downloadHandler)
	m.loggerComponent.RegisterAPI("/remove", m.removeHandler)
	return nil
}

func (m *apiManager) authenticate(r *http.Request) (account.UserID, bool) {
	username, password, _ := r.BasicAuth()
	return m.accountComponent.Authenticate(username, password)
}

func (m *apiManager) statsHandler(w http.ResponseWriter, _ *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m.loggerComponent.GetRecords()); err != nil {
		log.Println(err)
		http.Error(w, "failed to download stats", http.StatusNotFound)
	}
}

func (m *apiManager) logsHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := m.loggerComponent.DumpLog(w)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to download log", http.StatusNotFound)
	}
}

func (m *apiManager) registerHandler(w http.ResponseWriter, r *http.Request) {
	username, password, _ := r.BasicAuth()

	userID, ok := m.accountComponent.Register(username, password)
	if !ok {
		http.Error(w, "username is too short or too long or already exists", http.StatusNotFound)
	}

	if err := m.albumComponent.Register(userID); err != nil {
		log.Println(err)
		http.Error(w, "failed to allocate album storage", http.StatusNotFound)
	}
}

func (m *apiManager) loginHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := m.authenticate(r); ok {
		fmt.Fprintln(w, "login successfully")
	} else {
		http.Error(w, "failed to authenticate", http.StatusNotFound)
	}
}

func (m *apiManager) uploadHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(r)
	if !ok {
		http.Error(w, "failed to authenticate", http.StatusNotFound)
		return
	}

	//decode client data
	var album album.Album
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		log.Println(err)
		http.Error(w, "failed to decode album data", http.StatusNotFound)
		return
	}

	//upload to the storage
	if err := m.albumComponent.Upload(userID, album); err != nil {
		log.Println(err)
		http.Error(w, "failed to upload album data", http.StatusNotFound)
	}
}

func (m *apiManager) downloadHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(r)
	if !ok {
		http.Error(w, "failed to authenticate", http.StatusNotFound)
		return
	}

	//get the album key
	const kAlbumKeyParam = "albumKey"
	key, err := album.ParseAlbumKey(r.URL.Query().Get(kAlbumKeyParam))
	if err != nil {
		http.Error(w, "invalid album key", http.StatusNotFound)
		return
	}

	//get the album from the storage
	album, err := m.albumComponent.Download(userID, key)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to download album", http.StatusNotFound)
		return
	}

	//send to the client
	if err := json.NewEncoder(w).Encode(album); err != nil {
		log.Println(err)
		http.Error(w, "failed to download album", http.StatusNotFound)
	}
}

func (m *apiManager) downloadAllHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(r)
	if !ok {
		http.Error(w, "failed to authenticate", http.StatusNotFound)
		return
	}

	//get all albums from the storage
	albums, err := m.albumComponent.DownloadAll(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to download all albums", http.StatusNotFound)
		return
	}

	//send to the client
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		log.Println(err)
		http.Error(w, "failed to download album", http.StatusNotFound)
	}
}

func (m *apiManager) removeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := m.authenticate(r)
	if !ok {
		http.Error(w, "failed to authenticate", http.StatusNotFound)
		return
	}

	//get the album key
	const kAlbumKeyParam = "albumKey"
	key, err := album.ParseAlbumKey(r.URL.Query().Get(kAlbumKeyParam))
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid album key", http.StatusNotFound)
		return
	}

	if err := m.albumComponent.Remove(userID, key); err != nil {
		log.Println(err)
		http.Error(w, "failed to remove album", http.StatusNotFound)
		return
	}
}
