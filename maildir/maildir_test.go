package maildir_test

import (
	"reflect"
	"testing"

	"github.com/tennashi/goem/maildir"
)

func Test_ParseKey(t *testing.T) {
	cases := map[string]struct {
		input string
		want  maildir.Key
		err   bool
	}{
		"(valid)full format": {
			input: "123123.#123123X123123R123123I123123V123123M123123P123123Q123123.hostname,U=123123,S=123123:2,AB",
			want: maildir.Key{
				Raw:    "123123.#123123X123123R123123I123123V123123M123123P123123Q123123.hostname,U=123123,S=123123:2,AB",
				Second: 123123,
				DeliveryID: maildir.ID{
					UNIXSeq:     123123,
					Boot:        123123,
					Urandom:     123123,
					Inode:       123123,
					Dev:         123123,
					MicroSecond: 123123,
					PID:         123123,
					Seq:         123123,
				},
				HostName: "hostname",
				Params: map[string]string{
					"U": "123123",
					"S": "123123",
				},
				FlagType: maildir.FlagTypeNormal,
				Flags:    []string{"A", "B"},
			},
			err: false,
		},
	}
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			got, err := maildir.ParseKey(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("\n\tgot: %v\n\twant: %v", got, tt.want)
			}
		})
	}
}

func Test_ParseID(t *testing.T) {
	cases := map[string]struct {
		input string
		want  maildir.ID
		err   bool
	}{
		"(valid)full format": {
			input: "#123123X123123R123123I123123V123123M123123P123123Q123123",
			want: maildir.ID{
				UNIXSeq:     123123,
				Boot:        123123,
				Urandom:     123123,
				Inode:       123123,
				Dev:         123123,
				MicroSecond: 123123,
				PID:         123123,
				Seq:         123123,
			},
			err: false,
		},
		"(valid)specific parts": {
			input: "#123123",
			want: maildir.ID{
				UNIXSeq: 123123,
			},
			err: false,
		},
		"(valid)shorter old fashioned": {
			input: "123123",
			want: maildir.ID{
				PID: 123123,
			},
			err: false,
		},
		"(valid)longer old fashioned": {
			input: "123123_123123",
			want: maildir.ID{
				PID: 123123,
				Seq: 123123,
			},
			err: false,
		},
		"(invalid)empty": {
			input: "",
			want:  maildir.ID{},
			err:   true,
		},
		"(invalid)strings": {
			input: "hogehoge",
			want:  maildir.ID{},
			err:   true,
		},
		"(invalid)lower case": {
			input: "x123123r123123",
			want:  maildir.ID{},
			err:   true,
		},
	}
	for caseName, tt := range cases {
		t.Run(caseName, func(t *testing.T) {
			got, err := maildir.ParseID(tt.input)
			if !tt.err && err != nil {
				t.Fatalf("should not be error for %v but %v", caseName, err)
			}
			if tt.err && err == nil {
				t.Fatalf("should be error for %v but not", caseName)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("\n\tgot: %v\n\twant: %v", got, tt.want)
			}
		})
	}
}
