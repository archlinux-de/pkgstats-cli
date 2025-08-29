package cmd

import (
	"fmt"
	"io"

	"pkgstats-cli/internal/api/request"
)

type requestFunc func(client *request.Client, args []string) (request.PackagePopularityList, []error)

func handleRequest(writer io.Writer, args []string, fn requestFunc) error {
	client := request.NewClient(baseURL)

	ppl, errs := fn(client, args)

	if len(ppl.PackagePopularities) > 0 {
		request.PrintPackagePopularities(writer, ppl)
	}

	if len(errs) > 0 {
		fmt.Fprintln(writer, "\nErrors:")
		for _, err := range errs {
			fmt.Fprintln(writer, err)
		}
	}

	return nil
}
