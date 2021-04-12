package service

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"golang.org/x/crypto/bcrypt"
	"hackathon/config"
	"hackathon/dao/mysql"
	"hackathon/models"
	"hackathon/response"
	"hackathon/util"
	"math/rand"
	"time"
)

func Register(p *models.ParamRegister) int {
	//数据验证
	if len(p.Password) < 6 {
		return response.CodePwdLength
	}

	//如果名称没有传，就给名称一个随机的十位字符串
	if len(p.Name) == 0 {
		p.Name = util.RandomString(10)
	}
	//判断手机号是否存在
	DB := mysql.GetDB()
	if mysql.IsTelephoneExist(DB, p.Telephone) {
		return response.CodePhoneExist
	}
	var smsCode models.SmsCode
	mysql.RetrieveByStruct(&smsCode, p.Telephone)
	if p.Code != smsCode.Code {
		return 1
	}
	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		//response.Response(ctx,http.StatusUnprocessableEntity,422,nil,"加密错误")
		return response.CodeEncryptError
	}
	newUser := &models.User{
		Name:      p.Name,
		Telephone: p.Telephone,
		Password:  string(hasedPassword),
	}
	err = mysql.Create(newUser)
	if err != nil {
		return response.Error
	}
	//返回结果
	return response.OK
}

func Login(p *models.ParamLogin) (string, int) {
	//手机号是否存在
	DB := mysql.GetDB()
	var user models.User
	DB.Where("telephone = ?", p.Telephone).First(&user)
	if user.ID == 0 {

		return "", response.CodePhoneExist
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
		return "", response.CodePwdWrong
	}
	//发放token
	token, err := util.ReleaseToken(user)
	if err != nil {
		return "", response.Error
	}
	return token, response.OK
}

func SendSmsCode(phone string) bool {
	//产生验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	//调用sdk，完成发送
	client, err := dysmsapi.NewClientWithAccessKey(config.Conf.RegionArea, config.Conf.AccessKey, config.Conf.AccessSecret)
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = config.Conf.SignName
	request.TemplateCode = config.Conf.TemplateCode
	par, err := json.Marshal(map[string]string{"code": code})
	request.TemplateParam = string(par)
	response, err := client.SendSms(request)
	if err != nil {
		fmt.Print(err.Error())
		return false
	}
	fmt.Printf("response is %#v\n", response)
	if response.Code == "OK" {
		//将验证码保存到数据库中
		smsCode := &models.SmsCode{
			Phone: phone,
			BizId: response.BizId,
			Code:  code,
		}
		err := mysql.Create(smsCode)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
