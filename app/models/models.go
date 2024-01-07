package models

import (
	"gorm.io/datatypes"
)

type InternalTrack struct {
	ISRC           string         `gorm:"type:varchar(20);unique_index;primaryKey" json:"isrc"`
	SpotifyImage   string         `gorm:"column:spotify_image" json:"spotifyImage"`
	Title          string         `gorm:"type:varchar(255)" json:"title"`
	ArtistNameList datatypes.JSON `gorm:"type:json"`
	Popularity     int            `gorm:"type:integer" json:"popularity"`
	AccessToken    string         `gorm:"-" json:"accessToken"`
}

type SpotifyAPIResponse struct {
	Tracks Tracks `json:"tracks"`
}

type ExternalURLs struct {
	Spotify string `json:"spotify"`
}

type Artists struct {
	ExternalURLs ExternalURLs `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

type Album struct {
	AlbumType            string       `json:"album_type"`
	Artists              []Artists    `json:"artists"`
	AvailableMarkets     []string     `json:"available_markets"`
	ExternalURLs         ExternalURLs `json:"external_urls"`
	Href                 string       `json:"href"`
	ID                   string       `json:"id"`
	Images               []Image      `json:"images"`
	Name                 string       `json:"name"`
	ReleaseDate          string       `json:"release_date"`
	ReleaseDatePrecision string       `json:"release_date_precision"`
	TotalTracks          int          `json:"total_tracks"`
	Type                 string       `json:"type"`
	URI                  string       `json:"uri"`
}

type Image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type Items struct {
	Album        Album        `json:"album"`
	DiscNumber   int          `json:"disc_number"`
	DurationMs   int          `json:"duration_ms"`
	Explicit     bool         `json:"explicit"`
	ExternalIDs  ExternalIDs  `json:"external_ids"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	IsLocal      bool         `json:"is_local"`
	Name         string       `json:"name"`
	Popularity   int          `json:"popularity"`
	PreviewURL   string       `json:"preview_url"`
	TrackNumber  int          `json:"track_number"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}
type ExternalIDs struct {
	ISRC string `json:"isrc"`
}

type Tracks struct {
	Href     string  `json:"href"`
	Items    []Items `json:"items"`
	Limit    int     `json:"limit"`
	Next     string  `json:"next"`
	Offset   int     `json:"offset"`
	Previous string  `json:"previous"`
	Total    int     `json:"total"`
}
type CreateTrackReq struct {
	ISRC string `json:"isrc"`
}
type APIResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	ErrorCode int         `json:"error_code,omitempty"`
	Data      interface{} `json:"data"`
}
type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}
