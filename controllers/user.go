package controllers

import (

	// "encoding/json"
	// "fmt"
	// "math/rand"
	// "net/url"

	"log"
	"time"

	"net/http"
	// "github.com/astaxie/beego"

	"github.com/astaxie/beego"
	"github.com/chuanshuo843/12306_server/utils"
	"github.com/chuanshuo843/12306_server/utils/kyfw"
)

var (
	kyfwUser kyfw.User
)

// var request utils.Request

// //{"result_message":"验证码校验成功","result_code":"4"}
// type _VerifyRes struct {
// 	ResultMessage string `json:"result_message"`
// 	ResultCode    string `json:"result_code"`
// }

// //{"result_message":"登录成功","result_code":0,"uamtk":"tnRPMlCjrDGm3k5IbzlRKQrbmnKToZC_8WN4ePn32Mkhuc1c0"}
// type _LoginRes struct {
// 	ResultMessage string `json:"result_message"`
// 	ResultCode    int    `json:"result_code"`
// 	Uamtk         string `json:"uamtk"`
// }

// // {"result_message":"验证通过","result_code":0,"apptk":null,"newapptk":"P5e8H_FPPq-Q6kfa9uUsKC0PUdOyqGtE6OSTPKvol9Qhuc1c0"}
// type _UaMtkRes struct {
// 	ResultCode    int    `json:"result_code"`
// 	ResultMessage string `json:"result_message"`
// 	AppTk         string `json:"apptk"`
// 	NewAppTK      string `json:"newapptk"`
// }

// //{"apptk":"6fgxwb7avXwqubqIZr5kHbmHZY2wxV2RqUjDkX0xs8Etyc2c0","result_code":0,"result_message":"验证通过","username":"YouName"}
// type _AuthOk struct {
// 	ResultCode    int    `json:"result_code"`
// 	ResultMessage string `json:"result_message"`
// 	AppTk         string `json:"apptk"`
// 	UserName      string `json:"username"`
// }

// UserController Operations about Users
type UserController struct {
	BaseController
}

// Prepare .
func (u *UserController) Prepare() {
	req := u.req()
	if req == nil {
		sid := u.Ctx.Input.CruSession.SessionID()
		log.Println("sid = ", sid)
		if sid != "" {
			req := utils.NewRequest()
			kyfw.Store(sid, req)
		}
	}
}

// Login 登录12306
func (u *UserController) Login() {
	verify := u.GetString("verify")
	username := u.GetString("username")
	password := u.GetString("password")

	req := u.req()
	// key := u.GetString("key")
	errLogin := kyfwUser.Login(req, username, password, verify)
	if errLogin != nil {
		u.Fail().SetMsg(errLogin.Error()).Send()
	}
	//生成JWT
	jwt := &utils.Jwt{}
	jwt.InitJwt()
	jwt.Payload.Jti = time.Now().Unix()
	jwt.Payload.Iat = time.Now().Unix()
	jwt.Payload.Nbf = time.Now().Unix()
	jwt.Payload.Exp = time.Now().Unix() + 70000
	jwt.Payload.Data = `{"username":"` + kyfwUser.UserName + `"}`
	token := jwt.Encode()
	reJSON := map[string]string{"access_token": token}
	kyfw.Store(token, req)
	kyfw.Delete(u.Ctx.Input.Cookie(beego.BConfig.WebConfig.Session.SessionName))

	u.Success().SetMsg("登录成功").SetData(reJSON).Send()
}

// VerifyCode 获取12306登录验证码
func (u *UserController) VerifyCode() {
	req := u.req()
	//初始化登录页面
	_, errInit := kyfwUser.InitLogin(req)
	if errInit != nil {
		http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
		return
	}
	//获取验证码
	data, errVer := kyfwUser.GetVerifyImages(req)
	if errVer != nil {
		http.Error(u.Ctx.ResponseWriter, "Not Found", 404)
		return
	}
	u.Ctx.Output.ContentType("png")
	u.Ctx.Output.Body(data)
}
