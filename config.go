package components

import (
	"github.com/pkg/errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const tagName = "env"

// MapConfig 给所有可导出字段赋值
func MapConfig(dest interface{}) error {
	m := &mapper{strict: false}
	return m.mapConfig(reflect.TypeOf(dest), reflect.ValueOf(dest))
}

// MustMapConfig 所有可导出字段都必须赋值
// 每个可导出字段的值都不能为空，整个 mapping 过程中有 err 立即返回
func MustMapConfig(dest interface{}) error {
	m := &mapper{strict: true}
	return m.mapConfig(reflect.TypeOf(dest), reflect.ValueOf(dest))
}

type fieldTag struct {
	FieldName string
	Index     []int
	EnvName   string
	Default   string
	Value     string
}

// newTag 从 filed 中读取 env 标签的值，并从环境变量中读取，返回 *fieldTag
//
// 如果 tag 为空，只是 field name 的全大写
// 通过 "," 分割 tag 值，只分割第一个 ","，分割后的第二个值为字段的默认值
// 如果从环境变量中没有到对应的值，使用默认值代替
func newTag(field reflect.StructField) (t *fieldTag) {
	t = &fieldTag{
		FieldName: field.Name,
		Index:     field.Index,
	}
	tag := field.Tag.Get(tagName)
	if tag == "" {
		t.EnvName = strings.ToUpper(t.FieldName)
		return
	}
	tags := strings.SplitN(tag, ",", 2)
	t.EnvName = tags[0]
	if len(tags) > 1 {
		t.Default = tags[1]
	}
	t.Value = os.Getenv(t.EnvName)
	if t.Value == "" {
		t.Value = t.Default
	}
	t.Value = strings.Trim(t.Value, " ")
	return
}

type mapper struct {
	strict bool
}

func (m *mapper) mapConfig(t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.Ptr:
		// mapConfig 指针指向的值
		return m.mapConfig(t.Elem(), v.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			fieldType := t.Field(i)
			// 跳过非导出字段
			if fieldType.PkgPath != "" {
				continue
			}
			// 嵌入结构体处理
			if fieldType.Anonymous {
				if err := m.mapConfig(fieldType.Type, v.Field(i)); err != nil {
					return err
				}
				continue
			}
			// 处理标签
			tag := newTag(fieldType)
			if err := m.checkEmptyTag(tag, fieldType); err != nil {
				return err
			}
			fieldValue := v.FieldByIndex(tag.Index)
			// fieldValue 应该可寻址，可导出
			if !fieldValue.CanSet() {
				return ErrorNotWriteable
			}
			if err := m.setFieldValue(fieldValue, tag.Value); err != nil {
				if m.strict {
					return err
				}
				continue
			}
		}
	default:
		return ErrorUnsupportedType(t.String())
	}
	return nil
}

// checkEmptyTag 在严格模式下 tag 的值不允许为空，但是需要忽略 struct 和底层
// 类型指向一个 struct 的 ptr
func (m *mapper) checkEmptyTag(tag *fieldTag, fieldType reflect.StructField) error {
	if m.strict && tag.Value == "" {
		switch fieldType.Type.Kind() {
		case reflect.Struct:
			return nil
		case reflect.Ptr:
			if fieldType.Type.Elem().Kind() == reflect.Struct {
				return nil
			}
			return ErrorEmptyEnviron
		default:
			return ErrorEmptyEnviron
		}
	}
	return nil
}

// setFieldValue 给结构体字段赋值
// 如果 v 不可寻址或是不可导出字段（字段首字母小写），返回 ErrorNotWriteable 错误
// 给结构体赋值需要转换为对应类型，如果类型转换错误，返回相应的错误
func (m *mapper) setFieldValue(v reflect.Value, value string) error {
	switch v.Kind() {
	default:
		return ErrorUnsupportedType(v.String())
	case reflect.String:
		v.SetString(strings.Trim(value, " "))
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
		return errors.Wrap(err, "setInt8")
	}
	v.SetInt(i)
	return nil
}

// setInt16 设置 int16 类型
func (m *mapper) setInt16(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return errors.Wrap(err, "setInt16")
	}
	v.SetInt(i)
	return nil
}

// setInt32 设置 int32 类型
func (m *mapper) setInt32(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return errors.Wrap(err, "setInt32")
	}
	v.SetInt(i)
	return nil
}

// setInt64 设置 int64 和 int 类型
func (m *mapper) setInt64(v reflect.Value, value string) error {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return errors.Wrap(err, "setInt64")
	}
	v.SetInt(i)
	return nil
}

// setUint8 设置 uint8 类型
func (m *mapper) setUint8(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return errors.Wrap(err, "setUint8")
	}
	v.SetUint(i)
	return nil
}

// setUint16 设置 uint16 类型
func (m *mapper) setUint16(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return errors.Wrap(err, "setUint16")
	}
	v.SetUint(i)
	return nil
}

// setUint32 设置 uint32 类型
func (m *mapper) setUint32(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return errors.Wrap(err, "setUint32")
	}
	v.SetUint(i)
	return nil
}

// setUint64 设置 uint64 和 uint 类型
func (m *mapper) setUint64(v reflect.Value, value string) error {
	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return errors.Wrap(err, "setUint64")
	}
	v.SetUint(i)
	return nil
}

// setFloat32 设置 float32 类型
func (m *mapper) setFloat32(v reflect.Value, value string) error {
	i, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return errors.Wrap(err, "setFloat32")
	}
	v.SetFloat(i)
	return nil
}

// setFloat64 设置 float64 类型
func (m *mapper) setFloat64(v reflect.Value, value string) error {
	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return errors.Wrap(err, "setFloat64")
	}
	v.SetFloat(i)
	return nil
}

// setBool 设置 bool 类型
func (m *mapper) setBool(v reflect.Value, value string) error {
	i, err := strconv.ParseBool(value)
	if err != nil {
		return errors.Wrap(err, "setBool")
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
			return errors.Wrap(err, "setSlice")
		}
	}
	v.Set(slice)
	return nil
}

// setArray 设置 array 类型
func (m *mapper) setArray(v reflect.Value, value string) error {
	tags := strings.Split(value, ",")
	if len(tags) > v.Cap() {
		return ErrorArrayOutOfRange
	}
	for i, t := range tags {
		err := m.setFieldValue(v.Index(i), t)
		if err != nil {
			return errors.Wrap(err, "setArray")
		}
	}
	return nil
}

// setPtr 设置 ptr 类型
func (m *mapper) setPtr(v reflect.Value, value string) error {
	var s reflect.Value
	s = reflect.New(v.Type().Elem())
	if err := m.setFieldValue(s.Elem(), value); err != nil {
		return errors.Wrap(err, "setPtr")
	}
	v.Set(s)
	return nil
}

// setStruct 设置 struct 类型
func (m *mapper) setStruct(v reflect.Value) error {
	err := m.mapConfig(v.Type(), v)
	if err != nil {
		return err
	}
	return nil
}
