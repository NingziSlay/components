package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// MapConfig 可以从环境变量读取配置灌入到 dest 中，dest 必须是一个指向结构体的指针，
// 会解析结构体每个字段，通过结构体标签得到这个字段对应的环境变量名，自动将从环境变量读取
// 到的字符串值转为 dest 字段对应的类型并为其赋值
// 如果 dest 为 nil，返回 ErrorNilInput
// 如果 dest 不是一个指针，返回 ErrorNonPointer
// 如果 dest 指针指向的值不是一个结构体，返回 ErrorNonStruct
// 如果环境变量的值转为字段相应类型时发生错误，立即返回相应的错误
//
// MapConfig 允许字段值为空
func MapConfig(dest interface{}) error {
	return newMapper(false).mapper(dest)
}

// MustMapConfig 所有可导出字段都不允许为空
func MustMapConfig(dest interface{}) error {
	return newMapper(true).mapper(dest)
}

type mapper struct {
	strict bool
}

func newMapper(strict bool) *mapper {
	return &mapper{strict: strict}
}

/* dest 必须是一个指向结构体的指针

type Config struct{}
var config *Config (config 是一个空指针）
m.mapper(config) 传的是这个指针的值，也就是一个 nil，无法操作
m.mapStruct(&config) 传的是这个指向 *config 的指针，通过 reflect 可以设置这个指针指向的 *config 的值

                 检查入参
                    |  \
                  结构体 其他 --> error
                    |
     -------->  分解结构体
    |               |
    |           遍历每个字段
    |               |
    |           是否可导出  ---No--> 跳过
    |               |
    |            是否忽略  ---Yes--> 跳过
    |               |
    |             检查值
    |               |
    |             设置值 <----------------
    |               |                    |
    |             /   \                  |
    |          结构体  指针 ---behind()---
    |           /
     ----------
*/
func (m *mapper) mapper(dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return ErrorNonPointer
	}
	if v.IsNil() {
		return ErrorNilInput
	}
	v = behind(v)
	if v.Kind() != reflect.Struct {
		return ErrorNonStruct
	}
	return m.mapStruct(v)
}

// behind 返回 v 指针指向的最终值，如果 v 本身就是一个指针，就直接返回 v 自身
func behind(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}
	// 如果 v 是 nil，使用 reflect 为 v 创建一个值
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}

	return behind(reflect.Indirect(v))
}

// mapStruct 结构体内的字段处理
func (m *mapper) mapStruct(v reflect.Value) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		data := getData(ft)
		// 标签为 "-" 或非导出字段
		if data.shouldSkip() {
			continue
		}

		// 如果字段的 env 值为空，判断是否是结构体，如果是结构体则忽略，否则根据 strict 判断是否返回错误
		if !data.isValid() {
			if behind(fv).Kind() != reflect.Struct {
				if m.strict {
					return newE("missing value of %s.%s", t.String(), ft.Name)
				}
				continue
			}
		}
		// set value
		if err := m.setFieldValue(fv, data.val); err != nil {
			return err
		}
	}
	return nil
}

// setFieldValue 给结构体字段赋值
// 如果 v 不可寻址或是不可导出字段（字段首字母小写），返回 ErrorNotWritable 错误
// 给结构体赋值需要转换为对应类型，如果类型转换错误，返回相应的错误
func (m *mapper) setFieldValue(v reflect.Value, value string) error {
	switch v.Kind() {
	default:
		return newE("unsupported type: %s", v.String())
	case reflect.String:
		v.SetString(value)
	case reflect.Int8:
		return m.setInt8(v, value)
	case reflect.Int16:
		return m.setInt16(v, value)
	case reflect.Int32:
		return m.setInt32(v, value)
	case reflect.Int64, reflect.Int:
		return m.setInt64(v, value)
	case reflect.Uint8:
		return m.setUint8(v, value)
	case reflect.Uint16:
		return m.setUint16(v, value)
	case reflect.Uint32:
		return m.setUint32(v, value)
	case reflect.Uint64, reflect.Uint:
		return m.setUint64(v, value)
	case reflect.Float32:
		return m.setFloat32(v, value)
	case reflect.Float64:
		return m.setFloat64(v, value)
	case reflect.Bool:
		return m.setBool(v, value)
	case reflect.Slice:
		return m.setSlice(v, value)
	case reflect.Array:
		return m.setArray(v, value)
	case reflect.Ptr:
		return m.setPtr(v, value)
	case reflect.Struct:
		return m.setStruct(v)
	}
	return nil
}

// 下面的所有方法都是为了 setFieldValue 服务，针对结构体中不同的类型，
// 把从环境变量中得到的值转为对应的类型，并赋值给对应的 field 中
//
// 下面的所有方法都假定 v 是可写的，不会做可写判断
//
// 如果类型转换发生错误，直接返回
// 如果 field 是结构体、切片、列表、指针，需要递归调用的情况，可能会返回
// 方法 mapConfig 和方法 setFieldValue 中的错误

// setInt8 设置 int8 类型
func (m *mapper) setInt8(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return wrapE("setInt8", err)
	}
	v.SetInt(i)
	return nil
}

// setInt16 设置 int16 类型
func (m *mapper) setInt16(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return wrapE("setInt16", err)
	}
	v.SetInt(i)
	return nil
}

// setInt32 设置 int32 类型
func (m *mapper) setInt32(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return wrapE("setInt32", err)
	}
	v.SetInt(i)
	return nil
}

// setInt64 设置 int64 和 int 类型
func (m *mapper) setInt64(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return wrapE("setInt64", err)
	}
	v.SetInt(i)
	return nil
}

// setUint8 设置 uint8 类型
func (m *mapper) setUint8(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return wrapE("setUint8", err)
	}
	v.SetUint(i)
	return nil
}

// setUint16 设置 uint16 类型
func (m *mapper) setUint16(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return wrapE("setUint16", err)
	}
	v.SetUint(i)
	return nil
}

// setUint32 设置 uint32 类型
func (m *mapper) setUint32(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return wrapE("setUint32", err)
	}
	v.SetUint(i)
	return nil
}

// setUint64 设置 uint64 和 uint 类型
func (m *mapper) setUint64(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return wrapE("setUint64", err)
	}
	v.SetUint(i)
	return nil
}

// setFloat32 设置 float32 类型
func (m *mapper) setFloat32(v reflect.Value, value string) error {
	i, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return wrapE("setFloat32", err)
	}
	v.SetFloat(i)
	return nil
}

// setFloat64 设置 float64 类型
func (m *mapper) setFloat64(v reflect.Value, value string) error {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return wrapE("setFloat64", err)
	}
	v.SetFloat(i)
	return nil
}

// setBool 设置 bool 类型
func (m *mapper) setBool(v reflect.Value, value string) error {
	i, err := strconv.ParseBool(value)
	if err != nil {
		return wrapE("setBool", err)
	}
	v.SetBool(i)
	return nil
}

// setSlice 设置 slice 类型
func (m *mapper) setSlice(v reflect.Value, value string) error {
	tags := strings.Split(value, ",")
	slice := reflect.MakeSlice(v.Type(), len(tags), cap(tags))
	for i, t := range tags {
		elem := slice.Index(i)
		err := m.setFieldValue(elem, t)
		if err != nil {
			return wrapE("setSlice", err)
		}
	}
	v.Set(slice)
	return nil
}

// setArray 设置 array 类型
func (m *mapper) setArray(v reflect.Value, value string) error {
	tags := strings.Split(value, ",")
	if len(tags) > v.Cap() {
		return newE("array out of range, max: %d", v.Cap())
	}
	for i, t := range tags {
		err := m.setFieldValue(v.Index(i), t)
		if err != nil {
			return wrapE("setArray", err)
		}
	}
	return nil
}

// setPtr 设置 ptr 类型
func (m *mapper) setPtr(v reflect.Value, value string) error {
	v = behind(v)
	return m.setFieldValue(v, value)
}

// setStruct 设置 struct 类型
func (m *mapper) setStruct(v reflect.Value) error {
	return m.mapStruct(v)
}

const tagName = "env"

type data struct {
	typ      reflect.StructField // field type
	key      string              // env key
	val      string              // env value
	_default string              // default value, use replace when val is empty
	skip     bool                // - 则直接跳过
}

func getData(field reflect.StructField) (t *data) {
	t = &data{typ: field}
	// 非导出字段
	if field.PkgPath != "" {
		t.skip = true
		return
	}
	tag := field.Tag.Get(tagName)
	// 忽略符
	if tag == "-" {
		t.skip = true
		return
	}

	// 标签为空，默认使用字段名下划线大写命名作为默认环境变量名
	if tag == "" {
		t.key = camelCaseToUnderscoreUpper(field.Name)
	} else {
		tags := strings.SplitN(tag, ",", 2)
		t.key = tags[0]
		if len(tags) > 1 {
			t._default = tags[1]
		}
	}

	t.val = os.Getenv(t.key)
	if t.val == "" {
		t.val = t._default
	}
	t.val = strings.Trim(t.val, " ")
	return t
}

func (t *data) shouldSkip() bool {
	return t.skip
}

func (t *data) isValid() bool {
	return t.val != ""
}

// 驼峰单词转下划线单词
func camelCaseToUnderscoreUpper(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, r)
		} else {
			if unicode.IsUpper(r) {
				output = append(output, '_')
			}
			output = append(output, r)
		}
	}
	return strings.ToUpper(string(output))
}
