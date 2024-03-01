package model

import (
	"gorm.io/gorm"
	"time"
)

/*
不使用 uint
因为 uint 在类型转化中存在诸多不便 没必要用 uint
直接使用 int 即可
*/
type Model struct {
	ID        int       `gorm:"primary_key;auto_increment" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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
		&Article{},      // 文章
		&Tag{},          // 标签
		&Resource{},     // 资源
		&Category{},     // 分类
		&Comment{},      // 评论
		&Config{},       // 配置
		&FriendLink{},   // 友链
		&Message{},      // 消息
		&OperationLog{}, // 操作日志
		&Page{},
		&Tag{},
	)
	return err
}

// 分页
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// 数据列表
func List[T any](db *gorm.DB, data T, slt, order, query string, args ...any) (T, error) {
	db = db.Model(data).Select(slt).Order(order)
	if query != "" {
		db = db.Where(query, args...)
	}
	result := db.Find(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}
