package dispatch

import (
	"fmt"
	"sync"
	"time"
)

var tag = "JobDispatcher"

type Job interface {
	Start(updates <-chan interface{}, done chan<- bool)
	ShutDown()
}

type JobState struct {
	startTime   time.Time
	timeOut     time.Duration
	job         Job
	updatesChan chan interface{}
	doneChan    chan bool
}

type JobDispatcher struct {
	channels   map[int64]*JobState
	mu         sync.Mutex
	JobCreator func(update interface{}) (job Job, timout time.Duration)
}

func (d *JobDispatcher) Exists(id int64) bool {
	defer d.mu.Unlock()
	d.mu.Lock()
	if d.channels == nil {
		d.channels = make(map[int64]*JobState)
		return false
	}
	return d.channels[id] != nil
}

func (d *JobDispatcher) addJob(id int64, job Job, timemOut time.Duration) bool {
	defer d.mu.Unlock()
	d.mu.Lock()

	if d.channels == nil {
		d.channels = make(map[int64]*JobState)
	} else if d.channels[id] != nil {
		return false
	}
	d.channels[id] = &JobState{
		startTime:   time.Now(),
		timeOut:     timemOut,
		job:         job,
		updatesChan: make(chan interface{}),
		doneChan:    make(chan bool),
	}
	return true
}

func (d *JobDispatcher) killJob(id int64) bool {
	if jS := d.getJobState(id); jS == nil {
		fmt.Println(tag, "Job", id, "Not found, aborting kill")
		return false
	} else {
		fmt.Println(tag, "Shutting down Job", id)
		close(jS.updatesChan)
		jS.job.ShutDown()
		jS.job = nil
		delete(d.channels, id)
		d.channels[id] = nil
		return true
	}
}

func (d *JobDispatcher) getJobState(id int64) *JobState {
	defer d.mu.Unlock()
	d.mu.Lock()
	if jS := d.channels[id]; jS == nil {
		return nil
	} else {
		return jS
	}
}

func (d *JobDispatcher) startJob(id int64) {
	fmt.Println(tag, "Starting Job", id)
	jS := d.getJobState(id)
	if jS == nil {
		panic(fmt.Sprint(tag, "Trying to start a non existing Job", id))
	}
	fmt.Println(tag, "Starting Job", id)
	go jS.job.Start(jS.updatesChan, jS.doneChan)
	grace := time.After(jS.timeOut)
GRACE:
	select {
	case <-grace:
		fmt.Println(tag, "Job", id, "Exceeded timeout", jS.timeOut, "Killing..")
		d.killJob(id)
		fmt.Println(tag, "Killed", id)
		return
	case isDone := <-jS.doneChan:
		if !isDone {
			fmt.Println(tag, "Job needs more time")
			grace = time.After(jS.timeOut)
			goto GRACE
		} else {
			fmt.Println(tag, "Job is done")
			d.killJob(id)
			fmt.Println(tag, "Killed", id)
			return
		}
	}

}

func (d *JobDispatcher) sendUpdate(id int64, update interface{}) {
	if jS := d.getJobState(id); jS != nil {
		fmt.Println(tag, "Sending Update", id)
		jS.updatesChan <- update
	}
}
func (d *JobDispatcher) Dispatch(id int64, update interface{}) {
	if d.Exists(id) {
		fmt.Println(tag, "Found job for", id, " Sending update")
		go d.sendUpdate(id, update)
	} else {
		fmt.Println(tag, "No job found for", id, ". Creating new one")
		j, t := d.JobCreator(update)
		if j == nil {
			fmt.Println(tag, "Job Creator returned nil, abortsing dispatch")
			return
		}
		d.addJob(id, j, t)
		go d.startJob(id)
		fmt.Println(tag, "Sending first update", id)
		go d.sendUpdate(id, update)
	}
}
