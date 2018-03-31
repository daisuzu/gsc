package checker

import (
	"testing"

	"honnef.co/go/tools/lint/testutil"
)

func TestAll(t *testing.T) {
	c := New(WithAdditionalContexts("MyCtx"))
	testutil.TestAll(t, c, "")
}
