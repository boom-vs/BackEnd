package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"unicode"
)

func makeControllerBody(name string) string {

	shortName := "c"
	bytes := []byte(name)
	for _, thisByte := range bytes {
		if !unicode.IsUpper(rune(thisByte)) {
			continue
		}

		shortName += string(unicode.ToLower(rune(thisByte)))
	}

	buff := []byte(name)
	buff[0] = byte(unicode.ToLower(rune(buff[0])))
	lowCapitalCaseName := string(buff)

	//Structure
	body := fmt.Sprintf("type Controller%s struct {\n}\n\n", name)

	//log
	body += fmt.Sprintf("func (%s *Controller%s) log(message string) {\n\tlog.Println(message)\n}\n\n",
		shortName, name)

	//setOne
	body += fmt.Sprintf("func (%s *Controller%s) setOne(context *types.RequestContext,\n", shortName, name)
	body += fmt.Sprintf("\t%s *webSocketModels.%s) (error, *webSocketModels.%s) {\n",
		lowCapitalCaseName, name, name)
	body += fmt.Sprintf("\tresult := &databaseModels.%s{}\n", name)
	body += fmt.Sprintf("\tif %s.ID == 0 {\n", lowCapitalCaseName)
	body += fmt.Sprintf("\t\tif %s.SortNumber == 0 {\n", lowCapitalCaseName)
	body += fmt.Sprintf("\t\t\t%s.SortNumber = types.DefaultSortNumber\n", lowCapitalCaseName)
	body += fmt.Sprintf("\t\t}\n")
	body += fmt.Sprintf("\t\tinternal.Copier(%s, result)\n\n", lowCapitalCaseName)
	body += fmt.Sprintf("\t\tdb := context.Base.Create(&result)\n")
	body += fmt.Sprintf("\t\tif result.ID == 0 || db.Error != nil {\n")
	body += fmt.Sprintf("\t\t\treturn errors.New(fmt.Sprintf(\"Some problem with db\")), nil\n")
	body += fmt.Sprintf("\t\t}\n")
	body += fmt.Sprintf("\t} else {\n")
	body += fmt.Sprintf("\t\tdb := context.Base.First(result, \"id = ?\", %s.ID)\n\n",
		lowCapitalCaseName)
	body += fmt.Sprintf("\t\tif db.Error != nil {\n")
	body += fmt.Sprintf("\t\t\t%s.log(db.Error.Error())\n", shortName)
	body += fmt.Sprintf("\t\t\treturn errors.New(fmt.Sprintf(\"Some problem with db \" + db.Error.Error())), nil\n")
	body += fmt.Sprintf("\t\t}\n\n")
	body += fmt.Sprintf("\t\tinternal.Copier(%s, result)\n", lowCapitalCaseName)
	body += fmt.Sprintf("\t\tresult.ID = %s.ID\n\n", lowCapitalCaseName)
	body += fmt.Sprintf("\t\tdb = context.Base.Updates(result)\n\n")
	body += fmt.Sprintf("\t\tif db.Error != nil {\n")
	body += fmt.Sprintf("\t\t\t%s.log(db.Error.Error())\n", shortName)
	body += fmt.Sprintf("\t\t\treturn errors.New(fmt.Sprintf(\"Some problem with db\" + db.Error.Error())), nil\n")
	body += fmt.Sprintf("\t\t}\n")
	body += fmt.Sprintf("\t}\n\n")
	body += fmt.Sprintf("\tinternal.Copier(result, %s)\n", lowCapitalCaseName)
	body += fmt.Sprintf("\t%s.ID = result.ID\n", lowCapitalCaseName)
	body += fmt.Sprintf("\treturn nil, %s\n", lowCapitalCaseName)
	body += fmt.Sprintf("}\n\n")

	//Set
	body += fmt.Sprintf("func (%s *Controller%s) Set(context *types.RequestContext) {\n", shortName, name)
	body += fmt.Sprintf("\tvar request []webSocketModels.%s\n\n", name)
	body += fmt.Sprintf("\tmapToStructConfig := &mapstructure.DecoderConfig{\n")
	body += fmt.Sprintf("\t\tErrorUnset: true,\n")
	body += fmt.Sprintf("\t\tResult: &request,\n")
	body += fmt.Sprintf("\t}\n\n")
	body += fmt.Sprintf("\tmapToStruct, err := mapstructure.NewDecoder(mapToStructConfig)\n")
	body += fmt.Sprintf("\tif err != nil {\n")
	body += fmt.Sprintf("\t\t%s.log(err.Error())\n", shortName)
	body += fmt.Sprintf("\t\treturn\n")
	body += fmt.Sprintf("\t}\n\n")
	body += fmt.Sprintf("\terr = mapToStruct.Decode(context.ReceivedData)\n")
	body += fmt.Sprintf("\tif err != nil {\n")
	body += fmt.Sprintf("\t\tcontext.Response.Error = fmt.Sprintf(\"Invalid structure: `%%+v` %%s\"," +
		" context.ReceivedData, err.Error())\n")
	body += fmt.Sprintf("\t\treturn\n")
	body += fmt.Sprintf("\t}\n\n")
	body += fmt.Sprintf("\tvar updater []interface{}\n\n")
	body += fmt.Sprintf("\tfor _, %s := range request {\n", lowCapitalCaseName)
	body += fmt.Sprintf("\t\terr, item := %s.setOne(context, &%s)\n", shortName, lowCapitalCaseName)
	body += fmt.Sprintf("\t\tif err != nil {\n")
	body += fmt.Sprintf("\t\t\tcontext.Response.Error = err.Error()\n")
	body += fmt.Sprintf("\t\t\treturn\n")
	body += fmt.Sprintf("\t\t}\n")
	body += fmt.Sprintf("\t\tupdater = append(updater, item)\n")
	body += fmt.Sprintf("\t}\n")
	body += fmt.Sprintf("\tcontext.Updater(\"%s\", updater)\n", name)
	body += fmt.Sprintf("}\n\n")

	//get list
	body += fmt.Sprintf("func (%s *Controller%s) GetList(context *types.RequestContext) {\n", shortName, name)
	body += fmt.Sprintf("\tvar %ss []*databaseModels.%s\n\n", lowCapitalCaseName, name)
	body += fmt.Sprintf("\tcontext.Base.Order(\"sort_number, name\").Find(&%ss)\n",
		lowCapitalCaseName)
	body += fmt.Sprintf("\tfor _, %s := range %ss {\n", lowCapitalCaseName, lowCapitalCaseName)
	body += fmt.Sprintf("\t\tresponse := &webSocketModels.%s{}\n", name)
	body += fmt.Sprintf("\t\tinternal.Copier(%s, response)\n", lowCapitalCaseName)
	//body += fmt.Sprintf("\t\tresponse.ID = %s.ID\n", lowCapitalCaseName) // ???
	body += fmt.Sprintf("\t\tcontext.Response.Data = append(context.Response.Data, response)\n")
	body += fmt.Sprintf("\t}\n")
	body += fmt.Sprintf("}\n")

	return body
}

func AutoControllers(mvc []MVC) {

	/*
		import (
			"crm-backend/internal"
			"crm-backend/models/databaseModels"
			"crm-backend/models/webSocketModels"
			"crm-backend/types"
			"errors"
			"fmt"
			"github.com/mitchellh/mapstructure"
			"log"
		)
	*/

	header, err := getGoHeader(ControllerBasePath, []string{
		"crm-backend/internal",
		"crm-backend/models/databaseModels",
		"crm-backend/models/webSocketModels",
		"crm-backend/types",
		"errors",
		"fmt",
		"github.com/mitchellh/mapstructure",
		"log",
	})

	if err != nil {
		log.Println("[controller-builder] cant' make header " + err.Error())
		return
	}

	for _, mvcItem := range mvc {
		if len(mvcItem.ControllerPath) != 0 || len(mvcItem.GormModelPath) == 0 || len(mvcItem.WSModelPath) == 0 {
			continue
		}

		body := makeControllerBody(mvcItem.CommonName)

		err = os.WriteFile(filepath.Join(ControllerBasePath, mvcItem.CommonName+".go"), []byte(header+body), 0666)
		if err != nil {
			log.Println("[gen-auto-conroller]: " + err.Error())
		} else {
			log.Println("[gen-auto-conroller]: " + filepath.Join(ControllerBasePath, mvcItem.CommonName+".go"))
		}
	}

}
