# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: evochi/v1/evochi.proto
# Protobuf Python Version: 5.27.2
"""Generated protocol buffer code."""

from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder

_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC, 5, 27, 2, "", "evochi/v1/evochi.proto"
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(
    b'\n\x16\x65vochi/v1/evochi.proto\x12\tevochi.v1\x1a\x1fgoogle/protobuf/timestamp.proto"#\n\x05Slice\x12\r\n\x05start\x18\x01 \x01(\x05\x12\x0b\n\x03\x65nd\x18\x02 \x01(\x05">\n\nEvaluation\x12\x1f\n\x05slice\x18\x01 \x01(\x0b\x32\x10.evochi.v1.Slice\x12\x0f\n\x07rewards\x18\x02 \x03(\x02"z\n\nHelloEvent\x12\n\n\x02id\x18\x01 \x01(\t\x12\r\n\x05token\x18\x02 \x01(\t\x12\x17\n\x0fpopulation_size\x18\x03 \x01(\x05\x12\x1a\n\x12heartbeat_interval\x18\x04 \x01(\x05\x12\x12\n\x05state\x18\x05 \x01(\x0cH\x00\x88\x01\x01\x42\x08\n\x06_state"Q\n\rEvaluateEvent\x12\x0f\n\x07task_id\x18\x01 \x01(\t\x12\r\n\x05\x65poch\x18\x02 \x01(\x05\x12 \n\x06slices\x18\x03 \x03(\x0b\x32\x10.evochi.v1.Slice"@\n\rOptimizeEvent\x12\x0f\n\x07task_id\x18\x01 \x01(\t\x12\r\n\x05\x65poch\x18\x02 \x01(\x05\x12\x0f\n\x07rewards\x18\x03 \x03(\x02""\n\x0fInitializeEvent\x12\x0f\n\x07task_id\x18\x01 \x01(\t"1\n\x0fShareStateEvent\x12\x0f\n\x07task_id\x18\x01 \x01(\t\x12\r\n\x05\x65poch\x18\x02 \x01(\x05"!\n\x10SubscribeRequest\x12\r\n\x05\x63ores\x18\x01 \x01(\x05"\xa9\x02\n\x11SubscribeResponse\x12"\n\x04type\x18\x01 \x01(\x0e\x32\x14.evochi.v1.EventType\x12&\n\x05hello\x18\x02 \x01(\x0b\x32\x15.evochi.v1.HelloEventH\x00\x12,\n\x08\x65valuate\x18\x03 \x01(\x0b\x32\x18.evochi.v1.EvaluateEventH\x00\x12,\n\x08optimize\x18\x04 \x01(\x0b\x32\x18.evochi.v1.OptimizeEventH\x00\x12\x30\n\ninitialize\x18\x05 \x01(\x0b\x32\x1a.evochi.v1.InitializeEventH\x00\x12\x31\n\x0bshare_state\x18\x06 \x01(\x0b\x32\x1a.evochi.v1.ShareStateEventH\x00\x42\x07\n\x05\x65vent"Q\n\x10HeartbeatRequest\x12\x0e\n\x06seq_id\x18\x01 \x01(\x05\x12-\n\ttimestamp\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.Timestamp"\x1f\n\x11HeartbeatResponse\x12\n\n\x02ok\x18\x01 \x01(\x08"V\n\x17\x46inishEvaluationRequest\x12\x0f\n\x07task_id\x18\x01 \x01(\t\x12*\n\x0b\x65valuations\x18\x02 \x03(\x0b\x32\x15.evochi.v1.Evaluation"&\n\x18\x46inishEvaluationResponse\x12\n\n\x02ok\x18\x01 \x01(\x08",\n\x19\x46inishOptimizationRequest\x12\x0f\n\x07task_id\x18\x01 \x01(\t"(\n\x1a\x46inishOptimizationResponse\x12\n\n\x02ok\x18\x01 \x01(\x08"=\n\x1b\x46inishInitializationRequest\x12\x0f\n\x07task_id\x18\x01 \x01(\t\x12\r\n\x05state\x18\x02 \x01(\x0c"*\n\x1c\x46inishInitializationResponse\x12\n\n\x02ok\x18\x01 \x01(\x08"9\n\x17\x46inishShareStateRequest\x12\x0f\n\x07task_id\x18\x01 \x01(\t\x12\r\n\x05state\x18\x02 \x01(\x0c"&\n\x18\x46inishShareStateResponse\x12\n\n\x02ok\x18\x01 \x01(\x08*\xa6\x01\n\tEventType\x12\x1a\n\x16\x45VENT_TYPE_UNSPECIFIED\x10\x00\x12\x14\n\x10\x45VENT_TYPE_HELLO\x10\x01\x12\x17\n\x13\x45VENT_TYPE_EVALUATE\x10\x03\x12\x17\n\x13\x45VENT_TYPE_OPTIMIZE\x10\x04\x12\x19\n\x15\x45VENT_TYPE_INITIALIZE\x10\x05\x12\x1a\n\x16\x45VENT_TYPE_SHARE_STATE\x10\x06\x32\xb3\x04\n\rEvochiService\x12J\n\tSubscribe\x12\x1b.evochi.v1.SubscribeRequest\x1a\x1c.evochi.v1.SubscribeResponse"\x00\x30\x01\x12H\n\tHeartbeat\x12\x1b.evochi.v1.HeartbeatRequest\x1a\x1c.evochi.v1.HeartbeatResponse"\x00\x12]\n\x10\x46inishEvaluation\x12".evochi.v1.FinishEvaluationRequest\x1a#.evochi.v1.FinishEvaluationResponse"\x00\x12\x63\n\x12\x46inishOptimization\x12$.evochi.v1.FinishOptimizationRequest\x1a%.evochi.v1.FinishOptimizationResponse"\x00\x12i\n\x14\x46inishInitialization\x12&.evochi.v1.FinishInitializationRequest\x1a\'.evochi.v1.FinishInitializationResponse"\x00\x12]\n\x10\x46inishShareState\x12".evochi.v1.FinishShareStateRequest\x1a#.evochi.v1.FinishShareStateResponse"\x00\x62\x06proto3'
)

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, "evochi.v1.evochi_pb2", _globals)
if not _descriptor._USE_C_DESCRIPTORS:
    DESCRIPTOR._loaded_options = None
    _globals["_EVENTTYPE"]._serialized_start = 1405
    _globals["_EVENTTYPE"]._serialized_end = 1571
    _globals["_SLICE"]._serialized_start = 70
    _globals["_SLICE"]._serialized_end = 105
    _globals["_EVALUATION"]._serialized_start = 107
    _globals["_EVALUATION"]._serialized_end = 169
    _globals["_HELLOEVENT"]._serialized_start = 171
    _globals["_HELLOEVENT"]._serialized_end = 293
    _globals["_EVALUATEEVENT"]._serialized_start = 295
    _globals["_EVALUATEEVENT"]._serialized_end = 376
    _globals["_OPTIMIZEEVENT"]._serialized_start = 378
    _globals["_OPTIMIZEEVENT"]._serialized_end = 442
    _globals["_INITIALIZEEVENT"]._serialized_start = 444
    _globals["_INITIALIZEEVENT"]._serialized_end = 478
    _globals["_SHARESTATEEVENT"]._serialized_start = 480
    _globals["_SHARESTATEEVENT"]._serialized_end = 529
    _globals["_SUBSCRIBEREQUEST"]._serialized_start = 531
    _globals["_SUBSCRIBEREQUEST"]._serialized_end = 564
    _globals["_SUBSCRIBERESPONSE"]._serialized_start = 567
    _globals["_SUBSCRIBERESPONSE"]._serialized_end = 864
    _globals["_HEARTBEATREQUEST"]._serialized_start = 866
    _globals["_HEARTBEATREQUEST"]._serialized_end = 947
    _globals["_HEARTBEATRESPONSE"]._serialized_start = 949
    _globals["_HEARTBEATRESPONSE"]._serialized_end = 980
    _globals["_FINISHEVALUATIONREQUEST"]._serialized_start = 982
    _globals["_FINISHEVALUATIONREQUEST"]._serialized_end = 1068
    _globals["_FINISHEVALUATIONRESPONSE"]._serialized_start = 1070
    _globals["_FINISHEVALUATIONRESPONSE"]._serialized_end = 1108
    _globals["_FINISHOPTIMIZATIONREQUEST"]._serialized_start = 1110
    _globals["_FINISHOPTIMIZATIONREQUEST"]._serialized_end = 1154
    _globals["_FINISHOPTIMIZATIONRESPONSE"]._serialized_start = 1156
    _globals["_FINISHOPTIMIZATIONRESPONSE"]._serialized_end = 1196
    _globals["_FINISHINITIALIZATIONREQUEST"]._serialized_start = 1198
    _globals["_FINISHINITIALIZATIONREQUEST"]._serialized_end = 1259
    _globals["_FINISHINITIALIZATIONRESPONSE"]._serialized_start = 1261
    _globals["_FINISHINITIALIZATIONRESPONSE"]._serialized_end = 1303
    _globals["_FINISHSHARESTATEREQUEST"]._serialized_start = 1305
    _globals["_FINISHSHARESTATEREQUEST"]._serialized_end = 1362
    _globals["_FINISHSHARESTATERESPONSE"]._serialized_start = 1364
    _globals["_FINISHSHARESTATERESPONSE"]._serialized_end = 1402
    _globals["_EVOCHISERVICE"]._serialized_start = 1574
    _globals["_EVOCHISERVICE"]._serialized_end = 2137
# @@protoc_insertion_point(module_scope)
