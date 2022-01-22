package graphtypes

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/bwmarrin/snowflake"
	"io"
	"strconv"
)

type ID snowflake.ID

func MarshalMyID(id ID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		str := strconv.FormatInt(int64(id), 10)
		_, _ = w.Write([]byte(fmt.Sprintf(`"%v"`, str)))
	})
}

func UnmarshalMyID(v interface{}) (ID, error) {
	switch v := v.(type) {
	case string:
		id, err := snowflake.ParseString(v)
		if err != nil {
			return 0, fmt.Errorf("given string is not parable to snowflake-id")
		}
		return ID(id), nil
	case int64:
		return ID(snowflake.ParseInt64(v)), nil
	default:
		return 0, fmt.Errorf("snowflake only supports string and int64")
	}
}
