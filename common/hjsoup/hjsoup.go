package hjsoup

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ProductCodeDto 商品条码信息
type ProductCodeDto struct {
	ProductCode      string   // 商品条码
	ProductCodeImage string   // 商品条形码图片
	CompanyName      string   // 企业名称
	CompanyAddress   string   // 企业注册地址
	ProductName      string   // 产品名称
	ProductCategory  string   // 产品分类
	Brand            string   // 品牌
	ProductSpec      string   // 产品规格
	StandardNo       string   // 标准号
	StandardName     string   // 标准名称
	ProductExp       string   // 保质期
	UpMarketTime     string   // 上市日期
	DownMarketTime   string   // 下市日期
	ProductImageList []string // 图片列表
}

// SearchByProductCode 根据商品条码查询商品信息，警告：此方法仅供学习参考，禁止用于商业
//	@productCode	商品条码
//	@debug			是否开启debug
func SearchByProductCode(productCode string, debug bool) (*ProductCodeDto, error) {

	// 商品条码信息
	productCodeDto := ProductCodeDto{}

	// 1.获取验证码信息
	response, err := g.Client().
		Timeout(20 * time.Second).
		Header(map[string]string{
			"Host":       "www.chinatrace.org",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		Get("http://www.chinatrace.org/trace/verification/image?_=" + strconv.Itoa(time.Now().Second()*1000))
	if err != nil {
		return &productCodeDto, err
	}

	cookies := response.Cookies()

	// 会话sessionId
	sessionId := cookies[0].Value

	fmt.Println("当前会话：", sessionId)

	// 2.查询y值
	y := gjson.New(response.ReadAllString()).GetString("y")
	fmt.Println("验证码y值：", y)

	// 定义有效x值
	validX := ""

	// 定义验证码
	doorCode := ""

	// 3.解析验证码
	for x := 0; x < 350; x++ {
		body := g.Client().
			SetCookieMap(map[string]string{
				"JSESSIONID": sessionId,
			}).
			Header(map[string]string{
				"Host":       "www.chinatrace.org",
				"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
			}).
			PostContent("http://www.chinatrace.org/trace/verification/result?x=" + strconv.Itoa(x) + "&y=" + y)

		if gstr.LenRune(body) > 0 {
			doorCode = body
			validX = strconv.Itoa(x)
			fmt.Println("重试次数："+strconv.Itoa(x)+"次，", "成功获取验证码：", doorCode)
			break
		}
	}

	if gstr.LenRune(doorCode) == 0 {
		return &productCodeDto, err
	}

	// 此处传参是正确的
	formData1 := map[string]string{
		"productCode":  productCode,
		"batchNo":      "",
		"productCode1": "",
		"traceCode":    "",
		"doorCode":     doorCode,
		"validX":       validX,
	}

	fmt.Println("请求参数：", formData1)

	// 警告：不能按照这种方式传参->使用FormPost方法，此处cookie会丢失
	formData2 := url.Values{
		"productCode":  {productCode},
		"batchNo":      {""},
		"productCode1": {""},
		"traceCode":    {""},
		"doorCode":     {doorCode},
		"validX":       {validX},
	}

	fmt.Println("格式化参数：", formData2.Encode())

	formResp, err := g.Client().Timeout(20*time.Second).
		SetCookieMap(map[string]string{
			"JSESSIONID": sessionId,
		}).
		Header(map[string]string{
			"Host":       "www.chinatrace.org",
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		}).
		ContentType("application/x-www-form-urlencoded").
		Post("http://www.chinatrace.org/trace/door/controller/SearchController/searchByProductCode.do", formData1)

	if err != nil {
		return &productCodeDto, err
	}

	body, err := ioutil.ReadAll(formResp.Body)
	if err != nil {
		return &productCodeDto, err
	}

	//fmt.Println(string(body))

	//fmt.Println(formResp.Request.Cookies())

	// 设置商品条码
	productCodeDto.ProductCode = productCode

	// 4.解析数据
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return &productCodeDto, err
	}

	// ID选择器，处理企业信息
	dom.Find("#company").Each(func(i int, selection *goquery.Selection) {
		selection.Find("p").Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 0:
				// 设置企业名称
				productCodeDto.CompanyName = selection.Find("span").Text()
			case 1:
				// 设置企业注册地址
				productCodeDto.CompanyAddress = selection.Find("span").Text()
			}
		})
	})

	productDiv := dom.Find("#product")

	// class选择器
	table := productDiv.Find(".table")

	table.Find("tr").Each(func(t1 int, selection1 *goquery.Selection) {

		// 处理tr单元格
		selection1.Find("td").Each(func(t2 int, selection2 *goquery.Selection) {
			if t2 == 1 {
				switch t1 {
				case 1:
					// 设置产品名称
					productCodeDto.ProductName = TrimBlank(selection2.Text())
				case 2:
					// 设置产品分类
					productCodeDto.ProductCategory = TrimBlank(selection2.Text())
				case 3:
					// 设置品牌
					productCodeDto.Brand = TrimBlank(selection2.Text())
				case 4:
					// 设置商品规格
					productCodeDto.ProductSpec = TrimBlank(selection2.Text())
				case 5:
					// 设置标准号
					productCodeDto.StandardNo = TrimBlank(selection2.Text())
				case 6:
					// 设置标准名称
					productCodeDto.StandardName = TrimBlank(selection2.Text())
				case 7:
					// 设置保质期
					productCodeDto.ProductExp = TrimBlank(selection2.Text())
				case 8:
					// 设置上市日期
					productCodeDto.UpMarketTime = TrimBlank(selection2.Text())
				case 9:
					// 设置下市日期
					productCodeDto.DownMarketTime = TrimBlank(selection2.Text())
				}
			}

		})
	})

	// 实例化图片切片
	imgList := make([]string, 0)

	// 元素选择器，处理图片列表
	table.Find("img").Each(func(i int, selection *goquery.Selection) {
		src, _ := selection.Attr("src")
		if i == 0 {
			// 设置条形码图片
			productCodeDto.ProductCodeImage = src
		} else {
			imgList = append(imgList, src)
		}

	})
	// 设置图片列表
	productCodeDto.ProductImageList = imgList

	return &productCodeDto, nil
}

// TrimBlank 删除字符串首尾空白
func TrimBlank(str string) string {
	// 去除制表符
	trimStr := gstr.TrimStr(str, "\t")
	// 去除换行符
	trimStr = gstr.TrimStr(trimStr, "\n")
	// 去除首尾空白字符
	return gstr.Trim(trimStr)
}