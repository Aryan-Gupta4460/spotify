package lt

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"lt/app/client/cache"
	"lt/app/models"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Manager struct {
	httpClient *http.Client
	log        *zap.SugaredLogger
	db         *gorm.DB
	Dao        *DaoRepo
	Cache      cache.Cache
}

func NewManager(log *zap.SugaredLogger, db *gorm.DB, cache cache.Cache) *Manager {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
	manager := &Manager{
		httpClient: httpClient,
		log:        log,
		db:         db,
		Dao:        NewDaoRepo(log, db, cache),
		Cache:      cache,
	}
	return manager
}
func FindArtistWithHighestPopularity(artists []models.Items) (models.Items, error) {
	if len(artists) == 0 {
		return models.Items{}, errors.New("no artists found")
	}

	highestPopularity := -1
	var selectedArtist models.Items
	for _, artist := range artists {
		if artist.Popularity > highestPopularity {
			highestPopularity = artist.Popularity
			selectedArtist = artist
		}
	}

	return selectedArtist, nil
}
func FindTrackWithHighestPopularity(tracks []models.Items) ([]models.Items, error) {
	if len(tracks) == 0 {
		return []models.Items{}, errors.New("no tracks found")
	}

	highestPopularity := -1
	var selectedTrack []models.Items
	for _, track := range tracks {
		if track.Popularity > highestPopularity {
			highestPopularity = track.Popularity
			selectedTrack = append(selectedTrack, track)
		}

	}

	return selectedTrack, nil
}
func GetArtistNames(artists []models.Artists) (names []string) {
	for _, artist := range artists {
		names = append(names, artist.Name)
	}
	return names
}

func ParseSpotifyResponse(body []byte) (models.SpotifyAPIResponse, error) {
	var spotifyResponse models.SpotifyAPIResponse
	if err := json.Unmarshal(body, &spotifyResponse); err != nil {
		return models.SpotifyAPIResponse{}, err
	}
	return spotifyResponse, nil
}
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
func ValidateAuthKey(log *zap.SugaredLogger, cache cache.Cache, authKey, redisKey string) (valid bool) {
	authKey = `"` + authKey + `"`
	redisKeyValue := cache.Get(redisKey)
	if redisKeyValue == authKey {
		valid = true
		return
	}
	log.Infof("Auth token Check Failed : redisKey : %v, redisKeyValue : %v, authToken : %v", redisKey, redisKeyValue, authKey)
	return
}
