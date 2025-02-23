package service

import (
	"fmt"
	"pkg/timeout"
	"slices"
	"sync"
	"time"

	"github.com/roadtoseniors/apicalc/internal/orchestrator/config"
	"github.com/roadtoseniors/apicalc/internal/task"
)

type CalcService struct {
	locker        sync.RWMutex
	exprTable     map[string]*Expression
	taskID        int64
	tasks         []*task.Task
	taskTable     map[int64]ExprElement
	timeTable     map[string]time.Duration
	timeoutsTable map[int64]*timeout.Timeout
}

func NewCalcService(cfg config.TimeConfig) *CalcService {
	cs := CalcService{
		exprTable:     make(map[string]*Expression),
		taskTable:     make(map[int64]ExprElement),
		timeTable:     make(map[string]time.Duration),
		timeoutsTable: make(map[int64]*timeout.Timeout),
	}
	cs.timeTable["+"] = cfg.Add
	cs.timeTable["-"] = cfg.Sub
	cs.timeTable["*"] = cfg.Mul
	cs.timeTable["/"] = cfg.Div

	return &cs
}

func (cs *CalcService) AddExpression(id, expr string) error {
	if len(id) == 0 {
		return fmt.Errorf("empty ID")
	}
	if len(expr) == 0 {
		return fmt.Errorf("empty expression")
	}

	cs.locker.Lock()
	defer cs.locker.Unlock()

	if _, found := cs.exprTable[id]; found {
		return fmt.Errorf("not a unique ID: %q", id)
	}

	expression, err := NewExpression(id, expr)
	cs.exprTable[id] = expression
	//извлекаем задачи если выражение в процессе вычисления
	if err == nil && expression.Status == StatusInProcess {
		cs.extractTasksFromExpression(expression)
	}

	return nil
}

// возвращаем список всех выражений
func (cs *CalcService) ListAll() ExpressionList {
	cs.locker.RLock()
	defer cs.locker.RUnlock()

	lst := ExpressionList{}
	for _, expr := range cs.exprTable {
		lst.Exprs = append(lst.Exprs, *expr)
	}

	//сортируем по айди
	slices.SortFunc(lst.Exprs, func(a, b Expression) int {
		if a.ID > b.ID {
			return 1
		} else if a.ID < b.ID {
			return -1
		}
		return 0
	})

	return lst
}

// возвращаю выражение по айди
func (cs *CalcService) FindById(id string) (*ExpressionUnit, error) {
	cs.locker.RLock()
	defer cs.locker.RUnlock()

	// ищу в таблице
	expr, found := cs.exprTable[id]
	if !found {
		return nil, fmt.Errorf("id %q not found", id)
	}

	return &ExpressionUnit{Expr: *expr}, nil
}

// возврат для выполнения задачи
func (cs *CalcService) GetTask() *task.Task {
	cs.locker.Lock()
	defer cs.locker.Unlock()

	if len(cs.tasks) == 0 {
		return nil
	}

	newtask := cs.tasks[0]
	cs.tasks = cs.tasks[1:]

	cs.timeoutsTable[newtask.ID] = timeout.NewTimeout(
		5*time.Second + newtask.OperationTime,
	)

	// горутина обрабатывает таймаут
	go func(task task.Task) {
		cs.locker.Lock()
		timeout, found := cs.timeoutsTable[task.ID]
		cs.locker.Unlock()
		if !found {
			return
		}

		select {
		case <-timeout.Timer.C:
			cs.locker.Lock()
			cs.tasks = append(cs.tasks, &task)
			cs.locker.Unlock()
		case <-timeout.Ctx.Done():
			return
		}
	}(*newtask)

	return newtask
}

// сохраняю результат выполнения задачи
func (cs *CalcService) PutResult(id int64, value float64) error {
	cs.locker.Lock()
	defer cs.locker.Unlock()

	timeout, found := cs.timeoutsTable[id]
	if found {
		timeout.Cancel()
	}

	_, found = cs.taskTable[id]
	if !found {
		return fmt.Errorf("Task id %d not found", id)
	}

	el := cs.taskTable[id].Ptr
	exprID := cs.taskTable[id].ID
	delete(cs.taskTable, id)

	expr, found := cs.exprTable[exprID]
	if !found {
		return fmt.Errorf("Expression for task %d not found", id)
	}

	if expr.Len() == 1 {
		expr.Result = fmt.Sprintf("%g", value)
		expr.Status = StatusDone
		expr.Remove(el)
	} else {
		numToken := NumToken{value}
		expr.InsertBefore(numToken, el)
		expr.Remove(el)
		cs.extractTasksFromExpression(expr)
	}

	return nil
}

// извлекаю все задачи для выполнения
func (cs *CalcService) extractTasksFromExpression(expr *Expression) int {
	var taskCount int
	el := expr.Front()

	for el != nil {
		el1 := el
		if el1.Value.(Token).Type() != TokenTypeNumber {
			el = el.Next()
			continue
		}

		el2 := el1.Next()
		if el2 == nil || el2.Value.(Token).Type() != TokenTypeNumber {
			el = el.Next()
			continue
		}

		op := el2.Next()
		if op == nil || op.Value.(Token).Type() != TokenTypeOperation {
			el = el.Next()
			continue
		}

		// создаём новую задачу
		task := new(task.Task)
		task.ID = cs.taskID
		cs.taskID++
		taskToken := TaskToken{ID: task.ID}
		taskElement := expr.InsertBefore(&taskToken, el)
		cs.taskTable[task.ID] = ExprElement{expr.ID, taskElement}

		task.Arg1 = fmt.Sprintf("%f", el1.Value.(NumToken).Value)
		task.Arg2 = fmt.Sprintf("%f", el2.Value.(NumToken).Value)
		task.Operation = op.Value.(OpToken).Value
		task.OperationTime = cs.timeTable[task.Operation]

		taskCount++
		cs.tasks = append(cs.tasks, task)
		el = op.Next()

		expr.Remove(el1)
		expr.Remove(el2)
		expr.Remove(op)
	}

	return taskCount
}
