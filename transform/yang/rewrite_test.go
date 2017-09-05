package yang

import (
	"reflect"
	"testing"

	"github.com/oshothebig/pbast"
)

func TestLiftMessage(t *testing.T) {
	table := []struct {
		in       *pbast.File
		expected *pbast.File
	}{
		{
			&pbast.File{},
			&pbast.File{},
		},
		{
			&pbast.File{
				Messages: []*pbast.Message{
					&pbast.Message{
						Name: "top",
					},
				},
			},
			&pbast.File{
				Messages: []*pbast.Message{
					&pbast.Message{
						Name: "top",
					},
				},
			},
		},
		{
			&pbast.File{
				Messages: []*pbast.Message{
					&pbast.Message{
						Name: "top",
						Messages: []*pbast.Message{
							&pbast.Message{
								Name: "second",
							},
						},
					},
				},
			},
			&pbast.File{
				Messages: []*pbast.Message{
					&pbast.Message{
						Name: "top",
					},
					&pbast.Message{
						Name: "second",
					},
				},
			},
		},
	}

	for x, d := range table {
		if actual := LiftMessage(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}

func TestCountName(t *testing.T) {
	table := []struct {
		in       *pbast.File
		expected map[string]int
	}{
		{
			&pbast.File{},
			map[string]int{},
		},
		{
			&pbast.File{
				Messages: []*pbast.Message{
					&pbast.Message{
						Name: "top",
					},
				},
			},
			map[string]int{
				"top": 1,
			},
		},
		{
			&pbast.File{
				Messages: []*pbast.Message{
					&pbast.Message{
						Name: "top",
						Messages: []*pbast.Message{
							&pbast.Message{
								Name: "second",
							},
						},
					},
				},
			},
			map[string]int{
				"top":    1,
				"second": 1,
			},
		},
	}

	for x, d := range table {
		if actual := countName(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}

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
