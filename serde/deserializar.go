package serde

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
)

type deserializer struct{}

//NewDeserializer will create a new instance of the deserializer.
func NewDeserializer() *deserializer {
	return &deserializer{}
}

//CreateStruct deserialize the buffer and populate the fields of the received object
func (des *deserializer) CreateStruct(obj interface{}, buff []byte) (uint64, error) {
	buffer := NewSourceBuffer(buff)

	reflectedValue, err := des.getReflectedValue(obj)
	if err != nil {
		return buffer.Pos(), err
	}
	if reflectedValue.Kind() != reflect.Struct || reflectedValue.Type() == reflect.TypeOf(big.Int{}) {
		return buffer.Pos(), errors.New("invalid type")
	}

	usedBytes, err := des.setFields(reflectedValue, buffer)

	if err != nil {
		return usedBytes, err
	}

	return buffer.Pos(), nil
}

func (des *deserializer) setFields(reflectedValue reflect.Value, buffer *SourceBuffer) (uint64, error) {
	for fieldIndex := 0; fieldIndex < reflectedValue.NumField(); fieldIndex++ {
		reflectedValueField := reflectedValue.Field(fieldIndex)
		reflectedValueFieldName := reflectedValue.Type().Field(fieldIndex).Name

		valueFromBuffer, _ := des.getNextValueFromBuffer(buffer, reflectedValueField)

		if (reflectedValueField.Kind() == reflect.Struct || reflectedValueField.Kind() == reflect.Ptr) &&
			reflect.TypeOf(valueFromBuffer) == reflect.TypeOf([]byte{}) {
			usedBytes, err := des.CreateStruct(reflectedValueField, valueFromBuffer.([]byte))
			if err != nil {
				return usedBytes, err
			}
			buffer.Skip(usedBytes)
			continue
		}

		err := des.setField(reflectedValue, reflectedValueFieldName, valueFromBuffer)
		if err != nil {
			return buffer.Pos(), errors.New("can't set value for field")
		}
	}
	return buffer.Pos(), nil
}

func (des *deserializer) setField(structValue reflect.Value, name string, value interface{}) error {
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return errors.New(fmt.Sprintf("No such field: %s in obj", name))
	}

	if !structFieldValue.CanSet() {
		return errors.New(fmt.Sprintf("Cannot set %s field value", name))
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	valType := val.Type()
	if structFieldType != valType {
		return errors.New("provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func (des *deserializer) getReflectedValue(obj interface{}) (value reflect.Value, err error) {
	if _, ok := obj.(reflect.Value); ok {
		value = obj.(reflect.Value)
	} else {
		value = reflect.ValueOf(obj)
	}
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	if !value.CanSet() {
		err = errors.New("invalid type")
	}

	return
}

func (des *deserializer) getNextValueFromBuffer(buffer *SourceBuffer, v reflect.Value) (interface{}, bool) {
	t := v.Kind()
	switch t {
	case reflect.Int8:
		value, eof := buffer.NextUint8()
		return int8(value), eof
	case reflect.Int16:
		value, eof := buffer.NextInt16()
		return value, eof
	case reflect.Int32:
		value, eof := buffer.NextInt32()
		return value, eof
	case reflect.Int64:
		value, eof := buffer.NextInt64()
		return value, eof
	case reflect.Uint16:
		value, eof := buffer.NextUint16()
		return value, eof
	case reflect.Uint32:
		value, eof := buffer.NextUint32()
		return value, eof
	case reflect.Uint64:
		value, eof := buffer.NextUint64()
		return value, eof
	case reflect.Uint8:
		value, eof := buffer.NextUint8()
		return value, eof
	case reflect.Bool:
		value, eof := buffer.NextBool()
		return value, eof
	case reflect.String:
		value, eof := buffer.NextVarBytes()
		return string(value), eof
	case reflect.Struct:
		if v.Type() == reflect.ValueOf(big.Int{}).Type() {
			buff, eof := buffer.NextVarBytes()
			return *big.NewInt(0).SetBytes(buff), eof
		}
		return buffer.OffBytes(), true
	case reflect.Ptr:
		return buffer.OffBytes(), true
	default:
		return nil, true
	}
}
