package Database

type ID int64

/*
func (id *ID) str() string {
	return strconv.FormatInt(int64(*id), 10)
}

func ToID(v interface{}) ID {
	switch v := v.(type) {
	case string:
		id, err := snowflake.ParseString(v)
		if err != nil {
			return 0
		}
		return ID(id)
	case int64:
		return ID(snowflake.ParseInt64(v))
	default:
		return 0
	}
}
*/