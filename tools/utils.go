package tools

import (
	"bytes"
	"strings"
)

// BuildInsertValuePlaceHolder 生成批量写入 Mysql 时的占位符
// INSERT INTO `table`(col1, col2) values < 占位符... >
// 如果需要批量插入时 values 后面的占位符处理起来会很痛苦，这个方法可以根据字段
// 的数量和插入的数据的数量生成相应的占位符，业务层只需要 Fprintf 插入到自己的
// sql 中即可
// colCount 指的是字段的数量，控制每个括号里有几个占位符
// valCount 指的是需要插入的数据的数量，控制需要生成多少个占位符
//
// 	str := BuildInsertValuePlaceHolder(2, 3)
//	str == "(?, ?), (?, ?), (?, ?)"
func BuildInsertValuePlaceHolder(itemLength, valueLength int) string {
	s := strings.Repeat("?, ", itemLength)
	s = strings.Join([]string{"(", s[:len(s)-2], ")", ", "}, "")
	s = strings.Repeat(s, valueLength)
	return s[:len(s)-2]
}

func buildInsertValuePlaceHolderV1(itemLength, valueLength int) string {
	return buildBytes(valueLength, itemLength)
}

func buildBytes(out, in int) string {
	length := 3*in*out - 2*out - 2

	var buffer bytes.Buffer
	buffer.Grow(length)
	for i := 0; i < out; i++ {
		buffer.WriteString("(")
		for j := 0; j < in; j++ {
			buffer.WriteRune('?')
			if j+1 == in {
				continue
			}
			buffer.WriteString(", ")
		}
		buffer.WriteString(")")
		if i+1 == out {
			continue
		}
		buffer.WriteString(", ")
	}
	return buffer.String()
}

func buildInsertValuePlaceHolderV2(itemLength, valueLength int) string {
	return buildString(valueLength, itemLength)
}

func buildString(out, in int) string {
	buffer := strings.Builder{}
	buffer.Grow(3*in*out - 2*out - 2)
	for i := 0; i < out; i++ {
		buffer.WriteString("(")
		for j := 0; j < in; j++ {
			buffer.WriteString("?")
			if j+1 == in {
				continue
			}
			buffer.WriteString(", ")
		}
		buffer.WriteString(")")
		if i+1 == out {
			continue
		}
		buffer.WriteString(", ")
	}
	return buffer.String()
}

func buildInsertValuePlaceHolderV3(itemLength, valueLength int) string {
	return buildV3(buildV3("?", itemLength, true), valueLength, false)
}

func buildV3(str string, length int, brackets bool) string {
	values := make([]string, 0, length)
	for i := 0; i < length; i++ {
		values = append(values, str)
	}
	if brackets {
		return "(" + strings.Join(values, ", ") + ")"
	}
	return strings.Join(values, ", ")
}

func buildInsertValuePlaceHolderV4(colCount, valCount int) string {
	return buildV4(
		buildV4("?", colCount*3, colCount, true),
		3*colCount*valCount+2*valCount-2, valCount,
		false,
	)
}

func buildV4(str string, cap int, len int, brackets bool) string {
	buffer := strings.Builder{}
	buffer.Grow(cap)
	if brackets {
		buffer.Write([]byte("("))
	}
	for i := 1; i <= len; i++ {
		buffer.WriteString(str)
		if i == len {
			break
		}
		buffer.WriteString(", ")
	}
	if brackets {
		buffer.WriteString(")")
	}
	return buffer.String()
}

/*
(?, ?, ?, ?) 12	 4   4 + 2 + 3 + 3
(?, ?, ?)    9	 3	 3 + 2 + 2 + 2
(?, ?)       6	 2   2 + 2 + 1 + 1
(?)			 3   1   1 + 2

----[3x]----

(?, ?), (?, ?), (?, ?), (?, ?) 	   30  	4	6 * 4 + 2 * 3
(?, ?), (?, ?), (?, ?)  		   22  	3	6 * 3 + 2 * 2
(?, ?), (?, ?)  				   14  	2	6 * 2 + 2
(?, ?)							   6  	1	6 * 1

----[3xy + 2y - 2]---
*/
