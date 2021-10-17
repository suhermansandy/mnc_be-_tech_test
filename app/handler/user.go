package handler

import (
	"encoding/json"
	"errors"
	"mnc-be-tech-test/app/db"
	"mnc-be-tech-test/app/model"
	"strconv"
)

type UserHandler struct {
	RESTHandler
}

func (h UserHandler) DB(state string) db.DBHandler {
	return h.RESTHandler.DB(state).Model(model.User{})
}

func (h UserHandler) New() interface{} {
	return &model.User{}
}

func (h UserHandler) GetModel(data []byte) interface{} {
	model := []*model.User{}
	json.Unmarshal(data, &model)
	return model
}

func (h UserHandler) KeywordMap() map[string]*Keyword {
	return map[string]*Keyword{
		"user_name": {Type: StringContains, ColumnName: "user_name"},
		"password":  {Type: StringContains, ColumnName: "password"},
		"name":      {Type: StringContains, ColumnName: "name"},
	}
}

func (h UserHandler) Find(db db.DBHandler) interface{} {
	data := []model.User{}
	db.Order("id", true).Find(&data)
	return data
}

func (h UserHandler) CreateUserId(userID string, v interface{}) interface{} {
	data, ok := v.(*model.User)
	if !ok {
		return errors.New("Invalid Type")
	}

	data.CreatedBy = userID
	data.UpdatedBy = userID

	return data
}

func (h UserHandler) UpdateUserId(userID string, v interface{}) interface{} {
	data, ok := v.(*model.User)
	if !ok {
		return errors.New("Invalid Type")
	}

	data.UpdatedBy = userID

	return data
}

func (h UserHandler) DeleteUserId(userID string, v interface{}) interface{} {
	data, ok := v.(*model.User)
	if !ok {
		return errors.New("Invalid Type")
	}

	data.DeletedBy = userID
	data.UpdatedBy = userID

	return data
}

func (h UserHandler) Validate(state string, v interface{}) error {
	data, ok := v.(*model.User)
	if !ok {
		return errors.New("Invalid Type")
	}

	strFieldMap := map[string]strFieldInfo{
		"UserName": {Value: data.UserName, IsMandatory: true},
		"Password": {Value: data.Password, IsMandatory: true},
		"Name":     {Value: data.Name, IsMandatory: true},
	}

	if err := validateStrings(strFieldMap); err != nil {
		return err
	}

	dbUser := h.RESTHandler.DB(state).Model(model.User{})
	modelUser := model.User{}
	dbUser.Where("user_name = ?", data.UserName).Last(&modelUser)
	if modelUser.ID != 0 {
		return errors.New("User Name Alredy Exist")
	}

	return nil
}

func (h UserHandler) ModelName() string {
	return "users"
}

// GetID is get id from interface
func (h UserHandler) GetID(v interface{}) uint {
	result := v.(*model.User)
	return result.ID
}

func (h UserHandler) AfterCreate(state string, v interface{}) error {
	result, ok := v.(*model.User)
	if !ok {
		return errors.New("Invalid Type")
	}
	dbHis := h.RESTHandler.DB(state).Model(model.History{})
	modelsHis := model.History{}

	modelsHis.Note = "Create New User " + result.Name + " and ID = " + strconv.Itoa(int(result.ID))
	modelsHis.Module = "User"
	dbHis.Save(&modelsHis)

	return nil
}
