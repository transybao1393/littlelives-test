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

func UpdateQuotaUsedByUserIP(IP string) (bool, error) {
	//- get current quote + 1
	quota := &mongoRepository.UserInfo{
		UserIP: IP,
	}
	return quota.UpdateQuotaUsedByUserIP()
}

func IsOverQuota(IP string) (bool, error) {
	quota := &mongoRepository.UserInfo{
		UserIP: IP,
	}
	_, err := quota.GetQuotaByUserIP()
	if err != nil {
		return true, err
	}
	return false, nil
}
