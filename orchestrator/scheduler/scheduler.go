package scheduler

//determine a set of candidate workers on which a task could run
//Score the candidate workers from best to worst
//Pick the worker with the best score

type Scheduler interface {
	SelectCandiateNodes()
	Score()
	Pick()
}
