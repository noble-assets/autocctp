// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/autocctp/v1/event.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// AccountRegistered is emitted whenever a new AutoCCTP account is registered.
type AccountRegistered struct {
	Address           string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	DestinationDomain uint32 `protobuf:"varint,2,opt,name=destination_domain,json=destinationDomain,proto3" json:"destination_domain,omitempty"`
	MintRecipient     []byte `protobuf:"bytes,3,opt,name=mint_recipient,json=mintRecipient,proto3" json:"mint_recipient,omitempty"`
	FallbackRecipient string `protobuf:"bytes,4,opt,name=fallback_recipient,json=fallbackRecipient,proto3" json:"fallback_recipient,omitempty"`
	DestinationCaller []byte `protobuf:"bytes,5,opt,name=destination_caller,json=destinationCaller,proto3" json:"destination_caller,omitempty"`
	Signerlessly      bool   `protobuf:"varint,6,opt,name=signerlessly,proto3" json:"signerlessly,omitempty"`
}

func (m *AccountRegistered) Reset()         { *m = AccountRegistered{} }
func (m *AccountRegistered) String() string { return proto.CompactTextString(m) }
func (*AccountRegistered) ProtoMessage()    {}
func (*AccountRegistered) Descriptor() ([]byte, []int) {
	return fileDescriptor_c4b6599cb121ef2c, []int{0}
}
func (m *AccountRegistered) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AccountRegistered) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AccountRegistered.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AccountRegistered) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountRegistered.Merge(m, src)
}
func (m *AccountRegistered) XXX_Size() int {
	return m.Size()
}
func (m *AccountRegistered) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountRegistered.DiscardUnknown(m)
}

var xxx_messageInfo_AccountRegistered proto.InternalMessageInfo

func (m *AccountRegistered) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *AccountRegistered) GetDestinationDomain() uint32 {
	if m != nil {
		return m.DestinationDomain
	}
	return 0
}

func (m *AccountRegistered) GetMintRecipient() []byte {
	if m != nil {
		return m.MintRecipient
	}
	return nil
}

func (m *AccountRegistered) GetFallbackRecipient() string {
	if m != nil {
		return m.FallbackRecipient
	}
	return ""
}

func (m *AccountRegistered) GetDestinationCaller() []byte {
	if m != nil {
		return m.DestinationCaller
	}
	return nil
}

func (m *AccountRegistered) GetSignerlessly() bool {
	if m != nil {
		return m.Signerlessly
	}
	return false
}

func init() {
	proto.RegisterType((*AccountRegistered)(nil), "noble.autocctp.v1.AccountRegistered")
}

func init() { proto.RegisterFile("noble/autocctp/v1/event.proto", fileDescriptor_c4b6599cb121ef2c) }

var fileDescriptor_c4b6599cb121ef2c = []byte{
	// 270 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x4f, 0x4a, 0x03, 0x31,
	0x18, 0xc5, 0x9b, 0xaa, 0x55, 0x43, 0x2b, 0x4c, 0x56, 0xd9, 0x18, 0x86, 0x82, 0x30, 0x0b, 0x9d,
	0xa1, 0x78, 0x02, 0xff, 0x9c, 0x20, 0x4b, 0x37, 0x25, 0x93, 0x7c, 0x96, 0x60, 0x9a, 0x0c, 0x49,
	0x3a, 0xd0, 0x03, 0xb8, 0xf7, 0x58, 0x2e, 0xbb, 0x74, 0x29, 0x33, 0x17, 0x91, 0x06, 0x46, 0xa6,
	0xcb, 0xef, 0xfd, 0x1e, 0xef, 0x7b, 0x3c, 0x7c, 0x6b, 0x5d, 0x6d, 0xa0, 0x12, 0xbb, 0xe8, 0xa4,
	0x8c, 0x4d, 0xd5, 0xae, 0x2a, 0x68, 0xc1, 0xc6, 0xb2, 0xf1, 0x2e, 0x3a, 0x92, 0x25, 0x5c, 0x0e,
	0xb8, 0x6c, 0x57, 0xcb, 0xcf, 0x29, 0xce, 0x9e, 0xa4, 0x74, 0x3b, 0x1b, 0x39, 0x6c, 0x74, 0x88,
	0xe0, 0x41, 0x11, 0x8a, 0x2f, 0x85, 0x52, 0x1e, 0x42, 0xa0, 0x28, 0x47, 0xc5, 0x35, 0x1f, 0x4e,
	0xf2, 0x80, 0x89, 0x82, 0x10, 0xb5, 0x15, 0x51, 0x3b, 0xbb, 0x56, 0x6e, 0x2b, 0xb4, 0xa5, 0xd3,
	0x1c, 0x15, 0x0b, 0x9e, 0x8d, 0xc8, 0x6b, 0x02, 0xe4, 0x0e, 0xdf, 0x6c, 0xb5, 0x8d, 0x6b, 0x0f,
	0x52, 0x37, 0x1a, 0x6c, 0xa4, 0x67, 0x39, 0x2a, 0xe6, 0x7c, 0x71, 0x54, 0xf9, 0x20, 0x1e, 0x53,
	0xdf, 0x85, 0x31, 0xb5, 0x90, 0x1f, 0x23, 0xeb, 0x79, 0x7a, 0x9d, 0x0d, 0xe4, 0xc4, 0x3e, 0x2e,
	0x21, 0x85, 0x31, 0xe0, 0xe9, 0x45, 0x4a, 0x1e, 0x97, 0x78, 0x49, 0x80, 0x2c, 0xf1, 0x3c, 0xe8,
	0x8d, 0x05, 0x6f, 0x20, 0x04, 0xb3, 0xa7, 0xb3, 0x1c, 0x15, 0x57, 0xfc, 0x44, 0x7b, 0xbe, 0xff,
	0xee, 0x18, 0x3a, 0x74, 0x0c, 0xfd, 0x76, 0x0c, 0x7d, 0xf5, 0x6c, 0x72, 0xe8, 0xd9, 0xe4, 0xa7,
	0x67, 0x93, 0x37, 0xf2, 0x3f, 0x97, 0x82, 0xb6, 0x8a, 0xfb, 0x06, 0x42, 0x3d, 0x4b, 0x7b, 0x3e,
	0xfe, 0x05, 0x00, 0x00, 0xff, 0xff, 0x6f, 0xe3, 0x07, 0xb8, 0x70, 0x01, 0x00, 0x00,
}

func (m *AccountRegistered) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AccountRegistered) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AccountRegistered) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Signerlessly {
		i--
		if m.Signerlessly {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if len(m.DestinationCaller) > 0 {
		i -= len(m.DestinationCaller)
		copy(dAtA[i:], m.DestinationCaller)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.DestinationCaller)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.FallbackRecipient) > 0 {
		i -= len(m.FallbackRecipient)
		copy(dAtA[i:], m.FallbackRecipient)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.FallbackRecipient)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.MintRecipient) > 0 {
		i -= len(m.MintRecipient)
		copy(dAtA[i:], m.MintRecipient)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.MintRecipient)))
		i--
		dAtA[i] = 0x1a
	}
	if m.DestinationDomain != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.DestinationDomain))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AccountRegistered) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.DestinationDomain != 0 {
		n += 1 + sovEvent(uint64(m.DestinationDomain))
	}
	l = len(m.MintRecipient)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.FallbackRecipient)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.DestinationCaller)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.Signerlessly {
		n += 2
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AccountRegistered) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AccountRegistered: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AccountRegistered: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationDomain", wireType)
			}
			m.DestinationDomain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DestinationDomain |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MintRecipient", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MintRecipient = append(m.MintRecipient[:0], dAtA[iNdEx:postIndex]...)
			if m.MintRecipient == nil {
				m.MintRecipient = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FallbackRecipient", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FallbackRecipient = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationCaller", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationCaller = append(m.DestinationCaller[:0], dAtA[iNdEx:postIndex]...)
			if m.DestinationCaller == nil {
				m.DestinationCaller = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signerlessly", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Signerlessly = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
