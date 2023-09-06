package models

type CommonErrorResponse struct {
	Error AclError
}

type CommandResponse struct {
	IsSuccessd bool
	Message    string
}

type AclError struct {
	Code             string
	Message          string
	Details          string
	RequestStack     string
	ValidationErrors []ValidationError
}

type ValidationError struct {
	Message string
	Members []string
}
