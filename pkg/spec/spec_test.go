package spec

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestTriggerDecode(t *testing.T) {
	examples := map[string][]string{
		`message: plain1`:                     {"plain1"},
		`message: ["arr1", "arr2"]`:           {"arr1", "arr2"},
		`message: {text: text1}`:              {"text1"},
		`message: {text: ["tarr1", "tarr2"]}`: {"tarr1", "tarr2"},
	}
	var cnt int
	for src, exp := range examples {
		t.Run(fmt.Sprintf("case-%d", cnt), func(t *testing.T) {
			var tr Trigger
			err := yaml.Unmarshal([]byte(src), &tr)
			require.NoError(t, err)
			require.Equal(t, exp, tr.Message.Text)
		})
	}
}
