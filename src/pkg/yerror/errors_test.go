package yerror_test

import (
	"errors"
	"testing"

	"redistore/pkg/yerror"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestE(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want *yerror.Error
	}{
		{
			"Simple",
			args{args: []interface{}{yerror.Op("blog.findByID"), yerror.KindNotFound, errors.New("not found error")}},
			&yerror.Error{
				Op:   "blog.findByID",
				Kind: yerror.KindNotFound,
				Err:  errors.New("not found error"),
			},
		},
		{
			"Nested",
			args{args: []interface{}{
				yerror.Op("blog.create"),
				yerror.KindUnauthorized,
				&yerror.Error{
					Op:   "account.getUser",
					Kind: yerror.KindNotFound,
					Err:  errors.New("user not found error"),
				},
			}},
			&yerror.Error{
				Op:   "blog.create",
				Kind: yerror.KindUnauthorized,
				Err: &yerror.Error{
					Op:   "account.getUser",
					Kind: yerror.KindNotFound,
					Err:  errors.New("user not found error"),
				},
			},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, yerror.E(tt.args.args...), tt.name)
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Op   yerror.Op
		Kind codes.Code
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"Error with nested error",
			fields{
				Op:   "blog.findByID",
				Kind: yerror.KindNotFound,
				Err:  errors.New("blog not found"),
			},
			"blog not found",
		},
		{
			"Error with nested Error",
			fields{
				Op:   "blog.findByID",
				Kind: yerror.KindNotFound,
				Err:  yerror.E(yerror.Op("account.getUser"), errors.New("unexpected error")),
			},
			"unexpected error",
		},
	}
	for _, tt := range tests {
		e := &yerror.Error{
			Op:   tt.fields.Op,
			Kind: tt.fields.Kind,
			Err:  tt.fields.Err,
		}
		assert.Equal(t, tt.want, e.Error(), tt.name)
	}
}

func TestOps(t *testing.T) {
	type args struct {
		e *yerror.Error
	}
	tests := []struct {
		name string
		args args
		want []yerror.Op
	}{
		{
			"Nested Errors",
			args{e: yerror.E(yerror.Op("blog.findByID"), yerror.E(yerror.Op("account.getUser")))},
			[]yerror.Op{"blog.findByID", "account.getUser"},
		},
		{
			"Error with nested error",
			args{e: yerror.E(yerror.Op("blog.findByID"), errors.New("unexpected error"))},
			[]yerror.Op{"blog.findByID"},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, yerror.Ops(tt.args.e), tt.name)
	}
}

func TestLevel(t *testing.T) {
	const op yerror.Op = "core.yerror.TestLevel"

	testCases := []struct {
		desc string
		err  error
		want logrus.Level
	}{
		{
			desc: "native error",
			err:  errors.New("simple error has no level"),
			want: logrus.ErrorLevel,
		},
		{
			desc: "no level yerror",
			err:  yerror.E(op, errors.New("simple error")),
			want: logrus.ErrorLevel,
		},
		{
			desc: "no child yerror",
			err:  yerror.E(op),
			want: logrus.ErrorLevel,
		},
		{
			desc: "debug level yerror",
			err:  yerror.E(op, yerror.LevelDebug, errors.New("simple error")),
			want: logrus.DebugLevel,
		},
		{
			desc: "warn level yerror",
			err:  yerror.E(op, yerror.LevelWarn, errors.New("simple error")),
			want: logrus.WarnLevel,
		},
		{
			desc: "info level yerror",
			err:  yerror.E(op, yerror.LevelInfo, errors.New("simple error")),
			want: logrus.InfoLevel,
		},
		{
			desc: "error level yerror",
			err:  yerror.E(op, yerror.LevelError, errors.New("simple error")),
			want: logrus.ErrorLevel,
		},
		{
			desc: "nested yerror",
			err:  yerror.E(op, yerror.E(op, yerror.LevelDebug, errors.New("simple error"))),
			want: logrus.DebugLevel,
		},
	}
	for _, tC := range testCases {
		assert.EqualValues(t, tC.want, yerror.Level(tC.err), tC.desc)
	}
}

func TestKind(t *testing.T) {
	const op yerror.Op = "core.yerror.TestKind"

	simpleError := errors.New("simple error")

	testCases := []struct {
		desc string
		err  error
		want codes.Code
	}{
		{
			desc: "native error",
			err:  errors.New("simple error has no kind"),
			want: codes.Unknown,
		},
		{
			desc: "yerror with no kind",
			err:  yerror.E(op, simpleError),
			want: codes.Unknown,
		},
		{
			desc: "yerror without child",
			err:  yerror.E(op),
			want: codes.Unknown,
		},
		{
			desc: "not found",
			err:  yerror.E(op, yerror.KindNotFound, simpleError),
			want: codes.NotFound,
		},
		{
			desc: "invalid argument",
			err:  yerror.E(op, yerror.KindInvalidArgument, simpleError),
			want: codes.InvalidArgument,
		},
		{
			desc: "unauthenticated",
			err:  yerror.E(op, yerror.KindUnauthenticated, simpleError),
			want: codes.Unauthenticated,
		},
		{
			desc: "unauthorized",
			err:  yerror.E(op, yerror.KindUnauthorized, simpleError),
			want: codes.PermissionDenied,
		},
		{
			desc: "internal",
			err:  yerror.E(op, yerror.KindInternal, simpleError),
			want: codes.Internal,
		},
		{
			desc: "unexpected",
			err:  yerror.E(op, yerror.KindUnexpected, simpleError),
			want: codes.Unknown,
		},
		{
			desc: "nested yerror",
			err:  yerror.E(op, yerror.E(op, yerror.KindInternal, simpleError)),
			want: codes.Internal,
		},
	}
	for _, tC := range testCases {
		assert.EqualValues(t, tC.want, yerror.Kind(tC.err), tC.desc)
	}
}
