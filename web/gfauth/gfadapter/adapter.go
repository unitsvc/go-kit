package gfadapter

import (
	"fmt"
	"runtime"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

type CasbinRule struct {
	PType string `json:"pType"`
	V0    string `json:"v0"`
	V1    string `json:"v1"`
	V2    string `json:"v2"`
	V3    string `json:"v3"`
	V4    string `json:"v4"`
	V5    string `json:"v5"`
}

// Adapter 表示策略存储的gdb适配器。
type Adapter struct {
	driverName string
	dbLink     string
	tableName  string
	db         gdb.DB
}

// finalizer 是适配器的析构函数。
func finalizer(a *Adapter) {
	// 注意不用的时候不需要使用Close方法关闭数据库连接(并且gdb也没有提供Close方法)，
	// 数据库引擎底层采用了链接池设计，当链接不再使用时会自动关闭
	a.db = nil
}

// NewAdapter 是适配器的构造函数。
func NewAdapter(driverName string, dbLink string) (*Adapter, error) {
	a := &Adapter{}
	a.driverName = driverName
	a.dbLink = dbLink
	a.tableName = "casbin_rule"

	// 打开数据库，如果不存在就创建它。
	err := a.open()
	if err != nil {
		return nil, err
	}

	// 释放对象时调用析构函数。
	runtime.SetFinalizer(a, finalizer)

	return a, nil
}

// NewAdapterFromOptions 是已存在连接的适配器的构造函数。
func NewAdapterFromOptions(adapter *Adapter) (*Adapter, error) {

	if adapter.tableName == "" {
		adapter.tableName = "casbin_rule"
	}
	if adapter.db == nil {
		err := adapter.open()
		if err != nil {
			return nil, err
		}

		runtime.SetFinalizer(adapter, finalizer)
	}

	return adapter, nil
}

func (a *Adapter) open() error {
	var err error
	var db gdb.DB

	gdb.SetConfig(gdb.Config{
		"casbin": gdb.ConfigGroup{
			gdb.ConfigNode{
				Type:   a.driverName,
				Link:   a.dbLink,
				Role:   "master",
				Weight: 100,
			},
		},
	})
	db, err = gdb.New("casbin")

	if err != nil {
		return err
	}

	a.db = db

	return a.createTable()
}

func (a *Adapter) close() error {
	// 注意不用的时候不需要使用Close方法关闭数据库连接(并且gdb也没有提供Close方法)，
	// 数据库引擎底层采用了链接池设计，当链接不再使用时会自动关闭
	a.db = nil
	return nil
}

func (a *Adapter) createTable() error {
	_, err := a.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (p_type VARCHAR(10), v0 VARCHAR(256), v1 VARCHAR(256), v2 VARCHAR(256), v3 VARCHAR(256), v4 VARCHAR(256), v5 VARCHAR(256))", a.tableName))
	return err
}

func (a *Adapter) dropTable() error {
	_, err := a.db.Exec(fmt.Sprintf("DROP TABLE %s", a.tableName))
	return err
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	lineText := line.PType
	if line.V0 != "" {
		lineText += ", " + line.V0
	}
	if line.V1 != "" {
		lineText += ", " + line.V1
	}
	if line.V2 != "" {
		lineText += ", " + line.V2
	}
	if line.V3 != "" {
		lineText += ", " + line.V3
	}
	if line.V4 != "" {
		lineText += ", " + line.V4
	}
	if line.V5 != "" {
		lineText += ", " + line.V5
	}

	persist.LoadPolicyLine(lineText, model)
}

// LoadPolicy 从数据库加载所有策略规则。（必须实现此方法）
func (a *Adapter) LoadPolicy(model model.Model) error {
	// 打印日志
	g.Log().Line(false).Debug("从数据库加载所有策略规则")

	var lines []CasbinRule

	if err := a.db.Table(a.tableName).
		Fields("p_type", "v0", "v1", "v2", "v3", "v4", "v5").
		Scan(&lines); err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

func savePolicyLine(pType string, rule []string) CasbinRule {
	line := CasbinRule{}

	line.PType = pType
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

// SavePolicy 将所有策略保存到数据库。（必须实现此方法）（慎用）
//  备注：只有适配器才能调用此方法，执行器无法调用。
func (a *Adapter) SavePolicy(model model.Model) error {
	// 打印日志
	g.Log().Line(false).Warning("保存所有策略到数据库，执行流程：删除表，新建表，将所有策略规则保存到存储中，不建议使用此方法a.SavePolicy()。")

	// TODO 需要修改，不能直接删表
	err := a.dropTable()
	if err != nil {
		return err
	}
	err = a.createTable()
	if err != nil {
		return err
	}

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			_, err := a.db.Table(a.tableName).Data(&line).Insert()
			if err != nil {
				return err
			}
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			_, err := a.db.Table(a.tableName).Data(&line).Insert()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// AddPolicy 将策略规则添加到存储中。
func (a *Adapter) AddPolicy(sec string, pType string, rule []string) error {
	// 打印日志
	g.Log().Line(false).Debug("持久化策略：", sec, pType, rule)
	line := savePolicyLine(pType, rule)
	_, err := a.db.Table(a.tableName).Data(&line).Insert()
	return err
}

// RemovePolicy 从存储中删除策略规则。
func (a *Adapter) RemovePolicy(sec string, pType string, rule []string) error {
	g.Log().Line(false).Debug("删除持久化策略：", sec, pType, rule)
	line := savePolicyLine(pType, rule)
	err := rawDelete(a, line)
	return err
}

// RemoveFilteredPolicy 从存储中删除匹配过滤器的策略规则。
func (a *Adapter) RemoveFilteredPolicy(sec string, pType string, fieldIndex int, fieldValues ...string) error {
	// 打印日志
	g.Log().Line(false).Debug("过滤指定字段并删除相应策略：", sec, pType, fieldIndex, fieldValues)

	line := CasbinRule{}

	line.PType = pType
	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		line.V5 = fieldValues[5-fieldIndex]
	}
	err := rawDelete(a, line)
	return err
}

func rawDelete(a *Adapter, line CasbinRule) error {
	db := a.db.Table(a.tableName)

	db.Where("p_type = ?", line.PType)
	if line.V0 != "" {
		db.Where("v0 = ?", line.V0)
	}
	if line.V1 != "" {
		db.Where("v1 = ?", line.V1)
	}
	if line.V2 != "" {
		db.Where("v2 = ?", line.V2)
	}
	if line.V3 != "" {
		db.Where("v3 = ?", line.V3)
	}
	if line.V4 != "" {
		db.Where("v4 = ?", line.V4)
	}
	if line.V5 != "" {
		db.Where("v5 = ?", line.V5)
	}

	_, err := db.Delete()
	return err
}
