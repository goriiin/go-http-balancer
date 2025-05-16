package logger

import (
	"log/slog"
	"sync/atomic"
)

type atomicLevel struct{ v atomic.Int64 }

func (a *atomicLevel) Level() slog.Level { return slog.Level(a.v.Load()) }
func (a *atomicLevel) Set(l slog.Level)  { a.v.Store(int64(l)) }
