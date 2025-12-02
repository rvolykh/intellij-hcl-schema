package terraform

// typeValToStr converts GRPC response Type value from bytes to string with removing open/close quotes
func typeValToStr(v []byte) string {
	if len(v) < 2 {
		return ""
	}
	return string(v[1 : len(v)-1])
}

// deprecatedValToStr converts GRPC response Deprecated value from bool to string
func deprecatedValToStr(v bool) string {
	if !v {
		return ""
	}
	return "Deprecated! Please refer to documentation for more details"
}
