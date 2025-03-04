package task_heap

func (th TaskHeap) Less(i, j int) bool { return th[i].ExecutionTime.Before(th[j].ExecutionTime) }
