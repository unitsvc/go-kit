package hjsoup

import (
	"github.com/gogf/gf/frame/g"

	"fmt"
	"testing"
)

func TestHttpClient(t *testing.T) {

	productCode, err := SearchByProductCode("6928804011111", false)
	if err != nil {
		fmt.Println(err)
	}
	g.Dump(productCode)
}

func TestSearchByProductCodeCache(t *testing.T) {

	// 商品条码
	productCode := "6928804011111"

	productCodeInfo, err := SearchByProductCodeCache(productCode, true)
	if err != nil {
		g.Log().Line(false).Error(err)
	}

	g.Log().Line(false).Info(productCodeInfo)

}
