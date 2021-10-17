package handler

import (
	"encoding/json"
	"errors"
	"mnc-be-tech-test/app/db"
	"mnc-be-tech-test/app/model"
	"strconv"
)

type PaymentHandler struct {
	RESTHandler
}

func (h PaymentHandler) DB(state string) db.DBHandler {
	return h.RESTHandler.DB(state).Model(model.Payment{})
}

func (h PaymentHandler) New() interface{} {
	return &model.Payment{}
}

func (h PaymentHandler) GetModel(data []byte) interface{} {
	model := []*model.Payment{}
	json.Unmarshal(data, &model)
	return model
}

func (h PaymentHandler) KeywordMap() map[string]*Keyword {
	return map[string]*Keyword{
		"from":            {Type: StringContains, ColumnName: "from"},
		"to":              {Type: StringContains, ColumnName: "to"},
		"nominal":         {Type: Equal, ColumnName: "nominal"},
		"created_at_from": {Type: MoreThanOrEqual, ColumnName: "created_at"},
		"created_at_to":   {Type: LessThanOrEqual, ColumnName: "created_at"},
	}
}

func (h PaymentHandler) Find(db db.DBHandler) interface{} {
	data := []model.Payment{}
	db.Order("id", true).Find(&data)
	return data
}

func (h PaymentHandler) CreateUserId(userID string, v interface{}) interface{} {
	data, ok := v.(*model.Payment)
	if !ok {
		return errors.New("Invalid Type")
	}

	data.CreatedBy = userID
	data.UpdatedBy = userID

	return data
}

func (h PaymentHandler) UpdateUserId(userID string, v interface{}) interface{} {
	data, ok := v.(*model.Payment)
	if !ok {
		return errors.New("Invalid Type")
	}

	data.UpdatedBy = userID

	return data
}

func (h PaymentHandler) DeleteUserId(userID string, v interface{}) interface{} {
	data, ok := v.(*model.Payment)
	if !ok {
		return errors.New("Invalid Type")
	}

	data.DeletedBy = userID
	data.UpdatedBy = userID

	return data
}

func (h PaymentHandler) Validate(state string, v interface{}) error {
	data, ok := v.(*model.Payment)
	if !ok {
		return errors.New("Invalid Type")
	}

	strFieldMap := map[string]strFieldInfo{
		"From": {Value: data.From, IsMandatory: true},
		"To":   {Value: data.To, IsMandatory: true},
	}

	if err := validateStrings(strFieldMap); err != nil {
		return err
	}

	dbUserFrom := h.RESTHandler.DB(state).Model(model.User{})
	modelUserFrom := model.User{}
	dbUserFrom.Where("user_name = ?", data.From).Last(&modelUserFrom)
	if modelUserFrom.ID == 0 {
		return errors.New("User Send Not Exist")
	}

	dbUserTo := h.RESTHandler.DB(state).Model(model.User{})
	modelUserTo := model.User{}
	dbUserTo.Where("user_name = ?", data.To).Last(&modelUserTo)
	if modelUserTo.ID == 0 {
		return errors.New("User Not Exist")
	}

	return nil
}

func (h PaymentHandler) ModelName() string {
	return "payments"
}

// GetID is get id from interface
func (h PaymentHandler) GetID(v interface{}) uint {
	result := v.(*model.Payment)
	return result.ID
}

func (h PaymentHandler) AfterCreate(state string, v interface{}) error {
	result, ok := v.(*model.Payment)
	if !ok {
		return errors.New("Invalid Type")
	}
	dbHis := h.RESTHandler.DB(state).Model(model.History{})
	modelsHis := model.History{}

	modelsHis.Note = "Create Payment, send From " + result.From + " To " + result.To + " and ID = " + strconv.Itoa(int(result.ID))
	modelsHis.Module = "Payment"
	dbHis.Save(&modelsHis)

	return nil
}
