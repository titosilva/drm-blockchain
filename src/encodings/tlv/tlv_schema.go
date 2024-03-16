package tlv

import (
	errorutils "drm-blockchain/src/utils/error"
	"encoding/binary"
	"reflect"
)

type TLVSchema struct {
	Type uint8
	// Empty for struct and array
	// LittleEndian encoded
	LEBytes []uint8
	// Filled only for struct and array
	Children []TLVSchema
}

const (
	// Those numbers should NEVER change
	// This would cause previous conversions to become invalid
	TLVSchemaInt  uint8 = 1
	TLVSchemaByte uint8 = 2

	TLVSchemaArray  uint8 = 254
	TLVSchemaObject uint8 = 255
)

func BuildTLVSchema(st any) (schema TLVSchema, err error) {
	return buildTLVSchema(reflect.ValueOf(st))
}

func buildTLVSchema(stv reflect.Value) (schema TLVSchema, err error) {
	schema.Type, schema.LEBytes, err = getTLVSchemaTypeAndBytes(stv)

	if err != nil {
		return schema, err
	}

	if stv.Kind() == reflect.Struct {
		schema.Children = make([]TLVSchema, stv.NumField())

		for i := 0; i < stv.NumField(); i++ {
			f := stv.Field(i)
			fchild, _ := buildTLVSchema(f)
			schema.Children[i] = fchild
		}
	}

	return schema, err
}

func getTLVSchemaTypeAndBytes(t reflect.Value) (tlvType uint8, bs []byte, err error) {
	switch t.Kind() {
	case reflect.Int:
		tlvType, bs = getTLVInt32(t)
	case reflect.Struct:
		tlvType = TLVSchemaObject
	default:
		err = errorutils.Newf("Cannot convert kind %d to TLVSchemaType", t.Kind())
	}

	return tlvType, bs, err
}

func getTLVInt32(t reflect.Value) (uint8, []byte) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(t.Int()))
	return TLVSchemaInt, bs
}
