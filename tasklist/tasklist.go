package tasklist

/*
	This package will house the tasklist structure and provide a number of
	methods for it.
	- Tasklist functions as a FIFO tasks queue
*/

// channel for tasks
type Tasklist struct {
	tasks chan Task
}

// filepaths
type Task string

// create a new Task
func CreateTask(path string) Task {
	return Task(path)
}

// add a task to the TL
func (t *Tasklist) Enqueue(task Task) {
	t.tasks <- task
}

// grab the next task in line from the TL
func (t *Tasklist) Dequeue() Task{
	next := <-t.tasks
	return next
}

// create a buffered TL
func CreateTLChannel(bufSize int) Tasklist {
	return Tasklist{make(chan Task, bufSize)}
}