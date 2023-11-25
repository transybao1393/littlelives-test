package router

import (
	"bytes"
	"fmt"
	"io"
	"ll_test/app/logger"
	"ll_test/app/utils"
	"ll_test/domain"
	"net/http"

	usecases "ll_test/ll/usecase"

	"github.com/go-chi/render"
	"golang.org/x/exp/slices"
)

var log = logger.NewLogrusLogger()

const (
	MB        = 1 << 20 //- 1MB
	fileLimit = 5 * MB  // 5 MB
)

func LLAddNewUser(w http.ResponseWriter, r *http.Request) error {
	userIP := utils.UserIPHandling(r)
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

	//- response
	statusCode := http.StatusCreated
	w.WriteHeader(statusCode)
	render.JSON(w, r, &domain.Response{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
		Data:       nil,
	},
	)
	return nil
}

// - using form
func LLFileUpload(w http.ResponseWriter, r *http.Request) error {
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

	// validation media type is video
	//- FIXME: Move this to middleware
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
	userIP := utils.UserIPHandling(r)
	isOverLimit, _ := usecases.IsOverQuota(userIP)
	if !isOverLimit {
		usecases.SaveToMinIO(
			file,
			handler.Header.Get("Content-Type"),
			fileBuf,
			handler.Filename,
			handler.Size,
			userIP, //- user IP
		)
		_, err = usecases.UpdateQuotaUsedByUserIP(userIP)
		if err != nil {
			fields := logger.Fields{
				"service": "ll",
				"message": "Error when update quota used",
			}
			log.Fields(fields).Errorf(err, "Error when update quota used")
			return err
		}
	} else {
		//- response
		statusCode := http.StatusBadRequest
		w.WriteHeader(statusCode)
		render.JSON(w, r, &domain.Response{
			StatusCode: statusCode,
			Message:    http.StatusText(statusCode),
			Data:       "Quota exceeded! Please buy more quota",
		},
		)
		return nil
	}

	//- response
	statusCode := http.StatusOK
	w.WriteHeader(statusCode)
	render.JSON(w, r, &domain.Response{
		StatusCode: statusCode,
		Message:    http.StatusText(statusCode),
		Data:       nil,
	},
	)
	return nil
}
