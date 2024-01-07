package lt

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lt/app/models"
	"lt/configs"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"gorm.io/datatypes"
)

var Auth = &oauth2.Config{
	ClientID:     configs.CLIENT_ID,
	ClientSecret: configs.CLIENT_SECRET,
	RedirectURL:  configs.REDIRECT_URI,
	Scopes:       configs.SCOPES,
	Endpoint:     configs.ENDPOINT,
}

func (m *Manager) LoginHandler(w http.ResponseWriter, r *http.Request) {
	State := GenerateRandomString(16)
	configs.STATE = State
	url := Auth.AuthCodeURL(State)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (m *Manager) CallbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")

	tok, err := Auth.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		return
	}
	if st := r.FormValue("state"); st != configs.STATE {
		http.NotFound(w, r)
		return
	}
	if tok.AccessToken == "" {
		m.Dao.log.Error("Access Token not generated")
	}
	m.log.Info(tok.AccessToken)
	configs.ACCESS_TOKEN = tok.AccessToken
	m.log.Info(tok.AccessToken)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
func (m *Manager) CreateTrackHandler(w http.ResponseWriter, r *http.Request) {
	var (
		req  models.CreateTrackReq
		err  error
		resp models.APIResponse
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		m.log.Errorf("error %v", err)
		resp.Message = fmt.Sprintf("error %v", err)
		resp.ErrorCode = http.StatusInternalServerError
		resp.Status = "error"
		result, err := json.Marshal(resp)
		if err != nil {
			m.log.Errorf("unable to marshal to json %v", err)
		}
		_, _ = w.Write(result)
		return
	}
	track, err := m.GetMetadata(req.ISRC)
	if err != nil {
		m.log.Errorf("%v", err)
		resp.Message = fmt.Sprintf("error %v", err)
		resp.ErrorCode = http.StatusNotFound
		resp.Status = "error"
		result, err := json.Marshal(resp)
		if err != nil {
			m.log.Errorf("unable to marshal to json %v", err)
		}
		_, _ = w.Write(result)
		return
	}
	resp.Status = "success"
	resp.Data = track
	resp.Message = "Spotify Meta Data  inserted"
	result, err := json.Marshal(resp)
	if err != nil {
		m.log.Errorf("unable to marshal to json %v", err)
	}
	_, _ = w.Write(result)

}

func (m *Manager) GetMetadata(isrc string) (track models.InternalTrack, err error) {
	url := fmt.Sprintf("https://api.spotify.com/v1/search?type=track&q=isrc:%s", isrc)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		m.log.Errorf("%v", err)
		return track, err
	}

	accessToken := fmt.Sprintf("Bearer %v", configs.ACCESS_TOKEN)
	req.Header.Set("Authorization", accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		m.log.Errorf("%v", err)
		return track, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		m.log.Errorf("%v", err)
		return track, err
	}

	spotifyResponse, err := ParseSpotifyResponse(body)
	if err != nil {
		m.log.Errorf("%v", err)
		return track, err
	}

	selectedTrack, err := FindTrackWithHighestPopularity(spotifyResponse.Tracks.Items)
	if err != nil {
		m.log.Errorf("%v", err)
		return track, err
	}
	selectedArtist, err := FindArtistWithHighestPopularity(selectedTrack)
	if err != nil {
		m.log.Errorf("%v", err)
		return track, err
	}

	track = models.InternalTrack{
		ISRC:         isrc,
		SpotifyImage: selectedArtist.Album.Images[0].URL,
		Title:        selectedArtist.Album.Name,
		Popularity:   selectedArtist.Popularity,
	}
	artistNamesJSON, err := json.Marshal(GetArtistNames(selectedArtist.Album.Artists))
	if err != nil {
		m.log.Errorf("%v", err)
		return track, err

	}
	track.ArtistNameList = datatypes.JSON(artistNamesJSON)
	result := m.Dao.DB.Create(track)
	if result.Error != nil {
		m.log.Errorf("%v", result.Error.Error())
		return track, result.Error
	}
	key := fmt.Sprintf("INTERNAL_ACCESS_TOKEN_%v", configs.ACCESS_TOKEN)
	isExists := m.Cache.Exists(key)
	if isExists == 0 {
		m.Cache.Set(key, GenerateSecureToken(200), 45*time.Minute)
	}
	secToken := m.Cache.Get(key)
	secToken = strings.ReplaceAll(secToken, "/", "")
	secToken = strings.ReplaceAll(secToken, `"`, "")
	track.AccessToken = secToken
	return track, nil
}
func (m *Manager) GetMataDataByIsrc(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp models.APIResponse
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	isrc := r.URL.Query().Get("isrc")
	authToken := r.Header.Get("auth-token")
	m.log.Info(authToken, " HEADERPART")
	key := fmt.Sprintf("INTERNAL_ACCESS_TOKEN_%v", configs.ACCESS_TOKEN)
	validAuthKey := ValidateAuthKey(m.log, m.Cache, authToken, key)
	if !validAuthKey {
		resp.Message = "invalid authCode"
		resp.ErrorCode = http.StatusNotFound
		resp.Status = "error"
		result, err := json.Marshal(resp)
		if err != nil {
			m.log.Errorf("unable to marshal to json %v", err)
		}
		_, _ = w.Write(result)
		return
	}
	track, err := m.Dao.GetMetaDataByIsrc(isrc)
	if err != nil {
		m.log.Errorf("%v", err)
		resp.Message = fmt.Sprintf("error %v", err)
		resp.ErrorCode = http.StatusNotFound
		resp.Status = "error"
		result, err := json.Marshal(resp)
		if err != nil {
			m.log.Errorf("unable to marshal to json %v", err)
		}
		_, _ = w.Write(result)
		return
	}
	resp.Status = "success"
	resp.Data = track
	resp.Message = "Served  Spotify Meta Data based on Isrc"
	result, err := json.Marshal(resp)
	if err != nil {
		m.log.Errorf("unable to marshal to json %v", err)
	}
	_, _ = w.Write(result)
}
func (m *Manager) GetMataDataByArtist(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		resp models.APIResponse
	)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	artist := r.URL.Query().Get("artist")
	artistNames := strings.Split(artist, ",")
	key := fmt.Sprintf("INTERNAL_ACCESS_TOKEN_%v", configs.ACCESS_TOKEN)
	validAuthKey := ValidateAuthKey(m.log, m.Cache, r.Header.Get("auth-token"), key)
	if !validAuthKey {
		resp.Message = "invalid authCode"
		resp.ErrorCode = http.StatusNotFound
		resp.Status = "error"
		result, err := json.Marshal(resp)
		if err != nil {
			m.log.Errorf("unable to marshal to json %v", err)
		}
		_, _ = w.Write(result)
		return
	}
	track, err := m.Dao.GetMetaDataByArtist(artistNames)
	if err != nil {
		m.log.Errorf("%v", err)
		resp.Message = fmt.Sprintf("error %v", err)
		resp.ErrorCode = http.StatusNotFound
		resp.Status = "error"
		result, err := json.Marshal(resp)
		if err != nil {
			m.log.Errorf("unable to marshal to json %v", err)
		}
		_, _ = w.Write(result)
		return
	}
	resp.Status = "success"
	resp.Data = track
	resp.Message = "Served  Spotify Meta Data based on Artist"
	result, err := json.Marshal(resp)
	if err != nil {
		m.log.Errorf("unable to marshal to json %v", err)
	}
	_, _ = w.Write(result)
}
