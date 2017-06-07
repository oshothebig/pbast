package yang

import (
	"reflect"
	"testing"

	"github.com/oshothebig/pbast"
)

func TestRemoveRootMessage(t *testing.T) {
	table := []struct {
		in       *pbast.File
		expected *pbast.File
	}{
		{
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")),
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")),
		},
		{
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")).
				AddMessage(pbast.NewMessage("Root")),
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")),
		},
		{
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")).
				AddMessage(pbast.NewMessage("Root").
					AddField(pbast.NewMessageField(pbast.Bool, "enabled", 1))),
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")).
				AddMessage(pbast.NewMessage("Root").
					AddField(pbast.NewMessageField(pbast.Bool, "enabled", 1))),
		},
		{
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")).
				AddMessage(pbast.NewMessage("Root").
					AddMessage(pbast.NewMessage("Sub"))),
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")).
				AddMessage(pbast.NewMessage("Sub")),
		},
		{
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")).
				AddMessage(pbast.NewMessage("Root").
					AddMessage(pbast.NewMessage("Sub").
						AddField(pbast.NewMessageField(pbast.Bool, "enabled", 1)))),
			pbast.NewFile("org.example").
				AddMessage(pbast.NewMessage("Origin")).
				AddMessage(pbast.NewMessage("Sub").
					AddField(pbast.NewMessageField(pbast.Bool, "enabled", 1))),
		},
	}

	for x, d := range table {
		if actual := RemoveRootMessage(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}
