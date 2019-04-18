package fbsparser_test

import (
	"testing"

	"github.com/kazu/fbshelper/fbsparser"
	"github.com/stretchr/testify/assert"
)

func TestParserTableFixedField(t *testing.T) {

	data := `namespace mb_schema;
	table RegistGameServer {
		message_type:uint;
		object_id:uint;
		room_id:long;
		user_id:long;
		hoge_id:float64;
		uuid:[byte];
		hogas:[Hogo];
		Hogg:Hoga;
		
	  }
	struct Hoga {
		mope:int64;
	}  
	union Hoge { RegistGameServer, Hoga }  
	  
	`

	parser := &fbsparser.Parser{Buffer: data}

	parser.Init()
	err := parser.Parse()

	if err != nil {
		t.Error(err)
		//return
	}

	parser.Execute()
	parser.Fbs.FinilizeForFbs()

	assert.Equal(t, len(parser.Fbs.Structs), 2, parser.Fbs.Structs)
	assert.Equal(t, parser.Fbs.Structs[0].Name, "RegistGameServer")
	assert.Equal(t, len(parser.Fbs.Structs[0].Fields), 8, parser.Fbs.Structs)
	assert.Equal(t, len(parser.Fbs.Unions), 1, parser.Fbs.Unions)
	assert.Equal(t, len(parser.Fbs.Unions[0].Aliases), 2, parser.Fbs.Unions)

}
