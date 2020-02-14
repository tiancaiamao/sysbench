package sysbench

import (
	"fmt"
	"time"

	"github.com/VividCortex/gohistogram"
)

type Report struct {
	Hist *gohistogram.NumericHistogram
	Succ int
	Fail int
	time.Duration
	Error bool
}

func (r *Report) Report() {
	fmt.Println("========= hisgotram ==========")
	fmt.Println(r.Hist)
	fmt.Println("========= summary ==========")
	fmt.Println("Mean: ", r.Hist.Mean())
	fmt.Println("Quantile 80: ", r.Hist.Quantile(0.8))
	fmt.Println("Quantile 95: ", r.Hist.Quantile(0.95))
	fmt.Println("Quantile 99: ", r.Hist.Quantile(0.99))
	fmt.Println("Total Succ:", r.Succ)
	fmt.Println("Total Fail:", r.Fail)
	if r.Error {
		fmt.Println("Failed!!!")
	}
}
