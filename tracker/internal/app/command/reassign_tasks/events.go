package reassign_tasks

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/meta"
	task_event "github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_assigned"
	"github.com/ilyatos/aac-task-tracker/schema_registry/pkg/task/task_updated"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/broker"
)

func makeTaskAssignedMessage(task repository.Task) (kafka.Message, error) {
	taskAssigned, err := proto.Marshal(&task_assigned.TaskAssigned{
		Header: &meta.Header{
			Producer: "tracker.reassign_tasks",
		},
		Payload: &task_assigned.TaskAssigned_V1{
			V1: &task_assigned.V1{
				PublicId:     task.PublicID.String(),
				UserPublicId: task.UserPublicID.String(),
			},
		},
	})
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{
		Key:   []byte("TaskAssigned"),
		Value: taskAssigned,
	}, nil
}

func produceTaskAssignedEvents(ctx context.Context, msgs []kafka.Message) error {
	return broker.Produce(ctx, "tasks", msgs...)
}

func makeTaskUpdatedMessage(task repository.Task) (kafka.Message, error) {
	taskUpdated, err := proto.Marshal(&task_updated.TaskUpdated{
		Header: &meta.Header{Producer: "tracker.reassign_tasks"},
		Payload: &task_updated.TaskUpdated_V1{
			V1: &task_updated.V1{
				PublicId:     task.PublicID.String(),
				UserPublicId: task.UserPublicID.String(),
				Description:  task.Description,
				Status:       task_event.Status(task_event.Status_value[strings.ToUpper(string(task.Status))]),
			},
		},
	})
	if err != nil {
		return kafka.Message{}, err
	}

	return kafka.Message{
		Key:   []byte("TaskUpdated"),
		Value: taskUpdated,
	}, nil
}

func produceTaskUpdatedEvents(ctx context.Context, msgs []kafka.Message) error {
	return broker.Produce(ctx, "tasks-stream", msgs...)
}
