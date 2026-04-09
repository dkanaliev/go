package runtime_test

import (
	"runtime"
	"runtime/debug"
	"testing"
)

// Attempt to allocate pointer rich objects from many goroutines
// concurrently in an effort to  force GC assists, wiith GOMAXPRO CS 
//set ABOVE the physical core count
//
//	go test -run=^$ -bench=BenchmarkGCAssistFlush -benchtime=5s -count=10
func BenchmarkGCAssistFlush(b *testing.B) {
	for _, procs := range []int{8, 32, 64, 128} {
		b.Run("procs="+itoa(procs), func(b *testing.B) {
			oldProcs := runtime.GOMAXPROCS(procs)
			defer runtime.GOMAXPROCS(oldProcs)

			oldGC := debug.SetGCPercent(10)
			defer debug.SetGCPercent(oldGC)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					sl := make([]*byte, 128)
					for i := range sl {
						v := byte(i)
						sl[i] = &v
					}
					runtime.KeepAlive(sl)
				}
			})
		})
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(buf[pos:])
}

