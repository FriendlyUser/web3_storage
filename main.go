// author: David Li
// description: Utility script to put a file in web3.storage
// usage: go run main.go list -cid bafybeide43vps6vt2oo7nbqfwn5zz6l2alyi64mym3sb7reqhmypjnmej4

package main

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
	w3fs "github.com/web3-storage/go-w3s-client/fs"
)

// Create a new type for a list of Strings
type stringList []string

// Implement the flag.Value interface
func (s *stringList) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringList) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

// Usage:
// TOKEN="API_TOKEN" go run ./main.go
// see https://docs.web3.storage/how-tos/list#listing-your-uploads
// for details on how to retrieve all items
func main() {

	// parse command line arguments
	// Subcommands
	uploadCommand := flag.NewFlagSet("upload", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	// Use flag.Var to create a flag of our new flagType
	// Default value is the current value at countStringListPtr (currently a nil value)

	// List subcommand flag pointers
	listCidPtr := listCommand.String("cid", "", "Cid to list files for. (Required)")

	// upload subcommand flag pointers
	uploadPathPtr := uploadCommand.String("folder", "", "Folder to upload to web3 storage. (Required)")
	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("upload or list subcommand is required")
		os.Exit(1)
	}

	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "list":
		listCommand.Parse(os.Args[2:])
	case "upload":
		uploadCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
	c, _ := w3s.NewClient(w3s.WithEndpoint(os.Getenv("ENDPOINT")), w3s.WithToken(os.Getenv("TOKEN")))
	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	if listCommand.Parsed() {
		// List all items based on cid
		if *listCidPtr == "" {
			listCommand.PrintDefaults()
			os.Exit(1)
		}
		getFiles(c, *listCidPtr)
	}

	if uploadCommand.Parsed() {
		// // Required Flags
		if *uploadPathPtr == "" {
			uploadCommand.PrintDefaults()
			os.Exit(1)
		}
		p := filepath.FromSlash(*uploadPathPtr)
		putDirectory(c, p)
		// // If the metric flag is substring, the substring or substringList flag is required
		// if *countMetricPtr == "substring" && *countSubstringPtr == "" && (&countStringList).String() == "[]" {
		// 	countCommand.PrintDefaults()
		// 	os.Exit(1)
		// }
		// //If the metric flag is not substring, the substring flag must not be used
		// if *countMetricPtr != "substring" && (*countSubstringPtr != "" || (&countStringList).String() != "[]") {
		// 	fmt.Println("--substring and --substringList may only be used with --metric=substring.")
		// 	countCommand.PrintDefaults()
		// 	os.Exit(1)
		// }
		// //Choice flag
		// metricChoices := map[string]bool{"chars": true, "words": true, "lines": true, "substring": true}
		// if _, validChoice := metricChoices[*listMetricPtr]; !validChoice {
		// 	countCommand.PrintDefaults()
		// 	os.Exit(1)
		// }
		// //Print
		// fmt.Printf("textPtr: %s, metricPtr: %s, substringPtr: %v, substringListPtr: %v, uniquePtr: %t\n",
		// 	*countTextPtr,
		// 	*countMetricPtr,
		// 	*countSubstringPtr,
		// 	(&countStringList).String(),
		// 	*countUniquePtr,
		// )
	}
}

func putSingleFile(c w3s.Client) cid.Cid {
	file, err := os.Open("images/donotresist.jpg")
	if err != nil {
		panic(err)
	}
	return putFile(c, file)
}

func putMultipleFiles(c w3s.Client) cid.Cid {
	f0, err := os.Open("images/donotresist.jpg")
	if err != nil {
		panic(err)
	}
	f1, err := os.Open("images/pinpie.jpg")
	if err != nil {
		panic(err)
	}
	dir := w3fs.NewDir("comic", []fs.File{f0, f1})
	return putFile(c, dir)
}

func putMultipleFilesAndDirectories(c w3s.Client) cid.Cid {
	f0, err := os.Open("images/donotresist.jpg")
	if err != nil {
		panic(err)
	}
	f1, err := os.Open("images/pinpie.jpg")
	if err != nil {
		panic(err)
	}
	d0 := w3fs.NewDir("one", []fs.File{f0})
	d1 := w3fs.NewDir("two", []fs.File{f1})
	rootdir := w3fs.NewDir("comic", []fs.File{d0, d1})
	return putFile(c, rootdir)
}

func putDirectory(c w3s.Client, v string) cid.Cid {
	if v == "" {
		v = "images"
	}
	dir, err := os.Open(v)
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)
	return putFile(c, dir)
}

func putFile(c w3s.Client, f fs.File, opts ...w3s.PutOption) cid.Cid {
	cid, err := c.Put(context.Background(), f, opts...)
	if err != nil {
		panic(err)
	}
	fmt.Printf("https://%v.ipfs.dweb.link\n", cid)
	return cid
}

func getStatusForCid(c w3s.Client, cid cid.Cid) {
	s, err := c.Status(context.Background(), cid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Status: %+v", s)
}

func getStatusForKnownCid(c w3s.Client) {
	cid, _ := cid.Parse("bafybeig7qnlzyregxe2m63b4kkpx3ujqm5bwmn5wtvtftp7j27tmdtznji")
	getStatusForCid(c, cid)
}

func getFiles(c w3s.Client, v string) {
	if v == "" {
		v = "bafybeide43vps6vt2oo7nbqfwn5zz6l2alyi64mym3sb7reqhmypjnmej4"
	}
	cid, _ := cid.Parse(v)

	res, err := c.Get(context.Background(), cid)
	if err != nil {
		panic(err)
	}

	f, fsys, err := res.Files()
	if err != nil {
		panic(err)
	}

	info, err := f.Stat()
	if err != nil {
		panic(err)
	}

	if info.IsDir() {
		err = fs.WalkDir(fsys, "/", func(path string, d fs.DirEntry, err error) error {
			info, _ := d.Info()
			fmt.Printf("%s (%d bytes)\n", path, info.Size())
			return err
		})
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("%s (%d bytes)\n", cid.String(), info.Size())
	}
}
