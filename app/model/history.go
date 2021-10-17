package model

type History struct {
	Default
	Note   string `json:"Note"`
	Module string `json:"Module"`
}
