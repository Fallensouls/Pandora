package controllers

import "github.com/astaxie/beego"


// 建立base controller，可以用于检验用户是否登录或检查用户的权限，相当于request拦截器
type BaseController struct {
	beego.Controller
}

// Prepare会在controller的其他方法之前执行
func (this *BaseController) Prepare()  {

}