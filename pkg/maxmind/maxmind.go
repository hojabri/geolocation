package maxmind

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/hojabri/geolocation/pkg/config"
	"github.com/jasonlvhit/gocron"
	"github.com/oschwald/geoip2-golang"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	log        *zerolog.Logger
	licenseKey string
	DB         *geoip2.Reader
}

func New(log *zerolog.Logger, licenseKey string) Server {
	return Server{
		log:        log,
		licenseKey: licenseKey,
	}
}

func (s *Server) OpenDB() error {
	db, err := geoip2.Open(fmt.Sprintf("geodb/%s", config.Configuration.GetString("GEO_CITY_DB")))
	if err != nil {
		return err
	}
	s.DB = db
	return nil
}

func (s *Server) CloseDB() error {
	err := s.DB.Close()
	if err != nil {
		return err
	}
	return nil
}
func (s *Server) RunDownloadScheduler() error {
	if err := gocron.Every(1).Wednesday().Do(s.DownloadDB); err != nil {
		return err
	}
	return nil
}

func (s *Server) DownloadDB() error {
	s.log.Info().Msg("Downloading maxmind db...")
	dbUrl := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key={{LICENSE_KEY}}&suffix=tar.gz"
	dbUrl = strings.Replace(dbUrl, "{{LICENSE_KEY}}", s.licenseKey, -1)

	// Download
	if err := downloadFile("geodb/db.tar.gz", dbUrl); err != nil {
		s.log.Err(err).Msg("Error in downloading maxmind db")
		return err
	}

	// Decompress
	r, err := os.Open("geodb/db.tar.gz")
	if err != nil {
		s.log.Err(err).Msg("Error in opening db.tar.gz file")
		return err
	}
	defer r.Close()
	extractedDir, err := extractTarGz(r, "geodb")
	if err != nil {
		s.log.Err(err).Msg("Error in extracting db.tar.gz file")
		return err
	}

	// move extracted file to main directory
	err = os.Rename("geodb/"+extractedDir+"GeoLite2-City.mmdb", "geodb/GeoLite2-City.mmdb")
	if err != nil {
		s.log.Err(err).Msg("Error in moving extracted GeoLite2-City.mmdb file to geodb directory")
		return err
	}

	// delete temp decompressed directory
	err = os.RemoveAll("geodb/" + extractedDir)
	if err != nil {
		s.log.Err(err).Msg("Error in deleting the temp directory")
		return err
	}

	// delete temp compressed file
	err = os.Remove("geodb/db.tar.gz")
	if err != nil {
		s.log.Err(err).Msg("Error in deleting the compressed file")
		return err
	}

	return nil

}

func downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func extractTarGz(gzipStream io.Reader, target string) (string, error) {
	var extractedDir string
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return "", err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(target+"/"+header.Name, 0755); err != nil {
				return "", err
			}
			extractedDir = header.Name
		case tar.TypeReg:
			outFile, err := os.Create(target + "/" + header.Name)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return "", err
			}
			outFile.Close()

		default:
			return "", errors.New(fmt.Sprintf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name))
		}

	}

	return extractedDir, nil
}
