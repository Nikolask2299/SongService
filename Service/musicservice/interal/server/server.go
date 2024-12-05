package server

import (
	"encoding/json"
	"log/slog"
	"musicservice/interal/app"
	"musicservice/interal/models"
	"net/http"
	"strconv"
)

// Server represents the server interface
// @Description Server interface for interacting with the server
type MysicServer struct {
	logger *slog.Logger
	app   app.App
}

func NewMysicServer(logger *slog.Logger, app app.App) *MysicServer {
    return &MysicServer{logger: logger, app: app}
}

// GetData godoc
// @Summary      Get Data 
// @Description  get songs from database
// @Tags         data
// @Accept       json
// @Produce      json
// @Param        page query string true "first page"
// @Param        limit query string true "count page"
// @Param        input body models.FilterSong true "filter information"
// @Success      200  {array} models.Song
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /search [post]
func (s *MysicServer) GetData(w http.ResponseWriter, r *http.Request) {
    s.logger.Info("Getting data music from server" + r.URL.String())

    if r.Method!= "POST" {
        s.logger.Error("Error getting data from server" + " method not allowed")
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    defer r.Body.Close()

    page := r.URL.Query().Get("page")
    if page == "" {
        page = "1"
    }

    limit := r.URL.Query().Get("limit")
    if limit == "" {
        limit = "1000"
    }

    frstpg, err := strconv.Atoi(page)
    if err != nil || frstpg < 1 {
        s.logger.Error("Error converting page to integer", err)
        http.Error(w, "Invalid page", http.StatusBadRequest)
    }

    limcnt, err := strconv.Atoi(limit)
    if err != nil {
        s.logger.Error("Error converting limit to integer", err)
        http.Error(w, "Invalid limit", http.StatusBadRequest)
    }


	var filter models.FilterSong
	err = json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		s.logger.Error("Error decoding filter song from server" + err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

    songs, err := s.app.GetDataMusic(filter, frstpg, limcnt)
    if err!= nil {
        s.logger.Error("Error getting data from database" + err.Error())
        http.Error(w, "Failed to get data from database", http.StatusInternalServerError)
        return
    }

	json.NewEncoder(w).Encode(songs)
	w.Header().Set("Content-Type", "application/json")
	s.logger.Info("Data music returned to server" + r.URL.String())
}

// GetText godoc
// @Summary      Get Text 
// @Description  get text from database
// @Tags         text
// @Accept       json
// @Produce      json
// @Param        page query string true "first page"
// @Param        limit query string true "count page"
// @Param        song query string true "song name"
// @Success      200  {object} server.TextSong
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /text [post]
func (s *MysicServer) GetText(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Getting text from database " + r.URL.String())


    if r.Method!= "POST" {
        s.logger.Error("Error getting data from server" + " method not allowed")
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    defer r.Body.Close()

    page := r.URL.Query().Get("page")
    if page == "" {
        page = "1"
    }

    limit := r.URL.Query().Get("limit")
    if limit == "" {
        limit = "1000"
    }

    frstpg, err := strconv.Atoi(page)
    if err != nil || frstpg < 1 {
        s.logger.Error("Error converting page to integer", err)
        http.Error(w, "Invalid page", http.StatusBadRequest)
    }

    limcnt, err := strconv.Atoi(limit)
    if err != nil {
        s.logger.Error("Error converting limit to integer", err)
        http.Error(w, "Invalid limit", http.StatusBadRequest)
    }


	song := r.URL.Query().Get("song")
	if song == "" {
		s.logger.Debug("Error getting song from server" + " song not found")
        http.Error(w, "Song not found", http.StatusNotFound)
        return
	}

	text, err := s.app.GetTextSong(song, frstpg, limcnt)
	if err!= nil {
        s.logger.Error("Error getting text from database" + err.Error())
        http.Error(w, "Failed to get text from database", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(TextSong{Text: string(text)})
	w.Header().Set("Content-Type", "application/json")
	s.logger.Info("Text returned to server" + r.URL.String())
}

// Text is song 
// @Description Text song
type TextSong struct {
    Text string `json:"text"`
}

// DelSong godoc
// @Summary      Delete Song    
// @Description  delete song from database
// @Tags         deleted
// @Accept       json
// @Produce      json
// @Param        song query string true "song name"
// @Success      204 "success response"
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /delete [delete]
func (s *MysicServer) DeleteSong(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Deleting song from database " + r.URL.String())

    if r.Method != "DELETE" {
        s.logger.Error("Error deleting song from server" + " method not allowed")
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    song := r.URL.Query().Get("song")
    if song == "" {
        s.logger.Debug("Error getting song from server" + " song not found")
        http.Error(w, "Song not found", http.StatusNotFound)
        return
    }

    err := s.app.DeleteSong(song)
    if err!= nil {
        s.logger.Error("Error deleting song from database" + err.Error())
        http.Error(w, "Failed to delete song from database", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
    s.logger.Info("Song deleted from server" + r.URL.String())
}

// UpdateSong godoc
// @Summary      Update song 
// @Description  update song from database
// @Tags         update
// @Accept       json
// @Produce      json
// @Param        input body models.FilterSong true "update song"
// @Success      204 "success response"
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /update [post]
func (s *MysicServer) UpdateSong(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Update song from database  " + r.URL.String())

    if r.Method!= "POST" {
        s.logger.Error("Error updating song from server" + " method not allowed")
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	var song models.FilterSong
	err := json.NewDecoder(r.Body).Decode(&song)
	if err!= nil {
        s.logger.Debug("Error decoding song from server" + err.Error())
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	err = s.app.UpdateSong(song)
	if err!= nil {
        s.logger.Error("Error updating song from database" + err.Error())
        http.Error(w, "Failed to update song from database", http.StatusInternalServerError)
        return
    }

	w.WriteHeader(http.StatusNoContent)
	s.logger.Info("Song updated from server" + r.URL.String())
}

// CreateSong godoc
// @Summary      Create song 
// @Description  create song from database
// @Tags         create
// @Accept       json
// @Produce      json
// @Param        input body models.NewSong true "song struct"
// @Success      200 {object} server.NewID
// @Failure      400  "Bad request error"
// @Failure      404 "Not found error"
// @Failure      405 "Method not allowed"
// @Failure      500  "Internal server error"
// @Router       /create [post]
func (s *MysicServer) CreateSong(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Creating new song in database " + r.URL.String())
    
    if r.Method == "POST" {
        s.logger.Error("Error creating song from server" + " method not allowed")
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var newsong models.NewSong
    err := json.NewDecoder(r.Body).Decode(&newsong)
    if err != nil {
        s.logger.Debug("Error decoding new song from server" + err.Error())
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    id, err := s.app.CreateSong(newsong)
    if err != nil {
        s.logger.Error("Error creating song in database" + err.Error())
        http.Error(w, "Failed to create song in database", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(NewID{ID: id})
    w.Header().Set("Content-Type", "application/json")
    
	s.logger.Info("New song created in server " + r.URL.String())
}

// ID is song 
// @Description ID song
type NewID struct {
    ID uint64 `json:"id"`
}