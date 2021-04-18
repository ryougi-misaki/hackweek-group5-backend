package models

// ParamRegister 获取注册请求application/json 参数
type ParamRegister struct {
	Name      string `json:"name" form:"name"`                              //可以为空
	Telephone string `json:"telephone" form:"telephone" binding:"required"` //手机号
	Password  string `json:"password" form:"password" binding:"required"`   //密码
}

// ParamLogin 获取登录请求application/json参数
type ParamLogin struct {
	Telephone string `json:"telephone" form:"telephone" binding:"required"` //手机号
	Password  string `json:"password" form:"password" binding:"required"`   //密码
}

// ParamEditInfo 获取编辑用户信息请求application/json参数
type ParamEditInfo struct {
	Name        string `json:"name" form:"name" binding:"required" example:"shiki"` //名称
	Icon        string `json:"icon" form:"icon" binding:"required" example:"1.png"` //头像url
	Description string `json:"description" form:"description" example:"这是一个介绍"`     //简介
	Gender      string `json:"gender" form:"gender" example:"男"`                    //性别
	Birth       string `json:"birth" form:"birth" example:"20000101"`               //生日
}

// ParamChangePwd 获取编辑用户信息请求application/json参数
type ParamChangePwd struct {
	Password    string `json:"password" form:"password" binding:"required"`         //原密码
	NewPassword string `json:"new_password" form:"new_password" binding:"required"` //新密码
}

// ParamCreateTag 创建版块请求application/json参数
type ParamCreateTag struct {
	Name        string `json:"name" form:"name" binding:"required" example:"百年校区祝福"` //名称
	Description string `json:"description" form:"description" example:"介绍"`          //简介
}

// ParamCreatePost 创建版块请求application/json参数
type ParamCreatePost struct {
	TagId   int    `json:"tag_id" form:"tag_id" binding:"required" example:"2"`       //版块id
	Title   string `json:"title" form:"title" binding:"required" example:"这是个标题"`     //标题
	Content string `json:"content" form:"content" binding:"required" example:"杰哥不要啊"` //内容
}

type ParamAddChatRecord struct {
	FromId  int    `json:"from_id" form:"from_id" binding:"required"`   //聊天对象id
	LastMsg string `json:"last_msg" form:"last_msg" binging:"required"` //最后一条聊天记录
}
