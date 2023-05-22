package snowflake

import (
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/sealsee/web-base/public/utils/node"

	"time"

	sf "github.com/bwmarrin/snowflake"
)

var sfNode *sf.Node

func Init() {
	id := node.GetNodeId()
	var st time.Time
	st, err := time.Parse("2006-01-02", "2023-05-01")
	if err != nil {
		panic(err)
	}
	sf.Epoch = st.UnixNano() / 1000000
	sfNode, err = sf.NewNode(gconv.Int64(id))
	if err != nil {
		panic(err)
	}
}

func GenID() int64 {
	return sfNode.Generate().Int64()
}
