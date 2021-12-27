package apiparser

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type Parser struct {
	Types map[string]string
	Service map[string][]Router
}

type Router struct {
	Handler string `json:"handler"`
	Middle string `json:"middle"`
	Method string `json:"method"`
	Action string `json:"action"`
	Request string `json:"request"`
	Response string `json:"response"`
}

func NewParser() *Parser {
	return &Parser{}
}

func(m *Parser) Parse(filename string) error {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	ctx, err := ioutil.ReadFile(abs)
	if err != nil {
		return err
	}
	m.fillType(string(ctx))
	err = m.fillService(string(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (m *Parser) fillType(ctx string)  {
	typeList := map[string]string{}
	reg2 := regexp.MustCompile("type (\\w*) .* {\n(.*?\"`\n)*}")
	result2 := reg2.FindAllStringSubmatch(ctx, -1)
	for _, v := range result2 {
		typeList[v[1]] = v[0]
	}
	m.Types = typeList
}

func (m *Parser)  fillService(ctx string) (err error) {
	result2 := getRegexp("service (.*) {\\n([\\w@\\s/()]*\\n)*}", ctx)
	routerMap := map[string][]Router{}
	for _, v := range result2 {
		var arr []string
		arr = strings.Split(v[2], "\n\n")
		for _, v1 := range arr {
			lineArr := strings.Split(v1, "\n")
			for k, line := range lineArr {
				lineArr[k] = strings.Trim(line, " ")
				lineArr[k] = strings.Trim(lineArr[k], "\t")
			}
			str := strings.Join(lineArr, "\n")
			var router Router
			//匹配handler
			result := getRegexp("@handler (\\w*)", str)
			if len(result) == 0 {
				err = errors.New("缺失handler")
				return
			}
			router.Handler = result[0][1]
			//匹配middle中间件
			result = getRegexp("@middle (\\w*)", str)
			if len(result) > 0 {
				router.Middle = result[0][1]
			}
			//匹配方法
			result = getRegexp("(\\w*) ([\\w/]*) (.*) returns (.*)", str)

			if len(result) == 0 {
				err = errors.New("匹配方式缺失")
				return
			}
			router.Method = result[0][1]
			router.Action = result[0][2]
			router.Request = strings.TrimRight(strings.TrimLeft(result[0][3], "("), ")")
			router.Response =strings.TrimRight(strings.TrimLeft(result[0][4], "("), ")")
			//检测request以及response是否已定义
			if _, ok := m.Types[router.Request]; !ok {
				err = errors.New(router.Request+ "未定义")
				break
			}
			if _, ok := m.Types[router.Response]; !ok {
				err = errors.New(router.Response+ "未定义")
				break
			}
			routerMap[v[1]] = append(routerMap[v[1]], router)
		}
	}
	m.Service = routerMap
	return
}

func getRegexp(pattern, ctx string) [][]string {
	reg2 := regexp.MustCompile(pattern)
	result2 := reg2.FindAllStringSubmatch(ctx, -1)
	return result2
}

