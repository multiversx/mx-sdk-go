package codec

import (
	"errors"
	"io"

	"github.com/multiversx/mx-sdk-go/abi/values"
)

func (c *defaultCodec) encodeNestedOption(writer io.Writer, value values.OptionValue) error {
	if value.Value == nil {
		_, err := writer.Write([]byte{0})
		return err
	}

	_, err := writer.Write([]byte{1})
	if err != nil {
		return err
	}

	return c.doEncodeNested(writer, value.Value)
}

func (c *defaultCodec) decodeNestedOption(reader io.Reader, value *values.OptionValue) error {
	bytes, err := readBytesExactly(reader, 1)
	if err != nil {
		return err
	}

	if bytes[0] == 0 {
		value.Value = nil
		return nil
	}

	return c.doDecodeNested(reader, value.Value)
}

func (c *defaultCodec) encodeNestedList(writer io.Writer, value values.InputListValue) error {
	err := c.encodeLength(writer, uint32(len(value.Items)))
	if err != nil {
		return err
	}

	for _, item := range value.Items {
		err := c.doEncodeNested(writer, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultCodec) decodeNestedList(reader io.Reader, value *values.OutputListValue) error {
	if value.ItemCreator == nil {
		return errors.New("cannot deserialize list: item creator is nil")
	}

	length, err := c.decodeLength(reader)
	if err != nil {
		return err
	}

	value.Items = make([]any, 0, length)

	for i := uint32(0); i < length; i++ {
		newItem := value.ItemCreator()

		err := c.doDecodeNested(reader, newItem)
		if err != nil {
			return err
		}

		value.Items = append(value.Items, newItem)
	}

	return nil
}
