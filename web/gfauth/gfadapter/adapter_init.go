package gfadapter

import (
	"runtime"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
)

const (
	sqlite = "sqlite"
	mysql  = "mysql"
	pgsql  = "pgsql"
)

// NewAdapterByGdb 从自定义数据库连接中创建适配器
func NewAdapterByGdb(customDb gdb.DB) (*Adapter, error) {

	// 获取当前数据库类型
	dbType := customDb.GetConfig().Type

	//  TODO 需要对不同类型数据库进行处理

	// 判断当前数据库类型-Sqlite3
	if gstr.Equal(sqlite, dbType) { // 自动创建sqlite3数据库casbin_rule表
		sql := CreateSqlite3Table("casbin_rule")
		if _, err := customDb.Exec(sql); err != nil {
			return nil, err
		}
	}

	// 判断当前数据库类型-Mysql
	if gstr.Equal(mysql, dbType) { // 自动创建mysql数据库casbin_rule表
		sql := CreateMysqlTable("casbin_rule")
		if _, err := customDb.Exec(sql); err != nil {
			return nil, err
		}
	}

	// 判断当前数据库类型-Pgsql
	if gstr.Equal(pgsql, dbType) { // 自动创建pgsql数据库casbin_rule表
		sql := CreatePgsqlTable("casbin_rule")
		if _, err := customDb.Exec(sql); err != nil {
			return nil, err
		}
	}

	// TODO ------------------------------------------------------------
	// 1.添加事务操作
	// 2.修改自动保存逻辑
	// 3.完善不同数据库操作
	// 4.单例模式实例化对象
	// -----------------------------------------------------------------

	// 构造适配器对象
	a := &Adapter{
		db: customDb,
	}

	if customDb == nil {
		return nil, gerror.New("数据库默认连接不存在，无法实例化casbin执行器")
	}

	// 释放对象时调用析构函数。
	runtime.SetFinalizer(a, finalizer)

	return a, nil
}

// NewEnforcer 实例化gf-casbin执行器对象
//  1、支持自动注册，自动寻找gf框架default分组数据源（gdb.DB），无需关心数据源种类。
//  2、支持自定义分组数据源注册。
//
//  目前支持sqlite3、mysql5.7、postgresql数据库。
//
//  备注：1.sqlite3、mysql5.7数据库表新增主键自增，postgresql数据库无主键。
//       2.sqlite3、pgsql需要添加额外驱动
//
//	sqlite3驱动：	_ "github.com/lib/pq"
//
//	pgsql驱动：		_ "github.com/mattn/go-sqlite3"
//
//  示例：
//  e, err := gfadapter.NewEnforcer()
//  e, err := gfadapter.NewEnforcer(g.DB())
//  e, err := gfadapter.NewEnforcer(g.DB("sqlite"))
//  e, err := gfadapter.NewEnforcer(g.DB("mysql"))
//  e, err := gfadapter.NewEnforcer(g.DB("pgsql"))
func NewEnforcer(customDb ...gdb.DB) (*casbin.Enforcer, error) {

	// TODO 需要添加动态模型配置，从配置文件中读取，获取不到则使用默认配置
	// 加载casbin默认模型
	modelFromString, _ := getDefaultNewModel()

	// 定义数据源
	var db gdb.DB

	// 判断是否传参数
	num := len(customDb)

	if num == 0 {
		db = g.DB()
	} else {
		db = customDb[0]
	}

	// 创建gf默认数据源适配器
	if adapter, err := NewAdapterByGdb(db); err == nil {

		// 调用已有连接的适配器中的构造器
		if options, err := NewAdapterFromOptions(adapter); err == nil {
			// 返回casbin执行器
			return casbin.NewEnforcer(modelFromString, options)
		} else {
			return nil, err
		}

	} else {
		return nil, err
	}

}

// 加载casbin默认rbac模型
// rbac_model.conf
func getDefaultNewModel() (model.Model, error) {
	// 打印日志
	g.Log().Line(false).Debug("加载casbin模型")
	// 配置rbac_model.conf字符串
	rbacModelText := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	// 从字符串中加载模型
	return model.NewModelFromString(rbacModelText)
}
