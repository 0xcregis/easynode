package service

type TaskApiInterface interface {
	SendNodeTask(task *NodeTask) error
}
