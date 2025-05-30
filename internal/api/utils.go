package api

func checkMethod(allowedMethod, reqMethod string) (string, error) {
	if allowedMethod != reqMethod {
		return "method not allowed", nil
	}
	return "", nil
}
