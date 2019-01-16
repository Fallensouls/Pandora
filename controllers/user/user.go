package user

import (
	"Pandora/controllers"
	. "Pandora/models/user"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	"regexp"
	"strconv"
)

// Operations about Users
type UserController struct {
	controllers.BaseController
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (this *UserController) addUser() {
	username := this.GetString("username")
	password := this.GetString("password")
	email := this.GetString("email")
	tel := this.GetString("telphone")

	user := User{Username:username, Password:password, Email:email, Telphone:tel}
	by := ""

	valid := validation.Validation{}
	reg1 := regexp.MustCompile("^[a-zA-Z0-9_]{4,16}$")
	reg2 := regexp.MustCompile("^[a-zA-Z0-9]{8,16}$")
	valid.Match(user.Username, reg1, "username").Message("用户名只能包括数字、字母以及下划线，且在4-16位之间")
	valid.Match(user.Password, reg2, "password").Message("密码只能包括数字、字母，且在8-16位之间")
	if user.Email != ""{
		valid.Email(user.Email,"email")
		by = "email"
	}
	if user.Telphone != ""{
		valid.Mobile(user.Telphone,"telphone")
		by = "tel"
	}
	b, err := valid.Valid(&user)
	if err != nil {
		logs.Error(err)
		this.Abort("500")
	}
	if !b && user.Telphone == "" && user.Email == ""{
		// validation does not pass
		this.Data["json"] = "invalid data"
		this.Abort("400")  // bad request
	}
	if err := AddUser(user, by); err != nil{
		logs.Error(err)
		this.Abort("500")
	}
	this.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
//func (u *UserController) GetAll() {
//	users := models.GetAllUsers()
//	u.Data["json"] = users
//	u.ServeJSON()
//}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:id [get]
func (this *UserController) getUser() {
	uid := this.GetString(":id")
	if uid != "" {
		id, _ := strconv.ParseInt(uid, 10, 64) // string转int64
		user, err := GetUser(id)
		if err != nil {
			this.Data["json"] = err.Error()
		} else {
			this.Data["json"] = user
		}
	}
	this.ServeJSON()
}

// @Title Update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:id [put]
func (this *UserController) updateUser() {
	uid := this.GetString(":id")
	if uid != "" {
		var profile User
		if err := json.Unmarshal(this.Ctx.Input.RequestBody, &profile); err != nil{
			logs.Error(err)
			this.Abort("500")
		}

		valid := validation.Validation{}
		valid.Range(profile.Age,0,140,"age").Message("年龄应在0-140之间")
		valid.MaxSize(profile.Address,100,"address").Message("地址长度最大为100个字符")
		valid.MaxSize(profile.Description,200,"description").Message("个人说明最大为200个字符")
		reg := regexp.MustCompile("^[0-2]*$")
		valid.Match(profile.Gender, reg,"gender")
		b, err := valid.Valid(&profile)
		if err != nil {
			logs.Error(err)
			this.Abort("500")
		}
		if !b {
			// validation does not pass
			this.Data["json"] = "invalid data"
			this.Abort("400")  // bad request
		}

		if err := UpdateUserProfile(profile); err != nil{
			logs.Error(err)
			this.Abort("500")
		}
	}
	this.ServeJSON()
}


// @Title Login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UserController) Login() {
	username := u.GetString("username")
	password := u.GetString("password")
	if models.Login(username, password) {
		u.Data["json"] = "login success"
	} else {
		u.Data["json"] = "user not exist"
	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}

