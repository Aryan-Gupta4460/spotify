package lt

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type InternalController struct {
	log *zap.SugaredLogger
	Mgr *Manager
}

func NewController(logger *zap.SugaredLogger, mgr *Manager) *InternalController {
	return &InternalController{Mgr: mgr, log: logger}
}
func (c *InternalController) SetRouteHandlers(router *mux.Router) {
	v1Router := router.PathPrefix("/v1/lt/spotify").Subrouter()
	v1Router.HandleFunc("/login", c.Mgr.LoginHandler).Methods("GET")
	v1Router.HandleFunc("/callback", c.Mgr.CallbackHandler).Methods("GET")
	v1Router.HandleFunc("/get_byisrc", c.Mgr.GetMataDataByIsrc).Methods("GET")
	v1Router.HandleFunc("/get_byartist", c.Mgr.GetMataDataByArtist).Methods("GET")
	v1Router.HandleFunc("/create_track", c.Mgr.CreateTrackHandler).Methods("POST")

}
