package domain

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type money struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

var (
	testInstanceContext = map[string]interface{}{
		"instance_id": uuid.MustParse("88f60b3c-0a1b-11ef-8f21-0b0a3415e515"),
		"order_id":    uuid.MustParse("dc270aa0-0a1a-11ef-b099-2be64697d22e"),
		"pay_system":  "SBP",
		"order_data": map[string]interface{}{
			"items": []uuid.UUID{
				uuid.MustParse("011785ec-0a1b-11ef-8c55-7f8cc2b0e46f"),
				uuid.MustParse("04819948-0a1b-11ef-9ccd-e7b03a48f8ef"),
				uuid.MustParse("0b66506e-0a1b-11ef-be47-4391f102dbdf"),
			},
			"masked_pan": "1234 5XXX XXXX 6789",
			"total_amount": money{
				Amount:       "6666.6",
				CurrencyCode: "RUB",
			},
		},
	}
)

func TestNewInstanceContext(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "success",
			data: []byte(`{"name": "Name", "last_name": "Last Name"}`),
			want: map[string]interface{}{
				"name":      "Name",
				"last_name": "Last Name",
			},
			wantErr: false,
		},
		{
			name:    "fail",
			data:    []byte("INVALID JSON"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := NewInstanceContext(tc.data)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewInstanceContext() wantErr = %v, err = %v", tc.wantErr, err)
			}
			reflect.DeepEqual(actual, tc.want)
		})
	}
}

func TestInstanceContext_MergeData(t *testing.T) {
	tests := []struct {
		name        string
		source      []byte
		dataToMerge []byte
		want        map[string]interface{}
		wantErr     bool
	}{
		{
			name:        "success simple",
			source:      []byte(`{"value": "321"}`),
			dataToMerge: []byte(`{"merged_value": "123"}`),
			want: map[string]interface{}{
				"value":        "321",
				"merged_value": "123",
			},
			wantErr: false,
		},
		{
			name:        "success overwrite value",
			source:      []byte(`{"value": "321"}`),
			dataToMerge: []byte(`{"value": "123"}`),
			want: map[string]interface{}{
				"value": "123",
			},
			wantErr: false,
		},
		{
			name:        "success overwrite object",
			source:      []byte(`{"value": "321", "obj": {"item1": "val1", "item2": "val2"}}`),
			dataToMerge: []byte(`{"value": "123", "obj": {"item1": "new_val"}}`),
			want: map[string]interface{}{
				"value": "123",
				"obj": map[string]interface{}{
					"item1": "new_val",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := NewInstanceContext(tc.source)
			assert.NoError(t, err)
			err = ctx.MergeData(tc.dataToMerge)

			if (err != nil) != tc.wantErr {
				t.Errorf("MergeData() wantErr = %v, err = %v", tc.wantErr, err)
			}

			reflect.DeepEqual(ctx.data, tc.want)
		})
	}
}

func TestInstanceContext_GetValue(t *testing.T) {
	tests := []struct {
		path string
		want interface{}
	}{
		{
			path: "instance_id",
			want: testInstanceContext["instance_id"],
		},
		{
			path: "order_id",
			want: testInstanceContext["order_id"],
		},
		{
			path: "pay_system",
			want: testInstanceContext["pay_system"],
		},
		{
			path: "order_data",
			want: testInstanceContext["order_data"],
		},
		{
			path: "order_data.items",
			want: testInstanceContext["order_data"].(map[string]interface{})["items"],
		},
		{
			path: "order_data.masked_pan",
			want: testInstanceContext["order_data"].(map[string]interface{})["masked_pan"],
		},
		{
			path: "order_data.total_amount",
			want: testInstanceContext["order_data"].(map[string]interface{})["total_amount"],
		},
		{
			path: "order_data.total_amount.amount",
			want: testInstanceContext["order_data"].(map[string]interface{})["total_amount"].(money).Amount,
		},
		{
			path: "order_data.total_amount.currency_code",
			want: testInstanceContext["order_data"].(map[string]interface{})["total_amount"].(money).CurrencyCode,
		},
		{
			path: "non_existence",
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			marshalled, err := json.Marshal(testInstanceContext)
			assert.NoError(t, err)
			ctx, err := NewInstanceContext(marshalled)
			assert.NoError(t, err)

			actual, err := ctx.GetValue(tc.path)
			assert.NoError(t, err)
			reflect.DeepEqual(actual, tc.want)
		})
	}
}
