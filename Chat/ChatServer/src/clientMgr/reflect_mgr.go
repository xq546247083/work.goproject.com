package clientMgr

import (
	"fmt"
	"reflect"
	"strings"

	. "work.goproject.com/Chat/ChatServerModel/src"
	"work.goproject.com/goutil/debugUtil"
	"work.goproject.com/goutil/logUtil"
)

const (
	// 供客户端访问的模块的后缀
	con_ModuleSuffix = "Module"

	// 定义用于分隔模块名称和方法名称的分隔符
	con_DelimeterOfObjAndMethod = "_"
)

var (
	// 定义存放所有方法映射的变量
	methodMap = make(map[string]*methodAndInOutTypes)
)

// 获取结构体类型的名称
// structType：结构体类型
// 返回值：
// 结构体类型的名称
func getStructName(structType reflect.Type) string {
	reflectTypeStr := structType.String()
	reflectTypeArr := strings.Split(reflectTypeStr, ".")

	return reflectTypeArr[len(reflectTypeArr)-1]
}

// 获取完整的模块名称
// moduleName：模块名称
// 返回值：
// 完整的模块名称
func getFullModuleName(moduleName string) string {
	return moduleName + con_ModuleSuffix
}

// 获取完整的方法名称
// structName：结构体名称
// methodName：方法名称
// 返回值：
// 完整的方法名称
func getFullMethodName(structName, methodName string) string {
	return structName + con_DelimeterOfObjAndMethod + methodName
}

// 解析方法的输入输出参数
// method：方法对应的反射值
// 返回值：
// 输入参数类型集合
// 输出参数类型集合
func resolveMethodInOutParams(method reflect.Value) (inTypes []reflect.Type, outTypes []reflect.Type) {
	methodType := method.Type()
	for i := 0; i < methodType.NumIn(); i++ {
		inTypes = append(inTypes, methodType.In(i))
	}

	for i := 0; i < methodType.NumOut(); i++ {
		outTypes = append(outTypes, methodType.Out(i))
	}

	return
}

// 将需要对客户端提供方法的对象进行注册
// structObject：对象
func RegisterFunction(structObject interface{}) {
	// 获取structObject对应的反射 Type 和 Value
	reflectValue := reflect.ValueOf(structObject)
	reflectType := reflect.TypeOf(structObject)

	// 提取对象类型名称
	structName := getStructName(reflectType)

	// 获取structObject中返回值为responseObject的方法
	for i := 0; i < reflectType.NumMethod(); i++ {
		// 获得方法名称
		methodName := reflectType.Method(i).Name

		// 判断是否为导出的方法

		// 获得方法及其输入参数的类型列表
		method := reflectValue.MethodByName(methodName)
		inTypes, outTypes := resolveMethodInOutParams(method)

		// 判断输出参数数量是否正确
		if len(outTypes) != 1 {
			continue
		}

		// 判断返回值是否为responseObject
		outType := outTypes[0]
		if _, ok := outType.FieldByName("Code"); !ok {
			continue
		}
		if _, ok := outType.FieldByName("Message"); !ok {
			continue
		}
		if _, ok := outType.FieldByName("Data"); !ok {
			continue
		}

		// 添加到列表中
		methodMap[getFullMethodName(structName, methodName)] = newmethodAndInOutTypes(method, inTypes, outTypes)

		debugUtil.Println(fmt.Sprintf("%s_%s注册成功,当前共%d个方法", structName, methodName, len(methodMap)))
	}
}

// 调用方法
// requestObj：请求对象
func callFunction(requestObj *ServerRequestObject) *ServerResponseObject {
	responseObj := NewServerResponseObject()

	var methodAndInOutTypes *methodAndInOutTypes
	var exists bool

	// 根据传入的ModuleName和MethodName找到对应的方法对象
	key := getFullMethodName(getFullModuleName(requestObj.ModuleName), requestObj.MethodName)
	if methodAndInOutTypes, exists = methodMap[key]; !exists {
		logUtil.ErrorLog("找不到指定的方法：%v,%s", methodMap, key)
		return responseObj.SetResultStatus(NoTargetMethod)
	}

	// 判断参数数量是否相同
	inTypesLength := len(methodAndInOutTypes.InTypes)
	paramLength := len(requestObj.Parameters)
	if paramLength != inTypesLength {
		logUtil.ErrorLog("传入的参数数量不符，本地方法%s的参数数量：%d，传入的参数数量为：%d", key, inTypesLength, paramLength)
		return responseObj.SetResultStatus(ParamNotMatch)
	}

	// 构造参数
	in := make([]reflect.Value, inTypesLength)
	for i := 0; i < inTypesLength; i++ {
		inTypeItem := methodAndInOutTypes.InTypes[i]
		paramItem := requestObj.Parameters[i]

		// 已支持类型：Client,Player(非基本类型)
		// 已支持类型：Bool,Int,Int8,Int16,Int32,Int64,Uint,Uint8,Uint16,Uint32,Uint64,Float32,Float64,String
		// 已支持类型：以及上面所列出类型的Slice类型
		// 未支持类型：Uintptr,Complex64,Complex128,Array,Chan,Func,Interface,Map,Ptr,Struct,UnsafePointer
		// 由于byte与int8同义，rune与int32同义，所以并不需要单独处理
		// 枚举参数的类型，并进行类型转换
		switch inTypeItem.Kind() {
		case reflect.Interface:
			if param_client, ok := paramItem.(IClient); ok {
				in[i] = reflect.ValueOf(param_client)
			}
		case reflect.Ptr:
			if param_player, ok := paramItem.(*Player); ok {
				in[i] = reflect.ValueOf(param_player)
			}
		case reflect.Bool:
			if param_bool, ok := paramItem.(bool); ok {
				in[i] = reflect.ValueOf(param_bool)
			}
		case reflect.Int:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(int(param_float64))
			}
		case reflect.Int8:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(int8(param_float64))
			}
		case reflect.Int16:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(int16(param_float64))
			}
		case reflect.Int32:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(int32(param_float64))
			}
		case reflect.Int64:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(int64(param_float64))
			}
		case reflect.Uint:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(uint(param_float64))
			}
		case reflect.Uint8:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(uint8(param_float64))
			}
		case reflect.Uint16:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(uint16(param_float64))
			}
		case reflect.Uint32:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(uint32(param_float64))
			}
		case reflect.Uint64:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(uint64(param_float64))
			}
		case reflect.Float32:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(float32(param_float64))
			}
		case reflect.Float64:
			if param_float64, ok := paramItem.(float64); ok {
				in[i] = reflect.ValueOf(param_float64)
			}
		case reflect.String:
			if param_string, ok := paramItem.(string); ok {
				in[i] = reflect.ValueOf(param_string)
			}
		case reflect.Slice:
			// 如果是Slice类型，则需要对其中的项再次进行类型判断及类型转换
			if param_interface, ok := paramItem.([]interface{}); ok {
				switch inTypeItem.String() {
				case "[]bool":
					params_inner := make([]bool, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_bool, ok := param_interface[i].(bool); ok {
							params_inner[i] = param_bool
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]int":
					params_inner := make([]int, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = int(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]int8":
					params_inner := make([]int8, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = int8(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]int16":
					params_inner := make([]int16, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = int16(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]int32":
					params_inner := make([]int32, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = int32(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]int64":
					params_inner := make([]int64, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = int64(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]uint":
					params_inner := make([]uint, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = uint(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				// case "[]uint8": 特殊处理
				case "[]uint16":
					params_inner := make([]uint16, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = uint16(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]uint32":
					params_inner := make([]uint32, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = uint32(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]uint64":
					params_inner := make([]uint64, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = uint64(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]float32":
					params_inner := make([]float32, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = float32(param_float64)
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]float64":
					params_inner := make([]float64, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_float64, ok := param_interface[i].(float64); ok {
							params_inner[i] = param_float64
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				case "[]string":
					params_inner := make([]string, len(param_interface))
					for i := 0; i < len(param_interface); i++ {
						if param_string, ok := param_interface[i].(string); ok {
							params_inner[i] = param_string
						}
					}
					in[i] = reflect.ValueOf(params_inner)
				}
			} else if inTypeItem.String() == "[]uint8" { // 由于[]uint8在传输过程中会被转化成字符串，所以单独处理;
				if param_string, ok := paramItem.(string); ok {
					param_uint8 := ([]uint8)(param_string)
					in[i] = reflect.ValueOf(param_uint8)
				}
			}
		}
	}

	// 判断是否有无效的参数（传入的参数类型和方法定义的类型不匹配导致没有赋值）
	for _, item := range in {
		if reflect.Value.IsValid(item) == false {
			logUtil.ErrorLog("type:%v,value:%v.方法%s传入的参数%v无效", reflect.TypeOf(item), reflect.ValueOf(item), key, requestObj.Parameters)
			return responseObj.SetResultStatus(ParamInValid)
		}
	}

	// 传入参数，调用方法
	out := methodAndInOutTypes.Method.Call(in)

	// 由于只有一个返回值，所以取out[0]
	if retResponseObj, ok := (out[0]).Interface().(ServerResponseObject); ok {
		responseObj.SetMethodName(requestObj.MethodName)
		responseObj.SetResultStatus(retResponseObj.ResultStatus)
		responseObj.SetData(retResponseObj.Data)
	} else {
		logUtil.ErrorLog("(&out[0]).Interface()转换类型失败")
	}

	return responseObj
}
