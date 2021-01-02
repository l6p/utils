package json

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type EmptyKey struct{}

type NormalKey struct {
	Name string
}

type ArrayKey struct {
	Name  string
	Index int
}

type Value struct {
	Parent *interface{}
	Key    interface{}
}

func (v *Value) Get() interface{} {
	return v.do(func(parent, currentValue interface{}, keyName string) interface{} {
		return currentValue
	})
}

func (v *Value) Set(value interface{}) {
	v.do(func(parent, currentValue interface{}, keyName string) interface{} {
		return value
	})
}

func (v *Value) Append(value interface{}) {
	v.do(func(parent, currentValue interface{}, keyName string) interface{} {
		return append(currentValue.([]interface{}), value)
	})
}

func (v *Value) Merge(value interface{}) {
	v.do(func(parent, currentValue interface{}, keyName string) interface{} {
		current := currentValue.(map[string]interface{})
		other := value.(map[string]interface{})
		for key := range other {
			current[key] = other[key]
		}
		return current
	})
}

func (v *Value) Filter(selector func(interface{}) bool) {
	v.do(func(parent, currentValue interface{}, keyName string) interface{} {
		var slice []interface{}
		for _, obj := range currentValue.([]interface{}) {
			if selector(obj) {
				slice = append(slice, obj)
			}
		}
		return slice
	})
}

func (v *Value) Delete() {
	v.do(func(parent, currentValue interface{}, keyName string) interface{} {
		if parent == nil {
			return make(map[string]interface{})
		}

		if reflect.ValueOf(currentValue).Kind() == reflect.Slice {
			slice := currentValue.([]interface{})
			return slice[:len(slice)-1]
		} else {
			delete(parent.(map[string]interface{}), keyName)
			return nil
		}
	})
}

func (v *Value) do(fn func(parent, currentValue interface{}, keyName string) interface{}) interface{} {
	switch v.Key.(type) {
	case EmptyKey:
		newValue := fn(nil, *v.Parent, "")
		if newValue != nil {
			*v.Parent = newValue
		}
		return newValue
	case NormalKey:
		normalKey := v.Key.(NormalKey)
		newValue := fn(*v.Parent, (*v.Parent).(map[string]interface{})[normalKey.Name], normalKey.Name)
		if newValue != nil {
			(*v.Parent).(map[string]interface{})[normalKey.Name] = newValue
		}
		return newValue
	case ArrayKey:
		arrayKey := v.Key.(ArrayKey)
		var newValue interface{}

		if arrayKey.Name == "" {
			newValue = fn(*v.Parent, (*v.Parent).([]interface{})[arrayKey.Index], arrayKey.Name)
			if newValue != nil {
				(*v.Parent).([]interface{})[arrayKey.Index] = newValue
			}
		} else {
			newValue = fn(*v.Parent, (*v.Parent).(map[string]interface{})[arrayKey.Name].([]interface{})[arrayKey.Index], arrayKey.Name)
			if newValue != nil {
				(*v.Parent).(map[string]interface{})[arrayKey.Name].([]interface{})[arrayKey.Index] = newValue
			}
		}

		return newValue
	default:
		panic("invalid key")
	}
}

type Data struct {
	jsonObj *interface{}
}

func D(json string) *Data {
	var jsonObj interface{} = make(map[string]interface{})
	data := Data{jsonObj: &jsonObj}
	data.value("").Set(data.unmarshal(json))
	return &data
}

func (d *Data) GetJson(key string) string {
	bytes, err := json.Marshal(d.value(key).Get())
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (d *Data) GetString(key string) string {
	return d.value(key).Get().(string)
}

func (d *Data) GetInt(key string) int {
	return int(d.value(key).Get().(float64))
}

func (d *Data) GetFloat(key string) float64 {
	return d.value(key).Get().(float64)
}

func (d *Data) GetBool(key string) bool {
	return d.value(key).Get().(bool)
}

func (d *Data) SetJson(key string, value string) *Data {
	d.value(key).Set(d.unmarshal(value))
	return d
}

func (d *Data) SetString(key string, value string) *Data {
	d.value(key).Set(value)
	return d
}

func (d *Data) SetInt(key string, value int) *Data {
	d.value(key).Set(value)
	return d
}

func (d *Data) SetFloat(key string, value float64) *Data {
	d.value(key).Set(value)
	return d
}

func (d *Data) SetBool(key string, value bool) *Data {
	d.value(key).Set(value)
	return d
}

func (d *Data) Append(key string, json string) *Data {
	d.value(key).Append(d.unmarshal(json))
	return d
}

func (d *Data) AppendString(key string, value string) *Data {
	d.value(key).Append(value)
	return d
}

func (d *Data) AppendInt(key string, value int) *Data {
	d.value(key).Append(value)
	return d
}

func (d *Data) AppendFloat(key string, value float64) *Data {
	d.value(key).Append(value)
	return d
}

func (d *Data) Filter(key string, selector func(*Data) bool) *Data {
	d.value(key).Filter(func(obj interface{}) bool {
		return selector(&Data{jsonObj: &obj})
	})
	return d
}

func (d *Data) Merge(key string, json string) *Data {
	d.value(key).Merge(d.unmarshal(json))
	return d
}

func (d *Data) Delete(key string) *Data {
	d.value(key).Delete()
	return d
}

func (d *Data) value(key string) *Value {
	var value *Value

	for _, part := range strings.Split(key, ".") {
		if value == nil {
			value = &Value{
				Parent: d.jsonObj,
				Key:    d.parseKey(part),
			}
		} else {
			parent := value.Get()
			value = &Value{
				Parent: &parent,
				Key:    d.parseKey(part),
			}
		}
	}
	return value
}

func (d *Data) unmarshal(value string) interface{} {
	var obj interface{}
	if err := json.Unmarshal([]byte(value), &obj); err != nil {
		panic(err)
	}
	return obj
}

func (d *Data) parseKey(key string) interface{} {
	if key == "" {
		return EmptyKey{}
	}

	compiled := regexp.MustCompile("^(.*?)(?:\\[([0-9]+)])?$")
	match := compiled.FindStringSubmatch(key)

	if len(match) == 0 {
		panic(errors.New("invalid key"))
	}

	if len(match[2]) == 0 {
		return NormalKey{Name: match[1]}
	} else {
		value, err := strconv.Atoi(match[2])
		if err != nil {
			panic(err)
		}
		return ArrayKey{Name: match[1], Index: value}
	}
}
