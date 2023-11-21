package router

import (
	"fmt"
	"ll_test/app/logger"
	"ll_test/app/utils"
	"net/http"

	"golang.org/x/exp/slices"
)

var log = logger.NewLogrusLogger()

const (
	MB = 10 << 20 //- 10MB
)

// - using form
func LLVideoUploadFile(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("LLVideoUploadFile...")
	//- receive file
	// message := "Video upload success"

	//- limit to 5mb per file
	if err := r.ParseMultipartForm(5 * MB); err != nil {
		fields := logger.Fields{
			"service": "Youtube",
			"message": "Error when parse multipart form",
		}
		log.Fields(fields).Errorf(err, "Error when parse multipart form")
		return err
	}

	// Limit upload size
	r.Body = http.MaxBytesReader(w, r.Body, 5*MB) // 5 Mb
	fmt.Printf("r.Body: %+v\n", r.Body)

	file, handler, err := r.FormFile("file")
	if err != nil {
		fields := logger.Fields{
			"service": "ll",
			"message": "Error when receive form file from request",
		}
		log.Fields(fields).Errorf(err, "Error when receive form file from request")
		return err
	}
	defer file.Close()

	fmt.Printf("file content type: %+v\n", handler.Header.Get("Content-Type"))

	// validation media type is video
	if !slices.Contains(utils.FileContentType, handler.Header.Get("Content-Type")) {
		fields := logger.Fields{
			"service": "ll",
			"message": "This is not a valid file",
		}
		log.Fields(fields).Errorf(err, "This is not a valid file")
		return fmt.Errorf("not allowed file type")
	}

	//-

	return nil
}
