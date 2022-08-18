// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: diligent_boss.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BossPingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BossPingRequest) Reset() {
	*x = BossPingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossPingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossPingRequest) ProtoMessage() {}

func (x *BossPingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossPingRequest.ProtoReflect.Descriptor instead.
func (*BossPingRequest) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{0}
}

type BossPingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BossPingResponse) Reset() {
	*x = BossPingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossPingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossPingResponse) ProtoMessage() {}

func (x *BossPingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossPingResponse.ProtoReflect.Descriptor instead.
func (*BossPingResponse) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{1}
}

type BossRegisterMinionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *BossRegisterMinionRequest) Reset() {
	*x = BossRegisterMinionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossRegisterMinionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossRegisterMinionRequest) ProtoMessage() {}

func (x *BossRegisterMinionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossRegisterMinionRequest.ProtoReflect.Descriptor instead.
func (*BossRegisterMinionRequest) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{2}
}

func (x *BossRegisterMinionRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type BossRegisterMinionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *GeneralStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *BossRegisterMinionResponse) Reset() {
	*x = BossRegisterMinionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossRegisterMinionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossRegisterMinionResponse) ProtoMessage() {}

func (x *BossRegisterMinionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossRegisterMinionResponse.ProtoReflect.Descriptor instead.
func (*BossRegisterMinionResponse) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{3}
}

func (x *BossRegisterMinionResponse) GetStatus() *GeneralStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type BossUnregisterMinonRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *BossUnregisterMinonRequest) Reset() {
	*x = BossUnregisterMinonRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossUnregisterMinonRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossUnregisterMinonRequest) ProtoMessage() {}

func (x *BossUnregisterMinonRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossUnregisterMinonRequest.ProtoReflect.Descriptor instead.
func (*BossUnregisterMinonRequest) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{4}
}

func (x *BossUnregisterMinonRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type BossUnregisterMinionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *GeneralStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *BossUnregisterMinionResponse) Reset() {
	*x = BossUnregisterMinionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossUnregisterMinionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossUnregisterMinionResponse) ProtoMessage() {}

func (x *BossUnregisterMinionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossUnregisterMinionResponse.ProtoReflect.Descriptor instead.
func (*BossUnregisterMinionResponse) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{5}
}

func (x *BossUnregisterMinionResponse) GetStatus() *GeneralStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type BossShowMinionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BossShowMinionRequest) Reset() {
	*x = BossShowMinionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossShowMinionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossShowMinionRequest) ProtoMessage() {}

func (x *BossShowMinionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossShowMinionRequest.ProtoReflect.Descriptor instead.
func (*BossShowMinionRequest) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{6}
}

type BossShowMinionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MinionStatuses []*MinionStatus `protobuf:"bytes,1,rep,name=minion_statuses,json=minionStatuses,proto3" json:"minion_statuses,omitempty"`
}

func (x *BossShowMinionResponse) Reset() {
	*x = BossShowMinionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossShowMinionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossShowMinionResponse) ProtoMessage() {}

func (x *BossShowMinionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossShowMinionResponse.ProtoReflect.Descriptor instead.
func (*BossShowMinionResponse) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{7}
}

func (x *BossShowMinionResponse) GetMinionStatuses() []*MinionStatus {
	if x != nil {
		return x.MinionStatuses
	}
	return nil
}

type MinionStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url    string         `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	Status *GeneralStatus `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *MinionStatus) Reset() {
	*x = MinionStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionStatus) ProtoMessage() {}

func (x *MinionStatus) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionStatus.ProtoReflect.Descriptor instead.
func (*MinionStatus) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{8}
}

func (x *MinionStatus) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *MinionStatus) GetStatus() *GeneralStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type BossRunWorkloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataSpec *DataSpec     `protobuf:"bytes,1,opt,name=data_spec,json=dataSpec,proto3" json:"data_spec,omitempty"`
	DbSpec   *DBSpec       `protobuf:"bytes,2,opt,name=db_spec,json=dbSpec,proto3" json:"db_spec,omitempty"`
	WlSpec   *WorkloadSpec `protobuf:"bytes,3,opt,name=wl_spec,json=wlSpec,proto3" json:"wl_spec,omitempty"`
}

func (x *BossRunWorkloadRequest) Reset() {
	*x = BossRunWorkloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossRunWorkloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossRunWorkloadRequest) ProtoMessage() {}

func (x *BossRunWorkloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossRunWorkloadRequest.ProtoReflect.Descriptor instead.
func (*BossRunWorkloadRequest) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{9}
}

func (x *BossRunWorkloadRequest) GetDataSpec() *DataSpec {
	if x != nil {
		return x.DataSpec
	}
	return nil
}

func (x *BossRunWorkloadRequest) GetDbSpec() *DBSpec {
	if x != nil {
		return x.DbSpec
	}
	return nil
}

func (x *BossRunWorkloadRequest) GetWlSpec() *WorkloadSpec {
	if x != nil {
		return x.WlSpec
	}
	return nil
}

type BossRunWorkloadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OverallStatus  *GeneralStatus  `protobuf:"bytes,1,opt,name=overall_status,json=overallStatus,proto3" json:"overall_status,omitempty"`
	MinionStatuses []*MinionStatus `protobuf:"bytes,2,rep,name=minion_statuses,json=minionStatuses,proto3" json:"minion_statuses,omitempty"`
}

func (x *BossRunWorkloadResponse) Reset() {
	*x = BossRunWorkloadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossRunWorkloadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossRunWorkloadResponse) ProtoMessage() {}

func (x *BossRunWorkloadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossRunWorkloadResponse.ProtoReflect.Descriptor instead.
func (*BossRunWorkloadResponse) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{10}
}

func (x *BossRunWorkloadResponse) GetOverallStatus() *GeneralStatus {
	if x != nil {
		return x.OverallStatus
	}
	return nil
}

func (x *BossRunWorkloadResponse) GetMinionStatuses() []*MinionStatus {
	if x != nil {
		return x.MinionStatuses
	}
	return nil
}

type BossStopWorkloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BossStopWorkloadRequest) Reset() {
	*x = BossStopWorkloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossStopWorkloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossStopWorkloadRequest) ProtoMessage() {}

func (x *BossStopWorkloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossStopWorkloadRequest.ProtoReflect.Descriptor instead.
func (*BossStopWorkloadRequest) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{11}
}

type BossStopWorkloadResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OverallStatus  *GeneralStatus  `protobuf:"bytes,1,opt,name=overall_status,json=overallStatus,proto3" json:"overall_status,omitempty"`
	MinionStatuses []*MinionStatus `protobuf:"bytes,2,rep,name=minion_statuses,json=minionStatuses,proto3" json:"minion_statuses,omitempty"`
}

func (x *BossStopWorkloadResponse) Reset() {
	*x = BossStopWorkloadResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_boss_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BossStopWorkloadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BossStopWorkloadResponse) ProtoMessage() {}

func (x *BossStopWorkloadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_boss_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BossStopWorkloadResponse.ProtoReflect.Descriptor instead.
func (*BossStopWorkloadResponse) Descriptor() ([]byte, []int) {
	return file_diligent_boss_proto_rawDescGZIP(), []int{12}
}

func (x *BossStopWorkloadResponse) GetOverallStatus() *GeneralStatus {
	if x != nil {
		return x.OverallStatus
	}
	return nil
}

func (x *BossStopWorkloadResponse) GetMinionStatuses() []*MinionStatus {
	if x != nil {
		return x.MinionStatuses
	}
	return nil
}

var File_diligent_boss_proto protoreflect.FileDescriptor

var file_diligent_boss_proto_rawDesc = []byte{
	0x0a, 0x13, 0x64, 0x69, 0x6c, 0x69, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x62, 0x6f, 0x73, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x64, 0x69,
	0x6c, 0x69, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x11, 0x0a, 0x0f, 0x42, 0x6f, 0x73, 0x73, 0x50, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x12, 0x0a, 0x10, 0x42, 0x6f, 0x73, 0x73, 0x50, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2d, 0x0a, 0x19, 0x42, 0x6f,
	0x73, 0x73, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x4a, 0x0a, 0x1a, 0x42, 0x6f, 0x73,
	0x73, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x2e, 0x0a, 0x1a, 0x42, 0x6f, 0x73, 0x73, 0x55, 0x6e, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x4c, 0x0a, 0x1c, 0x42, 0x6f, 0x73, 0x73, 0x55, 0x6e, 0x72,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x22, 0x17, 0x0a, 0x15, 0x42, 0x6f, 0x73, 0x73, 0x53, 0x68, 0x6f, 0x77, 0x4d,
	0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x56, 0x0a, 0x16,
	0x42, 0x6f, 0x73, 0x73, 0x53, 0x68, 0x6f, 0x77, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x0f, 0x6d, 0x69, 0x6e, 0x69, 0x6f, 0x6e,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x0e, 0x6d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x65, 0x73, 0x22, 0x4e, 0x0a, 0x0c, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x22, 0x9c, 0x01, 0x0a, 0x16, 0x42, 0x6f, 0x73, 0x73, 0x52, 0x75, 0x6e,
	0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2c, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x53, 0x70, 0x65, 0x63, 0x12, 0x26, 0x0a,
	0x07, 0x64, 0x62, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x42, 0x53, 0x70, 0x65, 0x63, 0x52, 0x06, 0x64,
	0x62, 0x53, 0x70, 0x65, 0x63, 0x12, 0x2c, 0x0a, 0x07, 0x77, 0x6c, 0x5f, 0x73, 0x70, 0x65, 0x63,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57,
	0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x70, 0x65, 0x63, 0x52, 0x06, 0x77, 0x6c, 0x53,
	0x70, 0x65, 0x63, 0x22, 0x94, 0x01, 0x0a, 0x17, 0x42, 0x6f, 0x73, 0x73, 0x52, 0x75, 0x6e, 0x57,
	0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3b, 0x0a, 0x0e, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x6c, 0x6c, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0d, 0x6f,
	0x76, 0x65, 0x72, 0x61, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x3c, 0x0a, 0x0f,
	0x6d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69,
	0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0e, 0x6d, 0x69, 0x6e, 0x69,
	0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x22, 0x19, 0x0a, 0x17, 0x42, 0x6f,
	0x73, 0x73, 0x53, 0x74, 0x6f, 0x70, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x95, 0x01, 0x0a, 0x18, 0x42, 0x6f, 0x73, 0x73, 0x53, 0x74,
	0x6f, 0x70, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0e, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x6c, 0x6c, 0x5f, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x0d, 0x6f, 0x76, 0x65, 0x72, 0x61, 0x6c, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x3c, 0x0a, 0x0f, 0x6d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0e, 0x6d,
	0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x65, 0x73, 0x32, 0xdd, 0x03,
	0x0a, 0x04, 0x42, 0x6f, 0x73, 0x73, 0x12, 0x37, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x12, 0x16,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x50, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42,
	0x6f, 0x73, 0x73, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x55, 0x0a, 0x0e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x69, 0x6f,
	0x6e, 0x12, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x10, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x4d, 0x69, 0x6e, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x55, 0x6e, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0b, 0x53, 0x68, 0x6f, 0x77, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x53, 0x68,
	0x6f, 0x77, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x53, 0x68, 0x6f, 0x77,
	0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c,
	0x0a, 0x0b, 0x52, 0x75, 0x6e, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x52, 0x75, 0x6e, 0x57, 0x6f, 0x72,
	0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x52, 0x75, 0x6e, 0x57, 0x6f, 0x72, 0x6b,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4f, 0x0a, 0x0c,
	0x53, 0x74, 0x6f, 0x70, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x53, 0x74, 0x6f, 0x70, 0x57, 0x6f, 0x72,
	0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x6f, 0x73, 0x73, 0x53, 0x74, 0x6f, 0x70, 0x57, 0x6f, 0x72,
	0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a,
	0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_diligent_boss_proto_rawDescOnce sync.Once
	file_diligent_boss_proto_rawDescData = file_diligent_boss_proto_rawDesc
)

func file_diligent_boss_proto_rawDescGZIP() []byte {
	file_diligent_boss_proto_rawDescOnce.Do(func() {
		file_diligent_boss_proto_rawDescData = protoimpl.X.CompressGZIP(file_diligent_boss_proto_rawDescData)
	})
	return file_diligent_boss_proto_rawDescData
}

var file_diligent_boss_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_diligent_boss_proto_goTypes = []interface{}{
	(*BossPingRequest)(nil),              // 0: proto.BossPingRequest
	(*BossPingResponse)(nil),             // 1: proto.BossPingResponse
	(*BossRegisterMinionRequest)(nil),    // 2: proto.BossRegisterMinionRequest
	(*BossRegisterMinionResponse)(nil),   // 3: proto.BossRegisterMinionResponse
	(*BossUnregisterMinonRequest)(nil),   // 4: proto.BossUnregisterMinonRequest
	(*BossUnregisterMinionResponse)(nil), // 5: proto.BossUnregisterMinionResponse
	(*BossShowMinionRequest)(nil),        // 6: proto.BossShowMinionRequest
	(*BossShowMinionResponse)(nil),       // 7: proto.BossShowMinionResponse
	(*MinionStatus)(nil),                 // 8: proto.MinionStatus
	(*BossRunWorkloadRequest)(nil),       // 9: proto.BossRunWorkloadRequest
	(*BossRunWorkloadResponse)(nil),      // 10: proto.BossRunWorkloadResponse
	(*BossStopWorkloadRequest)(nil),      // 11: proto.BossStopWorkloadRequest
	(*BossStopWorkloadResponse)(nil),     // 12: proto.BossStopWorkloadResponse
	(*GeneralStatus)(nil),                // 13: proto.GeneralStatus
	(*DataSpec)(nil),                     // 14: proto.DataSpec
	(*DBSpec)(nil),                       // 15: proto.DBSpec
	(*WorkloadSpec)(nil),                 // 16: proto.WorkloadSpec
}
var file_diligent_boss_proto_depIdxs = []int32{
	13, // 0: proto.BossRegisterMinionResponse.status:type_name -> proto.GeneralStatus
	13, // 1: proto.BossUnregisterMinionResponse.status:type_name -> proto.GeneralStatus
	8,  // 2: proto.BossShowMinionResponse.minion_statuses:type_name -> proto.MinionStatus
	13, // 3: proto.MinionStatus.status:type_name -> proto.GeneralStatus
	14, // 4: proto.BossRunWorkloadRequest.data_spec:type_name -> proto.DataSpec
	15, // 5: proto.BossRunWorkloadRequest.db_spec:type_name -> proto.DBSpec
	16, // 6: proto.BossRunWorkloadRequest.wl_spec:type_name -> proto.WorkloadSpec
	13, // 7: proto.BossRunWorkloadResponse.overall_status:type_name -> proto.GeneralStatus
	8,  // 8: proto.BossRunWorkloadResponse.minion_statuses:type_name -> proto.MinionStatus
	13, // 9: proto.BossStopWorkloadResponse.overall_status:type_name -> proto.GeneralStatus
	8,  // 10: proto.BossStopWorkloadResponse.minion_statuses:type_name -> proto.MinionStatus
	0,  // 11: proto.Boss.Ping:input_type -> proto.BossPingRequest
	2,  // 12: proto.Boss.RegisterMinion:input_type -> proto.BossRegisterMinionRequest
	4,  // 13: proto.Boss.UnregisterMinion:input_type -> proto.BossUnregisterMinonRequest
	6,  // 14: proto.Boss.ShowMinions:input_type -> proto.BossShowMinionRequest
	9,  // 15: proto.Boss.RunWorkload:input_type -> proto.BossRunWorkloadRequest
	11, // 16: proto.Boss.StopWorkload:input_type -> proto.BossStopWorkloadRequest
	1,  // 17: proto.Boss.Ping:output_type -> proto.BossPingResponse
	3,  // 18: proto.Boss.RegisterMinion:output_type -> proto.BossRegisterMinionResponse
	5,  // 19: proto.Boss.UnregisterMinion:output_type -> proto.BossUnregisterMinionResponse
	7,  // 20: proto.Boss.ShowMinions:output_type -> proto.BossShowMinionResponse
	10, // 21: proto.Boss.RunWorkload:output_type -> proto.BossRunWorkloadResponse
	12, // 22: proto.Boss.StopWorkload:output_type -> proto.BossStopWorkloadResponse
	17, // [17:23] is the sub-list for method output_type
	11, // [11:17] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_diligent_boss_proto_init() }
func file_diligent_boss_proto_init() {
	if File_diligent_boss_proto != nil {
		return
	}
	file_diligent_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_diligent_boss_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossPingRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossPingResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossRegisterMinionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossRegisterMinionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossUnregisterMinonRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossUnregisterMinionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossShowMinionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossShowMinionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossRunWorkloadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossRunWorkloadResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossStopWorkloadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_diligent_boss_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BossStopWorkloadResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_diligent_boss_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_diligent_boss_proto_goTypes,
		DependencyIndexes: file_diligent_boss_proto_depIdxs,
		MessageInfos:      file_diligent_boss_proto_msgTypes,
	}.Build()
	File_diligent_boss_proto = out.File
	file_diligent_boss_proto_rawDesc = nil
	file_diligent_boss_proto_goTypes = nil
	file_diligent_boss_proto_depIdxs = nil
}
