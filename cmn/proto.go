package cmn

import (
	"encoding/json"
	"errors"
)

type ReplyProto struct {
	//Status, 0: success, others: fault
	Status int `json:"status"`

	//Msg, Action result describe by literal
	Msg string `json:"msg,omitempty"`

	//Data, operand
	Data json.RawMessage `json:"data,omitempty"`

	// RowCount, just row count
	RowCount int64 `json:"rowCount,omitempty"`

	//API, call target
	API string `json:"API,omitempty"`

	//Method, using http method
	Method string `json:"method,omitempty"`

	//SN, call order
	SN int `json:"SN,omitempty"`
}

func NewErrorReply(err error, API, Method string) ReplyProto {
	var appErr *AppError
	ok := errors.As(err, &appErr)
	if !ok {
		return ReplyProto{
			Status:   CommonError,
			Msg:      err.Error(),
			API:      API,
			Method:   Method,
			RowCount: 0,
			Data:     nil,
			SN:       0,
		}
	} else {
		return ReplyProto{
			Status:   appErr.StatusCode,
			Msg:      appErr.Message,
			API:      API,
			Method:   Method,
			RowCount: 0,
			Data:     nil,
			SN:       0,
		}
	}
}

type ReqProto struct {
	Action string `json:"action,omitempty"`

	Sets    []string            `json:"sets,omitempty"`
	OrderBy []map[string]string `json:"orderBy,omitempty"`

	//***页码从第零页开始***
	Page     int64 `json:"page,omitempty"`
	PageSize int64 `json:"pageSize,omitempty"`

	Data   json.RawMessage `json:"data,omitempty"`
	Filter interface{}     `json:"filter,omitempty"`

	AuthFilter interface{} `json:"authFilter,omitempty"`
}
