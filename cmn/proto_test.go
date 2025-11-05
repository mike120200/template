package cmn

import (
	"errors"
	"testing"
)

func TestReplyProto_NewErrorReply(t *testing.T) {
	appErr := NewAppError(-3, "success")
	r := NewErrorReply(appErr, "API", "Method")
	t.Logf("result: %+v", r)

	err := errors.New("error")
	r = NewErrorReply(err, "API", "Method")
	t.Logf("result: %+v", r)
}
