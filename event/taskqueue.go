package event

import (
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"context"
	"encoding/json"
	"fmt"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"log"
	"net/http"
)

type TaskQueue interface {
	CreateTask(task Task) error
}

type Task struct {
	Queue   string
	Path    string
	Object  interface{}
	Payload []byte
}

type taskQueue struct {
	cli        *cloudtasks.Client
	projectID  string
	locationID string
}

func (t *taskQueue) CreateTask(task Task) error {
	ctx := context.Background()

	//headers := map[string]string{
	//	"Content-Type":  "application/json",
	//	"Authorization": c.authToken,
	//}

	body, err := json.Marshal(task.Object)
	if err != nil {
		return err
	}

	log.Printf("task.Path: %v", string(task.Path))
	log.Printf("ReadAll: %v", string(body))

	aeReq := &taskspb.AppEngineHttpRequest{
		//AppEngineRouting: &taskspb.AppEngineRouting{
		//	Service: "serviceID",
		//},
		HttpMethod:  taskspb.HttpMethod_POST,
		RelativeUri: task.Path,
		//Headers:     headers,
		Body: body,
	}

	req := &taskspb.CreateTaskRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", t.projectID, t.locationID, task.Queue),
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: aeReq,
			},
		},
	}
	_, err = t.cli.CreateTask(ctx, req)

	return err
}

func ParseTask(r *http.Request, o interface{}) error {
	return json.NewDecoder(r.Body).Decode(o)
}

func NewTasksClient(projectID, locationID string) TaskQueue {
	ctx := context.Background()

	//gOpt := option.WithGRPCDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
	//	Time:                1 * time.Second,
	//	Timeout:             5 * time.Second,
	//	PermitWithoutStream: true,
	//}))

	cli, err := cloudtasks.NewClient(ctx)

	if err != nil {
		panic(err)
	}

	return &taskQueue{
		cli:        cli,
		projectID:  projectID,
		locationID: locationID,
	}
}
