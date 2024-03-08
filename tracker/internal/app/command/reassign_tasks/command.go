package reassign_tasks

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/segmentio/kafka-go"

	task_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/task"
	user_model "github.com/ilyatos/aac-task-tracker/tracker/internal/app/model/user"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/repository"
)

type Command struct {
	userRepository *repository.UserRepository
	taskRepository *repository.TaskRepository
}

func New(userRepository *repository.UserRepository, taskRepository *repository.TaskRepository) *Command {
	return &Command{userRepository: userRepository, taskRepository: taskRepository}
}

func (c *Command) Handle(ctx context.Context) error {
	users, err := c.userRepository.GetAllByRole(ctx, user_model.Employee)
	if err != nil {
		return fmt.Errorf("get users error: %w", err)
	}
	if users == nil || len(users) == 0 {
		return fmt.Errorf("users not found")
	}

	tasks, err := c.taskRepository.GetAllByStatus(ctx, task_model.New)
	if err != nil {
		return fmt.Errorf("get tasks error: %w", err)
	}
	if tasks == nil || len(tasks) == 0 {
		return fmt.Errorf("tasks not found")
	}

	taskAssignedEvents, taskUpdatedEvents, err := c.updateTasks(ctx, tasks, users)

	// TODO: outbox pattern
	err = produceTaskAssignedEvents(ctx, taskAssignedEvents)
	if err != nil {
		fmt.Printf("produce task created events error: %v", err)
	}

	err = produceTaskUpdatedEvents(ctx, taskUpdatedEvents)
	if err != nil {
		fmt.Printf("produce task updated events error: %v", err)
	}

	return nil
}

// TODO: transaction
func (c *Command) updateTasks(ctx context.Context, tasks []repository.Task, users []repository.UserPublicID) ([]kafka.Message, []kafka.Message, error) {
	var (
		taskAssignedEvents []kafka.Message
		taskUpdatedEvents  []kafka.Message
	)
	for _, task := range tasks {
		task.UserPublicID = users[rand.Intn(len(users))].PublicID

		err := c.taskRepository.Update(ctx, task)
		if err != nil {
			return nil, nil, fmt.Errorf("update task error: %w", err)
		}

		taskAssignedEvent, err := makeTaskAssignedMessage(
			ctx,
			taskAssignedEvent{
				PublicID:     task.PublicID,
				UserPublicID: task.UserPublicID,
			},
		)
		taskAssignedEvents = append(taskAssignedEvents, taskAssignedEvent)
		if err != nil {
			return nil, nil, fmt.Errorf("make task assigned message error: %w", err)
		}

		taskUpdatedEvent, err := makeTaskUpdatedMessage(
			ctx,
			taskUpdatedEvent{
				PublicID:     task.PublicID,
				UserPublicID: task.UserPublicID,
				Description:  task.Description,
				Status:       task.Status,
			},
		)
		taskUpdatedEvents = append(taskUpdatedEvents, taskUpdatedEvent)
	}

	return taskAssignedEvents, taskUpdatedEvents, nil
}
