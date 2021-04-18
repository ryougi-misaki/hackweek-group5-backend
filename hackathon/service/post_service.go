package service

import (
	"hackathon/dao/mysql"
	"hackathon/models"
	"hackathon/response"
)

func CreateTag(p *models.ParamCreateTag) int {
	newTag := &models.Tag{
		Name:        p.Name,
		Description: p.Description,
	}
	err := mysql.Create(newTag)
	if err != nil {
		return response.CodeServerBusy
	}
	//返回结果
	return response.OK
}

func CreatePost(p *models.ParamCreatePost, id uint) int {
	newPost := &models.Post{
		UserId:  id,
		TagId:   uint(p.TagId),
		Title:   p.Title,
		Content: p.Content,
		Status:  0,
	}
	err := mysql.Create(newPost)
	if err != nil {
		return response.CodeServerBusy
	}
	//返回结果
	return response.OK
}
