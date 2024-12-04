// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: inventory/inventory.proto

package inventory

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UnitRequest int32

const (
	UnitRequest_i          UnitRequest = 0
	UnitRequest_UNIT_ITEMS UnitRequest = 1
	UnitRequest_UNIT_G     UnitRequest = 2
	UnitRequest_UNIT_KG    UnitRequest = 3
	UnitRequest_UNIT_ML    UnitRequest = 4
	UnitRequest_UNIT_L     UnitRequest = 5
	UnitRequest_UNIT_TSP   UnitRequest = 6
	UnitRequest_UNIT_TBSP  UnitRequest = 7
	UnitRequest_UNIT_CUP   UnitRequest = 8
)

// Enum value maps for UnitRequest.
var (
	UnitRequest_name = map[int32]string{
		0: "i",
		1: "UNIT_ITEMS",
		2: "UNIT_G",
		3: "UNIT_KG",
		4: "UNIT_ML",
		5: "UNIT_L",
		6: "UNIT_TSP",
		7: "UNIT_TBSP",
		8: "UNIT_CUP",
	}
	UnitRequest_value = map[string]int32{
		"i":          0,
		"UNIT_ITEMS": 1,
		"UNIT_G":     2,
		"UNIT_KG":    3,
		"UNIT_ML":    4,
		"UNIT_L":     5,
		"UNIT_TSP":   6,
		"UNIT_TBSP":  7,
		"UNIT_CUP":   8,
	}
)

func (x UnitRequest) Enum() *UnitRequest {
	p := new(UnitRequest)
	*p = x
	return p
}

func (x UnitRequest) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UnitRequest) Descriptor() protoreflect.EnumDescriptor {
	return file_inventory_inventory_proto_enumTypes[0].Descriptor()
}

func (UnitRequest) Type() protoreflect.EnumType {
	return &file_inventory_inventory_proto_enumTypes[0]
}

func (x UnitRequest) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UnitRequest.Descriptor instead.
func (UnitRequest) EnumDescriptor() ([]byte, []int) {
	return file_inventory_inventory_proto_rawDescGZIP(), []int{0}
}

type PostIngredientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId string  `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
	Name   string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Amount float64 `protobuf:"fixed64,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Unit   string  `protobuf:"bytes,5,opt,name=unit,proto3" json:"unit,omitempty"`
}

func (x *PostIngredientRequest) Reset() {
	*x = PostIngredientRequest{}
	mi := &file_inventory_inventory_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PostIngredientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostIngredientRequest) ProtoMessage() {}

func (x *PostIngredientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_inventory_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostIngredientRequest.ProtoReflect.Descriptor instead.
func (*PostIngredientRequest) Descriptor() ([]byte, []int) {
	return file_inventory_inventory_proto_rawDescGZIP(), []int{0}
}

func (x *PostIngredientRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PostIngredientRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *PostIngredientRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PostIngredientRequest) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *PostIngredientRequest) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

type DeleteIngredientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IngredientId string `protobuf:"bytes,1,opt,name=ingredientId,proto3" json:"ingredientId,omitempty"`
	UserId       string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *DeleteIngredientRequest) Reset() {
	*x = DeleteIngredientRequest{}
	mi := &file_inventory_inventory_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteIngredientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteIngredientRequest) ProtoMessage() {}

func (x *DeleteIngredientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_inventory_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteIngredientRequest.ProtoReflect.Descriptor instead.
func (*DeleteIngredientRequest) Descriptor() ([]byte, []int) {
	return file_inventory_inventory_proto_rawDescGZIP(), []int{1}
}

func (x *DeleteIngredientRequest) GetIngredientId() string {
	if x != nil {
		return x.IngredientId
	}
	return ""
}

func (x *DeleteIngredientRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

// IngredientResponse is the response of a user inventory
type IngredientResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name      string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Amount    float64                `protobuf:"fixed64,5,opt,name=amount,proto3" json:"amount,omitempty"`
	Unit      string                 `protobuf:"bytes,6,opt,name=unit,proto3" json:"unit,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

func (x *IngredientResponse) Reset() {
	*x = IngredientResponse{}
	mi := &file_inventory_inventory_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *IngredientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IngredientResponse) ProtoMessage() {}

func (x *IngredientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_inventory_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IngredientResponse.ProtoReflect.Descriptor instead.
func (*IngredientResponse) Descriptor() ([]byte, []int) {
	return file_inventory_inventory_proto_rawDescGZIP(), []int{2}
}

func (x *IngredientResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *IngredientResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *IngredientResponse) GetAmount() float64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *IngredientResponse) GetUnit() string {
	if x != nil {
		return x.Unit
	}
	return ""
}

func (x *IngredientResponse) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *IngredientResponse) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

// User inventory is the inventory of a user containing an array of IngredientResponse
type GetUserInventoryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserInventory []*IngredientResponse `protobuf:"bytes,1,rep,name=userInventory,proto3" json:"userInventory,omitempty"`
}

func (x *GetUserInventoryResponse) Reset() {
	*x = GetUserInventoryResponse{}
	mi := &file_inventory_inventory_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserInventoryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserInventoryResponse) ProtoMessage() {}

func (x *GetUserInventoryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_inventory_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserInventoryResponse.ProtoReflect.Descriptor instead.
func (*GetUserInventoryResponse) Descriptor() ([]byte, []int) {
	return file_inventory_inventory_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserInventoryResponse) GetUserInventory() []*IngredientResponse {
	if x != nil {
		return x.UserInventory
	}
	return nil
}

type GetInventoryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetInventoryRequest) Reset() {
	*x = GetInventoryRequest{}
	mi := &file_inventory_inventory_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetInventoryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetInventoryRequest) ProtoMessage() {}

func (x *GetInventoryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_inventory_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetInventoryRequest.ProtoReflect.Descriptor instead.
func (*GetInventoryRequest) Descriptor() ([]byte, []int) {
	return file_inventory_inventory_proto_rawDescGZIP(), []int{4}
}

func (x *GetInventoryRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type GetIngredientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IngredientId string `protobuf:"bytes,1,opt,name=ingredientId,proto3" json:"ingredientId,omitempty"`
	UserId       string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *GetIngredientRequest) Reset() {
	*x = GetIngredientRequest{}
	mi := &file_inventory_inventory_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetIngredientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetIngredientRequest) ProtoMessage() {}

func (x *GetIngredientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_inventory_inventory_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetIngredientRequest.ProtoReflect.Descriptor instead.
func (*GetIngredientRequest) Descriptor() ([]byte, []int) {
	return file_inventory_inventory_proto_rawDescGZIP(), []int{5}
}

func (x *GetIngredientRequest) GetIngredientId() string {
	if x != nil {
		return x.IngredientId
	}
	return ""
}

func (x *GetIngredientRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

var File_inventory_inventory_proto protoreflect.FileDescriptor

var file_inventory_inventory_proto_rawDesc = []byte{
	0x0a, 0x19, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x69, 0x6e, 0x76, 0x65,
	0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x69, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7f, 0x0a, 0x15, 0x50, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x67, 0x72,
	0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x75, 0x6e, 0x69, 0x74, 0x22, 0x55, 0x0a, 0x17, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49,
	0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x22, 0x0a, 0x0c, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65,
	0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0xd8, 0x01, 0x0a,
	0x12, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x75, 0x6e, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75,
	0x6e, 0x69, 0x74, 0x12, 0x38, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x38, 0x0a,
	0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x5f, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x69, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x22, 0x2d, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x49,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x52, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x49, 0x6e,
	0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x22, 0x0a, 0x0c, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e,
	0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x2a, 0x81, 0x01, 0x0a, 0x0b,
	0x55, 0x6e, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x05, 0x0a, 0x01, 0x69,
	0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x49, 0x54, 0x45, 0x4d, 0x53,
	0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x47, 0x10, 0x02, 0x12, 0x0b,
	0x0a, 0x07, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x4b, 0x47, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x55,
	0x4e, 0x49, 0x54, 0x5f, 0x4d, 0x4c, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x55, 0x4e, 0x49, 0x54,
	0x5f, 0x4c, 0x10, 0x05, 0x12, 0x0c, 0x0a, 0x08, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x54, 0x53, 0x50,
	0x10, 0x06, 0x12, 0x0d, 0x0a, 0x09, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x54, 0x42, 0x53, 0x50, 0x10,
	0x07, 0x12, 0x0c, 0x0a, 0x08, 0x55, 0x4e, 0x49, 0x54, 0x5f, 0x43, 0x55, 0x50, 0x10, 0x08, 0x32,
	0xb9, 0x03, 0x0a, 0x09, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x59, 0x0a,
	0x10, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x12, 0x1e, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x47, 0x65,
	0x74, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x23, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x49,
	0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x69, 0x6e, 0x76, 0x65,
	0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x69, 0x6e, 0x76,
	0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x10, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x12,
	0x20, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x50, 0x6f, 0x73, 0x74,
	0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e,
	0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x55, 0x0a, 0x10, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x67, 0x72,
	0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x50, 0x0a, 0x10, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x22, 0x2e,
	0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x49, 0x6e, 0x67, 0x72, 0x65, 0x64, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x15, 0x5a, 0x13, 0x69,
	0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f,
	0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_inventory_inventory_proto_rawDescOnce sync.Once
	file_inventory_inventory_proto_rawDescData = file_inventory_inventory_proto_rawDesc
)

func file_inventory_inventory_proto_rawDescGZIP() []byte {
	file_inventory_inventory_proto_rawDescOnce.Do(func() {
		file_inventory_inventory_proto_rawDescData = protoimpl.X.CompressGZIP(file_inventory_inventory_proto_rawDescData)
	})
	return file_inventory_inventory_proto_rawDescData
}

var file_inventory_inventory_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_inventory_inventory_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_inventory_inventory_proto_goTypes = []any{
	(UnitRequest)(0),                 // 0: inventory.UnitRequest
	(*PostIngredientRequest)(nil),    // 1: inventory.PostIngredientRequest
	(*DeleteIngredientRequest)(nil),  // 2: inventory.DeleteIngredientRequest
	(*IngredientResponse)(nil),       // 3: inventory.IngredientResponse
	(*GetUserInventoryResponse)(nil), // 4: inventory.GetUserInventoryResponse
	(*GetInventoryRequest)(nil),      // 5: inventory.GetInventoryRequest
	(*GetIngredientRequest)(nil),     // 6: inventory.GetIngredientRequest
	(*timestamppb.Timestamp)(nil),    // 7: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),            // 8: google.protobuf.Empty
}
var file_inventory_inventory_proto_depIdxs = []int32{
	7, // 0: inventory.IngredientResponse.createdAt:type_name -> google.protobuf.Timestamp
	7, // 1: inventory.IngredientResponse.updatedAt:type_name -> google.protobuf.Timestamp
	3, // 2: inventory.GetUserInventoryResponse.userInventory:type_name -> inventory.IngredientResponse
	5, // 3: inventory.Inventory.GetUserInventory:input_type -> inventory.GetInventoryRequest
	6, // 4: inventory.Inventory.GetIngredient:input_type -> inventory.GetIngredientRequest
	1, // 5: inventory.Inventory.CreateIngredient:input_type -> inventory.PostIngredientRequest
	1, // 6: inventory.Inventory.UpdateIngredient:input_type -> inventory.PostIngredientRequest
	2, // 7: inventory.Inventory.DeleteIngredient:input_type -> inventory.DeleteIngredientRequest
	4, // 8: inventory.Inventory.GetUserInventory:output_type -> inventory.GetUserInventoryResponse
	3, // 9: inventory.Inventory.GetIngredient:output_type -> inventory.IngredientResponse
	3, // 10: inventory.Inventory.CreateIngredient:output_type -> inventory.IngredientResponse
	3, // 11: inventory.Inventory.UpdateIngredient:output_type -> inventory.IngredientResponse
	8, // 12: inventory.Inventory.DeleteIngredient:output_type -> google.protobuf.Empty
	8, // [8:13] is the sub-list for method output_type
	3, // [3:8] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_inventory_inventory_proto_init() }
func file_inventory_inventory_proto_init() {
	if File_inventory_inventory_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_inventory_inventory_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_inventory_inventory_proto_goTypes,
		DependencyIndexes: file_inventory_inventory_proto_depIdxs,
		EnumInfos:         file_inventory_inventory_proto_enumTypes,
		MessageInfos:      file_inventory_inventory_proto_msgTypes,
	}.Build()
	File_inventory_inventory_proto = out.File
	file_inventory_inventory_proto_rawDesc = nil
	file_inventory_inventory_proto_goTypes = nil
	file_inventory_inventory_proto_depIdxs = nil
}
