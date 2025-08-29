package entity

import "time"

// Event 事件实体（可作为模板）
type Event struct {
	ID          uint   // 主键ID
	IsTemplate  bool   // 是否为模板
	Title       string // 标题
	Description string // 描述
	Location    string // 地点
	Priority    int    // 优先级（1-5）
	Category    string // 分类
}

// Task 任务实体,描述了任务本身
type Task struct {
	ID          uint   // 主键ID
	Description string // 任务描述

	// 可用于该task的时间段
	AllowedTime TimeSpan

	// 计划时间段
	PlannedDuration TimeSpan

	// 状态信息
	Status    Status    // 状态（待定/进行中/完成）
	CreatedAt time.Time // 创建时间

	// 关联关系
	EventID      uint   // 事件ID(描述了该任务的内容)
	ParentTaskID *uint  // 父任务ID（可空）
	PreTaskIDs   []uint // 前置任务ID（可空）
}

// Todo 待办事项,描述了我如何做任务
type Todo struct {
	ID uint // 主键ID

	// 时间段
	PlannedTime TimeSpan
	ActualTime  TimeSpan

	// 状态信息
	Status Status // 状态

	CompletedTime *time.Time

	// 关联关系
	EventID *uint //默认为空(使用任务的事件,除了当Todo需要与Task不同,如临时不同的地点等)
	TaskID  uint  // 所属任务ID
}

type Status string

const (
	StatusPending    Status = "pending"   // 待办
	StatusInProgress Status = "doing"     // 进行中
	StatusCompleted  Status = "done"      // 已完成
	StatusCancelled  Status = "cancelled" // 已取消
)

type TimeSpan struct {
	Start *time.Time
	End   *time.Time
}

func (ts TimeSpan) Duration() time.Duration {
	if ts.Start == nil || ts.End == nil {
		return 0
	}
	return ts.End.Sub(*ts.Start)
}
