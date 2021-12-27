package gen

import (
	"fmt"
	"github.com/zjswh/quickTool/apiparser"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func getOldAPiFuncList(apiOldContent string) map[string]int {
	funcList := map[string]int{}
	reg2 := regexp.MustCompile("func (\\w*)\\(c \\*gin\\.Context\\)")
	result2 := reg2.FindAllStringSubmatch(apiOldContent, -1)
	for _, v := range result2 {
		funcList[v[1]] = 1
	}
	return funcList
}

func getOldServiceFuncList(apiServiceContent string) map[string]int {
	funcList := map[string]int{}
	reg2 := regexp.MustCompile("func (\\w*)")
	result2 := reg2.FindAllStringSubmatch(apiServiceContent, -1)
	for _, v := range result2 {
		funcList[v[1]] = 1
	}
	return funcList
}

func GenApi(routerMap map[string][]apiparser.Router,  projectName, destPath string) error {
	apiDir := destPath + "/api/v1"
	err := MkdirIfNotExist(apiDir)
	if err != nil {
		return err
	}

	for service, routerArr := range routerMap {
		filename := apiDir+"/"+service+".go"
		f, err := os.OpenFile(filename, os.O_CREATE| os.O_APPEND,0600)
		defer f.Close()
		if err !=nil {
			return err
		}
		apiOldContentByte, _ := ioutil.ReadFile(filename)
		apiOldContent := string(apiOldContentByte)
		oldApiFuncList := getOldAPiFuncList(apiOldContent)

		apiFunc := ""
		for _, v := range routerArr {
			if _, ok := oldApiFuncList[ucFirst(v.Handler)]; !ok {
				funcInfo := functionTemplate
				validInfo, varStruct, isDefine := "", "", ":"
				funcInfo = strings.ReplaceAll(funcInfo, "FUNC_NAME", ucFirst(v.Handler))
				if v.Request != "" {
					isDefine = ""
					validInfo = validTemplate
					varStruct = v.Request
					validInfo = strings.ReplaceAll(validInfo, "STRUCT_E", ucFirst(v.Request))
					validInfo = "\n" + validInfo
				}
				funcInfo = strings.ReplaceAll(funcInfo, "VALID_TEMP", validInfo)
				funcInfo = strings.ReplaceAll(funcInfo, "IS_DEFINE", isDefine)
				funcInfo = strings.ReplaceAll(funcInfo, "VAR_STRUCT", varStruct)
				apiFunc += funcInfo
			}
		}
		apiContent := ""
		if len(oldApiFuncList) == 0 && apiOldContent == "" {
			apiContent = apiTemp
			apiContent = strings.ReplaceAll(apiContent, "FUNC_LIST", apiFunc)
			apiContent = strings.ReplaceAll(apiContent, "TEMPLATE", projectName)
		} else {
			apiContent = apiFunc
		}
		apiContent = strings.ReplaceAll(apiContent, "SERVICE_NAME", service+"Service")
		f.Write([]byte(apiContent))
	}
	return err
}

func GenService(routerMap map[string][]apiparser.Router, projectName, destPath string) error {
	for service, routerArr := range routerMap {
		serviceName := service + "Service"
		servicePath := destPath + "/service/" + serviceName
		err := MkdirIfNotExist(servicePath)
		if err != nil {
			break
		}
		filename := servicePath+"/"+serviceName+".go"
		f, err := os.OpenFile(filename, os.O_CREATE| os.O_APPEND,0600)
		defer f.Close()
		if err !=nil {
			return err
		}
		serviceOldContentByte, _ := ioutil.ReadFile(filename)
		serviceOldContent := string(serviceOldContentByte)
		oldServiceFuncList := getOldServiceFuncList(serviceOldContent)
		serviceFunc := ""
		for _, r := range routerArr {
			if _, ok := oldServiceFuncList[ucFirst(r.Handler)]; !ok {
				funcInfo := serviceFunctionTemplate
				funcInfo = strings.ReplaceAll(funcInfo, "FUNC_NAME", ucFirst(r.Handler))
				paramTemplate := ""
				if r.Request != "" {
					paramTemplate = "req types." + ucFirst(r.Request)
				}
				funcInfo = strings.ReplaceAll(funcInfo, "PARAM_TEMP", paramTemplate)
				serviceFunc += funcInfo
			}
		}
		serviceContent := ""
		if len(oldServiceFuncList) == 0 && serviceOldContent == "" {
			serviceContent = serviceTemp
			serviceContent = strings.ReplaceAll(serviceContent, "SERVICE_NAME", serviceName)
			serviceContent = strings.ReplaceAll(serviceContent, "FUNC_LIST", serviceFunc)
			serviceContent = strings.ReplaceAll(serviceContent, "TEMPLATE", projectName)
		} else {
			serviceContent = serviceFunc
		}
		f.Write([]byte(serviceContent))
	}
	return nil
}

func GenRoutes(routerMap map[string][]apiparser.Router,  projectName, destPath string) error {
	//创建文件夹
	routerDir := destPath + "/router"
	err := MkdirIfNotExist(routerDir)
	useMiddleImport := 0
	routerContent := ""
	middlewareMap := map[string]int{}

	for service, routerArr := range routerMap {
		//根据是否使用中间件进行分组
		arrMap := map[string][]apiparser.Router{}
		for _, v := range routerArr {
			arrMap[v.Middle] = append(arrMap[v.Middle], v)
		}
		for middle, groupRouterArr := range arrMap {
			router := fmt.Sprintf("\t%sRouter := Router.Group(\"\")", service + middle)
			if middle != "" {
				router += fmt.Sprintf(".\r\tUse(middleware.%s())", middle)
				middlewareMap[middle] = 1
				useMiddleImport++
			}
			router += "\n\t{\n"
			for _, v := range groupRouterArr {
				router += fmt.Sprintf("\t\t%sRouter.%s(\"%s\", v1.%s)\n", service + middle, strings.ToUpper(v.Method), v.Action, ucFirst(v.Handler))
			}
			router += "\t}\n\n"
			routerContent += router
		}
	}
	middlewareImport := ""
	//判断是否使用了中间件
	if useMiddleImport > 0 {
		middlewareImport = "\""+projectName+"/middleware\""
		//生成中间件
		GenMiddleware(middlewareMap, destPath)
	}
	routerTemplate = strings.ReplaceAll(routerTemplate, "MIDDLEWARE_IMPORT", middlewareImport)
	routerTemplate = strings.ReplaceAll(routerTemplate, "ROUTER_TEMP", routerContent)
	routerTemplate = strings.ReplaceAll(routerTemplate, "TEMPLATE", projectName)
	err = ioutil.WriteFile(routerDir+"/router.go", []byte(routerTemplate), os.ModePerm)
	return err
}

func GenTypes(typesMap map[string]string, destPath string) error {
	//创建文件夹
	requestDir := destPath + "/types"
	err := MkdirIfNotExist(requestDir)
	if err != nil {
		return err
	}
	typesContent := ""
	for _, v := range typesMap {
		typesContent += v + "\r\n\r\n"
	}
	err = ioutil.WriteFile(requestDir+"/types.go", []byte("package types\n\n"+typesContent), os.ModePerm)
	return nil
}

func GenMiddleware(middlewareMap map[string]int, destPath string) error {
	//创建文件夹
	middlewareDir := destPath + "/middleware"
	err := MkdirIfNotExist(middlewareDir)
	if err != nil {
		return err
	}
	funcList := ""
	for k, _ := range middlewareMap {
		funcInfo := middlewareFuncTemplate
		funcInfo = strings.ReplaceAll(funcInfo, "FUNC_NAME", k)
		funcList += funcInfo + "\n\n"
	}
	middlewareContent := middlewareTemplate
	middlewareContent = strings.ReplaceAll(middlewareContent, "FUNC_LIST", funcList)
	err = ioutil.WriteFile(middlewareDir+"/middleware.go", []byte(middlewareContent), os.ModePerm)
	return nil
}

func ucFirst(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:]
}

func MkdirIfNotExist(dir string) error {
	if len(dir) == 0 {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}

