package yang

import (
	"reflect"
	"testing"

	"github.com/oshothebig/pbast"
)

func TestCompleteZeroInEnum(t *testing.T) {
	table := []struct {
		in       *pbast.File
		expected *pbast.File
	}{
		{
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))),
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))),
		},
		{
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 1))),
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("DEFAULT", 0)).
					AddField(pbast.NewEnumField("V1", 1))),
		},
		{
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("V2", 0))),
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("V2", 0))),
		},
		{
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 1))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("V2", 1))),
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("DEFAULT", 0)).
					AddField(pbast.NewEnumField("V1", 1))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("DEFAULT", 0)).
					AddField(pbast.NewEnumField("V2", 1))),
		},
		{
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("V2", 1))),
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("DEFAULT", 0)).
					AddField(pbast.NewEnumField("V2", 1))),
		},
		{
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("V2", 1))).
				AddMessage(pbast.NewMessage("M1").
					AddEnum(pbast.NewEnum("E1").
						AddField(pbast.NewEnumField("V1", 0)))),
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("DEFAULT", 0)).
					AddField(pbast.NewEnumField("V2", 1))).
				AddMessage(pbast.NewMessage("M1").
					AddEnum(pbast.NewEnum("E1").
						AddField(pbast.NewEnumField("V1", 0)))),
		},
		{
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("V2", 1))).
				AddMessage(pbast.NewMessage("M1").
					AddEnum(pbast.NewEnum("E1").
						AddField(pbast.NewEnumField("V1", 1)))),
			pbast.NewFile("org.example").
				AddEnum(pbast.NewEnum("E1").
					AddField(pbast.NewEnumField("V1", 0))).
				AddEnum(pbast.NewEnum("E2").
					AddField(pbast.NewEnumField("DEFAULT", 0)).
					AddField(pbast.NewEnumField("V2", 1))).
				AddMessage(pbast.NewMessage("M1").
					AddEnum(pbast.NewEnum("E1").
						AddField(pbast.NewEnumField("DEFAULT", 0)).
						AddField(pbast.NewEnumField("V1", 1)))),
		},
	}

	for x, d := range table {
		if actual := CompleteZeroInEnum(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}
