// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/config/accesslog/v2/file.proto

package v2

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/lyft/protoc-gen-validate/validate"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Custom configuration for an :ref:`AccessLog <envoy_api_msg_config.filter.accesslog.v2.AccessLog>`
// that writes log entries directly to a file. Configures the built-in *envoy.file_access_log*
// AccessLog.
type FileAccessLog struct {
	// A path to a local file to which to write the access log entries.
	Path string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	// Access log format. Envoy supports :ref:`custom access log formats
	// <config_access_log_format>` as well as a :ref:`default format
	// <config_access_log_default_format>`.
	Format               string   `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileAccessLog) Reset()         { *m = FileAccessLog{} }
func (m *FileAccessLog) String() string { return proto.CompactTextString(m) }
func (*FileAccessLog) ProtoMessage()    {}
func (*FileAccessLog) Descriptor() ([]byte, []int) {
	return fileDescriptor_bb42a04cfa71ce3c, []int{0}
}

func (m *FileAccessLog) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileAccessLog.Unmarshal(m, b)
}
func (m *FileAccessLog) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileAccessLog.Marshal(b, m, deterministic)
}
func (m *FileAccessLog) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileAccessLog.Merge(m, src)
}
func (m *FileAccessLog) XXX_Size() int {
	return xxx_messageInfo_FileAccessLog.Size(m)
}
func (m *FileAccessLog) XXX_DiscardUnknown() {
	xxx_messageInfo_FileAccessLog.DiscardUnknown(m)
}

var xxx_messageInfo_FileAccessLog proto.InternalMessageInfo

func (m *FileAccessLog) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *FileAccessLog) GetFormat() string {
	if m != nil {
		return m.Format
	}
	return ""
}

func init() {
	proto.RegisterType((*FileAccessLog)(nil), "envoy.config.accesslog.v2.FileAccessLog")
}

func init() {
	proto.RegisterFile("envoy/config/accesslog/v2/file.proto", fileDescriptor_bb42a04cfa71ce3c)
}

var fileDescriptor_bb42a04cfa71ce3c = []byte{
	// 160 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x49, 0xcd, 0x2b, 0xcb,
	0xaf, 0xd4, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x4f, 0x4c, 0x4e, 0x4e, 0x2d, 0x2e, 0xce,
	0xc9, 0x4f, 0xd7, 0x2f, 0x33, 0xd2, 0x4f, 0xcb, 0xcc, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x92, 0x04, 0xab, 0xd2, 0x83, 0xa8, 0xd2, 0x83, 0xab, 0xd2, 0x2b, 0x33, 0x92, 0x12, 0x2f,
	0x4b, 0xcc, 0xc9, 0x4c, 0x49, 0x2c, 0x49, 0xd5, 0x87, 0x31, 0x20, 0x7a, 0x94, 0xdc, 0xb8, 0x78,
	0xdd, 0x32, 0x73, 0x52, 0x1d, 0xc1, 0x8a, 0x7d, 0xf2, 0xd3, 0x85, 0x64, 0xb9, 0x58, 0x0a, 0x12,
	0x4b, 0x32, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x9d, 0x38, 0x77, 0xbd, 0x3c, 0xc0, 0xcc, 0x52,
	0xc4, 0xa4, 0xc0, 0x18, 0x04, 0x16, 0x16, 0x12, 0xe3, 0x62, 0x4b, 0xcb, 0x2f, 0xca, 0x4d, 0x2c,
	0x91, 0x60, 0x02, 0x29, 0x08, 0x82, 0xf2, 0x9c, 0x58, 0xa2, 0x98, 0xca, 0x8c, 0x92, 0xd8, 0xc0,
	0x86, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x86, 0xe4, 0x9b, 0xe1, 0xb0, 0x00, 0x00, 0x00,
}
