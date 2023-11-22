package router

import (
	"bytes"
	"fmt"
	"io"
	"ll_test/app/logger"
	"ll_test/app/utils"
	"net/http"

	minioUseCase "ll_test/ll/usecase"

	"golang.org/x/exp/slices"
)

var log = logger.NewLogrusLogger()

const (
	MB = 1 << 20 //- 10MB
)

// - using form
func LLVideoUploadFile(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("LLVideoUploadFile...")
	//- receive file
	// message := "Video upload success"

	//- limit to 5mb per file
	if err := r.ParseMultipartForm(10 * MB); err != nil {
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

	//- convert to buffer
	fileBuf := bytes.NewBuffer(nil)
	if _, err = io.Copy(fileBuf, file); err != nil {
		fields := logger.Fields{
			"service": "ll",
			"message": "Error when copy file to buffer",
		}
		log.Fields(fields).Errorf(err, "Error when copy file to buffer")
		return err
	}

	//- call to save to minio
	minioUseCase.SaveToMinIO(file, handler.Header.Get("Content-Type"), fileBuf, handler.Filename, handler.Size)

	//- save file information to mongodb
	//- return result

	return nil
}
