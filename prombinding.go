package prombinding

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/prometheus/common/model"
)

// BindError include details for binding error
type BindError struct {
	Messages []string
}

func (e *BindError) Error() string {
	return strings.Join(e.Messages, ";")
}

// Bind binds a sample to a struct
func Bind(sample *model.Sample, obj interface{}) error {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("expect obj to be a pointer to struct")
	}
	t = t.Elem()
	v = v.Elem()
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expect obj to be a pointer to struct")
	}
	bErr := &BindError{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		vv := v.Field(i)
		if field.Type.Kind() == reflect.Float64 {
			vv.SetFloat(float64(sample.Value))
			continue
		}
		if field.Type.Kind() == reflect.Int64 {
			vv.SetInt(int64(sample.Timestamp))
			continue
		}
		// TODO: allow bind to more kinds
		if field.Type.Kind() != reflect.String {
			continue
		}
		labelName := field.Tag.Get("prom")
		if labelName == "" {
			continue
		}
		if labelValue, exist := sample.Metric[model.LabelName(labelName)]; !exist {
			bErr.Messages = append(bErr.Messages, fmt.Sprintf("label %s is missing", labelName))
		} else {
			lv := string(labelValue)
			vv.SetString(lv)
		}
	}
	return bErr
}
