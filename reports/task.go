package reports

import "github.com/mrcook/time_warrior/timeslip"

type task struct {
	name       string
	project    string
	started    int
	finished   int
	timeWorked int
	feeRate    float64 // hourly rate in dollars
}

// Creates a new task from a timeslip JSON string.
func newTask(jsonData []byte) (*task, error) {
	slip := &timeslip.Slip{}
	if err := timeslip.Unmarshal(jsonData, slip); err != nil {
		return nil, err
	}

	var name string
	if slip.Task == "" {
		name = "."
	} else {
		name = slip.Task
	}

	// Default fee rate of 66$ per hour
	t := &task{
		name:       name,
		project:    slip.Project,
		started:    slip.Started,
		finished:   slip.Finished,
		timeWorked: slip.Worked,
		feeRate:    66.0,
	}

	return t, nil
}

// Name returns the task name
func (t *task) Name() string {
	return t.name
}

// Started returns the task start time
func (t *task) Started() int {
	return t.started
}

// Finished returns the task finish time
func (t *task) Finished() int {
	return t.finished
}

// TimeSpent returns the time worked in hours
func (t *task) TimeSpent() float64 {
	return float64(t.timeWorked) / 3600.0
}

// Fee returns the fee for this task in dollars
func (t *task) Fee() float64 {
	return t.TimeSpent() * t.feeRate
}
