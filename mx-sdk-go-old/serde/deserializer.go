package serde

import (
	"encoding/binary"
	"errors"
	"math/big"
	"reflect"
)

type deserializer struct{}

// NewDeserializer will create a new instance of the deserializer.
func NewDeserializer() *deserializer {
	return &deserializer{}
}

// CreatePrimitiveDataType deserializes the buffer and populates the received object
func (des *deserializer) CreatePrimitiveDataType(obj interface{}, buff []byte) error {
	reflectedValue, err := des.getReflectedValue(obj)
	if err != nil {
		return err
	}

	if reflectedValue.Kind() == reflect.String || reflectedValue.Type() == reflect.ValueOf(big.Int{}).Type() {
		length := make([]byte, 4)
		binary.BigEndian.PutUint32(length[:], uint32(len(buff)))
		buff = append(length[:], buff...)
	}
	buffer := NewSourceBuffer(buff)

	valueFromBuffer, eof := des.getNextValueFromBuffer(buffer, reflectedValue)
	if eof {
		return errors.New("empty buffer")
	}
	err = des.setValue(reflectedValue, valueFromBuffer)
	if err != nil {
		return err
	}

	return nil
}

// CreateStruct deserialize the buffer and populate the fields of the received object
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

		valueFromBuffer, eof := des.getNextValueFromBuffer(buffer, reflectedValueField)
		if eof {
			return buffer.Pos(), errors.New("empty buffer")
		}

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

	err := des.setValue(structFieldValue, value)
	if err != nil {
		return err
	}
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
	switch v.Kind() {
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
		return buffer.OffBytes(), false
	case reflect.Ptr:
		return buffer.OffBytes(), false
	default:
		return nil, true
	}
}

func (des *deserializer) setValue(obj reflect.Value, value interface{}) error {
	if !obj.IsValid() {
		return errors.New("invalid object")
	}
	if !obj.CanSet() {
		return errors.New("cannot set field value")
	}

	structFieldType := obj.Type()
	val := reflect.ValueOf(value)
	valType := val.Type()
	if structFieldType != valType {
		return errors.New("provided value type didn't match obj field type")
	}

	obj.Set(val)

	return nil
}
