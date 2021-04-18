package mysql

import (
	"hackathon/models"
)

func QueryHistory(ReceiveData []models.ChatHistory, Param models.ChatHistory, LastId int) []models.ChatHistory {

	MD := GetDB()
	if LastId == 0 {

		MD.Raw("select count(*) from chat_histories").Scan(&LastId)

	}
	MD.Where("id<?&&id>?", LastId, LastId-1000).Where(&Param).Find(&ReceiveData)
	return ReceiveData

}
