package handler

import (
	"encoding/json"
	"mnc-be-tech-test/app/db"
	"mnc-be-tech-test/app/model"
)

type HistoryHandler struct {
	RESTHandler
}

func (h HistoryHandler) DB(state string) db.DBHandler {
	return h.RESTHandler.DB(state).Model(model.History{})
}

func (h HistoryHandler) New() interface{} {
	return &model.History{}
}

func (h HistoryHandler) GetModel(data []byte) interface{} {
	model := []*model.History{}
	json.Unmarshal(data, &model)
	return model
}

func (h HistoryHandler) KeywordMap() map[string]*Keyword {
	return map[string]*Keyword{
		"Module": {Type: StringContains, ColumnName: "Module"},
	}
}

func (h HistoryHandler) Find(db db.DBHandler) interface{} {
	data := []model.History{}
	db.Order("id", true).Find(&data)
	return data
}

func (h HistoryHandler) ModelName() string {
	return "histories"
}

// GetID is get id from interface
func (h HistoryHandler) GetID(v interface{}) uint {
	result := v.(*model.History)
	return result.ID
}
