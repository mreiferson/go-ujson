package ujson

import (
	"errors"
	"strconv"
)

const (
	JT_NULL = iota
	JT_TRUE
	JT_FALSE
	JT_NUMERIC
	JT_UTF8
	JT_ARRAY
	JT_OBJECT
	JT_INVALID
)

type ObjectStore interface {
	NewObject() (interface{}, error)
	NewArray() (interface{}, error)
	ObjectAddKey(interface{}, interface{}, interface{}) error
	ArrayAddItem(interface{}, interface{}) error
	NewString([]byte) (interface{}, error)
	NewNumeric([]byte) (interface{}, error)
	NewTrue() (interface{}, error)
	NewFalse() (interface{}, error)
	NewNull() (interface{}, error)
}

type Decoder struct {
	store      ObjectStore
	data       []byte
	idx        int64
	lastTypeId int
}

func NewDecoder(store ObjectStore, data []byte) *Decoder {
	return &Decoder{
		store: store,
		data:  data,
	}
}

func (j *Decoder) Decode() (interface{}, error) {
	j.idx = 0
	j.lastTypeId = JT_INVALID
	return j.decodeAny()
}

func (j *Decoder) skipWhitespace() {
	maxLength := int64(len(j.data))

	for {
		if j.idx >= maxLength {
			return
		}

		switch j.data[j.idx] {
		case ' ', '\t', '\r', '\n':
			j.idx++
			continue
		}
		break
	}
}

func (j *Decoder) decodeAny() (interface{}, error) {
	for {
		c := j.data[j.idx]
		switch c {
		case '"':
			return j.decodeString()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			return j.decodeNumeric()
		case '[':
			return j.decodeArray()
		case '{':
			return j.decodeObject()
		case 't':
			return j.decodeTrue()
		case 'f':
			return j.decodeFalse()
		case 'n':
			return j.decodeNull()
		case ' ', '\t', '\r', '\n':
			j.idx++
			continue
		}
		break
	}

	return nil, errors.New("Expected object or value")
}

func (j *Decoder) decodeObject() (interface{}, error) {
	newObj, err := j.store.NewObject()
	if err != nil {
		return nil, err
	}

	maxLength := int64(len(j.data))
	j.idx++

	for {
		j.skipWhitespace()

		if j.idx >= maxLength {
			return nil, errors.New("Index exceeded maximum length of j.data")
		}

		if j.data[j.idx] == '}' {
			j.idx++
			return newObj, nil
		}

		j.lastTypeId = JT_INVALID
		itemName, err := j.decodeAny()
		if err != nil {
			return "", err
		}

		if j.lastTypeId != JT_UTF8 {
			return nil, errors.New("Key name of object must be 'string' when decoding 'object'")
		}

		j.skipWhitespace()

		nextChar := j.data[j.idx]
		j.idx++
		if nextChar != ':' {
			return nil, errors.New("No ':' found when decoding object value")
		}

		j.skipWhitespace()

		itemValue, err := j.decodeAny()
		if err != nil {
			return nil, err
		}

		err = j.store.ObjectAddKey(newObj, itemName, itemValue)
		if err != nil {
			return nil, err
		}

		j.skipWhitespace()

		nextChar = j.data[j.idx]
		j.idx++
		switch nextChar {
		case '}':
			return newObj, nil
		case ',':
			continue
		}
		break
	}

	return nil, errors.New("Unexpected character in found when decoding object value")
}

func (j *Decoder) decodeArray() (interface{}, error) {
	var length int

	newObj, err := j.store.NewArray()
	if err != nil {
		return nil, err
	}

	j.lastTypeId = JT_INVALID

	maxLength := int64(len(j.data))
	j.idx++

	for {
		j.skipWhitespace()

		if j.idx >= maxLength {
			return nil, errors.New("Index exceeded maximum length of j.data")
		}

		if j.data[j.idx] == ']' {
			if length == 0 {
				j.idx++
				return newObj, nil
			}
			return nil, errors.New(
				"Unexpected character found when decoding array value (" + string(length) + ")",
			)
		}

		itemValue, err := j.decodeAny()
		if err != nil {
			return nil, err
		}

		err = j.store.ArrayAddItem(newObj, itemValue)
		if err != nil {
			return nil, err
		}

		j.skipWhitespace()

		nextChar := j.data[j.idx]
		j.idx++
		switch nextChar {
		case ']':
			return newObj, nil
		case ',':
			length++
			continue
		}
		break
	}

	return nil, errors.New(
		"Unexpected character found when decoding array value (" + strconv.Itoa(length) + ")",
	)
}

const (
	SS_NORMAL = iota
	SS_ESC
)

func (j *Decoder) decodeString() (interface{}, error) {
	var c byte
	var escCount int

	j.lastTypeId = JT_INVALID
	j.idx++
	startIdx := j.idx
	state := SS_NORMAL
	maxLength := int64(len(j.data))

	for {
		if j.idx >= maxLength {
			return nil, errors.New("Index exceeded maximum length of j.data")
		}

		c = j.data[j.idx]
		j.idx++
		switch state {
		case SS_NORMAL:
			switch c {
			case '"':
				j.lastTypeId = JT_UTF8
				endIdx := j.idx - 1
				return j.store.NewString(j.data[startIdx:endIdx])
			case '\\':
				state = SS_ESC
				continue
			}
			if c >= 0x20 {
				continue
			}
		case SS_ESC:
			if escCount > 0 {
				if '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F' {
					escCount++
					if escCount > 4 {
						state = SS_NORMAL
						escCount = 0
					}
					continue
				}
				return nil, errors.New("Unexpected character " + string(c) + " in \\u hexadecimal character escape")
			}
			switch c {
			case 'b', 'f', 'n', 'r', 't', '\\', '/', '"':
				state = SS_NORMAL
				continue
			case 'u':
				escCount = 1
				continue
			}
		}
		break
	}
	return nil, errors.New("Unexpected character " + string(c) + " when decoding string value")
}

func (j *Decoder) decodeNumeric() (interface{}, error) {
	startIdx := j.idx
	for {
		c := j.data[j.idx]
		switch c {
		case '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			j.idx++
			continue
		}
		break
	}
	endIdx := j.idx
	j.lastTypeId = JT_NUMERIC
	return j.store.NewNumeric(j.data[startIdx:endIdx])
}

func (j *Decoder) decodeTrue() (interface{}, error) {
	j.idx++
	if j.data[j.idx] != 'r' {
		goto err
	}
	j.idx++
	if j.data[j.idx] != 'u' {
		goto err
	}
	j.idx++
	if j.data[j.idx] != 'e' {
		goto err
	}
	j.lastTypeId = JT_TRUE
	j.idx++
	return j.store.NewTrue()

err:
	return nil, errors.New("Unexpected character found when decoding 'true'")
}

func (j *Decoder) decodeFalse() (interface{}, error) {
	j.idx++
	if j.data[j.idx] != 'a' {
		goto err
	}
	j.idx++
	if j.data[j.idx] != 'l' {
		goto err
	}
	j.idx++
	if j.data[j.idx] != 's' {
		goto err
	}
	j.idx++
	if j.data[j.idx] != 'e' {
		goto err
	}
	j.lastTypeId = JT_FALSE
	j.idx++
	return j.store.NewFalse()

err:
	return nil, errors.New("Unexpected character found when decoding 'false'")
}

func (j *Decoder) decodeNull() (interface{}, error) {
	j.idx++
	if j.data[j.idx] != 'u' {
		goto err
	}
	j.idx++
	if j.data[j.idx] != 'l' {
		goto err
	}
	j.idx++
	if j.data[j.idx] != 'l' {
		goto err
	}
	j.lastTypeId = JT_NULL
	j.idx++
	return j.store.NewNull()

err:
	return nil, errors.New("Unexpected character found when decoding 'null'")
}
