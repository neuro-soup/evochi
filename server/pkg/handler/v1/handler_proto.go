package v1

import (
	"github.com/neuro-soup/evochi/server/internal/training/eval"
	evochiv1 "github.com/neuro-soup/evochi/server/pkg/proto/evochi/v1"
)

func slicesToProto(slices []eval.Slice) []*evochiv1.Slice {
	out := make([]*evochiv1.Slice, len(slices))
	for i, slice := range slices {
		out[i] = &evochiv1.Slice{
			Start: int32(slice.Start),
			End:   int32(slice.End),
		}
	}
	return out
}

func protoToSlice(proto *evochiv1.Slice) eval.Slice {
	return eval.Slice{
		Start: uint(proto.Start),
		End:   uint(proto.End),
	}
}

func protoToEval(proto *evochiv1.Evaluation) eval.Eval {
	return eval.Eval{
		Slice:   protoToSlice(proto.Slice),
		Rewards: f32ToF64(proto.Rewards),
	}
}

func protoToEvals(proto []*evochiv1.Evaluation) []eval.Eval {
	out := make([]eval.Eval, len(proto))
	for i, e := range proto {
		out[i] = protoToEval(e)
	}
	return out
}

func f32ToF64(f []float32) []float64 {
	out := make([]float64, len(f))
	for i, f := range f {
		out[i] = float64(f)
	}
	return out
}

func f64ToF32(f []float64) []float32 {
	out := make([]float32, len(f))
	for i, f := range f {
		out[i] = float32(f)
	}
	return out
}
