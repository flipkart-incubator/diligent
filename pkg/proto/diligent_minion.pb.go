// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: diligent_minion.proto

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

type MinionPingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MinionPingRequest) Reset() {
	*x = MinionPingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_minion_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionPingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionPingRequest) ProtoMessage() {}

func (x *MinionPingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_minion_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionPingRequest.ProtoReflect.Descriptor instead.
func (*MinionPingRequest) Descriptor() ([]byte, []int) {
	return file_diligent_minion_proto_rawDescGZIP(), []int{0}
}

type MinionPingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BuildInfo   *BuildInfo   `protobuf:"bytes,1,opt,name=build_info,json=buildInfo,proto3" json:"build_info,omitempty"`
	ProcessInfo *ProcessInfo `protobuf:"bytes,2,opt,name=process_info,json=processInfo,proto3" json:"process_info,omitempty"`
	JobInfo     *JobInfo     `protobuf:"bytes,3,opt,name=job_info,json=jobInfo,proto3" json:"job_info,omitempty"`
}

func (x *MinionPingResponse) Reset() {
	*x = MinionPingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_minion_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionPingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionPingResponse) ProtoMessage() {}

func (x *MinionPingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_minion_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionPingResponse.ProtoReflect.Descriptor instead.
func (*MinionPingResponse) Descriptor() ([]byte, []int) {
	return file_diligent_minion_proto_rawDescGZIP(), []int{1}
}

func (x *MinionPingResponse) GetBuildInfo() *BuildInfo {
	if x != nil {
		return x.BuildInfo
	}
	return nil
}

func (x *MinionPingResponse) GetProcessInfo() *ProcessInfo {
	if x != nil {
		return x.ProcessInfo
	}
	return nil
}

func (x *MinionPingResponse) GetJobInfo() *JobInfo {
	if x != nil {
		return x.JobInfo
	}
	return nil
}

type MinionPrepareJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobId   string   `protobuf:"bytes,1,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
	JobDesc string   `protobuf:"bytes,2,opt,name=job_desc,json=jobDesc,proto3" json:"job_desc,omitempty"`
	JobSpec *JobSpec `protobuf:"bytes,3,opt,name=job_spec,json=jobSpec,proto3" json:"job_spec,omitempty"`
}

func (x *MinionPrepareJobRequest) Reset() {
	*x = MinionPrepareJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_minion_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionPrepareJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionPrepareJobRequest) ProtoMessage() {}

func (x *MinionPrepareJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_minion_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionPrepareJobRequest.ProtoReflect.Descriptor instead.
func (*MinionPrepareJobRequest) Descriptor() ([]byte, []int) {
	return file_diligent_minion_proto_rawDescGZIP(), []int{2}
}

func (x *MinionPrepareJobRequest) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

func (x *MinionPrepareJobRequest) GetJobDesc() string {
	if x != nil {
		return x.JobDesc
	}
	return ""
}

func (x *MinionPrepareJobRequest) GetJobSpec() *JobSpec {
	if x != nil {
		return x.JobSpec
	}
	return nil
}

type MinionPrepareJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *GeneralStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Pid    string         `protobuf:"bytes,2,opt,name=pid,proto3" json:"pid,omitempty"`
	JobId  string         `protobuf:"bytes,3,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (x *MinionPrepareJobResponse) Reset() {
	*x = MinionPrepareJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_minion_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionPrepareJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionPrepareJobResponse) ProtoMessage() {}

func (x *MinionPrepareJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_minion_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionPrepareJobResponse.ProtoReflect.Descriptor instead.
func (*MinionPrepareJobResponse) Descriptor() ([]byte, []int) {
	return file_diligent_minion_proto_rawDescGZIP(), []int{3}
}

func (x *MinionPrepareJobResponse) GetStatus() *GeneralStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *MinionPrepareJobResponse) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *MinionPrepareJobResponse) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

type MinionRunJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MinionRunJobRequest) Reset() {
	*x = MinionRunJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_minion_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionRunJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionRunJobRequest) ProtoMessage() {}

func (x *MinionRunJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_minion_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionRunJobRequest.ProtoReflect.Descriptor instead.
func (*MinionRunJobRequest) Descriptor() ([]byte, []int) {
	return file_diligent_minion_proto_rawDescGZIP(), []int{4}
}

type MinionRunJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *GeneralStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Pid    string         `protobuf:"bytes,2,opt,name=pid,proto3" json:"pid,omitempty"`
	JobId  string         `protobuf:"bytes,3,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (x *MinionRunJobResponse) Reset() {
	*x = MinionRunJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_minion_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionRunJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionRunJobResponse) ProtoMessage() {}

func (x *MinionRunJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_minion_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionRunJobResponse.ProtoReflect.Descriptor instead.
func (*MinionRunJobResponse) Descriptor() ([]byte, []int) {
	return file_diligent_minion_proto_rawDescGZIP(), []int{5}
}

func (x *MinionRunJobResponse) GetStatus() *GeneralStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *MinionRunJobResponse) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *MinionRunJobResponse) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

type MinionStopJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MinionStopJobRequest) Reset() {
	*x = MinionStopJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_minion_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionStopJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionStopJobRequest) ProtoMessage() {}

func (x *MinionStopJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_minion_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionStopJobRequest.ProtoReflect.Descriptor instead.
func (*MinionStopJobRequest) Descriptor() ([]byte, []int) {
	return file_diligent_minion_proto_rawDescGZIP(), []int{6}
}

type MinionStopJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *GeneralStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	Pid    string         `protobuf:"bytes,2,opt,name=pid,proto3" json:"pid,omitempty"`
	JobId  string         `protobuf:"bytes,3,opt,name=job_id,json=jobId,proto3" json:"job_id,omitempty"`
}

func (x *MinionStopJobResponse) Reset() {
	*x = MinionStopJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_minion_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinionStopJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinionStopJobResponse) ProtoMessage() {}

func (x *MinionStopJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_minion_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinionStopJobResponse.ProtoReflect.Descriptor instead.
func (*MinionStopJobResponse) Descriptor() ([]byte, []int) {
	return file_diligent_minion_proto_rawDescGZIP(), []int{7}
}

func (x *MinionStopJobResponse) GetStatus() *GeneralStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *MinionStopJobResponse) GetPid() string {
	if x != nil {
		return x.Pid
	}
	return ""
}

func (x *MinionStopJobResponse) GetJobId() string {
	if x != nil {
		return x.JobId
	}
	return ""
}

var File_diligent_minion_proto protoreflect.FileDescriptor

var file_diligent_minion_proto_rawDesc = []byte{
	0x0a, 0x15, 0x64, 0x69, 0x6c, 0x69, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x6d, 0x69, 0x6e, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15,
	0x64, 0x69, 0x6c, 0x69, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x13, 0x0a, 0x11, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xa7, 0x01, 0x0a, 0x12, 0x4d,
	0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2f, 0x0a, 0x0a, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x35, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0b, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x29, 0x0a, 0x08, 0x6a, 0x6f, 0x62,
	0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x4a, 0x6f, 0x62, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x07, 0x6a, 0x6f, 0x62,
	0x49, 0x6e, 0x66, 0x6f, 0x22, 0x76, 0x0a, 0x17, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x50, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x15, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x64, 0x65,
	0x73, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x44, 0x65, 0x73,
	0x63, 0x12, 0x29, 0x0a, 0x08, 0x6a, 0x6f, 0x62, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4a, 0x6f, 0x62, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x07, 0x6a, 0x6f, 0x62, 0x53, 0x70, 0x65, 0x63, 0x22, 0x71, 0x0a, 0x18,
	0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x4a, 0x6f, 0x62,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x22,
	0x15, 0x0a, 0x13, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x75, 0x6e, 0x4a, 0x6f, 0x62, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x6d, 0x0a, 0x14, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e,
	0x52, 0x75, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03,
	0x70, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x15,
	0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x6a, 0x6f, 0x62, 0x49, 0x64, 0x22, 0x16, 0x0a, 0x14, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x53,
	0x74, 0x6f, 0x70, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x6e, 0x0a,
	0x15, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x6f, 0x70, 0x4a, 0x6f, 0x62, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x15, 0x0a, 0x06, 0x6a, 0x6f, 0x62, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x49, 0x64, 0x32, 0x9d, 0x02,
	0x0a, 0x06, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x12, 0x3b, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67,
	0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4d, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x4a, 0x6f, 0x62, 0x12, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69,
	0x6f, 0x6e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69,
	0x6f, 0x6e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x06, 0x52, 0x75, 0x6e, 0x4a, 0x6f, 0x62, 0x12, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x75, 0x6e,
	0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x52, 0x75, 0x6e, 0x4a, 0x6f, 0x62, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x70, 0x4a,
	0x6f, 0x62, 0x12, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x6f,
	0x6e, 0x53, 0x74, 0x6f, 0x70, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x69, 0x6e, 0x69, 0x6f, 0x6e, 0x53, 0x74,
	0x6f, 0x70, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a,
	0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_diligent_minion_proto_rawDescOnce sync.Once
	file_diligent_minion_proto_rawDescData = file_diligent_minion_proto_rawDesc
)

func file_diligent_minion_proto_rawDescGZIP() []byte {
	file_diligent_minion_proto_rawDescOnce.Do(func() {
		file_diligent_minion_proto_rawDescData = protoimpl.X.CompressGZIP(file_diligent_minion_proto_rawDescData)
	})
	return file_diligent_minion_proto_rawDescData
}

var file_diligent_minion_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_diligent_minion_proto_goTypes = []interface{}{
	(*MinionPingRequest)(nil),        // 0: proto.MinionPingRequest
	(*MinionPingResponse)(nil),       // 1: proto.MinionPingResponse
	(*MinionPrepareJobRequest)(nil),  // 2: proto.MinionPrepareJobRequest
	(*MinionPrepareJobResponse)(nil), // 3: proto.MinionPrepareJobResponse
	(*MinionRunJobRequest)(nil),      // 4: proto.MinionRunJobRequest
	(*MinionRunJobResponse)(nil),     // 5: proto.MinionRunJobResponse
	(*MinionStopJobRequest)(nil),     // 6: proto.MinionStopJobRequest
	(*MinionStopJobResponse)(nil),    // 7: proto.MinionStopJobResponse
	(*BuildInfo)(nil),                // 8: proto.BuildInfo
	(*ProcessInfo)(nil),              // 9: proto.ProcessInfo
	(*JobInfo)(nil),                  // 10: proto.JobInfo
	(*JobSpec)(nil),                  // 11: proto.JobSpec
	(*GeneralStatus)(nil),            // 12: proto.GeneralStatus
}
var file_diligent_minion_proto_depIdxs = []int32{
	8,  // 0: proto.MinionPingResponse.build_info:type_name -> proto.BuildInfo
	9,  // 1: proto.MinionPingResponse.process_info:type_name -> proto.ProcessInfo
	10, // 2: proto.MinionPingResponse.job_info:type_name -> proto.JobInfo
	11, // 3: proto.MinionPrepareJobRequest.job_spec:type_name -> proto.JobSpec
	12, // 4: proto.MinionPrepareJobResponse.status:type_name -> proto.GeneralStatus
	12, // 5: proto.MinionRunJobResponse.status:type_name -> proto.GeneralStatus
	12, // 6: proto.MinionStopJobResponse.status:type_name -> proto.GeneralStatus
	0,  // 7: proto.Minion.Ping:input_type -> proto.MinionPingRequest
	2,  // 8: proto.Minion.PrepareJob:input_type -> proto.MinionPrepareJobRequest
	4,  // 9: proto.Minion.RunJob:input_type -> proto.MinionRunJobRequest
	6,  // 10: proto.Minion.StopJob:input_type -> proto.MinionStopJobRequest
	1,  // 11: proto.Minion.Ping:output_type -> proto.MinionPingResponse
	3,  // 12: proto.Minion.PrepareJob:output_type -> proto.MinionPrepareJobResponse
	5,  // 13: proto.Minion.RunJob:output_type -> proto.MinionRunJobResponse
	7,  // 14: proto.Minion.StopJob:output_type -> proto.MinionStopJobResponse
	11, // [11:15] is the sub-list for method output_type
	7,  // [7:11] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_diligent_minion_proto_init() }
func file_diligent_minion_proto_init() {
	if File_diligent_minion_proto != nil {
		return
	}
	file_diligent_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_diligent_minion_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionPingRequest); i {
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
		file_diligent_minion_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionPingResponse); i {
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
		file_diligent_minion_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionPrepareJobRequest); i {
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
		file_diligent_minion_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionPrepareJobResponse); i {
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
		file_diligent_minion_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionRunJobRequest); i {
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
		file_diligent_minion_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionRunJobResponse); i {
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
		file_diligent_minion_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionStopJobRequest); i {
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
		file_diligent_minion_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinionStopJobResponse); i {
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
			RawDescriptor: file_diligent_minion_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_diligent_minion_proto_goTypes,
		DependencyIndexes: file_diligent_minion_proto_depIdxs,
		MessageInfos:      file_diligent_minion_proto_msgTypes,
	}.Build()
	File_diligent_minion_proto = out.File
	file_diligent_minion_proto_rawDesc = nil
	file_diligent_minion_proto_goTypes = nil
	file_diligent_minion_proto_depIdxs = nil
}
