package handlers

import "github.com/sirupsen/logrus"

var log *logrus.Logger

func InitLogger(l *logrus.Logger) {
	log = l
}
