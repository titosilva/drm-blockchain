package tlv_test

import (
	"drm-blockchain/src/encodings/tlv"
	"encoding/binary"
	"testing"
)

type STest struct {
	value int
}

type STest2 struct {
	value int
	s     STest
}

func Test__GetTLVSchema__Should__ReturnInt32Type__WhenProvidedAnInt(t *testing.T) {
	s, err := tlv.BuildTLVSchema(32)

	if err != nil {
		t.Errorf("Unexpected error while building TLV schema: %s", err.Error())
	}

	if s.Type != tlv.TLVSchemaInt {
		t.Errorf("Expected TLV schema type Int, but received other")
	}

	if binary.LittleEndian.Uint32(s.LEBytes) != 32 {
		t.Errorf("Expected TLV schema with value 32, but received other")
	}
}

func Test__GetTLVSchema__Should__ReturnStructTypeWithChildren__WhenProvidedAStruct(t *testing.T) {
	s, err := tlv.BuildTLVSchema(STest{
		value: 32,
	})

	if err != nil {
		t.Error("Unexpected error while building TLV schema")
	}

	if s.Type != tlv.TLVSchemaObject {
		t.Errorf("Expected TLV schema type Object, but received other")
	}

	if len(s.Children) != 1 {
		t.Errorf("Expected exactly one children, got %d", len(s.Children))
	}

	if s.Children[0].Type != tlv.TLVSchemaInt {
		t.Errorf("Expected TLV schema type Int on struct child, but received other")
	}

	if binary.LittleEndian.Uint32(s.Children[0].LEBytes) != 32 {
		t.Errorf("Expected TLV schema with value 32, but received other")
	}
}

func Test__GetTLVSchema__Should__ReturnStructTypeWithSubChildren__WhenProvidedAStruct(t *testing.T) {
	s, err := tlv.BuildTLVSchema(STest2{
		value: 32,
		s: STest{
			value: 123,
		},
	})

	if err != nil {
		t.Error("Unexpected error while building TLV schema")
	}

	if s.Type != tlv.TLVSchemaObject {
		t.Errorf("Expected TLV schema type Object, but received other")
	}

	if len(s.Children) != 2 {
		t.Errorf("Expected exactly one children, got %d", len(s.Children))
	}

	if s.Children[0].Type != tlv.TLVSchemaInt {
		t.Errorf("Expected TLV schema type Int on struct child, but received other")
	}

	if s.Children[1].Type != tlv.TLVSchemaObject {
		t.Errorf("Expected TLV schema type Object on struct child, but received other")
	}

	if s.Children[1].Children[0].Type != tlv.TLVSchemaInt {
		t.Errorf("Expected TLV schema type Int on struct sub child, but received other")
	}

	if binary.LittleEndian.Uint32(s.Children[0].LEBytes) != 32 {
		t.Errorf("Expected TLV schema with value 32, but received other")
	}

	if binary.LittleEndian.Uint32(s.Children[1].Children[0].LEBytes) != 123 {
		t.Errorf("Expected TLV schema with value 123, but received other")
	}
}

func Test__GetTLVSchema__Should__ReturnArrayStructTypeWithChildren__WhenProvidedAnArray(t *testing.T) {
	arr := [2]STest{
		{value: 1},
		{value: 4},
	}
	s, err := tlv.BuildTLVSchema(arr)

	if err != nil {
		t.Error("Unexpected error while building TLV schema")
	}

	if s.Type != tlv.TLVSchemaArray {
		t.Errorf("Expected TLV schema type Array, but received other")
	}

	if len(s.Children) != 2 {
		t.Errorf("Expected exactly one children, got %d", len(s.Children))
	}

	if s.Children[0].Type != tlv.TLVSchemaObject {
		t.Errorf("Expected TLV schema type Int on struct child, but received other")
	}

	if binary.LittleEndian.Uint32(s.Children[0].Children[0].LEBytes) != 1 {
		t.Errorf("Expected TLV schema with value 1, but received other")
	}
}
