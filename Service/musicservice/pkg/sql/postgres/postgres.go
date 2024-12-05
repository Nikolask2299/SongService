package postgres

import (
	"client"
	"database/sql"
	"fmt"
	"musicservice/interal/models"
	"strings"

	_ "github.com/lib/pq"
)

type Postgres struct {
    db *sql.DB
}

func NewPostgres(user, password, dbname, host, port string) (*Postgres, error) {
	psqlInfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)
    db, err := sql.Open("postgres", psqlInfo)
    if err!= nil {
        return nil, err
    }

    err = db.Ping()
    if err!= nil {
        return nil, err
    }

    return &Postgres{db: db}, nil
}


func (p *Postgres) Close() error {
    return p.db.Close()
}

func (p *Postgres) GetSongs(filter map[string]string) ([]models.Song ,error) {
    query := `SELECT songs.id, songs.group, songs.song, to_char(songs.releasedate, 'DD.MM.YYYY'), songs.text, songs.link FROM songs WHERE`
    
    for k, v := range filter {
       if k == "releasedate" {
            query += ` AND "` + k + `" = to_date('` + v + `', 'DD.MM.YYYY')`
            continue
       }

       if k == "text" {
            query += ` AND ` + `make_tsvector(songs.text) @@ plainto_tsquery($$` + v + `$$)`
            continue
       }

        query += ` AND "` + k + `"`
        query += ` = '` + v + `'`
    }
    query += ";"
    
    query = strings.Replace(query, "WHERE AND", "WHERE ", -1)

    rows, err := p.db.Query(query)
    if err!= nil {
        return nil, err
    }
    defer rows.Close()

    songs := make([]models.Song, 0, 10)
    for rows.Next() {
        var song models.Song
        err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
        if err != nil {
            return nil, err
        }
        songs = append(songs, song)
    }
    return songs, nil
}

func (p *Postgres) GetText(song string) ([]byte, error) {
    query := `SELECT "text" FROM songs WHERE "song" = '` + song + `';`
    var text []byte
    err := p.db.QueryRow(query).Scan(&text)
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("song not found")
    } else if err != nil {
        return nil, err
    }
    return text, nil
}

func (p *Postgres) UpdateSong(song map[string]string) error {
    query := `UPDATE songs SET `

    for i, v := range song {
        query += `, "` + i + `"`
        query += " = '" + v + "' "
    }

    query += ` WHERE "song" = '` + song["song"] + `';`

    query = strings.Replace(query, "SET ,", "SET ", -1)
    _, err := p.db.Exec(query)
    return err
}

func (p *Postgres) DeleteSong(song string) error {
    query := `DELETE FROM songs WHERE "song" = '`+ song +`';`
    _, err := p.db.Exec(query)
    return err
}

func (p *Postgres) SaveGroup(songs string) error {
    query := `INSERT INTO groups("group") VALUES ($1)
        ON CONFLICT ("group") DO NOTHING;`

    _, err := p.db.Exec(query, songs)
    return err
}

func (p *Postgres) SaveMusic(song models.NewSong, data client.SongDetail) (uint64, error) {
    query := `INSERT INTO songs("group", "song", "releasedate", "text", "link")
        VALUES ($1, $2, to_date($3, 'DD.MM.YYYY'), $4, $5)
        RETURNING id;
    `
    err := p.SaveGroup(song.Group)
    if err != nil {
        return 0, err
    }

    var id uint64
    err = p.db.QueryRow(query, song.Group, song.Song, data.ReleaseDate, data.Text, data.Link).Scan(&id)
    if err!= nil {
        return 0, err
    }
    return id, nil
}