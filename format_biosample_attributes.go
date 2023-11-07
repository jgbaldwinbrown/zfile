package csvh

import (
	"log"
	"encoding/csv"
	"fmt"
	"io"
)

func Handle0(format string) func(e error) error {
	return func(e error) error {
		return fmt.Errorf(format, e)
	}
}

func Handle1[T any](format string) func(e error) (T, error) {
	return func(e error) (t T, err error) {
		err = fmt.Errorf(format, e)
		return
	}
}

func Handle2[T, U any](format string) func(e error) (T, U, error) {
	return func(e error) (t T, u U, err error) {
		err = fmt.Errorf(format, e)
		return
	}
}

func Handle3[T, U, V any](format string) func(e error) (T, U, V, error) {
	return func(e error) (t T, u U, v V, err error) {
		err = fmt.Errorf(format, e)
		return
	}
}

func HandleWrap(format string) func(e error) error {
	return func(e error) error {
		if e != nil {
			return fmt.Errorf(format, e)
		}
		return nil
	}
}

func DeferE(dst *error, e error) {
	if *dst == nil {
		*dst = e
	}
}

func NamesToCols(names []string) map[string]int {
	m := make(map[string]int, len(names))

	for i, name := range names {
		m[name] = i
	}

	return m
}

func Extend[T any](in []T, n int) []T {
	var t T
	for i := 0; i < n; i++ {
		in = append(in, t)
	}
	return in
}

func CsvIn(r io.Reader) *csv.Reader {
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.FieldsPerRecord = -1
	cr.ReuseRecord = true
	cr.LazyQuotes = true
	return cr
}

func CsvOut(w io.Writer) *csv.Writer {
	cw := csv.NewWriter(w)
	cw.Comma = '\t'
	return cw
}

func buildMetadataExample(w io.Writer, r io.Reader) (err error) {
	h := Handle0("BuildMetadata: %w")

	cr := CsvIn(r)

	cw := CsvOut(w)
	cw.Comma = '\t'
	defer func() {
		cw.Flush()
		DeferE(&err, cw.Error())
	}()

	var outl []string

	l, e := cr.Read()
	if e != nil {
		return h(e)
	}
	icol := NamesToCols(l)

	ohead := []string{"Apple", "Banana", "Carrot"}
	ocol := NamesToCols(ohead)
	if e := cw.Write(ohead); e != nil {
		return h(e)
	}

	for {
		l, e := cr.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			return h(e)
		} else if len(l) < len(icol) {
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
