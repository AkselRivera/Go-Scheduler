package task_heap

func (th TaskHeap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
	th[i].Index = i
	th[j].Index = j
}
