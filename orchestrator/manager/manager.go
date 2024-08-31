package manager

import (
	"fmt"

	"github.com/dkr290/go-advanced-projects/orchestrator/task"
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	Pending        queue.Queue
	TaskDb         map[string][]*task.Task
	EventDb        map[string][]*task.TaskEvent
	Workers        []string
	WorkersTaskMap map[string][]uuid.UUID //the jobs that are assigned to each worker
	TaskWorkerMap  map[uuid.UUID]string   //TaskWorkerMap, which is a map of task UUIDs to strings,where the string is the name of the worker
}

//SelectWorker() - This method will be responsible for looking at the requirements
//specified in a Task and evaluating the resources available in the pool of workers to see
//which worker is best suited to run the task.
//the Manager must keep track of tasks, their states, and the machine on which they run

func (m *Manager) SelectWorker() {
	fmt.Println("I will select the appropriate worker")
}

// UpdateTasks() - this method triggewr call to workers CollectStats()

func (m *Manager) UpdateTasks() {
	fmt.Println("I will update tasks")
}

func (*Manager) SendWorkd() {
	fmt.Println("I will send work to workers")
}
