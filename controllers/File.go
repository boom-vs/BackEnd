package controllers

import (
	"crm-backend/internal"
	"crm-backend/models/databaseModels"
	"crm-backend/models/webSocketModels"
	"crm-backend/types"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"io"
	"log"
	"os"
	"path"
)

type ControllerFile struct {
}

func (cf *ControllerFile) log(message string) {
	log.Println(message)
}

func (cf *ControllerFile) Create(context *types.RequestContext) {
	var request []webSocketModels.FileCreate

	mapToStructConfig := &mapstructure.DecoderConfig{
		ErrorUnset: true,
		Result:     &request,
	}
	mapToStruct, err := mapstructure.NewDecoder(mapToStructConfig)
	if err != nil {
		cf.log("Login: " + err.Error())
	}
	err = mapToStruct.Decode(context.ReceivedData)

	if err != nil {
		context.Response.Error = fmt.Sprintf("Invalid structure: `%+v` %s", context.ReceivedData, err.Error())
		if err != nil {
			cf.log(err.Error())
		}
		return
	}

	if len(request) == 0 {
		cf.log(err.Error())
		return
	}

	createRequest := request[0]

	file := &databaseModels.File{
		Name:       createRequest.Name,
		Size:       createRequest.Size,
		HashSum:    createRequest.Hash,
		UploaderId: context.EmployeeId,
	}
	context.Base.Create(file)

	if file.ID == 0 {
		context.Response.Error = fmt.Sprintf("Can't create file, db reson")
		if err != nil {
			cf.log(err.Error())
		}
		return
	}

	context.Response.Data = append(context.Response.Data, webSocketModels.FileToken{
		Token: file.Token,
	})
}

func (cf *ControllerFile) getFileHander(fileName, tmpPath string) (*os.File, error) {

	filePath := path.Join(tmpPath, fileName)
	fileNandler, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	return fileNandler, err
}

func (cf *ControllerFile) Upload(context *types.RequestContext) {
	var request []webSocketModels.FileUpload

	mapToStructConfig := &mapstructure.DecoderConfig{
		ErrorUnset: true,
		Result:     &request,
	}

	mapToStruct, err := mapstructure.NewDecoder(mapToStructConfig)
	if err != nil {
		cf.log(err.Error())
		return
	}

	err = mapToStruct.Decode(context.ReceivedData)

	if len(request) == 0 {
		context.Response.Error = "Empty data"
		return
	}

	uploadRequest := request[0]

	dbFile := &databaseModels.File{}

	context.Base.First(dbFile, "token = ?", uploadRequest.Token)
	if dbFile.ID == 0 {
		context.Response.Error = fmt.Sprintf("Can't find token `%s`", uploadRequest.Token)
		return
	}

	if dbFile.UploaderId != context.EmployeeId {
		context.Response.Error = "user did not match"
		return
	}

	if uploadRequest.ByteShift > dbFile.RealSize+1 {
		context.Response.Error = "you must transfer the file from the beginning"
		return
	}

	if dbFile.Size == dbFile.RealSize {
		context.Response.Error = "file already uploaded"
		return
	}

	fileHandler, err := cf.getFileHander(dbFile.Token, context.Env["TMP"])
	if err != nil {
		context.Response.Error = fmt.Sprintf("Can't create dbFile %s", uploadRequest.Token)
		return
	}

	defer func() {
		err = fileHandler.Close()
		if err != nil {
			cf.log(err.Error())
		}
	}()

	stat, err := fileHandler.Stat()
	if err != nil {
		context.Response.Error = fmt.Sprintf("Can't get fileStat dbFile %s", uploadRequest.Token)
		return
	}

	dbFile.RealSize = stat.Size()

	context.Socket.WriteJSON(webSocketModels.FileInfo{
		Name:     dbFile.Name,
		Hash:     dbFile.HashSum,
		Size:     dbFile.Size,
		RealSize: dbFile.RealSize,
	})

	for dbFile.Size != dbFile.RealSize {
		messageType, bytes, err := context.Socket.ReadMessage()
		if messageType != websocket.BinaryMessage {
			context.Response.Error = "error message type"
			context.Response.Data = append(context.Response.Data, webSocketModels.FileInfo{
				Name:     dbFile.Name,
				Hash:     dbFile.HashSum,
				Size:     dbFile.Size,
				RealSize: dbFile.RealSize,
			})
			return
		}

		recorded, err := fileHandler.Write(bytes)
		if err != nil {
			context.Response.Error = fmt.Sprintf("Can't write date to dbFile %s", uploadRequest.Token)
			return
		}

		if recorded != len(bytes) {
			context.Response.Error =
				fmt.Sprintf("Size of the recorded buffer (%d) is smaller than the received buffer (%d)",
					recorded,
					len(bytes))
			return
		}

		err = fileHandler.Sync()
		if err != nil {
			context.Response.Error = "Sync err"
			return
		}

		dbFile.RealSize += int64(recorded)

		if dbFile.RealSize < dbFile.Size {
			context.Response.Error =
				fmt.Sprintf("Uploaded dbFile `%s` size `%d` more then decloreted `%d`",
					dbFile.Name, dbFile.RealSize, dbFile.Size)
		}

	}
	fileHandler.Seek(0, 0)

	dbFile.Data, err = io.ReadAll(fileHandler)
	if err != nil {
		context.Response.Error = fmt.Sprintf("Can't read date from dbFile %s",
			path.Join("/tmp", uploadRequest.Token))
		return
	}

	if dbFile.HashSum != internal.GetHash(dbFile.Data) {
		context.Response.Error = "hash sum does not match"
		return
	}

	context.Base.Updates(dbFile)

	err = os.Remove(path.Join(context.Env["TMP"], uploadRequest.Token))
	if dbFile.HashSum != internal.GetHash(dbFile.Data) {
		context.Response.Error = "can't remove file " + uploadRequest.Token
		return
	}

	context.Socket.Close()
}
