package array

//判断数组中是否存在值
func InArray(need string, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}
