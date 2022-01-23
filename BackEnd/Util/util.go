package Util

import (
	"context"
	"github.com/bwmarrin/snowflake"
)

func String(s string) *string {
	return &s
}

func GetSnowflakeFromCTX(ctx context.Context) *snowflake.Node {
	return ctx.Value("snowflake").(*snowflake.Node)
}
