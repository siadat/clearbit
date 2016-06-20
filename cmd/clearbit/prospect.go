package main

import (
	"github.com/codegangsta/cli"
	"github.com/thoughtbot/clearbit"
	"regexp"
)

var prospectCommand = cli.Command{
	Name:      "prospect",
	Usage:     "fetch contacts for a company",
	ArgsUsage: "domain",
	Action:    prospect,
	Flags: []cli.Flag{
		cli.StringFlag{Name: "name", Usage: "Filter by first or last name (case-insensitive)"},
		cli.StringFlag{Name: "role", Usage: "Filter by job role"},
		cli.StringFlag{Name: "seniority", Usage: "Filter by job seniority"},
		cli.StringSliceFlag{Name: "title", Usage: "Filter by job title"},
	},
}

func prospect(ctx *cli.Context) error {
	var (
		apiKey = apiKeyFromContext(ctx)
		client = clearbit.NewClient(apiKey, nil)

		domain    = ctx.Args().First()
		name      = ctx.String("name")
		role      = ctx.String("role")
		seniority = ctx.String("seniority")
		titlesCSV = ctx.StringSlice("title")
		titles    []string
		separator = regexp.MustCompile("\\s*,\\s*")
	)

	for _, title := range titlesCSV {
		titles = append(titles, separator.Split(title, -1)...)
	}

	if domain == "" {
		return requiredArgError("Usage: clearbit prospect <domain>")
	}

	prospects, err := client.Prospect(clearbit.ProspectQuery{
		Domain:    domain,
		Name:      name,
		Role:      role,
		Seniority: seniority,
		Titles:    titles,
	})
	if err != nil {
		return exitError(err)
	}

	display(prospects)
	return nil
}
