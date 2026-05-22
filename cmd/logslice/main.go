// Command logslice extracts time-range segments from large structured log files.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/reader"
	"github.com/yourorg/logslice/internal/slicer"
)

const timeLayout = time.RFC3339

func main() {
	var (
		from   = flag.String("from", "", "start of time range (RFC3339, inclusive)")
		to     = flag.String("to", "", "end of time range (RFC3339, exclusive)")
		input  = flag.String("input", "", "input log file (default: stdin); .gz files supported")
		output = flag.String("output", "", "output file (default: stdout)")
	)
	flag.Parse()

	if *from == "" || *to == "" {
		fmt.Fprintln(os.Stderr, "error: --from and --to are required")
		flag.Usage()
		os.Exit(1)
	}

	start, err := time.Parse(timeLayout, *from)
	if err != nil {
		log.Fatalf("invalid --from value: %v", err)
	}
	end, err := time.Parse(timeLayout, *to)
	if err != nil {
		log.Fatalf("invalid --to value: %v", err)
	}

	var r io.ReadCloser
	if *input == "" {
		r = io.NopCloser(os.Stdin)
	} else {
		r, err = reader.NewFileReader(*input)
		if err != nil {
			log.Fatalf("cannot open input: %v", err)
		}
	}
	defer r.Close()

	var w io.Writer
	if *output == "" {
		w = os.Stdout
	} else {
		f, ferr := os.Create(*output)
		if ferr != nil {
			log.Fatalf("cannot create output file: %v", ferr)
		}
		defer f.Close()
		w = f
	}

	lineReader := reader.NewReader(r)
	if err := slicer.Slice(lineReader, w, start, end); err != nil {
		log.Fatalf("slice error: %v", err)
	}
}
