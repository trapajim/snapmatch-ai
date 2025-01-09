package asset

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/trapajim/snapmatch-ai/snapmatchai"
	"github.com/trapajim/snapmatch-ai/snapmatchai/mocks"
	"io"
	"strings"
	"testing"
)

func TestService_Upload(t *testing.T) {
	type fields struct {
		appContext snapmatchai.Context
	}
	type args struct {
		ctx    context.Context
		file   io.Reader
		object string
	}
	tests := []struct {
		name             string
		mockExpectations func(fields, args)
		fields           fields
		args             args
		wantErr          assert.ErrorAssertionFunc
		targetErr        error
	}{
		{
			name: "Test Upload",
			fields: fields{
				appContext: snapmatchai.NewContextForTest(t),
			},
			mockExpectations: func(fields fields, args args) {
				fields.appContext.Storage.(*mocks.MockUploader).EXPECT().Upload(args.ctx, args.file, args.object).Return(nil).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				file:   strings.NewReader("test file content"),
				object: "test-object.txt",
			},
			wantErr: assert.NoError,
		},
		{
			name: "Test Upload - error handled",
			fields: fields{
				appContext: snapmatchai.NewContextForTest(t),
			},
			mockExpectations: func(fields fields, args args) {
				fields.appContext.Storage.(*mocks.MockUploader).EXPECT().Upload(args.ctx, args.file, args.object).Return(snapmatchai.NewError(errors.New("error"), "some error", 400)).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				file:   strings.NewReader("test file content"),
				object: "test-object.txt",
			},
			wantErr:   assert.Error,
			targetErr: &UploadError{},
		},
		{
			name: "Test Upload - unexpected error handled",
			fields: fields{
				appContext: snapmatchai.NewContextForTest(t),
			},
			mockExpectations: func(fields fields, args args) {
				fields.appContext.Storage.(*mocks.MockUploader).EXPECT().Upload(args.ctx, args.file, args.object).Return(errors.New("unexpected error")).Times(1)
			},
			args: args{
				ctx:    context.Background(),
				file:   strings.NewReader("test file content"),
				object: "test-object.txt",
			},
			wantErr:   assert.Error,
			targetErr: &snapmatchai.Error{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				appContext: tt.fields.appContext,
			}
			tt.mockExpectations(tt.fields, tt.args)
			err := s.Upload(tt.args.ctx, tt.args.file, tt.args.object)
			tt.wantErr(t, err)
			if tt.targetErr != nil {
				assert.ErrorAs(t, err, &tt.targetErr)
			}
		})
	}
}
