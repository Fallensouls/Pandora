package user

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id           int64
	Username     string          `valid:"Required"`
	password     string          `valid:"Required"`
	Profile      *UserProfile      // 用户的个人信息
	auth         *[]Authority      `orm:"rel(fk)"`     // 用户的角色权限
	status       int               // 用户的状态
	create       time.Time         // 创建用户的时间
	lastLogin    time.Time         // 记录最后登录的时间
	lastModify   time.Time         // 记录上次修改密码的时间
}

type UserProfile struct {
	Id           int64
	Avatar       string
	Address      string       `valid:"MaxSize(100)"`
	Description  string       `valid:"MaxSize(200)"`
	email        string       `valid:"Email; MaxSize(30)"`
	telphone     string       `valid:"Mobile"`
}

type Authority struct {
	Id			 int
	Role         string            // 用户的角色
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
	Normal       =   0      // 正常状态
	Restricted   =   1      // 受限状态,只能访问部分数据或进行有限的操作
	Banned       =   2      // 封禁状态,不允许登录
)

func init()  {
	orm.RegisterModel(new(User), new(UserProfile), new(Authority))
}

func GetUser(id int64) (User, error)  {
	user := User{Id: id}
	o := orm.NewOrm()
	err := o.Read(&user)
	if err == orm.ErrNoRows{
		return User{}, ErrUserNotFound
	}
	if user.status != Normal{  // 不处于正常状态的用户不可查询
		return User{},ErrUserRestricted
	}
	return user, err
}

func AddUser(user User, by string) error {  // by表示注册用户的方式,可以通过邮箱或者手机号注册
	o := orm.NewOrm()
	err := SetPassword(&user.password)
	if err != nil{
		return ErrEncodingPassword
	}
	user.status = Normal
	user.create = time.Now()
	user.lastLogin = time.Now()
	user.lastModify = time.Now()
	created, _, err := o.ReadOrCreate(&user, by)
	if err == nil{
		if created {   // 成功创建新用户
			return nil
		} else {
			switch by {
			case "email":
				return ErrEmailUsed
			case "tel":
				return ErrTelUsed
			default:
				return err
			}
		}
	}
	return err
}

func ChangeStatus (id int64, status int) error {
	o := orm.NewOrm()
	user := User{Id: id}
	if err := o.Read(&user); err != nil{
		if err == orm.ErrNoRows{
			return ErrUserNotFound
		}
		return err
	}
	user.status = status
	if _, err := o.Update(&user, "status"); err != nil{
		return ErrChangeStatus
	}
	return nil
}

func SetPassword(password *string) error {  // 利用加密方法生成加密的密码
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	*password = string(hash)
	return nil
}

func ChangePassword (id int64, old string, new string) error {  // 利用旧密码来修改密码
	o := orm.NewOrm()
	user := User{Id: id}
	if err := o.Read(&user); err != nil{
		if err == orm.ErrNoRows{
			return ErrUserNotFound
		}
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(old), []byte(user.password)); err != nil{
		return ErrWrongPassword
	}
	err := SetPassword(&user.password)
	if err != nil{
		return ErrEncodingPassword
	}
	user.lastModify = time.Now()
	if _, err := o.Update(&user, "password", "lastModify"); err != nil{
		return err
	}
	return nil
}

func GetProfile (id int64) (UserProfile, error) {
	o := orm.NewOrm()
	profile := UserProfile{Id: id}
	err := o.Read(&profile)
	if err == orm.ErrNoRows{
		return profile, ErrUserNotFound
	}
	return profile, err
}

func UpdateProfile (profile UserProfile) error {
	o := orm.NewOrm()
	if _, err := o.Update(&profile); err != nil{
		return err
	}
	return nil
}
