package logger

import (
	"context"
	"fmt"
	"log"
)

func Error(ctx context.Context, err error, msg string, args ...interface{}) {
	errorFmt := msg
	if len(args) > 0 {
		errorFmt = fmt.Sprintf(msg, args)
	}
	if err != nil {
		log.Print(fmt.Sprintf("ERROR [context=%v] [error=%s] %s", ctx, err.Error(), errorFmt))
	} else {
		log.Print(fmt.Sprintf("ERROR [context=%v] %s", ctx, errorFmt))
	}
}

func Info(ctx context.Context, msg string, args ...interface{}) {
	infoFmt := msg
	if len(args) > 0 {
		infoFmt = fmt.Sprintf(msg, args)
	}
	log.Print(fmt.Sprintf("INFO [context=%v] %s", ctx, infoFmt))
}
