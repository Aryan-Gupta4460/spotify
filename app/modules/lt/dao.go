package lt

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"lt/app/client/cache"
	"lt/app/models"
	"net/http"
	"time"

	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DaoRepo struct {
	httpClient *http.Client
	log        *zap.SugaredLogger
	DB         *gorm.DB
	Cache      cache.Cache
}

func NewDaoRepo(log *zap.SugaredLogger, db *gorm.DB, cache cache.Cache) *DaoRepo {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
	return &DaoRepo{
		log:        log,
		DB:         db,
		httpClient: httpClient,
		Cache:      cache,
	}
}

func (d *DaoRepo) GetMetaDataByIsrc(isrc string) (track models.InternalTrack, err error) {
	result := d.DB.Where("isrc = ?", isrc).First(&track)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return track, errors.New("track not found")
		}
		return track, result.Error
	}
	return track, nil
}
func (d *DaoRepo) GetMetaDataByArtist(artistNames []string) (tracks []models.InternalTrack, err error) {
	artistNamesJSON, err := json.Marshal(artistNames)
	if err != nil {
		d.log.Errorf("%v", err)
		return tracks, err

	}
	result := d.DB.Where("artist_name_list @> ?", datatypes.JSON(artistNamesJSON)).Find(&tracks)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return tracks, errors.New("track not found")
		}
		return tracks, result.Error
	}
	return tracks, nil
}
