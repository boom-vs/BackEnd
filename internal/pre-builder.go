package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

const (
	ControllerBasePath = "controllers"
	GormModelBasePath  = "models/databaseModels"
	WSModelBasePath    = "models/webSocketModels"
	autoGetDbPath      = "db"
	autoGetRoute       = "webserver"
)

type MVC struct {
	CommonName     string
	ControllerPath string
	GormModelPath  string
	WSModelPath    string
}

func isFileNameEntry(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	name := entry.Name()

	if !strings.HasSuffix(name, ".go") {
		return false
	}

	runes := []rune(name)

	if !unicode.IsUpper(runes[0]) || !unicode.IsLetter(runes[0]) {
		return false
	}
	return true
}

func getCommonName(pathString string) string {
	lastPathSeporator := strings.LastIndex(pathString, string(os.PathSeparator))
	lastDot := strings.LastIndex(pathString, ".")
	return pathString[lastPathSeporator+1 : lastDot]
}

func getItemList(localPath string) (list []string) {
	currPath, err := os.Getwd()
	if err != nil {
		log.Panic(err.Error())
	}

	selectedPath := path.Join(currPath, path.Join(strings.Split(localPath, "/")...))
	files, err := os.ReadDir(selectedPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !isFileNameEntry(file) {
			continue
		}
		list = append(list, path.Join(selectedPath, file.Name()))
	}
	return
}

func GetModels() []MVC {
	mvc := []MVC{}

	//Scan Controllers
	controllers := getItemList(ControllerBasePath)
	for _, controller := range controllers {

		mvc = append(mvc, MVC{
			CommonName:     getCommonName(controller),
			ControllerPath: controller,
		})
	}

	//Scan Gorm Models
	models := getItemList(GormModelBasePath)
	for _, model := range models {

		found := false
		for mvcIndex, mvcKey := range mvc {
			if mvcKey.CommonName == getCommonName(model) {
				mvc[mvcIndex].GormModelPath = model
				found = true
				break
			}
		}

		if !found {
			mvc = append(mvc, MVC{
				CommonName:    getCommonName(model),
				GormModelPath: model,
			})
		}
	}

	//Get WebSocket models
	models = getItemList(WSModelBasePath)
	for _, model := range models {

		found := false
		for mvcIndex, mvcKey := range mvc {
			if mvcKey.CommonName == getCommonName(model) {
				mvc[mvcIndex].WSModelPath = model
				found = true
				break
			}
		}

		if !found {
			mvc = append(mvc, MVC{
				WSModelPath: model,
			})
		}
	}

	return mvc
}

func getGoModulePath() string {
	currPath, _ := os.Getwd()

	for {
		if _, err := os.Stat(path.Join(currPath, "go.mod")); err == nil {
			return currPath
		}
		pathSplit := strings.Split(currPath, string(os.PathSeparator))
		if len(pathSplit) < 2 {
			return ""
		}
		currPath = strings.Join(pathSplit[:len(pathSplit)-2], string(os.PathSeparator))
	}
}

func getGoModuleName(localPath string) string {
	file, err := os.Open(path.Join(localPath, "go.mod"))
	if err != nil {
		return ""
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		if strings.Index(string(line), "module") > -1 {
			splitLine := strings.Split(string(line), " ")
			return splitLine[1]
		}
	}
	return ""
}

func getPackageNameByFilePath(filePath string) (string, error) {
	filePathSplit := strings.Split(filePath, string(os.PathSeparator))
	if len(filePathSplit) < 1 {
		return "", errors.New(fmt.Sprintf("path `%s` so short", filePath))
	}
	return filePathSplit[len(filePathSplit)-1], nil
}

func getGoHeader(getFilePath string, packages []string) (string, error) {
	pkgName, err := getPackageNameByFilePath(getFilePath)
	if err != nil {
		return "", err
	}

	header := fmt.Sprintf("package %s\n\n", pkgName)
	header += fmt.Sprintf("import (\n")
	for _, pkg := range packages {
		header += fmt.Sprintf("\t\"%s\"\n", pkg)
	}
	header += fmt.Sprintf(")\n\n")

	return header, nil
}

func getGormAutoMigration(migrateList []string) string {
	body := "func (db *DataBase) AutoMigrate() error {\n"
	body += "\tif db.Handler == nil {\n"
	body += "\t\treturn errors.New(\"handle is nil\")\n"
	body += "\t}\n"
	body += "\treturn db.Handler.AutoMigrate(\n"
	for _, item := range migrateList {
		body += fmt.Sprintf("\t\t&%s{},\n", item)
	}
	body += "\t)\n"
	body += "}\n"

	return body
}

func checkFileAsGormModel(filePath string, name string) (bool, error) {
	fh, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer fh.Close()

	data, err := io.ReadAll(fh)
	if err != nil {
		return false, err
	}

	res := strings.Index(string(data), "gorm.Model") > -1 &&
		strings.Index(string(data), fmt.Sprintf("type %s struct {", name)) > -1

	errPath := ""
	if !res {
		errPath = filePath
	}

	log.Printf("[check-gorm-model]\tIsControllerFound:%t\tName:%s\t%s\n", res, name, errPath)
	return res, nil
}

func checkFileAsController(filePath string, name string) (bool, error) {
	fh, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer fh.Close()

	data, err := io.ReadAll(fh)
	if err != nil {
		return false, err
	}

	res := strings.Index(string(data), fmt.Sprintf("type Controller%s struct {", name)) > -1

	errPath := ""
	if !res {
		errPath = filePath
	}
	log.Printf("[check-controller]\tIsControllerFound:%t\tName:%s\t%s\n", res, name, errPath)

	return res, nil
}

func getMethodsFromController(filePath string) ([]string, error) {
	fh, err := os.Open(filePath)
	if err != nil {
		return []string{}, err
	}
	defer fh.Close()

	methods := []string{}

	fileNameSplit := strings.Split(filepath.Base(filePath), ".")
	controllerName := fileNameSplit[0]

	reader := bufio.NewReader(fh)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		if strings.Index(string(line), "func") > -1 &&
			strings.Index(string(line), fmt.Sprintf("Controller%s", controllerName)) > -1 &&
			strings.Index(string(line), "*types.RequestContext") > -1 {

			beginCopy := strings.Index(string(line), ") ") + len(") ")
			endCopy := strings.LastIndex(string(line), "(")
			if endCopy <= beginCopy {
				continue
			}

			newLine := string(line)
			newLine = newLine[beginCopy:endCopy]
			line = []byte(newLine)
			if !unicode.IsUpper(rune(line[0])) {
				continue
			}
			log.Printf("[check-actions]\tAction:%s\n", newLine)

			methods = append(methods, newLine)
		}
	}

	return methods, nil
}

func getAutoRouteBody(routes map[string][]string) string {
	body := "func (sm *SocketManager) Enumerate() {\n"
	body += "\tsm.controllers = make(map[string]map[string]func(context *types.RequestContext))\n"

	for mapKey, mapValue := range routes {
		mapKeySplit := strings.Split(mapKey, ".")
		if len(mapKeySplit) < 2 {
			continue
		}

		controllerName := mapKeySplit[1][len("contoller")+1:]

		instanceName := []byte(controllerName)
		instanceName[0] = byte(unicode.ToLower(rune(instanceName[0])))
		actionMapName := "actions" + controllerName

		body += fmt.Sprintf("\t%s := &%s{}\n", instanceName, mapKey)
		body += fmt.Sprintf("\t%s := make(map[string]func(context *types.RequestContext))\n", actionMapName)
		for _, action := range mapValue {
			body += fmt.Sprintf("\t%s[\"%s\"] = %s.%s\n", actionMapName, action, instanceName, action)
		}
		body += fmt.Sprintf("\tsm.controllers[\"%s\"] = %s\n", controllerName, actionMapName)
		body += "\n"
	}

	body += "}\n"

	return body
}

func AutoGetDb(mvc []MVC) {

	goModFilePath := getGoModulePath()
	moduleName := getGoModuleName(goModFilePath)

	packages := make(map[string]bool)
	autoMigrateList := []string{}

	for _, mvcValue := range mvc {
		if len(mvcValue.GormModelPath) == 0 {
			continue
		}

		relPath, err := filepath.Rel(goModFilePath, mvcValue.GormModelPath)

		if err != nil {
			continue
		}

		relPath = filepath.Join(moduleName, relPath)
		relPath = filepath.Dir(relPath)

		if _, ok := packages[relPath]; !ok {
			packages[relPath] = true
		}

		ok, err := checkFileAsGormModel(mvcValue.GormModelPath, mvcValue.CommonName)
		if err != nil {
			log.Println("[gen-auto-migration]: " + err.Error())
			continue
		}

		if ok {
			pkgName, err := getPackageNameByFilePath(filepath.Dir(mvcValue.GormModelPath))
			if err != nil {
				log.Println("[gen-auto-migration]: " + err.Error())
				continue
			}
			autoMigrateList = append(autoMigrateList, fmt.Sprintf("%s.%s", pkgName, mvcValue.CommonName))
		} else {
			log.Printf("[gen-auto-migration]: File `%s.go` doesn't contain gorm model\n", mvcValue.CommonName)
			continue
		}

	}

	pkgList := []string{}
	pkgList = append(pkgList, "errors")
	for key := range packages {
		pkgList = append(pkgList, key)
	}

	header, err := getGoHeader(autoGetDbPath, pkgList)

	if err != nil {
		log.Println("[gen-auto-migration]: " + err.Error())
		return
	}

	body := getGormAutoMigration(autoMigrateList)

	err = os.WriteFile(filepath.Join(autoGetDbPath, "autoGen.go"), []byte(header+body), 0666)
	if err != nil {
		log.Println("[gen-auto-migration]: " + err.Error())
	}
}

func AutoRoute(mvc []MVC) {
	goModFilePath := getGoModulePath()
	moduleName := getGoModuleName(goModFilePath)

	packages := make(map[string]bool)
	autoRouteList := make(map[string][]string)

	for _, mvcValue := range mvc {
		if len(mvcValue.ControllerPath) == 0 {
			continue
		}

		relPath, err := filepath.Rel(goModFilePath, mvcValue.ControllerPath)

		if err != nil {
			continue
		}

		relPath = filepath.Join(moduleName, relPath)
		relPath = filepath.Dir(relPath)

		if _, ok := packages[relPath]; !ok {
			packages[relPath] = true
		}

		ok, err := checkFileAsController(mvcValue.ControllerPath, mvcValue.CommonName)
		if err != nil {
			log.Println("[gen-auto-route]: " + err.Error())
			continue
		}

		if ok {
			pkgName, err := getPackageNameByFilePath(filepath.Dir(mvcValue.ControllerPath))
			if err != nil {
				log.Println("[gen-auto-route]: " + err.Error())
				continue
			}

			methodsList, _ := getMethodsFromController(mvcValue.ControllerPath)

			autoRouteList[fmt.Sprintf("%s.Controller%s", pkgName, mvcValue.CommonName)] = methodsList

		} else {
			continue
		}
	}

	pkgList := []string{}
	pkgList = append(pkgList, "crm-backend/types")
	for key := range packages {
		pkgList = append(pkgList, key)
	}
	header, err := getGoHeader(autoGetRoute, pkgList)
	if err != nil {
		log.Println("[gen-auto-route]: " + err.Error())
	}

	body := getAutoRouteBody(autoRouteList)
	err = os.WriteFile(filepath.Join(autoGetRoute, "autoGen.go"), []byte(header+body), 0666)
	if err != nil {
		log.Println("[gen-auto-migration]: " + err.Error())
	}
}
