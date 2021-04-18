package service

import (
	"hackathon/dao/mysql"
	"hackathon/models"
	"hackathon/response"
)

func AddChatRecord(p *models.ParamAddChatRecord, userId uint) int {
	newRecord := &models.ChatRecord{
		UserId:  userId,
		FromId:  uint(p.FromId),
		LastMsg: p.LastMsg,
		HasBg:   0,
	}
	exit := &models.ChatRecord{}
	mysql.GetDB().Where("user_id = ? AND from_id = ?", userId, p.FromId).First(&exit)
	if exit.ID > 0 {
		err := mysql.Update(exit, newRecord)
		if err != nil {
			return response.CodeServerBusy
		}
	} else {
		err := mysql.Create(newRecord)
		if err != nil {
			return response.CodeServerBusy
		}
	}

	return response.OK
}
