package services

func IsValidStatus(status string) bool {
	return status == "Pending" ||
		status == "In Progress" ||
		status == "Completed"
}
