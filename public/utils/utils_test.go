package utils

import (
	"fmt"
	"testing"

	"github.com/sealsee/web-base/public/utils/snowflake"
)

func TestGenSnowId(t *testing.T) {
	snowflake.Init()

	fmt.Println(snowflake.GenID())
	fmt.Println(snowflake.GenID())
	fmt.Println(snowflake.GenID())
	fmt.Println(snowflake.GenID())
}
