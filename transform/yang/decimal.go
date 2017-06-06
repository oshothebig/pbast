package yang

import "github.com/oshothebig/pbast"

var decimal64Message = pbast.NewMessage("Decimal64").
	AddField(pbast.NewMessageField(pbast.Int64, "value", 1)).
	AddField(pbast.NewMessageField(pbast.UInt32, "fraction_digits", 2))
