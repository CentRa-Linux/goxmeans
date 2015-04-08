/*
An example of how to invoke goxmeans.  Data is stored in a file called "datastet" in
the working directory.  The data is only two dimensions for this example.

 Usage: go run ./xmeans_ex.go k kmax
 Where k and kmax are integers, and k <= kmax, that indicate the number of centroids to use.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/drewlanenga/goxmeans"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var c = flag.Int("centroids", 1234, "number of centroids")

func main() {
	usage := "usage: xmeans_ex k kmax"
	if len(os.Args) < 3 {
		fmt.Println(usage)
		return
	}

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Could not convert arg %s to int.\n", os.Args[1])
		return
	}

	kmax, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Could not convert arg %s to int.\n", os.Args[2])
		return
	}

	if kmax < k {
		fmt.Printf("k must be <= kmax\n")
		return
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Load data set.
	data, err := goxmeans.Load("dataset", ",")
	if err != nil {
		fmt.Println("Load: ", err)
		return
	}
	fmt.Println("Load complete")

	model, err := goxmeans.BestXmeans(data, k, kmax)
	if err != nil {
		fmt.Println("BestXmeans:", err)
		return
	}

	fmt.Printf("\nBest fit:[#centroids=%d BIC=%f]\n", model.Numcentroids(), model.Bic)
	for i, c := range model.Clusters {
		fmt.Printf("cluster-%d: numpoints=%d variance=%f\n", i, c.Numpoints(), c.Variance)
	}

	assignment := goxmeans.ExtractClusters(model)
	fmt.Printf("Cluster Assignment: %v...\n", assignment[0:5])
}
