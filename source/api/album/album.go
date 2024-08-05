package album

import (
	"encoding/json"
	"sync"
	"time"
)

// watch out path injection
type AlbumKey string
type MusicKey string

func (k AlbumKey) IsEmpty() bool { return len(k) == 0 }
func (k MusicKey) IsEmpty() bool { return len(k) == 0 }

type Album struct {
	key   AlbumKey
	date  time.Time
	title string
	music []Music
	cover []byte
}

type albumJson struct {
	Key   AlbumKey `json:"key"`
	Date  int64    `json:"date"` //epoch seconds
	Title string   `json:"title"`
	Music []Music  `json:"music"`
	Cover []byte   `json:"cover"`
}

func (a *Album) Key() AlbumKey   { return a.key }
func (a *Album) Date() time.Time { return a.date }
func (a Album) Title() string    { return a.title }
func (a *Album) Music() []Music  { return a.music }
func (a *Album) Cover() []byte   { return a.cover }
func (a *Album) Count() int      { return len(a.music) }

func (a *Album) MarshalJSON() ([]byte, error) {
	return json.Marshal(albumJson{a.key, a.date.Unix(), a.title, a.music, a.cover})
}

func (a *Album) UnmarshalJSON(data []byte) error {
	var buf albumJson
	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}
	*a = Album{buf.Key, time.Unix(buf.Date, 0), buf.Title, buf.Music, buf.Cover}
	return nil
}

type Music struct {
	date     time.Time
	title    string
	length   time.Duration
	platform string
	id       string
}

type musicJson struct {
	Date     int64  `json:"date"` //epoch seconds
	Title    string `json:"title"`
	Length   string `json:"length"`
	Platform string `json:"platform"`
	ID       string `json:"id"`
}

func (m *Music) Key() MusicKey         { return MusicKey(m.platform + m.id) }
func (m *Music) Date() time.Time       { return m.date }
func (m Music) Title() string          { return m.title }
func (m *Music) Length() time.Duration { return m.length }

func (m *Music) MarshalJSON() ([]byte, error) {
	return json.Marshal(musicJson{m.date.Unix(), m.title, m.length.String(), m.platform, m.id})
}

func (m *Music) UnmarshalJSON(data []byte) error {
	var buf musicJson
	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}
	length, err := time.ParseDuration(buf.Length)
	if err != nil {
		return err
	}

	*m = Music{time.Unix(buf.Date, 0), buf.Title, length, buf.Platform, buf.ID}
	return nil
}

// id: album
type albumMap struct {
	sync.Map
}

func (m *albumMap) MarshalJSON() ([]byte, error) {
	mp := map[string]Album{}
	m.Range(func(key, value any) bool {
		mp[key.(string)] = value.(Album)
		return true
	})
	return json.Marshal(mp)
}

func (m *albumMap) UnmarshalJSON(data []byte) error {
	mp := map[string]Album{}
	if err := json.Unmarshal(data, &mp); err != nil {
		return err
	}
	for key, val := range mp {
		m.Store(key, val)
	}
	return nil
}
