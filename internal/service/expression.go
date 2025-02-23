package service

import (
	"container/list"
	"strconv"
	"strings"

	"github.com/roadtoseniors/apicalc/pkg/rpn"
)

const (
	StatusError     = "Error"
	StatusDone      = "Done"
	StatusInProcess = "In process"
)

const (
	TokenTypeNumber = iota
	TokenTypeOperation
	TokenTypeTask
)

type Token interface {
	Type() int
}

type NumToken struct {
	Value float64
}

func (num NumToken) Type() int {
	return TokenTypeNumber
}

type OpToken struct {
	Value string
}

func (num OpToken) Type() int {
	return TokenTypeOperation
}

type TaskToken struct {
	ID int64
}

func (num TaskToken) Type() int {
	return TokenTypeTask
}

type Expression struct {
	*list.List        // Список токенов выражения
	ID         string `json:"id"`
	Status     string `json:"status"`
	Result     string `json:"result"`
	Source     string `json:"source"` // исходник
}

type ExpressionUnit struct {
	Expr Expression `json:"expression"`
}

type ExpressionList struct {
	Exprs []Expression `json:"expressions"`
}

func NewExpression(id, expr string) (*Expression, error) {
	// преобразуем выражение в обратную польскую запись
	rpn, err := rpn.NewRPN(expr)
	if err != nil {
		// если произошла ошибка
		expression := Expression{
			List:   list.New(),
			ID:     id,
			Status: StatusError,
			Result: "",
			Source: expr,
		}
		return &expression, err
	}

	// Если выражение состоит из одного числа, создаём выражение со статусом "Done".
	if len(rpn) == 1 {
		expression := Expression{
			List:   list.New(),
			ID:     id,
			Status: StatusDone,
			Result: rpn[0],
			Source: expr,
		}
		return &expression, nil
	}

	// Создаём выражение со статусом "In process".
	expression := Expression{
		List:   list.New(),
		ID:     id,
		Status: StatusInProcess,
		Result: "",
		Source: expr,
	}

	// Преобразуем RPN в список токенов.
	for _, val := range rpn {
		if strings.Contains("-+*/", val) {
			// Если это операция, добавляем OpToken.
			expression.PushBack(OpToken{val})
		} else {
			// Если это число, преобразуем его в float64 и добавляем NumToken.
			num, err := strconv.ParseFloat(val, 10)
			if err != nil {
				return nil, err
			}
			expression.PushBack(NumToken{num})
		}
	}

	return &expression, nil
}

type ExprElement struct {
	ID  string
	Ptr *list.Element // Указатель на элемент списка
}
