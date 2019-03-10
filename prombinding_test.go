package prombinding

import (
	"reflect"
	"testing"

	"github.com/prometheus/common/model"
)

func TestLabelMissing(t *testing.T) {
	type point struct {
		Value     float64
		Time      int64
		Namespace string `prom:"namespace"`
	}
	cases := []struct {
		sample *model.Sample
		err    bool
		result point
	}{
		{
			sample: &model.Sample{
				Metric: model.Metric{
					"namespace": "bookinfo",
				},
				Value:     model.SampleValue(1),
				Timestamp: model.Time(10),
			},
			err: false,
			result: point{
				Value:     1,
				Time:      10,
				Namespace: "bookinfo",
			},
		},
		{
			sample: &model.Sample{
				Metric: model.Metric{
					"namespacex": "bookinfo",
				},
				Value:     model.SampleValue(1),
				Timestamp: model.Time(10),
			},
			err: true,
			result: point{
				Value:     1,
				Time:      10,
				Namespace: "",
			},
		},
	}
	for i, cas := range cases {
		t.Log("Case", i)
		var target point
		err := Bind(cas.sample, &target)
		if cas.err && err == nil {
			t.Errorf("Expect error, got nil")
		}
		if !cas.err && err != nil {
			t.Errorf("Expect nil, got error")
		}
		if !reflect.DeepEqual(cas.result, target) {
			t.Errorf("Expect %#v,\n got %#v", cas.result, target)
		}
	}
}
