package app

import (
	"client"
	"context"
	"fmt"
	"net/http"

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

func (a *App) GetDataMusic(filter models.FilterSong) ([]models.Song, error) {
	a.logger.Info("GetDataMusic called with filter " + fmt.Sprintf(" %v", filter))
	
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
        a.logger.Debug("Error getting songs" + fmt.Sprintf(" %v", filter))
        return nil, fmt.Errorf("failed to get songs: %w", err)
    }

    if len(songs) == 0 {
        a.logger.Debug("No songs found with filter" + fmt.Sprintf(" %v", filter))
        return nil, fmt.Errorf("no songs found with filter: %w", err)
    }

    a.logger.Info("GetDataMusic complete with" + fmt.Sprintf(" %d songs", len(songs)))
    return songs, nil
}

func (a *App) GetTextSong(song string) ([]byte, error) {
	a.logger.Info("GetTextSong called with song " + fmt.Sprintf(" %s", song))

	text, err := a.db.GetText(song)
	if err!= nil {
        a.logger.Debug("Error getting text for song" + fmt.Sprintf(" %s", song))
        return nil, fmt.Errorf("failed to get text for song: %w", err)
    }

	if len(text) == 0 {
        a.logger.Debug("Text not found for song" + fmt.Sprintf(" %s", song))
        return nil, fmt.Errorf("text not found for song: %w", err)
    }

	a.logger.Info("GetTextSong complete" + fmt.Sprintf(" %s", song))
	return text, nil
}

func (a *App) DeleteSong(song string) (error) {
	a.logger.Info("DeleteSong called with song" + fmt.Sprintf(" %s", song))

	err := a.db.DeleteSong(song)
	if err!= nil {
        a.logger.Debug("Error deleting song" + fmt.Sprintf(" %s", song))
        return fmt.Errorf("failed to delete song: %w", err)
    }

	a.logger.Info("Song deleted" + fmt.Sprintf(" %s", song))
	return nil
}

func (a *App) UpdateSong(song models.FilterSong) (error) {
	a.logger.Info("UpdateSong called with song" + fmt.Sprintf(" %s", song))

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
        a.logger.Debug("Error updating song" + fmt.Sprintf(" %s", song.Song))
        return fmt.Errorf("failed to update song: %w", err)
    }

	a.logger.Info("Song updated" + fmt.Sprintf(" %s", song.Song))
	return nil
}

func (a *App) CreateSong(newsong models.NewSong) (uint64, error) {
	a.logger.Info("Creating new song" + fmt.Sprintf(" %s %s", newsong.Group, newsong.Song))

	resp, err := a.client.GetInfoWithResponse(context.TODO(), &client.GetInfoParams{Group: newsong.Group, Song: newsong.Song})
	if err != nil {
		a.logger.Debug("Error getting info" + fmt.Sprintf(" %s %s", newsong.Group, newsong.Song))
		return 0, fmt.Errorf("failed to get song info: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		a.logger.Debug("Expected HTTP 200 but recerved" + fmt.Sprintf(" %s %s %d", newsong.Group, newsong.Song, resp.StatusCode()))
        return 0, fmt.Errorf("failed to get song info: status code %d", resp.StatusCode())
	}

	id, err := a.db.SaveMusic(newsong, *resp.JSON200)
	if err != nil {
        a.logger.Debug("Error saving music" + fmt.Sprintf(" %s %s", newsong.Group, newsong.Song))
        return 0, fmt.Errorf("failed to save song: %w", err)
    }

	a.logger.Info("New song created" + fmt.Sprintf(" %s %s ID: %d", newsong.Group, newsong.Song, id))
	return id, nil
}