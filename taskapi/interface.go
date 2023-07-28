package taskapi

type TaskApiInterface interface {
	SendNodeTask(task *NodeTask) error
}
