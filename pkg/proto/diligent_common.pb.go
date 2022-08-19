// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: diligent_common.proto

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

// Proto representation of datagen.Spec
type DataSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SpecType       string      `protobuf:"bytes,1,opt,name=spec_type,json=specType,proto3" json:"spec_type,omitempty"`
	Version        int32       `protobuf:"varint,2,opt,name=version,proto3" json:"version,omitempty"`
	RecordSize     int32       `protobuf:"varint,3,opt,name=record_size,json=recordSize,proto3" json:"record_size,omitempty"`
	KeyGenSpec     *KeyGenSpec `protobuf:"bytes,4,opt,name=key_gen_spec,json=keyGenSpec,proto3" json:"key_gen_spec,omitempty"`
	UniqTrSpec     *TrSpec     `protobuf:"bytes,5,opt,name=uniq_tr_spec,json=uniqTrSpec,proto3" json:"uniq_tr_spec,omitempty"`
	SmallGrpTrSpec *TrSpec     `protobuf:"bytes,6,opt,name=small_grp_tr_spec,json=smallGrpTrSpec,proto3" json:"small_grp_tr_spec,omitempty"`
	LargeGrpTrSpec *TrSpec     `protobuf:"bytes,7,opt,name=large_grp_tr_spec,json=largeGrpTrSpec,proto3" json:"large_grp_tr_spec,omitempty"`
	FixedValue     string      `protobuf:"bytes,8,opt,name=fixed_value,json=fixedValue,proto3" json:"fixed_value,omitempty"`
}

func (x *DataSpec) Reset() {
	*x = DataSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataSpec) ProtoMessage() {}

func (x *DataSpec) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataSpec.ProtoReflect.Descriptor instead.
func (*DataSpec) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{0}
}

func (x *DataSpec) GetSpecType() string {
	if x != nil {
		return x.SpecType
	}
	return ""
}

func (x *DataSpec) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *DataSpec) GetRecordSize() int32 {
	if x != nil {
		return x.RecordSize
	}
	return 0
}

func (x *DataSpec) GetKeyGenSpec() *KeyGenSpec {
	if x != nil {
		return x.KeyGenSpec
	}
	return nil
}

func (x *DataSpec) GetUniqTrSpec() *TrSpec {
	if x != nil {
		return x.UniqTrSpec
	}
	return nil
}

func (x *DataSpec) GetSmallGrpTrSpec() *TrSpec {
	if x != nil {
		return x.SmallGrpTrSpec
	}
	return nil
}

func (x *DataSpec) GetLargeGrpTrSpec() *TrSpec {
	if x != nil {
		return x.LargeGrpTrSpec
	}
	return nil
}

func (x *DataSpec) GetFixedValue() string {
	if x != nil {
		return x.FixedValue
	}
	return ""
}

// Summary info of a data spec
type DataSpecInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SpecName   string `protobuf:"bytes,1,opt,name=spec_name,json=specName,proto3" json:"spec_name,omitempty"`
	SpecType   string `protobuf:"bytes,2,opt,name=spec_type,json=specType,proto3" json:"spec_type,omitempty"`
	Version    int32  `protobuf:"varint,3,opt,name=version,proto3" json:"version,omitempty"`
	NumRecs    int32  `protobuf:"varint,4,opt,name=num_recs,json=numRecs,proto3" json:"num_recs,omitempty"`
	RecordSize int32  `protobuf:"varint,5,opt,name=record_size,json=recordSize,proto3" json:"record_size,omitempty"`
	Hash       int32  `protobuf:"varint,6,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *DataSpecInfo) Reset() {
	*x = DataSpecInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataSpecInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataSpecInfo) ProtoMessage() {}

func (x *DataSpecInfo) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataSpecInfo.ProtoReflect.Descriptor instead.
func (*DataSpecInfo) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{1}
}

func (x *DataSpecInfo) GetSpecName() string {
	if x != nil {
		return x.SpecName
	}
	return ""
}

func (x *DataSpecInfo) GetSpecType() string {
	if x != nil {
		return x.SpecType
	}
	return ""
}

func (x *DataSpecInfo) GetVersion() int32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *DataSpecInfo) GetNumRecs() int32 {
	if x != nil {
		return x.NumRecs
	}
	return 0
}

func (x *DataSpecInfo) GetRecordSize() int32 {
	if x != nil {
		return x.RecordSize
	}
	return 0
}

func (x *DataSpecInfo) GetHash() int32 {
	if x != nil {
		return x.Hash
	}
	return 0
}

// Proto representation of keygen.LeveledKeyGenSpec - used by DataSpec
type KeyGenSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LevelSizes []int32  `protobuf:"varint,1,rep,packed,name=level_sizes,json=levelSizes,proto3" json:"level_sizes,omitempty"`
	SubKeys    []string `protobuf:"bytes,2,rep,name=sub_keys,json=subKeys,proto3" json:"sub_keys,omitempty"`
	Delim      string   `protobuf:"bytes,3,opt,name=delim,proto3" json:"delim,omitempty"`
}

func (x *KeyGenSpec) Reset() {
	*x = KeyGenSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyGenSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyGenSpec) ProtoMessage() {}

func (x *KeyGenSpec) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyGenSpec.ProtoReflect.Descriptor instead.
func (*KeyGenSpec) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{2}
}

func (x *KeyGenSpec) GetLevelSizes() []int32 {
	if x != nil {
		return x.LevelSizes
	}
	return nil
}

func (x *KeyGenSpec) GetSubKeys() []string {
	if x != nil {
		return x.SubKeys
	}
	return nil
}

func (x *KeyGenSpec) GetDelim() string {
	if x != nil {
		return x.Delim
	}
	return ""
}

// Proto representation of strtr.Spec - used by DataSpec
type TrSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Inputs       string `protobuf:"bytes,1,opt,name=inputs,proto3" json:"inputs,omitempty"`
	Replacements string `protobuf:"bytes,2,opt,name=replacements,proto3" json:"replacements,omitempty"`
}

func (x *TrSpec) Reset() {
	*x = TrSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TrSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TrSpec) ProtoMessage() {}

func (x *TrSpec) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TrSpec.ProtoReflect.Descriptor instead.
func (*TrSpec) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{3}
}

func (x *TrSpec) GetInputs() string {
	if x != nil {
		return x.Inputs
	}
	return ""
}

func (x *TrSpec) GetReplacements() string {
	if x != nil {
		return x.Replacements
	}
	return ""
}

// Specification for DB connection
type DBSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Driver string `protobuf:"bytes,1,opt,name=driver,proto3" json:"driver,omitempty"`
	Url    string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *DBSpec) Reset() {
	*x = DBSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBSpec) ProtoMessage() {}

func (x *DBSpec) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBSpec.ProtoReflect.Descriptor instead.
func (*DBSpec) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{4}
}

func (x *DBSpec) GetDriver() string {
	if x != nil {
		return x.Driver
	}
	return ""
}

func (x *DBSpec) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

// Specification of a Workload
type WorkloadSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	WorkloadName  string `protobuf:"bytes,1,opt,name=workload_name,json=workloadName,proto3" json:"workload_name,omitempty"`
	AssignedRange *Range `protobuf:"bytes,2,opt,name=assigned_range,json=assignedRange,proto3" json:"assigned_range,omitempty"`
	TableName     string `protobuf:"bytes,3,opt,name=table_name,json=tableName,proto3" json:"table_name,omitempty"`
	DurationSec   int32  `protobuf:"varint,4,opt,name=duration_sec,json=durationSec,proto3" json:"duration_sec,omitempty"`
	Concurrency   int32  `protobuf:"varint,5,opt,name=concurrency,proto3" json:"concurrency,omitempty"`
	BatchSize     int32  `protobuf:"varint,6,opt,name=batch_size,json=batchSize,proto3" json:"batch_size,omitempty"`
}

func (x *WorkloadSpec) Reset() {
	*x = WorkloadSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkloadSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkloadSpec) ProtoMessage() {}

func (x *WorkloadSpec) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkloadSpec.ProtoReflect.Descriptor instead.
func (*WorkloadSpec) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{5}
}

func (x *WorkloadSpec) GetWorkloadName() string {
	if x != nil {
		return x.WorkloadName
	}
	return ""
}

func (x *WorkloadSpec) GetAssignedRange() *Range {
	if x != nil {
		return x.AssignedRange
	}
	return nil
}

func (x *WorkloadSpec) GetTableName() string {
	if x != nil {
		return x.TableName
	}
	return ""
}

func (x *WorkloadSpec) GetDurationSec() int32 {
	if x != nil {
		return x.DurationSec
	}
	return 0
}

func (x *WorkloadSpec) GetConcurrency() int32 {
	if x != nil {
		return x.Concurrency
	}
	return 0
}

func (x *WorkloadSpec) GetBatchSize() int32 {
	if x != nil {
		return x.BatchSize
	}
	return 0
}

// A range of integers
type Range struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start int32 `protobuf:"varint,1,opt,name=start,proto3" json:"start,omitempty"`
	Limit int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
}

func (x *Range) Reset() {
	*x = Range{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Range) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Range) ProtoMessage() {}

func (x *Range) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Range.ProtoReflect.Descriptor instead.
func (*Range) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{6}
}

func (x *Range) GetStart() int32 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *Range) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

// A generic status information consisting of ok / not ok flag and a message to indicate reason for failure if any
type GeneralStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsOk          bool   `protobuf:"varint,1,opt,name=is_ok,json=isOk,proto3" json:"is_ok,omitempty"`
	FailureReason string `protobuf:"bytes,2,opt,name=failure_reason,json=failureReason,proto3" json:"failure_reason,omitempty"`
}

func (x *GeneralStatus) Reset() {
	*x = GeneralStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeneralStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeneralStatus) ProtoMessage() {}

func (x *GeneralStatus) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GeneralStatus.ProtoReflect.Descriptor instead.
func (*GeneralStatus) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{7}
}

func (x *GeneralStatus) GetIsOk() bool {
	if x != nil {
		return x.IsOk
	}
	return false
}

func (x *GeneralStatus) GetFailureReason() string {
	if x != nil {
		return x.FailureReason
	}
	return ""
}

// General build information
type BuildInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AppName    string `protobuf:"bytes,1,opt,name=app_name,json=appName,proto3" json:"app_name,omitempty"`
	AppVersion string `protobuf:"bytes,2,opt,name=app_version,json=appVersion,proto3" json:"app_version,omitempty"`
	CommitHash string `protobuf:"bytes,3,opt,name=commit_hash,json=commitHash,proto3" json:"commit_hash,omitempty"`
	GoVersion  string `protobuf:"bytes,4,opt,name=go_version,json=goVersion,proto3" json:"go_version,omitempty"`
	BuildTime  string `protobuf:"bytes,5,opt,name=build_time,json=buildTime,proto3" json:"build_time,omitempty"`
}

func (x *BuildInfo) Reset() {
	*x = BuildInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildInfo) ProtoMessage() {}

func (x *BuildInfo) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildInfo.ProtoReflect.Descriptor instead.
func (*BuildInfo) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{8}
}

func (x *BuildInfo) GetAppName() string {
	if x != nil {
		return x.AppName
	}
	return ""
}

func (x *BuildInfo) GetAppVersion() string {
	if x != nil {
		return x.AppVersion
	}
	return ""
}

func (x *BuildInfo) GetCommitHash() string {
	if x != nil {
		return x.CommitHash
	}
	return ""
}

func (x *BuildInfo) GetGoVersion() string {
	if x != nil {
		return x.GoVersion
	}
	return ""
}

func (x *BuildInfo) GetBuildTime() string {
	if x != nil {
		return x.BuildTime
	}
	return ""
}

// General uptime information
type UpTimeInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartTime string `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	Uptime    string `protobuf:"bytes,2,opt,name=uptime,proto3" json:"uptime,omitempty"`
}

func (x *UpTimeInfo) Reset() {
	*x = UpTimeInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_diligent_common_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpTimeInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpTimeInfo) ProtoMessage() {}

func (x *UpTimeInfo) ProtoReflect() protoreflect.Message {
	mi := &file_diligent_common_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpTimeInfo.ProtoReflect.Descriptor instead.
func (*UpTimeInfo) Descriptor() ([]byte, []int) {
	return file_diligent_common_proto_rawDescGZIP(), []int{9}
}

func (x *UpTimeInfo) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *UpTimeInfo) GetUptime() string {
	if x != nil {
		return x.Uptime
	}
	return ""
}

var File_diligent_common_proto protoreflect.FileDescriptor

var file_diligent_common_proto_rawDesc = []byte{
	0x0a, 0x15, 0x64, 0x69, 0x6c, 0x69, 0x67, 0x65, 0x6e, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdd,
	0x02, 0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x53, 0x70, 0x65, 0x63, 0x12, 0x1b, 0x0a, 0x09, 0x73,
	0x70, 0x65, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x73, 0x70, 0x65, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53,
	0x69, 0x7a, 0x65, 0x12, 0x33, 0x0a, 0x0c, 0x6b, 0x65, 0x79, 0x5f, 0x67, 0x65, 0x6e, 0x5f, 0x73,
	0x70, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4b, 0x65, 0x79, 0x47, 0x65, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0a, 0x6b, 0x65,
	0x79, 0x47, 0x65, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x12, 0x2f, 0x0a, 0x0c, 0x75, 0x6e, 0x69, 0x71,
	0x5f, 0x74, 0x72, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0a, 0x75,
	0x6e, 0x69, 0x71, 0x54, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x38, 0x0a, 0x11, 0x73, 0x6d, 0x61,
	0x6c, 0x6c, 0x5f, 0x67, 0x72, 0x70, 0x5f, 0x74, 0x72, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x53,
	0x70, 0x65, 0x63, 0x52, 0x0e, 0x73, 0x6d, 0x61, 0x6c, 0x6c, 0x47, 0x72, 0x70, 0x54, 0x72, 0x53,
	0x70, 0x65, 0x63, 0x12, 0x38, 0x0a, 0x11, 0x6c, 0x61, 0x72, 0x67, 0x65, 0x5f, 0x67, 0x72, 0x70,
	0x5f, 0x74, 0x72, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x53, 0x70, 0x65, 0x63, 0x52, 0x0e, 0x6c,
	0x61, 0x72, 0x67, 0x65, 0x47, 0x72, 0x70, 0x54, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x1f, 0x0a,
	0x0b, 0x66, 0x69, 0x78, 0x65, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x66, 0x69, 0x78, 0x65, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xb2,
	0x01, 0x0a, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x53, 0x70, 0x65, 0x63, 0x49, 0x6e, 0x66, 0x6f, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x70, 0x65, 0x63, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x70, 0x65, 0x63, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x73, 0x70, 0x65, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x73, 0x70, 0x65, 0x63, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x12, 0x19, 0x0a, 0x08, 0x6e, 0x75, 0x6d, 0x5f, 0x72, 0x65, 0x63, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6e, 0x75, 0x6d, 0x52, 0x65, 0x63, 0x73, 0x12, 0x1f,
	0x0a, 0x0b, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x69, 0x7a, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x68,
	0x61, 0x73, 0x68, 0x22, 0x5e, 0x0a, 0x0a, 0x4b, 0x65, 0x79, 0x47, 0x65, 0x6e, 0x53, 0x70, 0x65,
	0x63, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x05, 0x52, 0x0a, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x53, 0x69, 0x7a,
	0x65, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x75, 0x62, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x14, 0x0a,
	0x05, 0x64, 0x65, 0x6c, 0x69, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x64, 0x65,
	0x6c, 0x69, 0x6d, 0x22, 0x44, 0x0a, 0x06, 0x54, 0x72, 0x53, 0x70, 0x65, 0x63, 0x12, 0x16, 0x0a,
	0x06, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69,
	0x6e, 0x70, 0x75, 0x74, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x70,
	0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x32, 0x0a, 0x06, 0x44, 0x42, 0x53,
	0x70, 0x65, 0x63, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0xeb, 0x01,
	0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x70, 0x65, 0x63, 0x12, 0x23,
	0x0a, 0x0d, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x6c, 0x6f, 0x61, 0x64, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x33, 0x0a, 0x0e, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f,
	0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x0d, 0x61, 0x73, 0x73, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x61, 0x62, 0x6c,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x61,
	0x62, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x63, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x63, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f,
	0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0b, 0x63, 0x6f, 0x6e, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x1d, 0x0a, 0x0a,
	0x62, 0x61, 0x74, 0x63, 0x68, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x09, 0x62, 0x61, 0x74, 0x63, 0x68, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x33, 0x0a, 0x05, 0x52,
	0x61, 0x6e, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69,
	0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x22, 0x4b, 0x0a, 0x0d, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x13, 0x0a, 0x05, 0x69, 0x73, 0x5f, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x04, 0x69, 0x73, 0x4f, 0x6b, 0x12, 0x25, 0x0a, 0x0e, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72,
	0x65, 0x5f, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d,
	0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0xa6, 0x01,
	0x0a, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x61,
	0x70, 0x70, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x70, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x70, 0x70, 0x5f, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x70, 0x70,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1d, 0x0a, 0x0a, 0x67, 0x6f, 0x5f, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x67, 0x6f,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x75, 0x69, 0x6c, 0x64,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x43, 0x0a, 0x0a, 0x55, 0x70, 0x54, 0x69, 0x6d, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x70, 0x74, 0x69, 0x6d, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_diligent_common_proto_rawDescOnce sync.Once
	file_diligent_common_proto_rawDescData = file_diligent_common_proto_rawDesc
)

func file_diligent_common_proto_rawDescGZIP() []byte {
	file_diligent_common_proto_rawDescOnce.Do(func() {
		file_diligent_common_proto_rawDescData = protoimpl.X.CompressGZIP(file_diligent_common_proto_rawDescData)
	})
	return file_diligent_common_proto_rawDescData
}

var file_diligent_common_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_diligent_common_proto_goTypes = []interface{}{
	(*DataSpec)(nil),      // 0: proto.DataSpec
	(*DataSpecInfo)(nil),  // 1: proto.DataSpecInfo
	(*KeyGenSpec)(nil),    // 2: proto.KeyGenSpec
	(*TrSpec)(nil),        // 3: proto.TrSpec
	(*DBSpec)(nil),        // 4: proto.DBSpec
	(*WorkloadSpec)(nil),  // 5: proto.WorkloadSpec
	(*Range)(nil),         // 6: proto.Range
	(*GeneralStatus)(nil), // 7: proto.GeneralStatus
	(*BuildInfo)(nil),     // 8: proto.BuildInfo
	(*UpTimeInfo)(nil),    // 9: proto.UpTimeInfo
}
var file_diligent_common_proto_depIdxs = []int32{
	2, // 0: proto.DataSpec.key_gen_spec:type_name -> proto.KeyGenSpec
	3, // 1: proto.DataSpec.uniq_tr_spec:type_name -> proto.TrSpec
	3, // 2: proto.DataSpec.small_grp_tr_spec:type_name -> proto.TrSpec
	3, // 3: proto.DataSpec.large_grp_tr_spec:type_name -> proto.TrSpec
	6, // 4: proto.WorkloadSpec.assigned_range:type_name -> proto.Range
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_diligent_common_proto_init() }
func file_diligent_common_proto_init() {
	if File_diligent_common_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_diligent_common_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataSpec); i {
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
		file_diligent_common_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataSpecInfo); i {
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
		file_diligent_common_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyGenSpec); i {
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
		file_diligent_common_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TrSpec); i {
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
		file_diligent_common_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBSpec); i {
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
		file_diligent_common_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkloadSpec); i {
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
		file_diligent_common_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Range); i {
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
		file_diligent_common_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GeneralStatus); i {
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
		file_diligent_common_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildInfo); i {
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
		file_diligent_common_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpTimeInfo); i {
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
			RawDescriptor: file_diligent_common_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_diligent_common_proto_goTypes,
		DependencyIndexes: file_diligent_common_proto_depIdxs,
		MessageInfos:      file_diligent_common_proto_msgTypes,
	}.Build()
	File_diligent_common_proto = out.File
	file_diligent_common_proto_rawDesc = nil
	file_diligent_common_proto_goTypes = nil
	file_diligent_common_proto_depIdxs = nil
}
