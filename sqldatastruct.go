package pkg

import "time"

type RedisCache struct {
	User		  string
	Gmt_create    time.Time
	Gmt_modify	  time.Time
}

