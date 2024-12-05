package app

import (
	"client"
	"context"
	"fmt"
	"net/http"
	"strings"

	"log/slog"
	"musicservice/interal/models"
	"musicservice/pkg/sql/postgres"
)

type App struct {
	logger *slog.Logger
	db *postgres.Postgres
	client *client.ClientWithResponses
}

func NewApp(log *slog.Logger, db *postgres.Postgres, client *client.ClientWithResponses) *App {
    return &App{logger: log, db: db, client: client}
}

func (a *App) GetDataMusic(filter models.FilterSong, frstpg, limcnt int) ([]models.Song, error) {
	log := a.logger.With(
		slog.String("OP", "GetDataMusic"),
	)

	log.Info("GetDataMusic called with filter " + fmt.Sprintf(" %v", filter))
	
	filtermap := make(map[string]string)
	
	if filter.Group != "" {
        filtermap["group"] = filter.Group
    }
	if filter.Song != "" {
        filtermap["song"] = filter.Song
    } 
	if filter.Link != "" {
        filtermap["link"] = filter.Link
    }
	if filter.ReleaseDate != "" {
        filtermap["releasedate"] = filter.ReleaseDate
    }
	if filter.Text != "" {
        filtermap["text"] = filter.Text
    }

	

    songs, err := a.db.GetSongs(filtermap)
    if err!= nil {
        log.Error("Error getting songs" + fmt.Sprintf(" %v", filter))
        return nil, fmt.Errorf("failed to get songs: %w", err)
    }

    if len(songs) == 0 {
        log.Error("No songs found with filter" + fmt.Sprintf(" %v", filter))
        return nil, fmt.Errorf("no songs found with filter: %w", err)
    }

	for i, s := range songs {
		songs[i].Text, err = Pangination(s.Text, frstpg, limcnt)
		if err != nil {
			return nil, fmt.Errorf("failed to paginate song text: %w", err)
		}
    }
	
    log.Info("GetDataMusic complete with" + fmt.Sprintf(" %d songs", len(songs)))
    return songs, nil
}

func (a *App) GetTextSong(song string, page, limit int) ([]byte, error) {
	log := a.logger.With(
		slog.String("OP", "GetTextSong"),
	)
	log.Info("GetTextSong called with song " + fmt.Sprintf(" %s", song))

	text, err := a.db.GetText(song)
	if err!= nil {
        log.Error("Error getting text for song" + fmt.Sprintf(" %s", song))
        return nil, fmt.Errorf("failed to get text for song: %w", err)
    }

	if len(text) == 0 {
        log.Debug("Text not found for song" + fmt.Sprintf(" %s", song))
        return nil, fmt.Errorf("text not found for song: %w", err)
    }

	texts, err := Pangination(string(text), page, limit)
	if err!= nil {
        return nil, fmt.Errorf("failed to paginate song text: %w", err)
    }

	text = []byte(texts)

	log.Info("GetTextSong complete" + fmt.Sprintf(" %s", song))
	return text, nil
}

func (a *App) DeleteSong(song string) (error) {
	log := a.logger.With(
		slog.String("OP", "DeleteSong"),
	)
	log.Info("DeleteSong called with song" + fmt.Sprintf(" %s", song))

	err := a.db.DeleteSong(song)
	if err!= nil {
        log.Debug("Error deleting song" + fmt.Sprintf(" %s", song))
        return fmt.Errorf("failed to delete song: %w", err)
    }

	log.Info("Song deleted" + fmt.Sprintf(" %s", song))
	return nil
}

func (a *App) UpdateSong(song models.FilterSong) (error) {
	log := a.logger.With(
		slog.String("OP", "UpdateSong"),
	)

	log.Info("UpdateSong called with song" + fmt.Sprintf(" %s", song))

	songmap := make(map[string]string)
	
	if song.Group != "" {
        songmap["group"] = song.Group
    }
	if song.Song != "" {
        songmap["song"] = song.Song
    } else {
		return fmt.Errorf("song name is required")
	}
	if song.Link != "" {
        songmap["link"] = song.Link
    }
	if song.ReleaseDate != "" {
        songmap["releasedate"] = song.ReleaseDate
    }
	if song.Text != "" {
        songmap["text"] = song.Text
    }

	err := a.db.UpdateSong(songmap)
	if err != nil {
        log.Debug("Error updating song" + fmt.Sprintf(" %s", song.Song))
        return fmt.Errorf("failed to update song: %w", err)
    }

	log.Info("Song updated" + fmt.Sprintf(" %s", song.Song))
	return nil
}

func (a *App) CreateSong(newsong models.NewSong) (uint64, error) {
	log := a.logger.With(
		slog.String("OP", "CreateSong"),
	)

	log.Info("Creating new song" + fmt.Sprintf(" %s %s", newsong.Group, newsong.Song))

	resp, err := a.client.GetInfoWithResponse(context.TODO(), &client.GetInfoParams{Group: newsong.Group, Song: newsong.Song})
	if err != nil {
		log.Debug("Error getting info" + fmt.Sprintf(" %s %s", newsong.Group, newsong.Song))
		return 0, fmt.Errorf("failed to get song info: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		log.Debug("Expected HTTP 200 but recerved" + fmt.Sprintf(" %s %s %d", newsong.Group, newsong.Song, resp.StatusCode()))
        return 0, fmt.Errorf("failed to get song info: status code %d", resp.StatusCode())
	}

	id, err := a.db.SaveMusic(newsong, *resp.JSON200)
	if err != nil {
        log.Debug("Error saving music" + fmt.Sprintf(" %s %s", newsong.Group, newsong.Song))
        return 0, fmt.Errorf("failed to save song: %w", err)
    }

	log.Info("New song created" + fmt.Sprintf(" %s %s ID: %d", newsong.Group, newsong.Song, id))
	return id, nil
}

func Pangination(text string, page, limit int) (string, error) {
	paragr := strings.Split(text, "\n\n")

	if limit > len(paragr) {
		paragr = paragr[page-1:]
		return strings.Join(paragr, "\n\n") + "\n\n... (Page " + fmt.Sprintf("%d", page) + " of " + fmt.Sprintf("%d", len(paragr)) + ")", nil
	}
	
	return strings.Join(paragr[page-1:page+limit-1], "\n\n") + "\n\n... (Page " + fmt.Sprintf("%d", page) + " of " + fmt.Sprintf("%d", page + limit) + ")", nil
}