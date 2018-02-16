package main

import (
	"log"
	"strings"

	"github.com/LloydGriffiths/ecr-mop/mop"
	"github.com/namsral/flag"
)

var (
	ignoreTags   = flag.String("ignore-tags", "", "comma seperated list of tags to ignore")
	repository   = flag.String("repository", "", "the repository to mop")
	staleAfter   = flag.Int("stale-after", 365, "mop images older than number days")
	wipeUntagged = flag.Bool("untagged", true, "mop untagged images")
)

func main() {
	flag.Parse()

	m, err := mop.New(*repository, *staleAfter, *wipeUntagged, strings.Split(*ignoreTags, ","))
	if err != nil {
		log.Fatal(err)
	}
	r, err := m.Wipe()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("removed %d images with %d failures from %s repository", len(r.Removed), len(r.Failed), *repository)
}
