package utils

import "math/rand"

// RandomStr 生成随机长度的字符串，length控制长度，当randSet为空时，默认使用数字+大小写英文字符作为字符集，spark为是否添加分隔符(_和-),
// 当randSet不为空时，使用randSet作为字符集，spark无效
func RandomStr(length int, spark bool, randSet string) string {
	if length == 0 {
		return ""
	}
	if randSet == "" {
		randSet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		if spark {
			randSet += "-_"
		}
	}

	bytes := make([]byte, length)
	max := len(randSet)
	for i := 0; i < length; i++ {
		bytes[i] = randSet[RandomInt(0, max)]
	}
	return string(bytes)
}

// RandomInt 生成随机整数，介于[min,max)范围
func RandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}
