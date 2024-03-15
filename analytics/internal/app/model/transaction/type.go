package transaction

type Type string

const (
	TaskAssigned  Type = "task_assigned"
	TaskCompleted Type = "task_completed"
	Payout        Type = "payout"
)
