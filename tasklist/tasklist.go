package tasklist

/*
	This package will house the tasklist structure and provide a number of
	methods for it.
	- Tasklist functions as a FIFO tasks queue
*/

// channel for tasks
type Tasklist struct {
	Tasks chan Task
}

// filepaths
type Task string

// create a new Task
func CreateTask(path string) Task {
	return Task(path)
}

// add a task to the TL
func (t *Tasklist) Enqueue(task Task) {
	t.Tasks <- task
}

// create a buffered TL
func CreateTLChannel(bufSize int) Tasklist {
	return Tasklist{make(chan Task, bufSize)}
}
