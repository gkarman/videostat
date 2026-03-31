package mq

import "time"

type Config struct {
	User           string
	Password       string
	Host           string
	Port           string
	Exchange       string
	ReconnectDelay time.Duration
}
