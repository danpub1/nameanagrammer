package main

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
)

var primes = [...]uint64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41,
	43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101}

func hash(s string) Code {
	var low uint64 = 1
	var high uint64 = 0

	for _, ss := range s {
		h1, l1 := bits.Mul64(low, primes[ss-'A'])
		_, l2 := bits.Mul64(high, primes[ss-'A'])
		low = l1
		high = l2 + h1
	}

	return Code{low, high}
}

func remainder(source, first string) string {
	ll := len(source)

	if len(first) > ll {
		return ""
	}

	for _, ss := range first {
		orgLen := len(source)
		source = strings.Replace(source, string(ss), "", 1)
		if orgLen == len(source) {
			return ""
		}
	}

	return source
}

func title(source string) string {
	if len(source) > 1 {
		return strings.ToUpper(string(source[0])) + strings.ToLower(source[1:])
	} else if len(source) > 0 {
		return strings.ToUpper(source)
	} else {
		return source
	}
}

func findAnagrams(source string) []string {
	var last map[Code]string = make(map[Code]string)
	for idx, alast := range lastNames {
		prev, exists := last[alast.code]
		if exists {
			last[alast.code] = prev + "/" + idx
		} else {
			last[alast.code] = idx
		}
	}

	var first map[Code]string = make(map[Code]string)
	for idx, afirst := range firstNames {
		prev, exists := first[afirst.code]
		if exists {
			first[afirst.code] = prev + "/" + idx
		} else {
			first[afirst.code] = idx
		}
	}

	cleanSource := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r
		case r >= 'a' && r <= 'z':
			return r - 'a' + 'A'
		}
		return ' '
	}

	source = strings.Map(cleanSource, source)
	source = strings.ReplaceAll(source, " ", "")

	var found []Code = make([]Code, 0)

	for _, val := range first {
		afirst := strings.Split(val, "/")[0]
		alast := remainder(source, afirst)
		if len(alast) > 0 {
			acode := hash(alast)
			_, exists := last[acode]
			if exists {
				found = append(found, hash(afirst), acode)
			}
		}
	}

	var results []string = make([]string, 0)

	for idx := 0; idx < len(found); idx += 2 {
		fnarr := strings.Split(first[found[idx]], "/")
		lnarr := strings.Split(last[found[idx+1]], "/")
		for afn := range fnarr {
			for aln := range lnarr {
				results = append(results, fmt.Sprintf("%025.20f, %s %s", firstNames[fnarr[afn]].weight/100000*lastNames[lnarr[aln]].weight, title(fnarr[afn]), title(lnarr[aln])))
			}
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(results)))

	return results
}
