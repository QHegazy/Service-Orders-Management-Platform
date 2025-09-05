package utils

func SuccessResponse(message string, data interface{}) map[string]interface{} {
	apiResponse := map[string]interface{}{
		"message": message,
		"data":    data,
	}
	return apiResponse
}

func ErrorResponse(message string, errorstring interface{}) map[string]interface{} {
	apiResponse := map[string]interface{}{
		"message": message,
		"error":   errorstring,
	}
	return apiResponse
}
