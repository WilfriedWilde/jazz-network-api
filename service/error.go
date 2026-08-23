package service

type ServiceError struct {
	Status  int
	Code    string
	Message string
}

func (e *ServiceError) GetStatus() int {
	return e.Status
}

func (e *ServiceError) GetCode() string {
	return e.Code
}

func (e *ServiceError) Error() string {
	return e.Message
}
