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
		uuid:[byte];
	  }
	union Hoge { RegistGameServer }  
	  
	`

	parser := &fbsparser.Parser{Buffer: data}

	parser.Init()
	err := parser.Parse()

	if err != nil {
		t.Error(err)
		//return
	}

	parser.Execute()

	assert.Equal(t, len(parser.Fbs.Structs), 1)
	assert.Equal(t, parser.Fbs.Structs[0].Name, "RegistGameServer")
	assert.Equal(t, len(parser.Fbs.Structs[0].Fields), 5)
	assert.Equal(t, len(parser.Fbs.Structs[0].Fields), 4, parser.Fbs.Structs[0])
	assert.Equal(t, len(parser.Fbs.Structs[0].Fields), 4, parser.Fbs.Unions[0])

}
