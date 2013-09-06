package main

import "github.com/vgorin/ecgo/xorcoding"
import "flag"
import "os"
import "log"
import "runtime/pprof"

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

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

	// encode data block using b = 2 (renamed as bp)
	block_size := 1 << 22
	bp := byte(2)
	data_block := make([]byte, block_size)
	xorcoding.XorEncode(data_block, bp)
	log.Println("benchmark complete")
}
