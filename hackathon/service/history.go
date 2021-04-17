package service

import (
	"hackathon/dao/mysql"
	"hackathon/models"
)

func History(history models.ChatHistory, last int) ([]models.ChatHistory, []models.ChatHistory) {

	var From []models.ChatHistory
	var To []models.ChatHistory

	var data = models.ChatHistory{}

	data.From = history.From
	data.To = history.To

	From = mysql.QueryHistory(From,data,last)

	data.From = history.To
	data.To = history.From
	To = mysql.QueryHistory(To,data,last)

	return From , To

}
