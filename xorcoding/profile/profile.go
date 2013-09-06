package main

import "github.com/vgorin/ecgo/xorcoding"
import "flag"
import "os"
import "log"
import "runtime/pprof"
import "net/http"
import _ "net/http/pprof"
import "time"

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var web = flag.String("web", "", "enable web server on port")
var inf = flag.String("inf", "", "use infinite loop (useful for web server)")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		log.Printf("writing cpuprofile to %s", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *web != "" {
		go func() {
			log.Println(http.ListenAndServe("localhost:" + *web, nil))
		}()
	}

	// encode data block using b = 2 (renamed as bp)
	block_size := 1 << 28
	bp := byte(2)
	data_block := make([]byte, block_size)
	xorcoding.XorEncode(data_block, bp)
	for *inf != "" {
		xorcoding.XorEncode(data_block, bp)
		time.Sleep(1 << 20)
	}
	log.Println("profiling complete")
}
