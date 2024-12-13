package js_test

import (
	"io"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/aserto-dev/ds-load/sdk"
	"github.com/aserto-dev/ds-load/sdk/common/js"
	"github.com/aserto-dev/ds-load/sdk/common/msg"
)

func BenchmarkReadProtoMessage(b *testing.B) {
	assert := assert.New(b)

	b.ResetTimer()

	f, err := sdk.Assets().Open("assets/acmecorp.json")
	assert.NoError(err)
	defer f.Close()
	reader, err := js.NewJSONArrayReader(f)
	assert.NoError(err)

	for {
		var message msg.Transform
		err := reader.ReadProtoMessage(&message)
		if err == io.EOF {
			break
		}
		assert.NoError(err)
	}
}
