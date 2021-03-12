package pkg

import (
	"encoding/json"
	"math"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func toString(s interface{}) string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

func TestMapper_Parse(t *testing.T) {
	type sample struct{}
	var s sample
	// non pointer
	m := newMapper(true)
	if err := m.Parse(s); err == nil {
		t.Fatalf("expect error returned")
	}
	// nil
	var s1 *sample
	if err := m.Parse(s1); err == nil {
		t.Fatalf("expect error returned")
	}
	// non struct
	var s2 string
	if err := m.Parse(&s2); err == nil {
		t.Fatalf("expect error returned")
	}

	// no error
	if err := m.Parse(&s); err != nil {
		t.Fatalf("expect no error, got: %s", err)
	}
}

// 非导出字段测试
func TestMapperUnexported(t *testing.T) {
	type Config struct {
		url string `env:"URL,https://ningzi.club"`
	}
	var c Config
	var c1 Config
	_ = MustMapConfig(&c)
	_ = MapConfig(&c1)

	if c.url != "" {
		t.Fatalf("unexported failed should be ignore")
	}
	if c1.url != "" {
		t.Fatalf("unexported failed should be ignore")
	}
}

// 测试从环境变量读取
func TestMapperENV(t *testing.T) {
	defer os.Clearenv()
	type Config struct {
		Name string `env:"LIST,hello"`
	}
	var c Config
	if err := os.Setenv("LIST", "world"); err != nil {
		panic(err)
	}
	_ = MustMapConfig(&c)
	if c.Name != "world" {
		t.Fatalf("unexpected value")
	}
}

// 测试 "-" 标签
func TestMapperSkip(t *testing.T) {
	defer os.Clearenv()
	type Config struct {
		Name string `env:"-"`
	}
	if err := os.Setenv("NAME", "should be ignored"); err != nil {
		panic(err)
	}
	var c Config
	_ = MustMapConfig(&c)
	if c.Name != "" {
		t.Fatalf("tag with `-` should be ignored")
	}
}

// 测试数值类型
func TestMapperInteger(t *testing.T) {
	defer os.Clearenv()
	type Int8 struct {
		Int8 int8
	}
	type Int16 struct {
		Int16 int16
	}

	type Int32 struct {
		Int32 int32
	}

	type Int64 struct {
		Int64 int64
	}

	var (
		i8  Int8
		i16 Int16
		i32 Int32
		i64 Int64
	)

	var cases = []struct {
		input    interface{}
		err      bool
		setter   func()
		function func(interface{}) error
	}{
		{
			input:    &i8,
			err:      false,
			function: MapConfig, // 忽略空值
		},
		{
			input:    &i8,
			err:      true,
			function: MustMapConfig, // 空值返回错误
		},
		{
			input: &i8,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT8", strconv.Itoa(math.MaxInt8+1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &i8,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT8", strconv.Itoa(math.MaxInt8+1)) // out of range
			},
			function: MapConfig,
		},
		{
			input: &i8,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT8", strconv.Itoa(math.MinInt8-1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &i8,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT8", strconv.Itoa(math.MinInt8-1)) // out of range
			},
			function: MapConfig,
		},

		{
			input:    &i16,
			err:      false,
			function: MapConfig, // 忽略空值
		},
		{
			input:    &i16,
			err:      true,
			function: MustMapConfig, // 空值返回错误
		},
		{
			input: &i16,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT16", strconv.Itoa(math.MaxInt16+1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &i16,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT16", strconv.Itoa(math.MaxInt16+1)) // out of range
			},
			function: MapConfig,
		},
		{
			input: &i16,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT16", strconv.Itoa(math.MinInt16-1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &i16,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT16", strconv.Itoa(math.MinInt16-1)) // out of range
			},
			function: MapConfig,
		},

		{
			input:    &i32,
			err:      false,
			function: MapConfig, // 忽略空值
		},
		{
			input:    &i32,
			err:      true,
			function: MustMapConfig, // 空值返回错误
		},
		{
			input: &i32,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT32", strconv.Itoa(math.MaxInt32+1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &i32,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT32", strconv.Itoa(math.MaxInt32+1)) // out of range
			},
			function: MapConfig,
		},
		{
			input: &i32,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT32", strconv.Itoa(math.MinInt32-1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &i32,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT32", strconv.Itoa(math.MinInt32-1)) // out of range
			},
			function: MapConfig,
		},

		{
			input:    &i64,
			err:      false,
			function: MapConfig, // 忽略空值
		},
		{
			input:    &i64,
			err:      true,
			function: MustMapConfig, // 空值返回错误
		},
		{
			input: &i64,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT64", "9223372036854775809") // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &i64,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT64", "9223372036854775809") // out of range
			},
			function: MapConfig,
		},
		{
			input: &i64,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT64", "-9223372036854775810") // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &i64,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT64", "-9223372036854775810") // out of range
			},
			function: MapConfig,
		},
	}
	for i, c := range cases {
		os.Clearenv()
		if c.setter != nil {
			c.setter()
		}
		err := c.function(c.input)
		if c.err {
			if err == nil {
				t.Fatalf("expect error return, got nil, index: %d", i)
			}
		} else {
			if err != nil {
				t.Fatalf("expect no error return, got: %s", err)
			}
		}
	}
}

// 测试无符号数值类型
func TestMapperUInteger(t *testing.T) {
	defer os.Clearenv()
	type Uint8 struct {
		Int8 uint8
	}
	type Uint16 struct {
		Int16 uint16
	}

	type Uint32 struct {
		Int32 uint32
	}

	type Uint64 struct {
		Int64 uint64
	}

	var (
		u8  Uint8
		u16 Uint16
		u32 Uint32
		u64 Uint64
	)

	var cases = []struct {
		input    interface{}
		err      bool
		setter   func()
		function func(interface{}) error
	}{
		{
			input:    &u8,
			err:      false,
			function: MapConfig, // 忽略空值
		},
		{
			input:    &u8,
			err:      true,
			function: MustMapConfig, // 空值返回错误
		},
		{
			input: &u8,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT8", strconv.Itoa(math.MaxUint8+1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &u8,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT8", strconv.Itoa(math.MaxUint8+1)) // out of range
			},
			function: MapConfig,
		},
		{
			input: &u8,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT8", "-1") // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &u8,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT8", "-1") // out of range
			},
			function: MapConfig,
		},

		{
			input:    &u16,
			err:      false,
			function: MapConfig, // 忽略空值
		},
		{
			input:    &u16,
			err:      true,
			function: MustMapConfig, // 空值返回错误
		},
		{
			input: &u16,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT16", strconv.Itoa(math.MaxUint16+1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &u16,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT16", strconv.Itoa(math.MaxUint16+1)) // out of range
			},
			function: MapConfig,
		},
		{
			input: &u16,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT16", "-1") // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &u16,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT16", "-1") // out of range
			},
			function: MapConfig,
		},

		{
			input:    &u32,
			err:      false,
			function: MapConfig, // 忽略空值
		},
		{
			input:    &u32,
			err:      true,
			function: MustMapConfig, // 空值返回错误
		},
		{
			input: &u32,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT32", strconv.Itoa(math.MaxUint32+1)) // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &u32,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT32", strconv.Itoa(math.MaxUint32+1)) // out of range
			},
			function: MapConfig,
		},
		{
			input: &u32,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT32", "-1") // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &u32,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT32", "-1") // out of range
			},
			function: MapConfig,
		},

		{
			input:    &u64,
			err:      false,
			function: MapConfig, // 忽略空值
		},
		{
			input:    &u64,
			err:      true,
			function: MustMapConfig, // 空值返回错误
		},
		{
			input: &u64,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT64", "18446744073709551617") // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &u64,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT64", "18446744073709551617") // out of range
			},
			function: MapConfig,
		},
		{
			input: &u64,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT64", "-1") // out of range
			},
			function: MustMapConfig,
		},
		{
			input: &u64,
			err:   true,
			setter: func() {
				_ = os.Setenv("INT64", "-1") // out of range
			},
			function: MapConfig,
		},
	}
	for i, c := range cases {
		os.Clearenv()
		if c.setter != nil {
			c.setter()
		}
		err := c.function(c.input)
		if c.err {
			if err == nil {
				t.Fatalf("expect error return, got nil, index: %d", i)
			}
		} else {
			if err != nil {
				t.Fatalf("expect no error return, got: %s", err)
			}
		}
	}
}

// 测试 bool
func TestMapperBool(t *testing.T) {
	type Config struct {
		Debug bool `env:"DEBUG,true"`
		Bool  bool `env:"BOOL,false"`
	}
	var c Config
	if err := MapConfig(&c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := Config{Debug: true, Bool: false}
	if !reflect.DeepEqual(c, expect) {
		t.Fatalf("error mapper")
	}
}

// 测试 float
func TestMapperFloat(t *testing.T) {
	type Config struct {
		Float64 float64 `env:"Float64,0.123"`
		Float32 float32 `env:"Float32,0.123"`
	}
	var c Config
	if err := MapConfig(&c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f1, _ := strconv.ParseFloat("0.123", 64)
	f2, _ := strconv.ParseFloat("0.123", 32)
	expect := Config{Float64: f1, Float32: float32(f2)}
	if !reflect.DeepEqual(c, expect) {
		t.Fatalf("error mapper")
	}
}

// 测试切片
func TestMapperSlice(t *testing.T) {
	type Config struct {
		Slice []int `env:"SLICE,1,2,3"`
	}
	var c Config
	if err := MapConfig(&c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expect := Config{Slice: []int{1, 2, 3}}
	if !reflect.DeepEqual(c, expect) {
		t.Fatalf("error mapper")
	}
}

// 测试 array
func TestMapperArray(t *testing.T) {
	type Config struct {
		Slice [3]int `env:"SLICE,1,2,3"`
	}
	var c Config
	if err := MapConfig(&c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expect := Config{Slice: [3]int{1, 2, 3}}
	if !reflect.DeepEqual(c, expect) {
		t.Fatalf("error mapper")
	}
}

// 测试内嵌结构体
func TestMapperEmbedded(t *testing.T) {
	type Mysql struct {
		Url string `env:"URL,some url"`
	}
	type Book struct {
		Name string `env:"NAME,活着"`
	}
	type Config struct {
		Ptr    *Mysql
		Struct Mysql
		Mysql
		*Book
	}

	var c Config
	_ = MustMapConfig(&c)
	if c.Struct.Url == "" {
		t.Fatal("embedded failed")
	}
	if c.Ptr.Url == "" {
		t.Fatal("embedded failed")
	}
	if c.Url == "" {
		t.Fatal("embedded failed")
	}
	if c.Name == "" {
		t.Fatal("embedded failed")
	}
}

// **************************** Benchmark ****************************

func BenchmarkMustMapConfig(b *testing.B) {
	type Embed struct {
		U [5]int64 `env:"U,1,2,3,4,5"`
	}
	type Embed1 struct {
		V string `env:"hello again"`
	}
	type foo struct {
		A int8     `env:"A,-1"`
		B int16    `env:"B,-1"`
		C int32    `env:"C,-1"`
		D int64    `env:"C,-1"`
		E int      `env:"C,-1"`
		F uint8    `env:"F,1"`
		G uint16   `env:"G,1"`
		H uint32   `env:"H,1"`
		I uint64   `env:"I,1"`
		J uint     `env:"J,1"`
		K float32  `env:"K,0.123"`
		L float64  `env:"L,0.123"`
		M bool     `env:"M,True"`
		N string   `env:"N,hello"`
		O []string `env:"O,hello,world"`
		P [2]int64 `env:"P,1,2"`
		Q struct {
			R uint `env:"R,1"`
		}
		S *struct {
			T []bool `env:"T,T,F,0,1"`
		}
		Embed
		*Embed1
	}

	var s foo
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MustMapConfig(&s)
	}
}

func BenchmarkMapConfig(b *testing.B) {
	type Embed struct {
		U [5]int64 `env:"U,1,2,3,4,5"`
	}
	type Embed1 struct {
		V string `env:"hello again"`
	}
	type foo struct {
		A int8     `env:"A,-1"`
		B int16    `env:"B,-1"`
		C int32    `env:"C,-1"`
		D int64    `env:"C,-1"`
		E int      `env:"C,-1"`
		F uint8    `env:"F,1"`
		G uint16   `env:"G,1"`
		H uint32   `env:"H,1"`
		I uint64   `env:"I,1"`
		J uint     `env:"J,1"`
		K float32  `env:"K,0.123"`
		L float64  `env:"L,0.123"`
		M bool     `env:"M,True"`
		N string   `env:"N,hello"`
		O []string `env:"O,hello,world"`
		P [2]int64 `env:"P,1,2"`
		Q struct {
			R uint `env:"R,1"`
		}
		S *struct {
			T []bool `env:"T,T,F,0,1"`
		}
		Embed
		*Embed1
	}

	var s foo
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapConfig(&s)
	}
}
