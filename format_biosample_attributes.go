package main

import (
	"log"
	"flag"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"bufio"
)

func deferE(dst *error, e error) {
	if *dst == nil {
		*dst = e
	}
}

func OldInCols() map[string]int {
	names := []string {
		"Experiment",
		"Multiplex Group",
		"ID",
		"Sample Name",
		"Organism",
		"Bird_Breed",
		"Individual",
		"Replicate",
		"Sample Name 2",
		"Time",
		"Treatment",
	}

	m := map[string]int{}
	for i, name := range names {
		m[name] = i
	}
	log.Printf("Names: %v\n", names)
	return m
}

func InCols() map[string]int {
	names := []string {
		"Unnamed: 0",
		"Experiment",
		"Multiplex Group",
		"ID",
		"Sample Name",
		"Organism",
		"Bird_Breed",
		"Individual_x",
		"Replicate",
		"Sample Name.1",
		"Time",
		"Treatment",
		"Individual_y",
	}

	m := map[string]int{}
	for i, name := range names {
		m[name] = i
	}
	log.Printf("Names: %v\n", names)
	return m
}

func Extend[T any](in []T, n int) []T {
	var t T
	for i := 0; i < n; i++ {
		in = append(in, t)
	}
	return in
}

func OutHeader() []string {
	return []string {
		"*sample_name",
		"sample_title",
		"bioproject_accession",
		"*organism",
		"isolate",
		"breed",
		"host",
		"isolation_source",
		"*collection_date",
		"*geo_loc_name",
		"*tissue",
		"biomaterial_provider",
		"dev_stage",
		"specimen_voucher",
		"host breed",
		"host treatment",
		"months since start of experiment",
		"experimental replicate",
		"pooled?",
	}
}

func OutCols() map[string]int {
	names := OutHeader()

	m := map[string]int{}
	for i, name := range names {
		m[name] = i
	}
	log.Printf("OutCols: %v\n", m)
	return m
}

func BuildMetadata(w io.Writer, r io.Reader) (err error) {
	h := func(e error) error {
		return fmt.Errorf("BuildMetadata: %w", e)
	}

	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.FieldsPerRecord = -1
	cr.ReuseRecord = true
	cr.LazyQuotes = true

	cw := csv.NewWriter(w)
	cw.Comma = '\t'
	defer func() {
		cw.Flush()
		deferE(&err, cw.Error())
	}()

	var outl []string

	icol := InCols()
	ocol := OutCols()

	if _, e := cr.Read(); e != nil {
		return h(e)
	}
	if e := cw.Write(OutHeader()); e != nil {
		return h(e)
	}

	for l, e := cr.Read(); e != io.EOF; l, e = cr.Read() {
		if e != nil {
			return h(e)
		}
		if len(l) < len(icol) {
			log.Printf("len(l) %v < len(icol) %v; l: %v\n", len(l), len(icol), l)
			continue
		}

		outl = Extend(outl[:0], len(ocol))

		id := l[icol["ID"]]
		if id == "NA" {
			id = l[icol["Individual"]]
		}

		outl[ocol["*sample_name"]] = l[icol["Sample Name"]]
		outl[ocol["sample_title"]] = id
		outl[ocol["bioproject_accession"]] = "PJRNA..."
		outl[ocol["*organism"]] = "Columbicola columbae"
		outl[ocol["isolate"]] = l[icol["Sample Name"]]
		outl[ocol["breed"]] = "not applicable"
		outl[ocol["host"]] = "Columba livia"
		outl[ocol["isolation_source"]] = "not applicable"
		outl[ocol["*collection_date"]] = "2018-01-01"
		outl[ocol["*geo_loc_name"]] = "United States: Salt Lake City"
		outl[ocol["*tissue"]] = "whole animal"
		outl[ocol["biomaterial_provider"]] = "Clayton-Bush laboratory"
		outl[ocol["dev_stage"]] = "adult"
		outl[ocol["specimen_voucher"]] = id
		outl[ocol["host breed"]] = l[icol["Bird_Breed"]]
		outl[ocol["host treatment"]] = l[icol["Treatment"]]
		outl[ocol["months since start of experiment"]] = l[icol["Time"]]
		outl[ocol["experimental replicate"]] = l[icol["Replicate"]]

		pooled := "pooled"
		if l[icol["Time"]] == "36" {
			pooled = "single individual"
			log.Printf("single individual\n")
		}
		outl[ocol["pooled?"]] = pooled

		if e := cw.Write(outl); e != nil {
			return h(e)
		}
	}

	return nil
}

func main() {
	identpathp := flag.String("i", "", "identities file")
	flag.Parse()

	var r io.Reader = os.Stdin
	if (*identpathp != "") {
		ir, e := os.Open(*identpathp)
		if e != nil {
			log.Fatal(e)
		}
		defer func() {
			if e := ir.Close(); e != nil {
				log.Fatal(e)
			}
		}()
		r = ir
	}

	br := bufio.NewReader(r)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	e := BuildMetadata(w, br)
	if e != nil {
		log.Fatal(e)
	}
}

// *sample_name	sample_title	bioproject_accession	*organism	isolate	breed	host	isolation_source	*collection_date	*geo_loc_name	*tissue	age	altitude	biomaterial_provider	collected_by	depth	dev_stage	env_broad_scale	host_tissue_sampled	identified_by	lat_lon	sex	specimen_voucher	temp	description
// 
// Experiment	Multiplex Group	ID	Sample Name	Organism	Bird_Breed	Individual	Replicate	Sample Name	Time	Treatment
// NA	1	15515X1	EMW1	Pigeon louse	feral	EMW1	1	EMW1	0	yes
// NA	1	15515X2	EMW2	Pigeon louse	feral	EMW2	1	EMW2	0	yes
// NA	1	15515X3	EMW3	Pigeon louse	feral	EMW3	1	EMW3	0	yes
// NA	1	15515X4	EMW5	Pigeon louse	feral	EMW5	1	EMW5	0	yes
// NA	1	15515X5	EMW6	Pigeon louse	feral	EMW6	1	EMW6	0	yes
// NA	1	15515X6	EMW7	Pigeon louse	feral	EMW7	1	EMW7	0	yes
// NA	1	15515X7	EMW9	Pigeon louse	feral	EMW9	1	EMW9	0	yes
// NA	1	15515X8	EMW10	Pigeon louse	feral	EMW10	1	EMW10	0	yes
// NA	1	15515X9	EMW11	Pigeon louse	feral	EMW11	1	EMW11	0	yes
