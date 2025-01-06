# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""

import grpc
import warnings

from evochi.v1 import evochi_pb2 as evochi_dot_v1_dot_evochi__pb2

GRPC_GENERATED_VERSION = "1.67.0"
GRPC_VERSION = grpc.__version__
_version_not_supported = False

try:
    from grpc._utilities import first_version_is_lower

    _version_not_supported = first_version_is_lower(GRPC_VERSION, GRPC_GENERATED_VERSION)
except ImportError:
    _version_not_supported = True

if _version_not_supported:
    raise RuntimeError(
        f"The grpc package installed is at version {GRPC_VERSION},"
        + f" but the generated code in evochi/v1/evochi_pb2_grpc.py depends on"
        + f" grpcio>={GRPC_GENERATED_VERSION}."
        + f" Please upgrade your grpc module to grpcio>={GRPC_GENERATED_VERSION}"
        + f" or downgrade your generated code using grpcio-tools<={GRPC_VERSION}."
    )


class EvochiServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Subscribe = channel.unary_stream(
            "/evochi.v1.EvochiService/Subscribe",
            request_serializer=evochi_dot_v1_dot_evochi__pb2.SubscribeRequest.SerializeToString,
            response_deserializer=evochi_dot_v1_dot_evochi__pb2.SubscribeResponse.FromString,
            _registered_method=True,
        )
        self.Heartbeat = channel.unary_unary(
            "/evochi.v1.EvochiService/Heartbeat",
            request_serializer=evochi_dot_v1_dot_evochi__pb2.HeartbeatRequest.SerializeToString,
            response_deserializer=evochi_dot_v1_dot_evochi__pb2.HeartbeatResponse.FromString,
            _registered_method=True,
        )
        self.FinishEvaluation = channel.unary_unary(
            "/evochi.v1.EvochiService/FinishEvaluation",
            request_serializer=evochi_dot_v1_dot_evochi__pb2.FinishEvaluationRequest.SerializeToString,
            response_deserializer=evochi_dot_v1_dot_evochi__pb2.FinishEvaluationResponse.FromString,
            _registered_method=True,
        )
        self.FinishOptimization = channel.unary_unary(
            "/evochi.v1.EvochiService/FinishOptimization",
            request_serializer=evochi_dot_v1_dot_evochi__pb2.FinishOptimizationRequest.SerializeToString,
            response_deserializer=evochi_dot_v1_dot_evochi__pb2.FinishOptimizationResponse.FromString,
            _registered_method=True,
        )
        self.FinishInitialization = channel.unary_unary(
            "/evochi.v1.EvochiService/FinishInitialization",
            request_serializer=evochi_dot_v1_dot_evochi__pb2.FinishInitializationRequest.SerializeToString,
            response_deserializer=evochi_dot_v1_dot_evochi__pb2.FinishInitializationResponse.FromString,
            _registered_method=True,
        )
        self.FinishShareState = channel.unary_unary(
            "/evochi.v1.EvochiService/FinishShareState",
            request_serializer=evochi_dot_v1_dot_evochi__pb2.FinishShareStateRequest.SerializeToString,
            response_deserializer=evochi_dot_v1_dot_evochi__pb2.FinishShareStateResponse.FromString,
            _registered_method=True,
        )


class EvochiServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Subscribe(self, request, context):
        """join the work force and subscribe to events"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def Heartbeat(self, request, context):
        """send heartbeat to the server to keep the connection alive"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def FinishEvaluation(self, request, context):
        """finish the evaluation"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def FinishOptimization(self, request, context):
        """finish the optimization"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def FinishInitialization(self, request, context):
        """finish the initialization"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")

    def FinishShareState(self, request, context):
        """finish the state sharing"""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details("Method not implemented!")
        raise NotImplementedError("Method not implemented!")


def add_EvochiServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
        "Subscribe": grpc.unary_stream_rpc_method_handler(
            servicer.Subscribe,
            request_deserializer=evochi_dot_v1_dot_evochi__pb2.SubscribeRequest.FromString,
            response_serializer=evochi_dot_v1_dot_evochi__pb2.SubscribeResponse.SerializeToString,
        ),
        "Heartbeat": grpc.unary_unary_rpc_method_handler(
            servicer.Heartbeat,
            request_deserializer=evochi_dot_v1_dot_evochi__pb2.HeartbeatRequest.FromString,
            response_serializer=evochi_dot_v1_dot_evochi__pb2.HeartbeatResponse.SerializeToString,
        ),
        "FinishEvaluation": grpc.unary_unary_rpc_method_handler(
            servicer.FinishEvaluation,
            request_deserializer=evochi_dot_v1_dot_evochi__pb2.FinishEvaluationRequest.FromString,
            response_serializer=evochi_dot_v1_dot_evochi__pb2.FinishEvaluationResponse.SerializeToString,
        ),
        "FinishOptimization": grpc.unary_unary_rpc_method_handler(
            servicer.FinishOptimization,
            request_deserializer=evochi_dot_v1_dot_evochi__pb2.FinishOptimizationRequest.FromString,
            response_serializer=evochi_dot_v1_dot_evochi__pb2.FinishOptimizationResponse.SerializeToString,
        ),
        "FinishInitialization": grpc.unary_unary_rpc_method_handler(
            servicer.FinishInitialization,
            request_deserializer=evochi_dot_v1_dot_evochi__pb2.FinishInitializationRequest.FromString,
            response_serializer=evochi_dot_v1_dot_evochi__pb2.FinishInitializationResponse.SerializeToString,
        ),
        "FinishShareState": grpc.unary_unary_rpc_method_handler(
            servicer.FinishShareState,
            request_deserializer=evochi_dot_v1_dot_evochi__pb2.FinishShareStateRequest.FromString,
            response_serializer=evochi_dot_v1_dot_evochi__pb2.FinishShareStateResponse.SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler("evochi.v1.EvochiService", rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))
    server.add_registered_method_handlers("evochi.v1.EvochiService", rpc_method_handlers)


# This class is part of an EXPERIMENTAL API.
class EvochiService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Subscribe(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_stream(
            request,
            target,
            "/evochi.v1.EvochiService/Subscribe",
            evochi_dot_v1_dot_evochi__pb2.SubscribeRequest.SerializeToString,
            evochi_dot_v1_dot_evochi__pb2.SubscribeResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True,
        )

    @staticmethod
    def Heartbeat(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/evochi.v1.EvochiService/Heartbeat",
            evochi_dot_v1_dot_evochi__pb2.HeartbeatRequest.SerializeToString,
            evochi_dot_v1_dot_evochi__pb2.HeartbeatResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True,
        )

    @staticmethod
    def FinishEvaluation(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/evochi.v1.EvochiService/FinishEvaluation",
            evochi_dot_v1_dot_evochi__pb2.FinishEvaluationRequest.SerializeToString,
            evochi_dot_v1_dot_evochi__pb2.FinishEvaluationResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True,
        )

    @staticmethod
    def FinishOptimization(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/evochi.v1.EvochiService/FinishOptimization",
            evochi_dot_v1_dot_evochi__pb2.FinishOptimizationRequest.SerializeToString,
            evochi_dot_v1_dot_evochi__pb2.FinishOptimizationResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True,
        )

    @staticmethod
    def FinishInitialization(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/evochi.v1.EvochiService/FinishInitialization",
            evochi_dot_v1_dot_evochi__pb2.FinishInitializationRequest.SerializeToString,
            evochi_dot_v1_dot_evochi__pb2.FinishInitializationResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True,
        )

    @staticmethod
    def FinishShareState(
        request,
        target,
        options=(),
        channel_credentials=None,
        call_credentials=None,
        insecure=False,
        compression=None,
        wait_for_ready=None,
        timeout=None,
        metadata=None,
    ):
        return grpc.experimental.unary_unary(
            request,
            target,
            "/evochi.v1.EvochiService/FinishShareState",
            evochi_dot_v1_dot_evochi__pb2.FinishShareStateRequest.SerializeToString,
            evochi_dot_v1_dot_evochi__pb2.FinishShareStateResponse.FromString,
            options,
            channel_credentials,
            insecure,
            call_credentials,
            compression,
            wait_for_ready,
            timeout,
            metadata,
            _registered_method=True,
        )
