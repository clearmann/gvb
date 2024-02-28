package model

import (
	"gorm.io/gorm"
	"time"
)

type UserAuth struct {
	Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"`
	Password      string     `gorm:"type:varchar(100)" json:"password"`
	Email         string     `json:"email" gorm:"type:varchar(30)"`
	LoginType     int        `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`
	IpAddress     string     `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"`
	IpSource      string     `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`
	LastLoginTime *time.Time `json:"last_login_time"`
	IsDisable     bool       `json:"is_disable"`
	IsSuper       bool       `json:"is_super"` // 超级管理员只能后台设置

	UserInfoId int       `json:"user_info_id"`
	UserInfo   *UserInfo `json:"info"`
	Roles      []Role    `json:"roles" gorm:"many2many:user_auth_roles"`
}
type Role struct {
	gorm.Model
	Name      string `gorm:"unique" json:"name"`
	Label     string `gorm:"unique" json:"label"`
	IsDisable bool   `json:"is_disable"`

	Resources []Resource `json:"resources" gorm:"many2many:role_resources"`
	Menus     []Menu     `json:"menus" gorm:"many2many:role_menus"`
	Users     []UserAuth `json:"users" gorm:"many2many:user_auth_roles"`
}
type Resource struct {
	gorm.Model
	Name      string `gorm:"unique;type:varchar(50)" json:"name"`
	ParentId  int    `json:"parent_id"`
	Url       string `gorm:"type:varchar(255)" json:"url"`
	Method    string `gorm:"type:varchar(10)" json:"request_method"`
	Anonymous bool   `json:"is_anonymous"`

	Roles []Role `json:"roles" gorm:"many2many:role_resources"`
}
type Menu struct {
	gorm.Model
	ParentId     int    `json:"parent_id"`
	Name         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(20)" json:"name"` // 菜单名称
	Path         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(50)" json:"path"` // 路由地址
	Component    string `gorm:"type:varchar(50)" json:"component"`                          // 组件路径
	Icon         string `gorm:"type:varchar(50)" json:"icon"`                               // 图标
	OrderNum     int8   `json:"order_num"`                                                  // 排序
	Redirect     string `gorm:"type:varchar(50)" json:"redirect"`                           // 重定向地址
	Catalogue    bool   `json:"is_catalogue"`                                               // 是否为目录
	Hidden       bool   `json:"is_hidden"`                                                  // 是否隐藏
	KeepAlive    bool   `json:"keep_alive"`                                                 // 是否缓存
	External     bool   `json:"is_external"`                                                // 是否外链
	ExternalLink string `gorm:"type:varchar(255)" json:"external_link"`                     // 外链地址

	Roles []Role `json:"roles" gorm:"many2many:role_menus"`
}
type RoleMenu struct {
	RoleId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
	UserId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
}
type RoleResource struct {
	RoleId     int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
	ResourceId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
}
type UserAuthRole struct {
	UserAuthId int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
	RoleId     int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
}

func GetUserAuthInfoByName(db *gorm.DB, username string) (*UserAuth, error) {
	var userAuth UserAuth
	result := db.Where(&UserAuth{Username: username}).First(&userAuth)
	return &userAuth, result.Error
}
func GetUserAuthById(db *gorm.DB, id int) (*UserAuth, error) {
	var userAuth UserAuth
	result := db.Where(&UserAuth{Model: Model{ID: id}}).First(&userAuth)
	return &userAuth, result.Error
}
func CreateUser(db *gorm.DB, email string, username string, passwordHash string) error {
	result := db.Create(&UserAuth{
		Username:      username,
		Password:      passwordHash,
		LoginType:     0,
		IpAddress:     "",
		IpSource:      "",
		Email:         email,
		LastLoginTime: nil,
		IsDisable:     false,
		IsSuper:       false,
		UserInfoId:    1,
		UserInfo: &UserInfo{
			Nickname: "",
			Avatar:   "",
			Intro:    "",
			Website:  "",
		},
		Roles: []Role{
			{
				Name:      "root",
				Label:     "",
				IsDisable: true,
			},
		},
	})
	return result.Error
}
func GetRoleIdsByUserId(db *gorm.DB, userAuthId int) (ids []int, err error) {
	result := db.Model(&UserAuthRole{UserAuthId: userAuthId}).Pluck("role_id", &ids)
	return ids, result.Error
}
