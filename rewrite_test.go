package pbast

import (
	"reflect"
	"testing"
)

func TestLiftMessage(t *testing.T) {
	table := []struct {
		in       *File
		expected *File
	}{
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{Name: "A"},
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{Name: "A"},
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "A",
						Messages: []*Message{
							&Message{
								Name: "B",
							},
						},
					},
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{Name: "A"},
					&Message{Name: "B"},
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "A",
						Messages: []*Message{
							&Message{
								Name: "B",
								Messages: []*Message{
									&Message{Name: "C"},
								},
							},
							&Message{Name: "D"},
						},
					},
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					&Message{Name: "A"},
					&Message{Name: "B"},
					&Message{Name: "D"},
					&Message{Name: "C"},
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("C"))).
						AddMessage(NewMessage("D")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B"),
					NewMessage("D"),
					NewMessage("C"),
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("B"))).
						AddMessage(NewMessage("C")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B").
						AddMessage(NewMessage("B")),
					NewMessage("C"),
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("B").
								AddMessage(NewMessage("C")))).
						AddMessage(NewMessage("D")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B").AddMessage(NewMessage("B")),
					NewMessage("D"),
					NewMessage("C"),
				},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("B").
								AddMessage(NewMessage("C")))).
						AddMessage(NewMessage("C")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B").
						AddMessage(NewMessage("B")).
						AddMessage(NewMessage("C")),
					NewMessage("C")},
			},
		},
		{
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A").
						AddMessage(NewMessage("B").
							AddMessage(NewMessage("B").
								AddMessage(NewMessage("C").
									AddMessage(NewMessage("D"))))).
						AddMessage(NewMessage("C")),
				},
			},
			&File{
				Package: "com.example",
				Messages: []*Message{
					NewMessage("A"),
					NewMessage("B").
						AddMessage(NewMessage("B")).
						AddMessage(NewMessage("C")),
					NewMessage("C"),
					NewMessage("D"),
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

func TestLiftEnum(t *testing.T) {
	table := []struct {
		in       File
		expected File
	}{
		{
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "Root",
						Enums: []*Enum{
							&Enum{Name: "E1"},
							&Enum{Name: "E2"},
						},
					},
				},
			},
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "Root",
					},
				},
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				}},
		},
		{
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
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
				},
			},
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "Root",
						Messages: []*Message{
							&Message{Name: "Inner"},
						},
					},
				},
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
					&Enum{Name: "E3"},
				},
			},
		},
		{
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
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
			},
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "Root",
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
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				},
			},
		},
		{
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
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
				},
			},
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "Root",
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
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
					&Enum{Name: "E3"},
				},
			},
		},
		{
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
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
				},
			},
			File{
				Package: "com.example",
				Messages: []*Message{
					&Message{
						Name: "Root",
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
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
					&Enum{Name: "E3"},
					&Enum{Name: "E4"},
					&Enum{Name: "E5"},
				},
			},
		},
	}

	for x, d := range table {
		if actual := LiftEnum(&d.in); !reflect.DeepEqual(actual, &d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}

func TestBfsEnum(t *testing.T) {
	table := []struct {
		in       *Message
		expected []*Enum
	}{
		{
			&Message{
				Name: "Root",
				Enums: []*Enum{
					&Enum{Name: "E1"},
					&Enum{Name: "E2"},
				},
			},
			[]*Enum{
				&Enum{Name: "E1"},
				&Enum{Name: "E2"},
			},
		},
		{
			&Message{
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
			[]*Enum{
				&Enum{Name: "E1"},
				&Enum{Name: "E2"},
				&Enum{Name: "E3"},
			},
		},
		{
			&Message{
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
			[]*Enum{
				&Enum{Name: "E1"},
				&Enum{Name: "E2"},
				&Enum{Name: "E1"},
			},
		},
		{
			&Message{
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
			[]*Enum{
				&Enum{Name: "E1"},
				&Enum{Name: "E2"},
				&Enum{Name: "E3"},
				&Enum{Name: "E3"},
			},
		},
		{
			&Message{
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
			[]*Enum{
				&Enum{Name: "E1"},
				&Enum{Name: "E2"},
				&Enum{Name: "E3"},
				&Enum{Name: "E4"},
				&Enum{Name: "E5"},
				&Enum{Name: "E5"},
			},
		},
	}

	for x, d := range table {
		if actual := bfsEnum(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}

}
