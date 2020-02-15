// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/type/calendar_period.proto

package calendarperiod

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// A `CalendarPeriod` represents the abstract concept of a time period that has
// a canonical start. Grammatically, "the start of the current
// `CalendarPeriod`." All calendar times begin at midnight UTC.
type CalendarPeriod int32

const (
	// Undefined period, raises an error.
	CalendarPeriod_CALENDAR_PERIOD_UNSPECIFIED CalendarPeriod = 0
	// A day.
	CalendarPeriod_DAY CalendarPeriod = 1
	// A week. Weeks begin on Monday, following
	// [ISO 8601](https://en.wikipedia.org/wiki/ISO_week_date).
	CalendarPeriod_WEEK CalendarPeriod = 2
	// A fortnight. The first calendar fortnight of the year begins at the start
	// of week 1 according to
	// [ISO 8601](https://en.wikipedia.org/wiki/ISO_week_date).
	CalendarPeriod_FORTNIGHT CalendarPeriod = 3
	// A month.
	CalendarPeriod_MONTH CalendarPeriod = 4
	// A quarter. Quarters start on dates 1-Jan, 1-Apr, 1-Jul, and 1-Oct of each
	// year.
	CalendarPeriod_QUARTER CalendarPeriod = 5
	// A half-year. Half-years start on dates 1-Jan and 1-Jul.
	CalendarPeriod_HALF CalendarPeriod = 6
	// A year.
	CalendarPeriod_YEAR CalendarPeriod = 7
)

var CalendarPeriod_name = map[int32]string{
	0: "CALENDAR_PERIOD_UNSPECIFIED",
	1: "DAY",
	2: "WEEK",
	3: "FORTNIGHT",
	4: "MONTH",
	5: "QUARTER",
	6: "HALF",
	7: "YEAR",
}

var CalendarPeriod_value = map[string]int32{
	"CALENDAR_PERIOD_UNSPECIFIED": 0,
	"DAY":                         1,
	"WEEK":                        2,
	"FORTNIGHT":                   3,
	"MONTH":                       4,
	"QUARTER":                     5,
	"HALF":                        6,
	"YEAR":                        7,
}

func (x CalendarPeriod) String() string {
	return proto.EnumName(CalendarPeriod_name, int32(x))
}

func (CalendarPeriod) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_52aec33199a3da0e, []int{0}
}

func init() {
	proto.RegisterEnum("google.type.CalendarPeriod", CalendarPeriod_name, CalendarPeriod_value)
}

func init() { proto.RegisterFile("google/type/calendar_period.proto", fileDescriptor_52aec33199a3da0e) }

var fileDescriptor_52aec33199a3da0e = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x8f, 0xb1, 0x4f, 0x83, 0x40,
	0x14, 0x87, 0x6d, 0x69, 0x8b, 0x7d, 0x8d, 0x7a, 0x39, 0x47, 0x07, 0xe3, 0xea, 0x00, 0x83, 0xa3,
	0xd3, 0x15, 0x8e, 0x42, 0xac, 0x70, 0x9e, 0xd7, 0x98, 0xba, 0x10, 0x6c, 0x2f, 0x97, 0x26, 0xc8,
	0x23, 0xd8, 0x41, 0x27, 0xff, 0x17, 0xff, 0x52, 0x73, 0xc0, 0x50, 0xb6, 0xbb, 0xbc, 0xef, 0x97,
	0x7c, 0x1f, 0xdc, 0x19, 0x44, 0x53, 0x6a, 0xff, 0xf8, 0x53, 0x6b, 0x7f, 0x57, 0x94, 0xba, 0xda,
	0x17, 0x4d, 0x5e, 0xeb, 0xe6, 0x80, 0x7b, 0xaf, 0x6e, 0xf0, 0x88, 0x74, 0xd1, 0x21, 0x9e, 0x45,
	0xee, 0x7f, 0xe1, 0x32, 0xe8, 0x29, 0xd1, 0x42, 0xf4, 0x16, 0x6e, 0x02, 0xb6, 0xe6, 0x69, 0xc8,
	0x64, 0x2e, 0xb8, 0x4c, 0xb2, 0x30, 0xdf, 0xa4, 0xaf, 0x82, 0x07, 0x49, 0x94, 0xf0, 0x90, 0x9c,
	0x51, 0x17, 0x9c, 0x90, 0x6d, 0xc9, 0x88, 0x9e, 0xc3, 0xe4, 0x8d, 0xf3, 0x27, 0x32, 0xa6, 0x17,
	0x30, 0x8f, 0x32, 0xa9, 0xd2, 0x64, 0x15, 0x2b, 0xe2, 0xd0, 0x39, 0x4c, 0x9f, 0xb3, 0x54, 0xc5,
	0x64, 0x42, 0x17, 0xe0, 0xbe, 0x6c, 0x98, 0x54, 0x5c, 0x92, 0xa9, 0x1d, 0xc4, 0x6c, 0x1d, 0x91,
	0x99, 0x7d, 0x6d, 0x39, 0x93, 0xc4, 0x5d, 0x7e, 0xc3, 0xd5, 0x0e, 0x3f, 0xbd, 0x13, 0xa7, 0xe5,
	0xf5, 0xd0, 0x48, 0x58, 0x6b, 0x31, 0x7a, 0x8f, 0x7b, 0xc6, 0x60, 0x59, 0x54, 0xc6, 0xc3, 0xc6,
	0xf8, 0x46, 0x57, 0x6d, 0x93, 0xdf, 0x9d, 0x8a, 0xfa, 0xf0, 0x35, 0x2c, 0xef, 0xc2, 0x1f, 0x87,
	0xdf, 0xbf, 0xb1, 0xb3, 0x52, 0xe2, 0x63, 0xd6, 0x4e, 0x1f, 0xfe, 0x03, 0x00, 0x00, 0xff, 0xff,
	0x91, 0x18, 0xaa, 0x3f, 0x33, 0x01, 0x00, 0x00,
}
