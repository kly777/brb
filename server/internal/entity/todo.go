package entity

import "time"

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
	StartTime time.Time `json:"startTime"` // 开始时间
	EndTime   time.Time `json:"endTime"`   // 结束时间

	// 预估时间
	EstimateTime time.Duration `json:"estimateTime"`
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

	CompletedTime time.Time `json:"completedTime"`

	// 计划执行Task的时间
	StartTime time.Time `json:"startTime"` // 开始时间
	EndTime   time.Time `json:"endTime"`
}
