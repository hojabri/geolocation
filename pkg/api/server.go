package api

import (
	"github.com/oschwald/geoip2-golang"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Server struct {
	log *zerolog.Logger
	db  *geoip2.Reader
}

func New(log *zerolog.Logger, db *geoip2.Reader) *Server {
	return &Server{
		log: log,
		db:  db,
	}
}

func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/geo/{ip}", s.getGeolocationInformation)
	return r
}

func (s *Server) getGeolocationInformation(writer http.ResponseWriter, request *http.Request) {
	ip := chi.URLParam(request, "ip")
	var err error
	if !validateIP(ip) {
		s.log.Debug().Msg("invalid IP")
		err = SendHTTPError(writer, "Invalid IP address", http.StatusBadRequest)
		if err != nil {
			s.log.Error().Err(err).Msg("can't send error message")
		}
		return
	}

	record, err := s.GetIpLocationInfo(ip)
	if err != nil {
		s.log.Error().Err(err).Str("IP", ip).Msg("failed to get information about IP")
		err = SendHTTPError(writer, "failed to get information about IP", http.StatusInternalServerError)
		if err != nil {
			s.log.Error().Err(err).Msg("can't send error message")
		}
		return
	}
	err = SendHTTPMessage(writer, record, http.StatusOK)
	if err != nil {
		s.log.Error().Err(err).Msg("can't send message")
	}
}

func validateIP(ip string) bool {
	if net.ParseIP(ip) != nil {
		return true
	}
	return false
}

func (s *Server) GetIpLocationInfo(ipStr string) (*geoip2.City, error) {
	ip := net.ParseIP(ipStr)
	var record *geoip2.City
	record, err := s.db.City(ip)
	if err != nil {
		return nil, err
	}

	return record, nil

}
