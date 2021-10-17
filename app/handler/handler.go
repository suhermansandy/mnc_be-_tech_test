package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mnc-be-tech-test/app/db"
	"mnc-be-tech-test/app/excel"
	"mnc-be-tech-test/app/model"
	"mnc-be-tech-test/util"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type RESTHandler interface {
	DefaultDB(state string) db.DBHandler
	DB(state string) db.DBHandler
	ChangeLogDB(state string) db.DBHandler
	Artemis() *model.Artemis
	Excel() excel.ExcelHandler
	New() interface{}
	GetModel(d []byte) interface{}
	// Mapping(x []interface{}, y []interface{}) []interface{}
	GetID(v interface{}) uint
	KeywordMap() map[string]*Keyword
	SuggestionKeywords() []Keyword
	Find(db.DBHandler) interface{}
	Validate(state string, v interface{}) error

	BeforeDelete(state string, v interface{}) error

	AfterCreate(state string, v interface{}) error
	AfterCreateUpload(state string, v []interface{}) error
	AfterUpdateUpload(state string, v []interface{}) error

	ModelName() string
	CsvName() string
	ExcelHeaderMap() map[string]ExcelHeader

	MappingToModel(state string, line []string) (interface{}, string, error)
	StructIndex() int

	GetFieldNames() []string
	GetDBNames() []string
	GetFieldNamesUpdate() []string
	GetDBNamesUpdate() []string
	GetKeyFieldNames() []string

	CreateUserId(userID string, y interface{}) interface{}
	UpdateUserId(userID string, y interface{}) interface{}
	DeleteUserId(userID string, y interface{}) interface{}

	MappingFirmUpdate(x interface{}, y interface{}) interface{}
}

type DefaultHandler struct {
	DBHandler      map[string]db.DBHandler
	ExcelHandler   excel.ExcelHandler
	ArtemisHandler *model.Artemis
}

func (h DefaultHandler) DefaultDB(state string) db.DBHandler {
	return h.DBHandler[state]
}

func (h DefaultHandler) DB(state string) db.DBHandler {
	return h.DBHandler[state]
}

func (h DefaultHandler) ChangeLogDB(state string) db.DBHandler {
	return h.DBHandler[state].Model(model.ChangeLog{})
}

func (h DefaultHandler) Artemis() *model.Artemis {
	return h.ArtemisHandler
}

func (h DefaultHandler) Excel() excel.ExcelHandler {
	return h.ExcelHandler
}

func (h DefaultHandler) New() interface{} {
	return gorm.Model{}
}

func (h DefaultHandler) GetModel(p []byte) interface{} {
	return make([]*gorm.Model, 0)
}

// func (h DefaultHandler) Mapping(v []interface{}, u []interface{}) []interface{} {
// 	return []*gorm.Model{}
// }

func (h DefaultHandler) GetID(v interface{}) uint {
	return 0
}

func (h DefaultHandler) KeywordMap() map[string]*Keyword {
	return make(map[string]*Keyword, 0)
}

func (h DefaultHandler) SuggestionKeywords() []Keyword {
	return make([]Keyword, 0)
}

func (h DefaultHandler) Find(db db.DBHandler) interface{} {
	return make([]gorm.Model, 0)
}

func (h DefaultHandler) Validate(state string, v interface{}) error {
	return nil
}

func (h DefaultHandler) BeforeDelete(state string, v interface{}) error {
	return nil
}

func (h DefaultHandler) AfterCreate(state string, v interface{}) error {
	return nil
}

func (h DefaultHandler) AfterCreateUpload(state string, v []interface{}) error {
	return nil
}

func (h DefaultHandler) AfterUpdateUpload(state string, v []interface{}) error {
	return nil
}

func (h DefaultHandler) ModelName() string {
	return "empty"
}

func (h DefaultHandler) CsvName() string {
	return "empty"
}

func (h DefaultHandler) ExcelHeaderMap() map[string]ExcelHeader {
	return make(map[string]ExcelHeader)
}

func New(db map[string]db.DBHandler, excel excel.ExcelHandler, artemis *model.Artemis) RESTHandler {
	return DefaultHandler{DBHandler: db, ExcelHandler: excel, ArtemisHandler: artemis}
}

// MappingToModel excel to struct
func (h DefaultHandler) MappingToModel(state string, line []string) (interface{}, string, error) {
	return nil, "", nil
}

// StructIndex file
func (h DefaultHandler) StructIndex() int {
	return 1
}

// GetFieldNames data
func (h DefaultHandler) GetFieldNames() []string {
	return []string{}
}

// GetDBNames data
func (h DefaultHandler) GetDBNames() []string {
	return []string{}
}

// GetFieldNames data
func (h DefaultHandler) GetFieldNamesUpdate() []string {
	return []string{}
}

// GetDBNames data
func (h DefaultHandler) GetDBNamesUpdate() []string {
	return []string{}
}

// GetKeyFieldNames data
func (h DefaultHandler) GetKeyFieldNames() []string {
	return []string{}
}

// CreateUserId data
func (h DefaultHandler) CreateUserId(userID string, u interface{}) interface{} {
	return gorm.Model{}
}

// UpdateUserId data
func (h DefaultHandler) UpdateUserId(userID string, u interface{}) interface{} {
	return gorm.Model{}
}

// DeleteUserId data
func (h DefaultHandler) DeleteUserId(userID string, u interface{}) interface{} {
	return gorm.Model{}
}

func (h DefaultHandler) MappingFirmUpdate(v interface{}, u interface{}) interface{} {
	return gorm.Model{}
}

// KeywordType is naming data type query action
type KeywordType string

// const is for query action
const (
	StringContains  KeywordType = "STRING_CONTAINS"
	IntContains     KeywordType = "INT_CONTAINS"
	Equal           KeywordType = "EQUAL"
	MoreThanOrEqual KeywordType = "MORE_THAN_OR_EQUAL"
	LessThanOrEqual KeywordType = "LESS_THAN_OR_EQUAL"
)

// Keyword is type for query action to a column
type Keyword struct {
	Type       KeywordType
	ColumnName string
	Value      string
}

// ExcelHeader is for mapping excel column
type ExcelHeader struct {
	FieldName   string
	DisplayName string
}

// Query Keyword used to filter
type ResultQueryKeyword struct {
	Field string
	Value string
}

// ResultQuery model for query paging
type ResultQuery struct {
	Offset  int
	Limit   int
	Total   int
	Keyword []ResultQueryKeyword
}

// Result model for result paging
type Result struct {
	Query ResultQuery
	Data  interface{}
}

// CreateExcel is function to create an excel file
func CreateExcel(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	if state != "mock" {
		state = "db"
	}

	skipColumnStr := q.Get("hideColumn")
	skipColumnArr := strings.Split(skipColumnStr, ",")
	skipColumnMap := make(map[string]struct{})
	for _, col := range skipColumnArr {
		skipColumnMap[col] = struct{}{}
	}

	xlsx := h.Excel().NewFile()
	sheet1Name := "Sheet One"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

	excelHeaderMap := h.ExcelHeaderMap()
	headerColumns := make([]string, len(excelHeaderMap))
	index := 0
	for _, headerMap := range excelHeaderMap {
		if _, ok := skipColumnMap[headerMap.FieldName]; !ok {
			headerColumn := util.GetExcelColumnName(index + 1)
			headerCell := headerColumn + "1"
			xlsx.SetCellValue(sheet1Name, headerCell, headerMap.DisplayName)
			headerColumns[index] = headerColumn
			index++
		}
	}

	dbData := h.Find(h.DB(state))
	reflectedDbData := reflect.ValueOf(dbData)
	if reflectedDbData.Kind() != reflect.Slice {
		respondError(w, http.StatusInternalServerError, errors.New("Invalid Data").Error())
		return
	}

	results := make([]interface{}, reflectedDbData.Len())
	for i := 0; i < reflectedDbData.Len(); i++ {
		results[i] = reflectedDbData.Index(i).Interface()
	}

	index = 0
	for _, headerMap := range excelHeaderMap {
		if _, ok := skipColumnMap[headerMap.FieldName]; !ok {
			headerColumn := util.GetExcelColumnName(index + 1)
			for i, result := range results {
				resultVal := reflect.ValueOf(result)
				field := reflect.Indirect(resultVal.FieldByName(headerMap.FieldName))
				cell := fmt.Sprintf("%s%d", headerColumn, i+2)
				if !field.IsValid() {
					xlsx.SetCellValue(sheet1Name, cell, nil)
				} else {
					fieldVal := fmt.Sprintf("%v", field.Interface())
					xlsx.SetCellValue(sheet1Name, cell, fieldVal)
				}
			}
			index++
		}
	}

	sort.Slice(headerColumns, func(i, j int) bool {
		return len(headerColumns[i]) < len(headerColumns[j]) ||
			(len(headerColumns[i]) == len(headerColumns[j]) && headerColumns[i] < headerColumns[j])
	})
	firstHeaderCell := headerColumns[0] + "1"
	lastHeaderCell := fmt.Sprintf("%s%d", headerColumns[len(headerColumns)-1], len(results)+1)

	err := xlsx.AutoFilter(sheet1Name, firstHeaderCell, lastHeaderCell, "")
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	downloadName := time.Now().Local().Format(h.ModelName() + "-20060102150405.xlsx")
	w.Header().Set("Content-Type", "application/octet-stream")
	contentDisposition := fmt.Sprintf("attachment; filename=%s", downloadName)
	w.Header().Set("Content-Disposition", contentDisposition)
	w.Header().Set("File-Name", downloadName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	err = xlsx.Write(w)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// GetList is a function to get a list of data
func GetList(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	if state != "mock" {
		state = "db"
	}

	queryResult := Result{
		Query: ResultQuery{
			Keyword: []ResultQueryKeyword{},
		},
	}

	qOffset := q.Get("offset")
	if !util.String.IsEmpty(qOffset) {
		queryResult.Query.Offset, _ = strconv.Atoi(qOffset)
	}
	qLimit := q.Get("limit")
	if !util.String.IsEmpty(qLimit) {
		queryResult.Query.Limit, _ = strconv.Atoi(qLimit)
	}
	keywordMap := h.KeywordMap()
	keyStringWhere := ""
	for key, x := range q {
		if key == "limit" || key == "offset" {
			continue
		}
		if len(x) > 1 {
			keyword, ok := keywordMap[key]
			if ok {
				keyword.Value = q.Get(key)
			}
			delete(keywordMap, key)
			for z := range x {
				if z > 0 {
					keyStringWhere += " OR "
				}
				switch keyword.Type {
				case StringContains:
					keyStringWhere += key + " ILIKE '%" + x[z] + "%'"
				case IntContains:
					keyStringWhere += "CAST(" + key + " AS TEXT) ILIKE '%" + x[z] + "%'"
				}
				modelKeyword := ResultQueryKeyword{Field: key, Value: x[z]}
				queryResult.Query.Keyword = append(queryResult.Query.Keyword, modelKeyword)
			}
		} else {
			keyword, ok := keywordMap[key]
			if ok {
				keyword.Value = q.Get(key)
			}
			modelKeyword := ResultQueryKeyword{Field: key, Value: q.Get(key)}
			queryResult.Query.Keyword = append(queryResult.Query.Keyword, modelKeyword)
		}
	}

	chain := h.DB(state)
	qSuggestion := q.Get("suggestion")
	if !util.String.IsEmpty(qSuggestion) {
		suggestionKeywords := h.SuggestionKeywords()
		if suggestionKeywords != nil && len(suggestionKeywords) > 0 {
			suggestionWhere := ""
			for i, keyword := range suggestionKeywords {
				if i > 0 {
					suggestionWhere += " OR "
				}

				switch keyword.Type {
				case StringContains:
					suggestionWhere += keyword.ColumnName + " ILIKE '%" + qSuggestion + "%'"
				case IntContains:
					suggestionWhere += "CAST(" + keyword.ColumnName + " AS TEXT) ILIKE '%" + qSuggestion + "%'"
				}
			}
			chain = chain.Where(suggestionWhere)
		}
	}
	chain = chain.Where(keyStringWhere)
	for key, keyword := range keywordMap {
		if !util.String.IsEmpty(keyword.Value) {
			if util.String.IsEmpty(keyword.ColumnName) {
				keyword.ColumnName = key
			}
			switch keyword.Type {
			case StringContains:
				chain = chain.Where(keyword.ColumnName+" ILIKE ?", "%"+keyword.Value+"%")
			case IntContains:
				chain = chain.Where("CAST("+keyword.ColumnName+" AS TEXT) ILIKE ?", "%"+keyword.Value+"%")
			case Equal:
				chain = chain.Where(keyword.ColumnName+" = ?", keyword.Value)
			case MoreThanOrEqual:
				chain = chain.Where(keyword.ColumnName+" >= ?", keyword.Value)
			case LessThanOrEqual:
				chain = chain.Where(keyword.ColumnName+" <= ?", keyword.Value)
			}
		}
	}

	chain.Count(&queryResult.Query.Total)

	if queryResult.Query.Limit > 0 {
		chain = chain.Limit(queryResult.Query.Limit)
	}
	chain = chain.Offset(queryResult.Query.Offset)
	data := h.Find(chain)
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		queryResult.Data = data
	}
	respondJSON(w, http.StatusOK, queryResult)
}

// GetByID is get a data by id
func GetByID(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	if state != "mock" {
		state = "db"
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := h.New()
	if err := h.DB(state).First(result, id).Error(); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
	} else {
		respondJSON(w, http.StatusOK, result)
	}
}

// Create new data to database
func Create(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	var path = "history"
	if state == "mock" {
		path += "/mock"
	} else {
		state = "db"
	}

	result := h.New()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	token := r.Header.Get("Authorization")
	if token != "" {
		userID, _, err := util.GetUserAndFirmIDFromJWT(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		h.CreateUserId(userID, result)
	}

	if err := h.Validate(state, result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	tx := h.DB(state).Begin()
	if err := tx.Save(result).Error(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.AfterCreate(state, result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	rowID := h.GetID(result)
	changeLog := &model.ChangeLog{TableName: h.ModelName(), RowID: &rowID}
	if err := tx.Save(changeLog).Error(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()

	respondJSON(w, http.StatusCreated, result)
}

// Update data in database with id as key
func Update(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	var path = "history"
	if state == "mock" {
		path += "/mock"
	} else {
		state = "db"
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := h.New()
	if err := h.DB(state).First(result, id).Error(); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	x := r.Body
	// decoder := json.NewDecoder(r.Body)
	// if err := decoder.Decode(&result); err != nil {
	// 	respondError(w, http.StatusBadRequest, err.Error())
	// 	return
	// }
	body, _ := ioutil.ReadAll(x)
	err = json.Unmarshal(body, &result)

	//khusus firm
	rUp := h.New()
	err = json.Unmarshal(body, &rUp)

	defer r.Body.Close()

	token := r.Header.Get("Authorization")
	if token != "" {
		userID, _, err := util.GetUserAndFirmIDFromJWT(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		h.UpdateUserId(userID, result)
	}

	if err := h.Validate(state, result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.MappingFirmUpdate(result, rUp)

	tx := h.DB(state).Begin()
	if err := tx.Save(result).Error(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rowID := uint(id)
	changeLog := &model.ChangeLog{TableName: h.ModelName(), RowID: &rowID}
	if err := tx.Save(changeLog).Error(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()

	respondJSON(w, http.StatusOK, result)
}

// Delete data in database with id as key
func Delete(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	var path = "history"
	if state == "mock" {
		path += "/mock"
	} else {
		state = "db"
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result := h.New()
	if err := h.DB(state).First(result, id).Error(); err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	token := r.Header.Get("Authorization")
	if token != "" {
		userID, _, err := util.GetUserAndFirmIDFromJWT(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}

		h.DeleteUserId(userID, result)
	}

	// check for existing relations
	if err := h.BeforeDelete(state, result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	tx := h.DB(state).Begin()
	if err := tx.Delete(result).Error(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rowID := uint(id)
	changeLog := &model.ChangeLog{TableName: h.ModelName(), RowID: &rowID}
	if err := tx.Save(changeLog).Error(); err != nil {
		tx.Rollback()
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()

	respondJSON(w, http.StatusNoContent, nil)
}

func combineErrors(errs []error) string {
	errString := []string{}
	for _, err := range errs {
		errString = append(errString, err.Error())
	}
	return strings.Join(errString, ", ")
}

type strFieldInfo struct {
	Value       string
	MaxLength   int
	IsMandatory bool
}

func validateString(name, val string, maxLength int, isMandatory bool) error {
	if isMandatory && util.String.IsEmpty(val) {
		msg := fmt.Sprintf("%s Required.", name)
		return errors.New(msg)
	}

	if maxLength > 0 && len(val) > maxLength {
		msg := fmt.Sprintf("Invalid %s Length (Max. %d).", name, maxLength)
		return errors.New(msg)
	}

	return nil
}

func validateStrings(strFieldMap map[string]strFieldInfo) error {
	for name, info := range strFieldMap {
		if err := validateString(name, info.Value, info.MaxLength, info.IsMandatory); err != nil {
			return err
		}
	}

	return nil
}

func validatePointer(name string, p interface{}) error {
	pVal := reflect.ValueOf(p)
	if pVal.Kind() != reflect.Ptr || pVal.IsNil() {
		msg := fmt.Sprintf("%s Required.", name)
		return errors.New(msg)
	}

	return nil
}

func validatePointers(pMap map[string]interface{}) error {
	for name, p := range pMap {
		if err := validatePointer(name, p); err != nil {
			return err
		}
	}

	return nil
}

func LoginUser(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	if state != "mock" {
		state = "db"
	}

	result := &model.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	var tokenFix string
	dbUser := h.DefaultDB(state)
	modelUser := model.User{}
	dbUser.Where("user_name = ?", result.UserName).Last(&modelUser)
	if modelUser.UserName == "" {
		respondError(w, http.StatusBadRequest, "User Not Exist")
	} else {
		if result.Password != modelUser.Password {
			respondError(w, http.StatusBadRequest, "Wrong Password")
		} else {
			sign := jwt.New(jwt.GetSigningMethod("HS256"))
			token, _ := sign.SignedString([]byte("secret"))
			tokenFix = token
			dbHis := h.DefaultDB(state).Model(model.History{})
			modelsHis := model.History{}

			modelsHis.Note = "Login User " + result.UserName
			modelsHis.Module = "User Login"
			dbHis.Save(&modelsHis)
			respondJSON(w, http.StatusOK, tokenFix)
		}
	}
}

func LogOutUser(h RESTHandler, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	state := q.Get("state")
	if state != "mock" {
		state = "db"
	}

	result := &model.User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	dbUser := h.DefaultDB(state)
	modelUser := model.User{}
	dbUser.Where("user_name = ?", result.UserName).Last(&modelUser)
	if modelUser.UserName == "" {
		respondError(w, http.StatusBadRequest, "User Not Exist")
	} else {
		if result.Password != modelUser.Password {
			respondError(w, http.StatusBadRequest, "Wrong Password")
		} else {
			dbHis := h.DefaultDB(state).Model(model.History{})
			modelsHis := model.History{}

			modelsHis.Note = "Logout User " + result.UserName
			modelsHis.Module = "User Logout"
			dbHis.Save(&modelsHis)
			respondJSON(w, http.StatusOK, nil)
		}
	}
}
