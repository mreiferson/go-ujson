package ujson

/*
#include <stdlib.h>
#include "ultrajson.h"
#include "ultrajsondec.c"
// #include "ultrajsonenc.c"

extern void go_objectAddKey(void *obj, void *name, void *value);

static void objectAddKey(JSOBJ obj, JSOBJ name, JSOBJ value)
{
	go_objectAddKey(obj, name, value);
}

extern void go_arrayAddItem(void *obj, void *value);

static void arrayAddItem(JSOBJ obj, JSOBJ value)
{
	go_arrayAddItem(obj, value);
}

extern JSOBJ go_newString(void *start, int32_t len, int32_t sz);

static JSOBJ newString(wchar_t *start, wchar_t *end)
{
	return go_newString(start, (end - start) * sizeof(wchar_t), sizeof(wchar_t));
}

extern JSOBJ go_newTrue(void);

static JSOBJ newTrue(void)
{ 
	return go_newTrue();
}

extern JSOBJ go_newFalse(void);

static JSOBJ newFalse(void)
{
	return go_newFalse();
}

extern JSOBJ go_newNull(void);

static JSOBJ newNull(void)
{
	return go_newNull();
}

extern JSOBJ go_newObject(void);

static JSOBJ newObject(void)
{
	return go_newObject();
}

extern JSOBJ go_newArray(void);

static JSOBJ newArray(void)
{
	return go_newArray();
}

extern JSOBJ go_newInteger(int32_t value);

static JSOBJ newInteger(JSINT32 value)
{
	return go_newInteger((int32_t)value);
}

extern JSOBJ go_newLong(int64_t value);

static JSOBJ newLong(JSINT64 value)
{
	return go_newLong((int64_t)value);
}

extern JSOBJ go_newDouble(double value);

static JSOBJ newDouble(double value)
{
	return go_newDouble(value);
}

static void releaseObject(JSOBJ obj) {}

extern void go_finalize(void *obj);

static void *decodeString(char *str, size_t len)
{
	void *ret;

	JSONObjectDecoder decoder = {
		newString,
		objectAddKey,
		arrayAddItem,
		newTrue,
		newFalse,
		newNull,
		newObject,
		newArray,
		newInteger,
		newLong,
		newDouble,
		releaseObject,
		NULL,
		NULL,
		NULL
	};

	decoder.preciseFloat = 0;
	decoder.errorStr = NULL;
	decoder.errorOffset = NULL;

	ret = JSON_DecodeObject(&decoder, str, len);
	if (decoder.errorStr) {
		// TODO: figure out a way to return error string
		return NULL;
	}

	return ret;
}

*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"unsafe"
)

func Version() string {
	return "0.1.0"
}

func Unmarshal(d []byte) (map[string]interface{}, error) {
	cData := (*C.char)(unsafe.Pointer(&d[0]))
	ret := C.decodeString(cData, C.size_t(len(d)))
	if ret == nil {
		return nil, errors.New("failed to decode JSON")
	}
	return *(*map[string]interface{})(ret), nil
}

//export go_objectAddKey
func go_objectAddKey(obj unsafe.Pointer, name unsafe.Pointer, value unsafe.Pointer) {
	m := *(*map[string]interface{})(obj)
	niface := *(*interface{})(name)
	key := niface.(string)
	viface := *(*interface{})(value)
	m[key] = viface
	log.Printf("objectAddKey: %p, %p, %p %s=%v", obj, name, value, key, viface)
}

//export go_arrayAddItem
func go_arrayAddItem(obj unsafe.Pointer, value unsafe.Pointer) {
	iface := *(*interface{})(obj)
	sa := iface.(*staticArray)
	sa.sl = append(sa.sl, value)
	log.Printf("arrayAddItem: %p %p %v", obj, value, sa)
}

//export go_newString
func go_newString(start unsafe.Pointer, length C.int, sz C.int) unsafe.Pointer {
	r := make([]rune, length / sz)
	b := C.GoBytes(start, length)
	err := binary.Read(bytes.NewBuffer(b), binary.LittleEndian, r)
	if err != nil {
		log.Printf("ERROR: %s", err.Error())
	}
	s := string(r)
	log.Printf("newString: %s", s)
	iface := (interface{})(s)
	return unsafe.Pointer(&iface)
}

//export go_newTrue
func go_newTrue() unsafe.Pointer {
	b := (interface{})(true)
	return unsafe.Pointer(&b)
}

//export go_newFalse
func go_newFalse() unsafe.Pointer {
	b := (interface{})(false)
	return unsafe.Pointer(&b)
}

//export go_newNull
func go_newNull() unsafe.Pointer {
	return nil
}

//export go_newObject
func go_newObject() unsafe.Pointer {
	m := make(map[string]interface{})
	p := unsafe.Pointer(&m)
	log.Printf("newObject: %+v %p", m, p)
	return p
}

type staticArray struct {
	sl []interface{}
}

//export go_newArray
func go_newArray() unsafe.Pointer {
	sa := &staticArray{make([]interface{}, 0)}
	iface := (interface{})(sa)
	p := unsafe.Pointer(&iface)
	log.Printf("newArray: %p", p)
	return p
}

//export go_newInteger
func go_newInteger(v int32) unsafe.Pointer {
	log.Printf("newInteger: %d", v)
	iface := (interface{})(v)
	return unsafe.Pointer(&iface)
}

//export go_newLong
func go_newLong(v int64) unsafe.Pointer {
	log.Printf("newLong: %d", v)
	iface := (interface{})(v)
	return unsafe.Pointer(&iface)
}

//export go_newDouble
func go_newDouble(v float64) unsafe.Pointer {
	log.Printf("newDouble: %f", v)
	iface := (interface{})(v)
	return unsafe.Pointer(&iface)
}
