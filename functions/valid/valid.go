package valid

import "regexp"

//验证手机号
func IsChineseMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "86" {
		return false
	}
	reg := `^` + preCode + `1((33|49|53|73|77|80|81|89|91|99|30|31|32|45|55|56|66|71|75|76|85|86|34|35|36|37|38|39|47|50|51|52|57|58|59|72|78|82|83|84|87|88|98|62|65|67)\d{8})|((700|701|702|703|705|706|704|707|708|709|349)\d{7})$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证台湾手机号
func IsTwMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "886" {
		return false
	}
	reg := `^` + preCode + `0?9\d{8}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证香港手机号
func IsHkMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "852" {
		return false
	}
	reg := `^` + preCode + `[569]\d{3}\d{4}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证日本手机号
func IsJapaneseMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "81" {
		return false
	}
	reg := `^` + preCode + `0\d{1,4}\d{1,4}\d{4}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证印度手机号
func IsIndiaMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "91" {
		return false
	}
	reg := `^` + preCode + `0?[789]\d{9}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证印度尼西亚手机号
func IsIndonesiaMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "62" {
		return false
	}
	reg := `^` + preCode + `0?[8]\d{10}\d{4}?$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证马来西亚手机号
func IsMalaysiaiaMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "60" {
		return false
	}
	return false
}

//验证韩国手机号
func IsKoreaMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "82" {
		return false
	}
	reg := `^` + preCode + `0?10\d{8}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证蒙古手机号
func IsMongoliaMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "976" {
		return false
	}
	reg := `^` + preCode + `^(91|95|96|99|88|77|94|81)\d{6}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证俄罗斯手机号
func IsRussiaMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "7" {
		return false
	}
	reg := `^` + preCode + `9\d{9}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证新加坡手机号
func IsSingaporeMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "65" {
		return false
	}
	reg := `^` + preCode + `^(8|9)\d{7}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证泰国手机号
func IsThailandMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "66" {
		return false
	}

	reg := `^` + preCode + `0?[689]\d{8}$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证越南手机号
func IsVietnamMobile(validStr string, preCode string) bool {
	if preCode != "" && preCode != "84" {
		return false
	}
	reg := `^` + preCode + `0?((1(2([0-9])|6([2-9])|88|99))|(9((?!5)[0-9])))([0-9]{7})$`
	flag, _ := regexp.Match(reg, []byte(preCode+validStr))
	return flag
}

//验证邮箱
func IsEmail(validStr string) bool {
	reg := `^[a-zA-Z0-9]+[a-zA-Z0-9_-]*@[a-zA-Z0-9]+\.[a-zA-Z]{2,5}$`
	flag, _ := regexp.Match(reg, []byte(validStr))
	return flag
}

//验证url
func IsUrl(validStr string) bool {
	reg := `^(?:(?:https?|ftp):\/\/|www\.)[-a-z0-9+&@#\/%?=~_|!:,.;]*[-a-z0-9+&@#\/%=~_|]$`
	flag, _ := regexp.Match(reg, []byte(validStr))
	return flag
}

//判断是否是手机号
func IsMobile(mobile string, areaCode string) bool {
	var flag = false
	switch areaCode {
	case "86":
		flag = IsChineseMobile(mobile, areaCode)
		break
	case "886":
		flag = IsTwMobile(mobile, areaCode)
		break
	case "852":
		flag = IsHkMobile(mobile, areaCode)
		break
	case "81":
		flag = IsJapaneseMobile(mobile, areaCode)
		break
	case "91":
		flag = IsIndiaMobile(mobile, areaCode)
		break
	case "7":
		flag = IsRussiaMobile(mobile, areaCode)
		break
	case "84":
		flag = IsVietnamMobile(mobile, areaCode)
		break
	case "82":
		flag = IsKoreaMobile(mobile, areaCode)
		break
	case "66":
		flag = IsThailandMobile(mobile, areaCode)
		break
	case "62":
		flag = IsIndonesiaMobile(mobile, areaCode)
		break
	case "976":
		flag = IsMongoliaMobile(mobile, areaCode)
		break
	case "65":
		flag = IsSingaporeMobile(mobile, areaCode)
		break
	}

	return flag
}
