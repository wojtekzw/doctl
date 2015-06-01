package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/digitalocean/godo"

	"golang.org/x/oauth2"
)

var ImageCommand = cli.Command{
	Name:    "image",
	Aliases: []string{"i"},
	Usage:   "Image commands. Lists by default.",
	Action:  imageList,
	Subcommands: []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List images.",
			Action:  imageList,
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "", Usage: "List only "},
			},
		},
		{
			Name:    "destroy",
			Aliases: []string{"d"},
			Usage:   "[--id | <name>] Destroy an droplet.",
			Action:  imageDestroy,
			Flags: []cli.Flag{
				cli.IntFlag{Name: "id", Usage: "ID for image. (e.g. 1234567)"},
			},
		},
	},
}

func imageList(ctx *cli.Context) {
	if ctx.BoolT("help") == true {
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	tokenSource := &TokenSource{
		AccessToken: APIKey,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	opt := &godo.ListOptions{}
	imageList := []godo.Image{}

	for { // TODO make all optional
		imagePage, resp, err := client.Images.List(opt)
		if err != nil {
			fmt.Printf("Unable to list images: %s\n", err)
			os.Exit(1)
		}

		// append the current page's images to our list
		for _, d := range imagePage {
			imageList = append(imageList, d)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			fmt.Printf("Unable to get pagination: %s\n", err)
			os.Exit(1)
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	cliOut := NewCLIOutput()
	defer cliOut.Flush()
	cliOut.Header("ID", "Name", "Slug", "Type", "Distribution", "Disk Required", "Regions")
	for _, image := range imageList {
		cliOut.Writeln("%d\t%s\t%s\t%s\t%s\t%dGB\t%v\n",
			image.ID, image.Name, image.Slug, image.Type, image.Distribution, image.MinDiskSize, image.Regions)
	}
}

func imageDestroy(ctx *cli.Context) {
	if ctx.Int("id") == 0 && len(ctx.Args()) != 1 {
		fmt.Printf("Error: Must provide ID or name for Image to destroy.\n")
		os.Exit(1)
	}

	tokenSource := &TokenSource{
		AccessToken: APIKey,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	id := ctx.Int("id")
	if id == 0 {
		image, err := FindImageByName(client, ctx.Args()[0])
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(64)
		} else {
			id = image.ID
		}
	}

	image, _, err := client.Images.GetByID(id)
	if err != nil {
		fmt.Printf("Unable to find image: %s\n", err)
		os.Exit(1)
	}

	_, err = client.Images.Delete(id)
	if err != nil {
		fmt.Printf("Unable to destroy image: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Image %s destroyed.\n", image.Name)
}
