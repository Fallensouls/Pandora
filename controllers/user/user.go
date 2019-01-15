package user

import (
	"Pandora/controllers"
	. "Pandora/models/user"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
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
	var user User
	valid := validation.Validation{}
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &user); err != nil{
		logs.Error(err)
		this.Abort("500")
	}
	b, err := valid.Valid(&user)
	if err != nil {
		logs.Error(err)
		this.Abort("500")
	}
	if !b {
		// validation does not pass
		this.Data["json"] = "invalid data"
		this.Abort("400")  // bad request
	}
	if err := AddUser(user, "email"); err != nil{
		logs.Error(err)
	}
	//this.Data["json"] = map[string]string{"uid": uid}
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
		var profile UserProfile
		if err := json.Unmarshal(this.Ctx.Input.RequestBody, &profile); err != nil{
			logs.Error(err)
			this.Abort("500")
		}
		valid := validation.Validation{}
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
		if err := UpdateProfile(profile); err != nil{
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

