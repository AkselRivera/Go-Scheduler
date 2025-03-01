package task_heap

func (th *TaskHeap) Pop() interface{} {
	old := *th
	n := len(old)
	task := old[n-1]
	old[n-1] = nil  // Evitar memory leak
	task.Index = -1 // Por seguridad
	*th = old[0 : n-1]
	return task
}
