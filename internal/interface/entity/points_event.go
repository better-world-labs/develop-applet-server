package entity

type PointsChangeEvent struct {
	OperationId string `json:"operationId"`
	Points
}
