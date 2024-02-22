package model

import "gorm.io/gorm"

// 迁移数据表，在没有数据表结构结构变更时候，建议注释不执行
// 只支持创建表，增加表中没有的字段和索引
// 为了保护数据，并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate(db *gorm.DB) error {
	// 设置表关联
	//db.SetupJoinTable(&UserAuth{}, "Roles", &UserAuthRole{})
	//db.SetupJoinTable(&Role{}, "Users", &UserAuthRole{})
	//db.SetupJoinTable(&Role{}, "Menus", &RoleMenu{})
	//db.SetupJoinTable(&Role{}, "Resources", &RoleResource{})
	err := db.AutoMigrate(
		&UserAuth{}, // 用户验证
		&UserInfo{}, // 用户信息
		&Role{},     // 角色
		&Menu{},     // 菜单
		//&RoleMenu{},     // 角色-菜单 关联
		//&RoleResource{}, // 角色资源关联
		//&UserAuthRole{}, // 用户-角色 关联
		&Article{}, //文章
		&Tag{},     //标签
	)
	return err
}
