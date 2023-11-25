package usecase

import (
	"ll_test/app/logger"
	mongoRepository "ll_test/ll/repository/mongodb"
)

func UserQuotaSet(IP string, quotaNumber int) error {
	//- add new user with IP
	user := &mongoRepository.UserInfo{
		UserIP:    IP,
		Quota:     quotaNumber, //- 10 files
		QuotaUsed: 0,
	}
	err := user.AddNewUser()
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

func UpdateQuotaUsedByUserIP(IP string, quotaUsed int) bool {
	//- get current quote + 1
	quota := &mongoRepository.UserInfo{
		UserIP:    IP,
		QuotaUsed: quotaUsed,
	}
	return quota.UpdateQuotaUsedByUserIP()
}
