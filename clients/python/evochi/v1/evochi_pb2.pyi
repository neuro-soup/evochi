from google.protobuf import timestamp_pb2 as _timestamp_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import (
    ClassVar as _ClassVar,
    Iterable as _Iterable,
    Mapping as _Mapping,
    Optional as _Optional,
    Union as _Union,
)

DESCRIPTOR: _descriptor.FileDescriptor

class EventType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    EVENT_TYPE_UNSPECIFIED: _ClassVar[EventType]
    EVENT_TYPE_HELLO: _ClassVar[EventType]
    EVENT_TYPE_EVALUATE: _ClassVar[EventType]
    EVENT_TYPE_OPTIMIZE: _ClassVar[EventType]
    EVENT_TYPE_INITIALIZE: _ClassVar[EventType]
    EVENT_TYPE_SHARE_STATE: _ClassVar[EventType]
    EVENT_TYPE_STOP: _ClassVar[EventType]

EVENT_TYPE_UNSPECIFIED: EventType
EVENT_TYPE_HELLO: EventType
EVENT_TYPE_EVALUATE: EventType
EVENT_TYPE_OPTIMIZE: EventType
EVENT_TYPE_INITIALIZE: EventType
EVENT_TYPE_SHARE_STATE: EventType
EVENT_TYPE_STOP: EventType

class Slice(_message.Message):
    __slots__ = ("start", "end")
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    start: int
    end: int
    def __init__(self, start: _Optional[int] = ..., end: _Optional[int] = ...) -> None: ...

class Evaluation(_message.Message):
    __slots__ = ("slice", "rewards")
    SLICE_FIELD_NUMBER: _ClassVar[int]
    REWARDS_FIELD_NUMBER: _ClassVar[int]
    slice: Slice
    rewards: _containers.RepeatedScalarFieldContainer[float]
    def __init__(
        self, slice: _Optional[_Union[Slice, _Mapping]] = ..., rewards: _Optional[_Iterable[float]] = ...
    ) -> None: ...

class HelloEvent(_message.Message):
    __slots__ = ("id", "token", "state", "population_size", "heartbeat_interval", "max_epochs")
    ID_FIELD_NUMBER: _ClassVar[int]
    TOKEN_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    POPULATION_SIZE_FIELD_NUMBER: _ClassVar[int]
    HEARTBEAT_INTERVAL_FIELD_NUMBER: _ClassVar[int]
    MAX_EPOCHS_FIELD_NUMBER: _ClassVar[int]
    id: str
    token: str
    state: bytes
    population_size: int
    heartbeat_interval: int
    max_epochs: int
    def __init__(
        self,
        id: _Optional[str] = ...,
        token: _Optional[str] = ...,
        state: _Optional[bytes] = ...,
        population_size: _Optional[int] = ...,
        heartbeat_interval: _Optional[int] = ...,
        max_epochs: _Optional[int] = ...,
    ) -> None: ...

class EvaluateEvent(_message.Message):
    __slots__ = ("task_id", "epoch", "slices")
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    EPOCH_FIELD_NUMBER: _ClassVar[int]
    SLICES_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    epoch: int
    slices: _containers.RepeatedCompositeFieldContainer[Slice]
    def __init__(
        self,
        task_id: _Optional[str] = ...,
        epoch: _Optional[int] = ...,
        slices: _Optional[_Iterable[_Union[Slice, _Mapping]]] = ...,
    ) -> None: ...

class OptimizeEvent(_message.Message):
    __slots__ = ("task_id", "epoch", "rewards")
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    EPOCH_FIELD_NUMBER: _ClassVar[int]
    REWARDS_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    epoch: int
    rewards: _containers.RepeatedScalarFieldContainer[float]
    def __init__(
        self, task_id: _Optional[str] = ..., epoch: _Optional[int] = ..., rewards: _Optional[_Iterable[float]] = ...
    ) -> None: ...

class InitializeEvent(_message.Message):
    __slots__ = ("task_id",)
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    def __init__(self, task_id: _Optional[str] = ...) -> None: ...

class ShareStateEvent(_message.Message):
    __slots__ = ("task_id", "epoch")
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    EPOCH_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    epoch: int
    def __init__(self, task_id: _Optional[str] = ..., epoch: _Optional[int] = ...) -> None: ...

class StopEvent(_message.Message):
    __slots__ = ("task_id",)
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    def __init__(self, task_id: _Optional[str] = ...) -> None: ...

class SubscribeRequest(_message.Message):
    __slots__ = ("cores",)
    CORES_FIELD_NUMBER: _ClassVar[int]
    cores: int
    def __init__(self, cores: _Optional[int] = ...) -> None: ...

class SubscribeResponse(_message.Message):
    __slots__ = ("type", "hello", "evaluate", "optimize", "initialize", "share_state", "stop")
    TYPE_FIELD_NUMBER: _ClassVar[int]
    HELLO_FIELD_NUMBER: _ClassVar[int]
    EVALUATE_FIELD_NUMBER: _ClassVar[int]
    OPTIMIZE_FIELD_NUMBER: _ClassVar[int]
    INITIALIZE_FIELD_NUMBER: _ClassVar[int]
    SHARE_STATE_FIELD_NUMBER: _ClassVar[int]
    STOP_FIELD_NUMBER: _ClassVar[int]
    type: EventType
    hello: HelloEvent
    evaluate: EvaluateEvent
    optimize: OptimizeEvent
    initialize: InitializeEvent
    share_state: ShareStateEvent
    stop: StopEvent
    def __init__(
        self,
        type: _Optional[_Union[EventType, str]] = ...,
        hello: _Optional[_Union[HelloEvent, _Mapping]] = ...,
        evaluate: _Optional[_Union[EvaluateEvent, _Mapping]] = ...,
        optimize: _Optional[_Union[OptimizeEvent, _Mapping]] = ...,
        initialize: _Optional[_Union[InitializeEvent, _Mapping]] = ...,
        share_state: _Optional[_Union[ShareStateEvent, _Mapping]] = ...,
        stop: _Optional[_Union[StopEvent, _Mapping]] = ...,
    ) -> None: ...

class HeartbeatRequest(_message.Message):
    __slots__ = ("seq_id", "timestamp")
    SEQ_ID_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    seq_id: int
    timestamp: _timestamp_pb2.Timestamp
    def __init__(
        self, seq_id: _Optional[int] = ..., timestamp: _Optional[_Union[_timestamp_pb2.Timestamp, _Mapping]] = ...
    ) -> None: ...

class HeartbeatResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...

class FinishEvaluationRequest(_message.Message):
    __slots__ = ("task_id", "evaluations")
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    EVALUATIONS_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    evaluations: _containers.RepeatedCompositeFieldContainer[Evaluation]
    def __init__(
        self, task_id: _Optional[str] = ..., evaluations: _Optional[_Iterable[_Union[Evaluation, _Mapping]]] = ...
    ) -> None: ...

class FinishEvaluationResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...

class FinishOptimizationRequest(_message.Message):
    __slots__ = ("task_id",)
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    def __init__(self, task_id: _Optional[str] = ...) -> None: ...

class FinishOptimizationResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...

class FinishInitializationRequest(_message.Message):
    __slots__ = ("task_id", "state")
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    state: bytes
    def __init__(self, task_id: _Optional[str] = ..., state: _Optional[bytes] = ...) -> None: ...

class FinishInitializationResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...

class FinishShareStateRequest(_message.Message):
    __slots__ = ("task_id", "state")
    TASK_ID_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    task_id: str
    state: bytes
    def __init__(self, task_id: _Optional[str] = ..., state: _Optional[bytes] = ...) -> None: ...

class FinishShareStateResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: bool = ...) -> None: ...
