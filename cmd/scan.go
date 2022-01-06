/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/goccy/go-graphviz"
	"github.com/google/go-github/v41/github"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

var (
	all       bool
	verbose   bool
	graph     bool
	graphPath string
)

func parseRepoToString(repo string) (string, string, error) {

	split := strings.Split(repo, "/")
	if len(split) != 2 {
		return "", "", errors.New("invalid repository name. Please use the format: owner/repo")
	}
	return split[0], split[1], nil
}

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a repository to find out who it's contributors work for.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			color.Red("GITHUB_TOKEN is not set. Please set it and try again.")
			return
		}

		owner, repoName, err := parseRepoToString(cmd.Flag("repo").Value.String())
		if err != nil {
			color.Red(err.Error())
			return
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		var contributorsList []*github.Contributor
		idx := 0
		for {

			contributors, resp, err := client.Repositories.ListContributors(ctx, owner, repoName,
				&github.ListContributorsOptions{
					Anon: "true",
					ListOptions: github.ListOptions{
						Page:    idx,
						PerPage: 250,
					},
				})
			if err != nil {
				color.Red(err.Error())
				return
			}

			contributorsList = append(contributorsList, contributors...)

			if resp.NextPage == 0 {
				break
			}
			idx++
		}

		fmt.Println(color.BlueString("%d contributors found", len(contributorsList)))

		var organisationalUsers map[string][]string = make(map[string][]string)

		color.Yellow("This will take a while...")
		for userIdx, contributor := range contributorsList {
			orgs, resp, err := client.Organizations.List(ctx, contributor.GetLogin(), &github.ListOptions{})
			if err != nil {
				color.Red(err.Error())
				continue
			}
			if resp.StatusCode != 200 {
				color.Red(fmt.Sprintf("%s returned %d", contributor.GetLogin(), resp.StatusCode))
				continue
			}

			for _, org := range orgs {
				if *org.Login == "" || contributor.GetLogin() == "" {
					continue
				}
				if _, ok := organisationalUsers[*org.Login]; !ok {
					organisationalUsers[*org.Login] = []string{}
				}
				organisationalUsers[*org.Login] = append(organisationalUsers[*org.Login], contributor.GetLogin())

			}
			fmt.Printf("\033[2K\r%d/%d Users scanned", userIdx+1, len(contributorsList))
		}

		fmt.Println(color.BlueString("\n%d unique organisations found", len(organisationalUsers)))

		n := map[int][]string{}
		var a []int
		for k, v := range organisationalUsers {
			n[len(v)] = append(n[len(v)], k)
		}

		for k := range n {
			a = append(a, k)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(a)))

		if graph {
			g := graphviz.New()
			graph, err := g.Graph()
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := graph.Close(); err != nil {
					log.Fatal(err)
				}
				g.Close()
			}()

			// Generate graph
			topLevelRepo, err := graph.CreateNode(repoName)
			if err != nil {
				log.Fatal(err)
			}

			for _, k := range a {
				for _, s := range n[k] {
					org, err := graph.CreateNode(s)
					if err != nil {
						log.Fatal(err)
					}
					tmpName := fmt.Sprintf("%s-%s", repoName, s)
					_, err = graph.CreateEdge(tmpName, topLevelRepo, org)
					if err != nil {
						log.Fatal(err)
					}

					// Generate the user to org edges
					if verbose {
						for _, u := range organisationalUsers[s] {

							user, err := graph.CreateNode(u)
							if err != nil {
								log.Fatal(err)
							}
							tmpName := fmt.Sprintf("%s-%s", s, u)
							_, err = graph.CreateEdge(tmpName, org, user)
							if err != nil {
								log.Fatal(err)
							}

						}
					}

				}
			}

			var buf bytes.Buffer
			if err := g.Render(graph, "dot", &buf); err != nil {
				log.Fatal(err)
			}

			if err := g.RenderFilename(graph, graphviz.PNG, graphPath); err != nil {
				log.Fatal(err)
			}

		} else {
			headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgYellow).SprintfFunc()
			var tbl table.Table
			if !verbose {
				tbl = table.New("Organisation", "Contributors")
			} else {
				tbl = table.New("Organisation", "Contributors", "Users")
			}
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			for _, k := range a {
				for idx, s := range n[k] {
					if !all && idx >= 4 {
						break
					}
					if verbose {
						tbl.AddRow(s, k, strings.Join(organisationalUsers[s], ", "))
					} else {
						tbl.AddRow(s, k)
					}
				}
			}
			tbl.Print()
		}
	},
}

func init() {
	scanCmd.Flags().BoolVarP(&all, "all", "a", false, "Print all results")
	scanCmd.Flags().StringP("repo", "r", "", "The GitHub repository to scan e.g. AlexsJones/prop-rep")
	scanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print verbose output")
	scanCmd.Flags().BoolVarP(&graph, "graph", "g", false, "Print graph output")
	scanCmd.Flags().StringVarP(&graphPath, "graph-path", "p", "graph.png", "Path to save graph to. Defaults to graph.png")
	scanCmd.MarkFlagRequired("repo")
	rootCmd.AddCommand(scanCmd)

}
