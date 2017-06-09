package pbast

import (
	"reflect"
	"testing"
)

func TestFlatten(t *testing.T) {
	table := []struct {
		in       *Message
		expected []*Message
	}{
		{
			NewMessage("A"),
			[]*Message{NewMessage("A")},
		},
		{
			NewMessage("A").
				AddMessage(NewMessage("B")),
			[]*Message{NewMessage("A"), NewMessage("B")},
		},
		{
			NewMessage("A").
				AddMessage(NewMessage("B").
					AddMessage(NewMessage("C"))).
				AddMessage(NewMessage("D")),
			[]*Message{NewMessage("A"), NewMessage("B"), NewMessage("D"), NewMessage("C")},
		},
		{
			NewMessage("A").
				AddMessage(NewMessage("B").
					AddMessage(NewMessage("B"))).
				AddMessage(NewMessage("C")),
			[]*Message{
				NewMessage("A"),
				NewMessage("B").
					AddMessage(NewMessage("B")),
				NewMessage("C"),
			},
		},
		{
			NewMessage("A").
				AddMessage(NewMessage("B").
					AddMessage(NewMessage("B").
						AddMessage(NewMessage("C")))).
				AddMessage(NewMessage("D")),
			[]*Message{
				NewMessage("A"),
				NewMessage("B").AddMessage(NewMessage("B")),
				NewMessage("D"),
				NewMessage("C"),
			},
		},
		{
			NewMessage("A").
				AddMessage(NewMessage("B").
					AddMessage(NewMessage("B").
						AddMessage(NewMessage("C")))).
				AddMessage(NewMessage("C")),
			[]*Message{
				NewMessage("A"),
				NewMessage("B").
					AddMessage(NewMessage("B")).
					AddMessage(NewMessage("C")),
				NewMessage("C"),
			},
		},
		{
			NewMessage("A").
				AddMessage(NewMessage("B").
					AddMessage(NewMessage("B").
						AddMessage(NewMessage("C").
							AddMessage(NewMessage("D"))))).
				AddMessage(NewMessage("C")),
			[]*Message{
				NewMessage("A"),
				NewMessage("B").
					AddMessage(NewMessage("B")).
					AddMessage(NewMessage("C")),
				NewMessage("C"),
				NewMessage("D"),
			},
		},
	}

	for x, d := range table {
		if actual := flatten(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}

func TestEnumFlatten(t *testing.T) {
	table := []struct {
		in       Message
		expected Message
	}{
		{
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				},
			},
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				},
			},
		},
		{
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				},
				Messages: []*Message{
					&Message{
						Name: "Inner",
						Enums: []*Enum{
							&Enum{Name: "E3"},
						},
					},
				},
			},
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
					&Enum{Name: "E3"},
				},
				Messages: []*Message{
					&Message{Name: "Inner"},
				},
			},
		},
		{
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				},
				Messages: []*Message{
					&Message{
						Name: "Inner",
						Enums: []*Enum{
							&Enum{Name: "E1"},
						},
					},
				},
			},
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				},
				Messages: []*Message{
					&Message{
						Name: "Inner",
						Enums: []*Enum{
							&Enum{Name: "E1"},
						},
					},
				},
			},
		},
		{
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				},
				Messages: []*Message{
					&Message{
						Name: "Inner1",
						Enums: []*Enum{
							&Enum{Name: "E3"},
						},
					},
					&Message{
						Name: "Inner2",
						Enums: []*Enum{
							&Enum{Name: "E3"},
						},
					},
				},
			},
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
					&Enum{Name: "E3"},
				},
				Messages: []*Message{
					&Message{
						Name: "Inner1",
					},
					&Message{
						Name: "Inner2",
						Enums: []*Enum{
							&Enum{Name: "E3"},
						},
					},
				},
			},
		},
		{
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
					&Enum{Name: "E3"},
				},
				Messages: []*Message{
					&Message{
						Name: "Inner1",
						Enums: []*Enum{
							&Enum{Name: "E4"},
						},
						Messages: []*Message{
							&Message{
								Name: "InnerMost",
								Enums: []*Enum{
									&Enum{Name: "E5"},
								},
							},
						},
					},
					&Message{
						Name: "Inner2",
						Enums: []*Enum{
							&Enum{Name: "E5"},
						},
					},
				},
			},
			Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
					&Enum{Name: "E3"},
					&Enum{Name: "E4"},
					&Enum{Name: "E5"},
				},
				Messages: []*Message{
					&Message{
						Name: "Inner1",
						Messages: []*Message{
							&Message{
								Name: "InnerMost",
							},
						},
					},
					&Message{
						Name: "Inner2",
						Enums: []*Enum{
							&Enum{Name: "E5"},
						},
					},
				},
			},
		},
	}

	for x, d := range table {
		if actual := enumFlatten(&d.in); !reflect.DeepEqual(actual, &d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}
