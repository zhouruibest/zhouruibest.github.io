# ParseInt

n, err := strconv.ParseInt(str1, 10, 8) //10是进制； 8指的是 转换结果最大值不超过 int8 即 127, 否则会报错