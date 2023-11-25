package router

import (
	"bytes"
	"fmt"
	"io"
	"ll_test/app/logger"
	"ll_test/app/utils"
	"net/http"

	usecases "ll_test/ll/usecase"

	"golang.org/x/exp/slices"
)

var log = logger.NewLogrusLogger()

const (
	MB        = 1 << 20 //- 1MB
	fileLimit = 5 * MB  // 5 MB
)

func LLAddNewUser(w http.ResponseWriter, r *http.Request) error {
	userIP := r.RemoteAddr
	userQuota := 10 //- 10 files
	err := usecases.UserQuotaSet(userIP, userQuota)
	if err != nil {
		fields := logger.Fields{
			"service": "littlelives",
			"message": "Error when add new user",
		}
		log.Fields(fields).Errorf(err, "Error when add new user")
		return err
	}
	return nil
}

// - using form
func LLVideoUploadFile(w http.ResponseWriter, r *http.Request) error {
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
	r.Body = http.MaxBytesReader(w, r.Body, fileLimit) // 5 Mb
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
	usecases.SaveToMinIO(
		file,
		handler.Header.Get("Content-Type"),
		fileBuf,
		handler.Filename,
		handler.Size,
		r.RemoteAddr, //- user IP
	)

	//- save file information to mongodb
	//- return result

	return nil
}
