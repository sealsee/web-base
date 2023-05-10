package node

import (
	"math/rand"
	"time"
)

func GetNodeId() int {
	rand.Seed(time.Now().UnixNano())
	var id int = 1
	// for true {
	// 	id = rand.Intn(1023)
	// 	val := ds.GetRedisClient().Get("snowflake:" + gconv.String(id)).Val()
	// 	if val == "" {
	// 		break
	// 	}
	// }
	// s, err := ipUtils.GetLocalIP()
	// if err != nil {
	// 	panic(err)
	// }

	// go func(id string, data string) {
	// 	for true {
	// 		ds.GetRedisClient().Set("snowflake:"+id, s, time.Hour+time.Minute).Err()
	// 		t := time.NewTimer(time.Hour)
	// 		<-t.C
	// 	}
	// }(gconv.String(id), s)
	return id
}
