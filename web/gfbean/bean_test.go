package gfbean

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/frame/g"
	"go.uber.org/dig"
)

type Version struct {
	Id   string
	Name string
}

func (p *Version) getId() string {
	return p.Id
}

func newPerson() Version {
	return Version{
		Id:   "1.0.6",
		Name: "gf-extend",
	}
}

// 测试注入依赖
func TestDig(t *testing.T) {

	// 创建容器
	container := dig.New()
	// 注入依赖
	err := container.Provide(newPerson)

	if err != nil {
		fmt.Println(err)
	}

	// 使用依赖
	err = container.Invoke(func(p Version) {
		// 业务逻辑
		fmt.Println(p.getId())
	})

	if err != nil {
		fmt.Println(err)
	}
}

// 参数对象
type needPerson1 struct {
	dig.In         // 打包依赖
	P      Version `name:"v"` // 命名依赖
}

// 注入命名依赖
func TestDigName(t *testing.T) {

	// 创建容器
	container := dig.New()
	// 注入命名依赖
	err := container.Provide(newPerson, dig.Name("v"))

	if err != nil {
		fmt.Println(err)
	}

	// 使用依赖
	err = container.Invoke(func(n needPerson1) {
		// 业务逻辑
		g.Dump(n)
	})

	if err != nil {
		fmt.Println(err)
	}

}

// 参数对象
type needPerson2 struct {
	dig.In           // 打包依赖
	P      []Version `group:"v"` // 组依赖，必须是个切片
}

// 注入组依赖
func TestDigGroup(t *testing.T) {

	// 创建容器
	container := dig.New()
	// 注入命名依赖
	err := container.Provide(newPerson, dig.Group("v"))

	if err != nil {
		fmt.Println(err)
	}

	// 使用依赖
	err = container.Invoke(func(n needPerson2) {
		// 业务逻辑
		g.Dump(n)
	})

	if err != nil {
		fmt.Println(err)
	}

}
