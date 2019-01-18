package user

import (
	"Pandora/models"
	"errors"
	_ "golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	models.BaseModel
	Username     string          `json:"username"`
	Password     string          `json:"-"`
	Avatar       string			 `json:"avatar"`
	Age          int             `json:"age"`
	Gender       int			 `json:"gender"`
	Address      string          `json:"address"`
	Description  string          `json:"description"`
	Email        string          `json:"-"`
	Telphone     string          `json:"-"`
	Auth         *[]Authority    `json:"-"`     // 用户的角色权限
	Status       int             `json:"-"`     // 用户的状态
	LastLogin    time.Time       `json:"-"`     // 记录最后登录的时间
	LastModify   time.Time       `json:"-"`     // 记录上次修改密码的时间
}

type Authority struct {
	models.BaseModel
	Role         string            // 用户的角色
	User         *[]User
}

// 定义可能发生的异常
var (
	ErrUserNotFound      =   errors.New("<User Service> the user does not exist")
	ErrUserRestricted    =   errors.New("<User Service> the user has been restricted")
	ErrWrongPassword     =   errors.New("<User Service> the password is wrong")
	ErrEncodingPassword  =   errors.New("<User Service> fail to encode the password")
	ErrEmailUsed         =   errors.New("<User Service> the email address has already been used")
	ErrTelUsed           =   errors.New("<User Service> the telphone number has already been used")
	ErrChangeStatus      =   errors.New("<User Service> fail to change status")
)

// 定义用户账号的状态
const (
	Inactive     =   0      // 未激活状态
	Normal       =   1      // 正常状态
	Restricted   =   2      // 受限状态,只能访问部分数据或进行有限的操作
	Banned       =   3      // 封禁状态,不允许登录
)

// 性别
const(
	Male         =   0
	Female       =   1
	Secret       =   2
)


func GetUser(id int64) (User, error)  {
	//user := User{Id: id}
	//o := orm.NewOrm()
	//err := o.Read(&user)
	//if err == orm.ErrNoRows{
	//	return User{}, ErrUserNotFound
	//}
	//if user.status != Normal{   // 不处于正常状态的用户不可查询
	//	return User{},ErrUserRestricted
	//}
	//return user, err
}

// by表示注册用户的方式,可以通过邮箱或者手机号注册
func AddUser(user User, by string) error {
	//o := orm.NewOrm()
	//err := SetPassword(&user.Password)
	//if err != nil{
	//	return ErrEncodingPassword
	//}
	//user.status = Inactive
	//user.create = time.Now()
	//user.lastLogin = time.Now()
	//user.lastModify = time.Now()
	//created, _, err := o.ReadOrCreate(&user, by)
	//if err == nil{
	//	if created {   // 成功创建新用户
	//		return nil
	//	} else {
	//		switch by {
	//		case "email":
	//			return ErrEmailUsed
	//		case "tel":
	//			return ErrTelUsed
	//		default:
	//			return err
	//		}
	//	}
	//}
	//return err
}

// 修改用户的非隐私信息
func UpdateUserProfile(user User) error {
	//o := orm.NewOrm()
	//if _, err := o.Update(&user,"Avatar","Age","Gender","Address","Description"); err != nil{
	//	return err
	//}
	//return nil
}

// 特殊字段需要另外提供更改的方法
func ChangeStatus(id int64, status int) error {
	//o := orm.NewOrm()
	//user := User{Id: id, status: status}
	//if _, err := o.Update(&user, "status"); err != nil{
	//	return ErrChangeStatus
	//}
	//return nil
}

func SetPassword(password *string) error {  // 利用加密方法生成加密的密码
	//hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	//if err != nil{
	//	return err
	//}
	//*password = string(hash)
	//return nil
}

func ChangePassword(id int64, old string, new string) error {  // 利用旧密码来修改密码
	//o := orm.NewOrm()
	//user := User{Id: id}
	//if err := o.Read(&user); err != nil{
	//	if err == orm.ErrNoRows{
	//		return ErrUserNotFound
	//	}
	//	return err
	//}
	//if err := bcrypt.CompareHashAndPassword([]byte(old), []byte(user.Password)); err != nil{
	//	return ErrWrongPassword
	//}
	//err := SetPassword(&user.Password)
	//if err != nil{
	//	return ErrEncodingPassword
	//}
	//user.lastModify = time.Now()
	//if _, err := o.Update(&user, "password", "lastModify"); err != nil{
	//	return err
	//}
	//return nil
}

func ChangeEmail(id int64, email string) error {
	//o := orm.NewOrm()
	//user := User{Id: id, Email: email}
	//if _, err := o.Update(&user, "Email"); err != nil{
	//	return err
	//}
	//return nil
}

func ChangeTel(id int64, tel string) error {
	//o := orm.NewOrm()
	//user := User{Id: id, Telphone: tel}
	//if _, err := o.Update(&user, "Telphone"); err != nil{
	//	return err
	//}
	//return nil
}


