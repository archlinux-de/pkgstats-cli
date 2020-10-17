package request

import (
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"text/tabwriter"
)

func PrintPackagePopularities(writer io.Writer, ppl PackagePopularityList) {
	w := tabwriter.NewWriter(writer, 15, 0, 2, ' ', 0)
	for _, pkg := range ppl.PackagePopularities {
		fmt.Fprintf(w, "%s\t%.2f\n", pkg.Name, pkg.Popularity)
	}
	w.Flush()

	if ppl.Count > 0 && ppl.Total > 0 {
		fmt.Fprintf(writer, "\n%d of %d results\n", ppl.Count, ppl.Total)
	}
}

func PrintSearchURL(writer io.Writer, baseURL string, query string) {
	if len(query) > 0 {
		fmt.Fprintf(writer, "See more results at %s/packages#query=%s\n", baseURL, query)
	}
}

func PrintShowURL(writer io.Writer, baseURL string, packages []string) {
	if len(packages) > 0 {
		for i, p := range packages {
			packages[i] = url.QueryEscape(p)
		}

		sort.StringSlice.Sort(packages)
		fmt.Fprintf(writer, "See more results at %s/compare/packages#packages=%s\n", baseURL, strings.Join(packages, ","))
	}
}
