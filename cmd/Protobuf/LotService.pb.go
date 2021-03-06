// Code generated by protoc-gen-go. DO NOT EDIT.
// source: LotService.proto

package lotspb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type LotID struct {
	LotID                int64    `protobuf:"varint,1,opt,name=lotID,proto3" json:"lotID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LotID) Reset()         { *m = LotID{} }
func (m *LotID) String() string { return proto.CompactTextString(m) }
func (*LotID) ProtoMessage()    {}
func (*LotID) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{0}
}

func (m *LotID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LotID.Unmarshal(m, b)
}
func (m *LotID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LotID.Marshal(b, m, deterministic)
}
func (m *LotID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LotID.Merge(m, src)
}
func (m *LotID) XXX_Size() int {
	return xxx_messageInfo_LotID.Size(m)
}
func (m *LotID) XXX_DiscardUnknown() {
	xxx_messageInfo_LotID.DiscardUnknown(m)
}

var xxx_messageInfo_LotID proto.InternalMessageInfo

func (m *LotID) GetLotID() int64 {
	if m != nil {
		return m.LotID
	}
	return 0
}

type UserLotID struct {
	LotID                int64    `protobuf:"varint,1,opt,name=lotID,proto3" json:"lotID,omitempty"`
	UserID               int64    `protobuf:"varint,2,opt,name=userID,proto3" json:"userID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserLotID) Reset()         { *m = UserLotID{} }
func (m *UserLotID) String() string { return proto.CompactTextString(m) }
func (*UserLotID) ProtoMessage()    {}
func (*UserLotID) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{1}
}

func (m *UserLotID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserLotID.Unmarshal(m, b)
}
func (m *UserLotID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserLotID.Marshal(b, m, deterministic)
}
func (m *UserLotID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserLotID.Merge(m, src)
}
func (m *UserLotID) XXX_Size() int {
	return xxx_messageInfo_UserLotID.Size(m)
}
func (m *UserLotID) XXX_DiscardUnknown() {
	xxx_messageInfo_UserLotID.DiscardUnknown(m)
}

var xxx_messageInfo_UserLotID proto.InternalMessageInfo

func (m *UserLotID) GetLotID() int64 {
	if m != nil {
		return m.LotID
	}
	return 0
}

func (m *UserLotID) GetUserID() int64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

type Status struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{2}
}

func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (m *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(m, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type LotCreateUpdate struct {
	Title                string               `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description          string               `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	MinPrice             float64              `protobuf:"fixed64,3,opt,name=min_price,json=minPrice,proto3" json:"min_price,omitempty"`
	PriceStep            float64              `protobuf:"fixed64,4,opt,name=price_step,json=priceStep,proto3" json:"price_step,omitempty"`
	EndAt                *timestamp.Timestamp `protobuf:"bytes,5,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
	Status               string               `protobuf:"bytes,6,opt,name=status,proto3" json:"status,omitempty"`
	UserID               int64                `protobuf:"varint,7,opt,name=userID,proto3" json:"userID,omitempty"`
	LotID                int64                `protobuf:"varint,8,opt,name=lotID,proto3" json:"lotID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *LotCreateUpdate) Reset()         { *m = LotCreateUpdate{} }
func (m *LotCreateUpdate) String() string { return proto.CompactTextString(m) }
func (*LotCreateUpdate) ProtoMessage()    {}
func (*LotCreateUpdate) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{3}
}

func (m *LotCreateUpdate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LotCreateUpdate.Unmarshal(m, b)
}
func (m *LotCreateUpdate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LotCreateUpdate.Marshal(b, m, deterministic)
}
func (m *LotCreateUpdate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LotCreateUpdate.Merge(m, src)
}
func (m *LotCreateUpdate) XXX_Size() int {
	return xxx_messageInfo_LotCreateUpdate.Size(m)
}
func (m *LotCreateUpdate) XXX_DiscardUnknown() {
	xxx_messageInfo_LotCreateUpdate.DiscardUnknown(m)
}

var xxx_messageInfo_LotCreateUpdate proto.InternalMessageInfo

func (m *LotCreateUpdate) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *LotCreateUpdate) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *LotCreateUpdate) GetMinPrice() float64 {
	if m != nil {
		return m.MinPrice
	}
	return 0
}

func (m *LotCreateUpdate) GetPriceStep() float64 {
	if m != nil {
		return m.PriceStep
	}
	return 0
}

func (m *LotCreateUpdate) GetEndAt() *timestamp.Timestamp {
	if m != nil {
		return m.EndAt
	}
	return nil
}

func (m *LotCreateUpdate) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *LotCreateUpdate) GetUserID() int64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

func (m *LotCreateUpdate) GetLotID() int64 {
	if m != nil {
		return m.LotID
	}
	return 0
}

type BuyLot struct {
	Price                float64  `protobuf:"fixed64,1,opt,name=price,proto3" json:"price,omitempty"`
	UserID               int64    `protobuf:"varint,2,opt,name=userID,proto3" json:"userID,omitempty"`
	LotID                int64    `protobuf:"varint,3,opt,name=lotID,proto3" json:"lotID,omitempty"`
	IsWS                 bool     `protobuf:"varint,4,opt,name=isWS,proto3" json:"isWS,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BuyLot) Reset()         { *m = BuyLot{} }
func (m *BuyLot) String() string { return proto.CompactTextString(m) }
func (*BuyLot) ProtoMessage()    {}
func (*BuyLot) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{4}
}

func (m *BuyLot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BuyLot.Unmarshal(m, b)
}
func (m *BuyLot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BuyLot.Marshal(b, m, deterministic)
}
func (m *BuyLot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BuyLot.Merge(m, src)
}
func (m *BuyLot) XXX_Size() int {
	return xxx_messageInfo_BuyLot.Size(m)
}
func (m *BuyLot) XXX_DiscardUnknown() {
	xxx_messageInfo_BuyLot.DiscardUnknown(m)
}

var xxx_messageInfo_BuyLot proto.InternalMessageInfo

func (m *BuyLot) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *BuyLot) GetUserID() int64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

func (m *BuyLot) GetLotID() int64 {
	if m != nil {
		return m.LotID
	}
	return 0
}

func (m *BuyLot) GetIsWS() bool {
	if m != nil {
		return m.IsWS
	}
	return false
}

type ShortUser struct {
	ID                   int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	FirstName            string   `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName             string   `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShortUser) Reset()         { *m = ShortUser{} }
func (m *ShortUser) String() string { return proto.CompactTextString(m) }
func (*ShortUser) ProtoMessage()    {}
func (*ShortUser) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{5}
}

func (m *ShortUser) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShortUser.Unmarshal(m, b)
}
func (m *ShortUser) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShortUser.Marshal(b, m, deterministic)
}
func (m *ShortUser) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShortUser.Merge(m, src)
}
func (m *ShortUser) XXX_Size() int {
	return xxx_messageInfo_ShortUser.Size(m)
}
func (m *ShortUser) XXX_DiscardUnknown() {
	xxx_messageInfo_ShortUser.DiscardUnknown(m)
}

var xxx_messageInfo_ShortUser proto.InternalMessageInfo

func (m *ShortUser) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *ShortUser) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *ShortUser) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

type Lot struct {
	ID                   int64                `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Title                string               `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description          string               `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	MinPrice             float64              `protobuf:"fixed64,4,opt,name=min_price,json=minPrice,proto3" json:"min_price,omitempty"`
	PriceStep            float64              `protobuf:"fixed64,5,opt,name=price_step,json=priceStep,proto3" json:"price_step,omitempty"`
	EndAt                *timestamp.Timestamp `protobuf:"bytes,6,opt,name=end_at,json=endAt,proto3" json:"end_at,omitempty"`
	Status               string               `protobuf:"bytes,7,opt,name=status,proto3" json:"status,omitempty"`
	BuyPrice             float64              `protobuf:"fixed64,8,opt,name=buy_price,json=buyPrice,proto3" json:"buy_price,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt            *timestamp.Timestamp `protobuf:"bytes,11,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
	Creator              *ShortUser           `protobuf:"bytes,12,opt,name=creator,proto3" json:"creator,omitempty"`
	Buyer                *ShortUser           `protobuf:"bytes,13,opt,name=buyer,proto3" json:"buyer,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Lot) Reset()         { *m = Lot{} }
func (m *Lot) String() string { return proto.CompactTextString(m) }
func (*Lot) ProtoMessage()    {}
func (*Lot) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{6}
}

func (m *Lot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Lot.Unmarshal(m, b)
}
func (m *Lot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Lot.Marshal(b, m, deterministic)
}
func (m *Lot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Lot.Merge(m, src)
}
func (m *Lot) XXX_Size() int {
	return xxx_messageInfo_Lot.Size(m)
}
func (m *Lot) XXX_DiscardUnknown() {
	xxx_messageInfo_Lot.DiscardUnknown(m)
}

var xxx_messageInfo_Lot proto.InternalMessageInfo

func (m *Lot) GetID() int64 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Lot) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Lot) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Lot) GetMinPrice() float64 {
	if m != nil {
		return m.MinPrice
	}
	return 0
}

func (m *Lot) GetPriceStep() float64 {
	if m != nil {
		return m.PriceStep
	}
	return 0
}

func (m *Lot) GetEndAt() *timestamp.Timestamp {
	if m != nil {
		return m.EndAt
	}
	return nil
}

func (m *Lot) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *Lot) GetBuyPrice() float64 {
	if m != nil {
		return m.BuyPrice
	}
	return 0
}

func (m *Lot) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Lot) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

func (m *Lot) GetDeletedAt() *timestamp.Timestamp {
	if m != nil {
		return m.DeletedAt
	}
	return nil
}

func (m *Lot) GetCreator() *ShortUser {
	if m != nil {
		return m.Creator
	}
	return nil
}

func (m *Lot) GetBuyer() *ShortUser {
	if m != nil {
		return m.Buyer
	}
	return nil
}

type Lots struct {
	Lots                 []*Lot   `protobuf:"bytes,1,rep,name=lots,proto3" json:"lots,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Lots) Reset()         { *m = Lots{} }
func (m *Lots) String() string { return proto.CompactTextString(m) }
func (*Lots) ProtoMessage()    {}
func (*Lots) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{7}
}

func (m *Lots) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Lots.Unmarshal(m, b)
}
func (m *Lots) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Lots.Marshal(b, m, deterministic)
}
func (m *Lots) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Lots.Merge(m, src)
}
func (m *Lots) XXX_Size() int {
	return xxx_messageInfo_Lots.Size(m)
}
func (m *Lots) XXX_DiscardUnknown() {
	xxx_messageInfo_Lots.DiscardUnknown(m)
}

var xxx_messageInfo_Lots proto.InternalMessageInfo

func (m *Lots) GetLots() []*Lot {
	if m != nil {
		return m.Lots
	}
	return nil
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{8}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type UserLots struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Type                 string   `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserLots) Reset()         { *m = UserLots{} }
func (m *UserLots) String() string { return proto.CompactTextString(m) }
func (*UserLots) ProtoMessage()    {}
func (*UserLots) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f82c460469dc35b, []int{9}
}

func (m *UserLots) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserLots.Unmarshal(m, b)
}
func (m *UserLots) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserLots.Marshal(b, m, deterministic)
}
func (m *UserLots) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserLots.Merge(m, src)
}
func (m *UserLots) XXX_Size() int {
	return xxx_messageInfo_UserLots.Size(m)
}
func (m *UserLots) XXX_DiscardUnknown() {
	xxx_messageInfo_UserLots.DiscardUnknown(m)
}

var xxx_messageInfo_UserLots proto.InternalMessageInfo

func (m *UserLots) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UserLots) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func init() {
	proto.RegisterType((*LotID)(nil), "login.LotID")
	proto.RegisterType((*UserLotID)(nil), "login.UserLotID")
	proto.RegisterType((*Status)(nil), "login.Status")
	proto.RegisterType((*LotCreateUpdate)(nil), "login.LotCreateUpdate")
	proto.RegisterType((*BuyLot)(nil), "login.BuyLot")
	proto.RegisterType((*ShortUser)(nil), "login.ShortUser")
	proto.RegisterType((*Lot)(nil), "login.Lot")
	proto.RegisterType((*Lots)(nil), "login.Lots")
	proto.RegisterType((*Empty)(nil), "login.Empty")
	proto.RegisterType((*UserLots)(nil), "login.UserLots")
}

func init() { proto.RegisterFile("LotService.proto", fileDescriptor_6f82c460469dc35b) }

var fileDescriptor_6f82c460469dc35b = []byte{
	// 680 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x5d, 0x6b, 0xdb, 0x4a,
	0x10, 0x8d, 0x2c, 0x4b, 0x96, 0xc6, 0xf9, 0x62, 0x09, 0x41, 0x38, 0xe4, 0x5e, 0xa3, 0x87, 0x60,
	0x02, 0x57, 0x21, 0xc9, 0x53, 0x1e, 0xed, 0xeb, 0x52, 0x0c, 0xa2, 0x14, 0xb9, 0x21, 0xd0, 0x17,
	0x57, 0xb6, 0x36, 0xae, 0xa8, 0xa4, 0x15, 0xda, 0x51, 0x41, 0xbf, 0xa3, 0xff, 0xb5, 0x8f, 0xa5,
	0xec, 0xae, 0x64, 0x2b, 0x09, 0x89, 0xc9, 0xdb, 0xee, 0xcc, 0x19, 0xcd, 0xd9, 0x73, 0x46, 0x03,
	0xc7, 0x3e, 0xc3, 0x39, 0x2d, 0x7e, 0xc6, 0x2b, 0xea, 0xe5, 0x05, 0x43, 0x46, 0x8c, 0x84, 0xad,
	0xe3, 0x6c, 0xf0, 0xef, 0x9a, 0xb1, 0x75, 0x42, 0xaf, 0x64, 0x70, 0x59, 0x3e, 0x5e, 0x61, 0x9c,
	0x52, 0x8e, 0x61, 0x9a, 0x2b, 0x9c, 0x7b, 0x0e, 0x86, 0xcf, 0x70, 0x36, 0x25, 0x27, 0x60, 0x24,
	0xe2, 0xe0, 0x68, 0x43, 0x6d, 0xa4, 0x07, 0xea, 0xe2, 0xde, 0x81, 0x7d, 0xcf, 0x69, 0xf1, 0x06,
	0x84, 0x9c, 0x82, 0x59, 0x72, 0x5a, 0xcc, 0xa6, 0x4e, 0x47, 0x86, 0xeb, 0x9b, 0x3b, 0x04, 0x73,
	0x8e, 0x21, 0x96, 0x5c, 0x20, 0xb8, 0x3c, 0xc9, 0x42, 0x3b, 0xa8, 0x6f, 0xee, 0x1f, 0x0d, 0x8e,
	0x7c, 0x86, 0xff, 0x17, 0x34, 0x44, 0x7a, 0x9f, 0x47, 0x21, 0x52, 0xd1, 0x03, 0x63, 0x4c, 0x68,
	0x0d, 0x55, 0x17, 0x32, 0x84, 0x7e, 0x44, 0xf9, 0xaa, 0x88, 0x73, 0x8c, 0x59, 0x26, 0x1b, 0xd9,
	0x41, 0x3b, 0x44, 0xce, 0xc0, 0x4e, 0xe3, 0x6c, 0x91, 0x17, 0xf1, 0x8a, 0x3a, 0xfa, 0x50, 0x1b,
	0x69, 0x81, 0x95, 0xc6, 0xd9, 0x67, 0x71, 0x27, 0xe7, 0x00, 0x32, 0xb1, 0xe0, 0x48, 0x73, 0xa7,
	0x2b, 0xb3, 0xb6, 0x8c, 0xcc, 0x91, 0xe6, 0xe4, 0x1a, 0x4c, 0x9a, 0x45, 0x8b, 0x10, 0x1d, 0x63,
	0xa8, 0x8d, 0xfa, 0x37, 0x03, 0x4f, 0xa9, 0xe6, 0x35, 0xaa, 0x79, 0x5f, 0x1a, 0xd5, 0x02, 0x83,
	0x66, 0xd1, 0x18, 0x5b, 0x4f, 0x32, 0xdb, 0x4f, 0x6a, 0x89, 0xd1, 0x6b, 0x8b, 0xb1, 0x95, 0xce,
	0x6a, 0xab, 0xfb, 0x0d, 0xcc, 0x49, 0x59, 0xf9, 0x0c, 0x45, 0x5e, 0x51, 0xd7, 0x24, 0x39, 0x75,
	0x79, 0x4d, 0xda, 0xed, 0xd7, 0xf4, 0xb6, 0x11, 0x04, 0xba, 0x31, 0x7f, 0x98, 0xcb, 0xf7, 0x59,
	0x81, 0x3c, 0xbb, 0x0f, 0x60, 0xcf, 0xbf, 0xb3, 0x02, 0x85, 0x89, 0xe4, 0x10, 0x3a, 0x1b, 0xf3,
	0x3a, 0xb3, 0xa9, 0x90, 0xe5, 0x31, 0x2e, 0x38, 0x2e, 0xb2, 0x30, 0xa5, 0xb5, 0xa8, 0xb6, 0x8c,
	0x7c, 0x0a, 0x53, 0x2a, 0x24, 0x4d, 0xc2, 0x26, 0xab, 0xcb, 0xac, 0x25, 0x02, 0x22, 0xe9, 0xfe,
	0xd6, 0x41, 0x17, 0xc4, 0x9f, 0x7f, 0x73, 0xe3, 0x5f, 0xe7, 0x0d, 0xff, 0xf4, 0x1d, 0xfe, 0x75,
	0xdf, 0xf4, 0xcf, 0x78, 0xdd, 0x3f, 0xf3, 0xfd, 0xfe, 0xf5, 0x9e, 0xf8, 0x77, 0x06, 0xf6, 0xb2,
	0xac, 0x6a, 0x1a, 0x96, 0xa2, 0xb1, 0x2c, 0x2b, 0x45, 0xe3, 0x0e, 0x60, 0x25, 0x67, 0x55, 0xf6,
	0xb2, 0x77, 0xf6, 0xb2, 0x6b, 0xf4, 0x18, 0x45, 0x69, 0x29, 0x07, 0x5c, 0x96, 0xc2, 0xee, 0xd2,
	0x1a, 0xad, 0x4a, 0x23, 0x9a, 0xd0, 0xba, 0xb4, 0xbf, 0xbb, 0xb4, 0x46, 0x8f, 0x91, 0x5c, 0x42,
	0x4f, 0x52, 0x60, 0x85, 0xb3, 0x2f, 0xeb, 0x8e, 0x3d, 0xb9, 0x16, 0xbc, 0xcd, 0x4c, 0x04, 0x0d,
	0x80, 0x5c, 0x80, 0xb1, 0x2c, 0x2b, 0x5a, 0x38, 0x07, 0xaf, 0x20, 0x55, 0xda, 0xbd, 0x80, 0xae,
	0xcf, 0x90, 0x93, 0x7f, 0xa0, 0x9b, 0x30, 0x14, 0xbf, 0xb4, 0x3e, 0xea, 0xdf, 0x40, 0x0d, 0xf7,
	0x19, 0x06, 0x32, 0xee, 0xf6, 0xc0, 0xf8, 0x90, 0xe6, 0x58, 0xb9, 0x1e, 0x58, 0xf5, 0x0a, 0xe1,
	0x62, 0x5a, 0xe2, 0xa8, 0x99, 0x96, 0x38, 0x12, 0x23, 0x8b, 0x55, 0xde, 0x0c, 0x8b, 0x3c, 0xdf,
	0xfc, 0xd2, 0xa1, 0x2f, 0xc0, 0xf5, 0x3e, 0x23, 0x97, 0x00, 0x1f, 0x29, 0x8e, 0x93, 0x44, 0x7e,
	0xe1, 0xa0, 0xe1, 0x25, 0xfd, 0x1a, 0xf4, 0xb7, 0x7d, 0xb9, 0xbb, 0x47, 0x46, 0x12, 0xeb, 0x33,
	0x9c, 0x54, 0xb3, 0x29, 0xd9, 0xdf, 0x26, 0x67, 0xd3, 0x41, 0x8b, 0xa2, 0xbb, 0x47, 0xae, 0xc1,
	0x56, 0x7b, 0x47, 0x0c, 0xf1, 0xe9, 0x36, 0xd5, 0x5e, 0x46, 0x2f, 0x4b, 0x54, 0xfc, 0x3d, 0x25,
	0x07, 0x53, 0xe9, 0x46, 0x43, 0xa9, 0x91, 0x75, 0xb3, 0x54, 0x07, 0x0d, 0x49, 0x25, 0xd6, 0x1e,
	0xf9, 0x0f, 0x0e, 0x37, 0x5d, 0xd4, 0xd8, 0x35, 0x4f, 0x56, 0xab, 0xe2, 0x45, 0x87, 0x23, 0xf5,
	0x62, 0x3e, 0xa9, 0xee, 0xd5, 0x76, 0x38, 0x7a, 0xda, 0xe3, 0x85, 0x48, 0xb7, 0x70, 0x32, 0x09,
	0x57, 0x3f, 0xd6, 0x05, 0x2b, 0xb3, 0x68, 0xd3, 0x8b, 0x93, 0x27, 0x4c, 0x9e, 0x15, 0x4d, 0xac,
	0xaf, 0xa6, 0xb0, 0x35, 0x5f, 0x2e, 0x4d, 0x39, 0x73, 0xb7, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff,
	0xad, 0x79, 0xb9, 0x0b, 0x74, 0x06, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LotsServiceClient is the client API for LotsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LotsServiceClient interface {
	GetAllLots(ctx context.Context, in *Status, opts ...grpc.CallOption) (*Lots, error)
	GetLotByID(ctx context.Context, in *LotID, opts ...grpc.CallOption) (*Lot, error)
	CreateLot(ctx context.Context, in *LotCreateUpdate, opts ...grpc.CallOption) (*Lot, error)
	UpdateLot(ctx context.Context, in *LotCreateUpdate, opts ...grpc.CallOption) (*Lot, error)
	DeleteLotByID(ctx context.Context, in *UserLotID, opts ...grpc.CallOption) (*Empty, error)
	UpdateLotPrice(ctx context.Context, in *BuyLot, opts ...grpc.CallOption) (*Lot, error)
	GetLotsByUserID(ctx context.Context, in *UserLots, opts ...grpc.CallOption) (*Lots, error)
	BackgroundUpdateLots(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Lots, error)
}

type lotsServiceClient struct {
	cc *grpc.ClientConn
}

func NewLotsServiceClient(cc *grpc.ClientConn) LotsServiceClient {
	return &lotsServiceClient{cc}
}

func (c *lotsServiceClient) GetAllLots(ctx context.Context, in *Status, opts ...grpc.CallOption) (*Lots, error) {
	out := new(Lots)
	err := c.cc.Invoke(ctx, "/login.LotsService/GetAllLots", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotsServiceClient) GetLotByID(ctx context.Context, in *LotID, opts ...grpc.CallOption) (*Lot, error) {
	out := new(Lot)
	err := c.cc.Invoke(ctx, "/login.LotsService/GetLotByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotsServiceClient) CreateLot(ctx context.Context, in *LotCreateUpdate, opts ...grpc.CallOption) (*Lot, error) {
	out := new(Lot)
	err := c.cc.Invoke(ctx, "/login.LotsService/CreateLot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotsServiceClient) UpdateLot(ctx context.Context, in *LotCreateUpdate, opts ...grpc.CallOption) (*Lot, error) {
	out := new(Lot)
	err := c.cc.Invoke(ctx, "/login.LotsService/UpdateLot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotsServiceClient) DeleteLotByID(ctx context.Context, in *UserLotID, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/login.LotsService/DeleteLotByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotsServiceClient) UpdateLotPrice(ctx context.Context, in *BuyLot, opts ...grpc.CallOption) (*Lot, error) {
	out := new(Lot)
	err := c.cc.Invoke(ctx, "/login.LotsService/UpdateLotPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotsServiceClient) GetLotsByUserID(ctx context.Context, in *UserLots, opts ...grpc.CallOption) (*Lots, error) {
	out := new(Lots)
	err := c.cc.Invoke(ctx, "/login.LotsService/GetLotsByUserID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotsServiceClient) BackgroundUpdateLots(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Lots, error) {
	out := new(Lots)
	err := c.cc.Invoke(ctx, "/login.LotsService/BackgroundUpdateLots", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LotsServiceServer is the server API for LotsService service.
type LotsServiceServer interface {
	GetAllLots(context.Context, *Status) (*Lots, error)
	GetLotByID(context.Context, *LotID) (*Lot, error)
	CreateLot(context.Context, *LotCreateUpdate) (*Lot, error)
	UpdateLot(context.Context, *LotCreateUpdate) (*Lot, error)
	DeleteLotByID(context.Context, *UserLotID) (*Empty, error)
	UpdateLotPrice(context.Context, *BuyLot) (*Lot, error)
	GetLotsByUserID(context.Context, *UserLots) (*Lots, error)
	BackgroundUpdateLots(context.Context, *Empty) (*Lots, error)
}

// UnimplementedLotsServiceServer can be embedded to have forward compatible implementations.
type UnimplementedLotsServiceServer struct {
}

func (*UnimplementedLotsServiceServer) GetAllLots(ctx context.Context, req *Status) (*Lots, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllLots not implemented")
}
func (*UnimplementedLotsServiceServer) GetLotByID(ctx context.Context, req *LotID) (*Lot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLotByID not implemented")
}
func (*UnimplementedLotsServiceServer) CreateLot(ctx context.Context, req *LotCreateUpdate) (*Lot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLot not implemented")
}
func (*UnimplementedLotsServiceServer) UpdateLot(ctx context.Context, req *LotCreateUpdate) (*Lot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLot not implemented")
}
func (*UnimplementedLotsServiceServer) DeleteLotByID(ctx context.Context, req *UserLotID) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteLotByID not implemented")
}
func (*UnimplementedLotsServiceServer) UpdateLotPrice(ctx context.Context, req *BuyLot) (*Lot, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLotPrice not implemented")
}
func (*UnimplementedLotsServiceServer) GetLotsByUserID(ctx context.Context, req *UserLots) (*Lots, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLotsByUserID not implemented")
}
func (*UnimplementedLotsServiceServer) BackgroundUpdateLots(ctx context.Context, req *Empty) (*Lots, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BackgroundUpdateLots not implemented")
}

func RegisterLotsServiceServer(s *grpc.Server, srv LotsServiceServer) {
	s.RegisterService(&_LotsService_serviceDesc, srv)
}

func _LotsService_GetAllLots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Status)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotsServiceServer).GetAllLots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.LotsService/GetAllLots",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotsServiceServer).GetAllLots(ctx, req.(*Status))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotsService_GetLotByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LotID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotsServiceServer).GetLotByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.LotsService/GetLotByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotsServiceServer).GetLotByID(ctx, req.(*LotID))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotsService_CreateLot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LotCreateUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotsServiceServer).CreateLot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.LotsService/CreateLot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotsServiceServer).CreateLot(ctx, req.(*LotCreateUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotsService_UpdateLot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LotCreateUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotsServiceServer).UpdateLot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.LotsService/UpdateLot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotsServiceServer).UpdateLot(ctx, req.(*LotCreateUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotsService_DeleteLotByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserLotID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotsServiceServer).DeleteLotByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.LotsService/DeleteLotByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotsServiceServer).DeleteLotByID(ctx, req.(*UserLotID))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotsService_UpdateLotPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuyLot)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotsServiceServer).UpdateLotPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.LotsService/UpdateLotPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotsServiceServer).UpdateLotPrice(ctx, req.(*BuyLot))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotsService_GetLotsByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserLots)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotsServiceServer).GetLotsByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.LotsService/GetLotsByUserID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotsServiceServer).GetLotsByUserID(ctx, req.(*UserLots))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotsService_BackgroundUpdateLots_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotsServiceServer).BackgroundUpdateLots(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.LotsService/BackgroundUpdateLots",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotsServiceServer).BackgroundUpdateLots(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _LotsService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "login.LotsService",
	HandlerType: (*LotsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllLots",
			Handler:    _LotsService_GetAllLots_Handler,
		},
		{
			MethodName: "GetLotByID",
			Handler:    _LotsService_GetLotByID_Handler,
		},
		{
			MethodName: "CreateLot",
			Handler:    _LotsService_CreateLot_Handler,
		},
		{
			MethodName: "UpdateLot",
			Handler:    _LotsService_UpdateLot_Handler,
		},
		{
			MethodName: "DeleteLotByID",
			Handler:    _LotsService_DeleteLotByID_Handler,
		},
		{
			MethodName: "UpdateLotPrice",
			Handler:    _LotsService_UpdateLotPrice_Handler,
		},
		{
			MethodName: "GetLotsByUserID",
			Handler:    _LotsService_GetLotsByUserID_Handler,
		},
		{
			MethodName: "BackgroundUpdateLots",
			Handler:    _LotsService_BackgroundUpdateLots_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "LotService.proto",
}
