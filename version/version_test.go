package version

import (
	"bytes"
	"testing"

	"github.com/go-openapi/testify/v2/require"
)

func TestInfoOutput(t *testing.T) {
	t.Parallel()

	info := Info{
		Build: Build{
			Version: "foo",
			Commit:  "commit",
		},
		Runtime: Runtime{},
	}
	buf := &bytes.Buffer{}
	err := info.Output(buf, OutputFormatText)
	require.NoError(t, err)
	require.EqualT(t, "version.test version foo commit\n", buf.String())

	info = Info{}
	err = info.Output(nil, "foo")
	require.EqualError(t, err, "unknown output format foo")
}
