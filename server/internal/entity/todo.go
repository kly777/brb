package entity

import (
	"encoding/json"
	"strings"
	"time"
)

// Event 事件结构体,一个Event可以对应多个Task
type Event struct {
	ID          string `json:"id"`          // 事件唯一标识
	Title       string `json:"title"`       // 标题
	Description string `json:"description"` // 描述

	Recurrence string `json:"recurrence"` // 重复规则（如：daily/weekly/monthly）
}

// Task 任务结构体,一个Task可以对应一个Event
// 任务是对事件的具体执行计划,startTime和endTime是任务的执行期限
// 与Todo不同,在这个期限内完成即可,而不是时刻在做
type Task struct {
	ID         string `json:"id"`        // 任务唯一标识
	EventID    string `json:"eventId"`   // 关联的事件ID
	MainTaskID *string `json:"mainTaskId"` // 父任务ID，此任务是子任务

	Description string `json:"description"` // 任务描述

	// 在开始结束时间中完成即可
	StartTime *time.Time `json:"startTime"` // 开始时间
	EndTime   *time.Time `json:"endTime"`   // 结束时间

	// 预估时间
	EstimateTime *time.Duration `json:"estimateTime"`
}

// UnmarshalJSON 自定义JSON解码，处理空字符串时间字段
func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		StartTime    string `json:"startTime"`
		EndTime      string `json:"endTime"`
		EstimateTime string `json:"estimateTime"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 处理StartTime
	if strings.TrimSpace(aux.StartTime) == "" {
		t.StartTime = nil
	} else {
		parsedTime, err := time.Parse(time.RFC3339, aux.StartTime)
		if err != nil {
			return err
		}
		t.StartTime = &parsedTime
	}

	// 处理EndTime
	if strings.TrimSpace(aux.EndTime) == "" {
		t.EndTime = nil
	} else {
		parsedTime, err := time.Parse(time.RFC3339, aux.EndTime)
		if err != nil {
			return err
		}
		t.EndTime = &parsedTime
	}

	// 处理EstimateTime
	if strings.TrimSpace(aux.EstimateTime) == "" {
		t.EstimateTime = nil
	} else {
		parsedDuration, err := time.ParseDuration(aux.EstimateTime)
		if err != nil {
			return err
		}
		t.EstimateTime = &parsedDuration
	}

	return nil
}

// ValidationError 表示验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NotFoundError 表示资源未找到错误
type NotFoundError struct {
	ID string `json:"id"`
}

func (e *NotFoundError) Error() string {
	return "resource with ID " + e.ID + " not found"
}

type Status string

const (
	StatusPending    Status = "pending"   // 待办
	StatusInProgress Status = "doing"     // 进行中
	StatusCompleted  Status = "done"      // 已完成
	StatusCancelled  Status = "cancelled" // 已取消
)

// Todo 待办事项结构体,一个Todo对应一个Task
// Todo是对Task的具体执行计划,需要在Task的StartTime和EndTime之间
// 在Todo的StartTime和EndTime之间进行
type Todo struct {
	ID       string `json:"id"`     // 待办事项唯一标识
	TaskID   string `json:"taskId"` // 待办任务ID
	Status   Status `json:"status"`
	Priority int    `json:"priority"`

	CompletedTime *time.Time `json:"completedTime"`

	// 计划执行Task的时间
	StartTime *time.Time `json:"startTime"` // 开始时间
	EndTime   *time.Time `json:"endTime"`
}

// UnmarshalJSON 自定义JSON解码，处理空字符串时间字段
func (t *Todo) UnmarshalJSON(data []byte) error {
	type Alias Todo
	aux := &struct {
		StartTime      string `json:"startTime"`
		EndTime        string `json:"endTime"`
		CompletedTime  string `json:"completedTime"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 处理StartTime
	if strings.TrimSpace(aux.StartTime) == "" {
		t.StartTime = nil
	} else {
		parsedTime, err := time.Parse(time.RFC3339, aux.StartTime)
		if err != nil {
			return err
		}
		t.StartTime = &parsedTime
	}

	// 处理EndTime
	if strings.TrimSpace(aux.EndTime) == "" {
		t.EndTime = nil
	} else {
		parsedTime, err := time.Parse(time.RFC3339, aux.EndTime)
		if err != nil {
			return err
		}
		t.EndTime = &parsedTime
	}

	// 处理CompletedTime
	if strings.TrimSpace(aux.CompletedTime) == "" {
		t.CompletedTime = nil
	} else {
		parsedTime, err := time.Parse(time.RFC3339, aux.CompletedTime)
		if err != nil {
			return err
		}
		t.CompletedTime = &parsedTime
	}

	return nil
}
