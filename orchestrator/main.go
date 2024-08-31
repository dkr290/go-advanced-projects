package main

import (
	"fmt"
	"time"

	"github.com/dkr290/go-advanced-projects/orchestrator/manager"
	"github.com/dkr290/go-advanced-projects/orchestrator/node"
	"github.com/dkr290/go-advanced-projects/orchestrator/task"
	"github.com/dkr290/go-advanced-projects/orchestrator/worker"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

func main() {
	t := task.Task{
		ID:        uuid.New(),
		Name:      "Task-1",
		State:     task.Pending,
		Image:     "image-1",
		CPU:       1.0,
		Memory:    1024,
		Disk:      1,
		StartTime: time.Now(),
	}

	te := task.TaskEvent{

		ID:        uuid.New(),
		State:     task.Running,
		TimeStamp: time.Now(),
		Task:      t,
	}

	fmt.Printf("task: %v\n", t)
	fmt.Printf("task event: %v\n", te)

	w := worker.Worker{
		Name:  "worker-1",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	fmt.Printf("worker: %v\n", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StartTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}
	fmt.Printf("manager: %v\n", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWorkd()

	n := node.Node{
		Name:  "Node-1",
		Ip:    "192.168.1.1",
		Cores: 4,
		Disk:  25,
		Role:  "worker",
	}

	fmt.Printf("node: %v\n", n)

}
