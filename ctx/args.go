package ctx

import (
	"fmt"
	"io"
	"os"

	"github.com/pborman/getopt/v2"
)

const (
	confPathDefault = "/nxs-data-anonymizer.conf"
)

// Args contains arguments value read from command line
type Args struct {
	ConfigPath string
	Input      io.Reader
	Output     io.Writer
	Cleanup    bool
	DBType     DBType
}

type DBType string

const (
	DBTypeMySQL DBType = "mysql"
	DBTypePgSQL DBType = "pgsql"
)

// ArgsRead reads arguments from command line
func ArgsRead() Args {

	var a Args

	args := getopt.New()

	helpFlag := args.BoolLong(
		"help",
		'h',
		"Show help")

	versionFlag := args.BoolLong(
		"version",
		'v',
		"Show program version")

	confPath := args.StringLong(
		"conf",
		'c',
		"",
		"Config file path")

	input := args.StringLong(
		"input",
		'i',
		"",
		"Input file")

	output := args.StringLong(
		"output",
		'o',
		"",
		"Output file")

	dbType := args.EnumLong(
		"type",
		't',
		[]string{
			string(DBTypeMySQL),
			string(DBTypePgSQL),
		},
		"",
		"Database type you need to operate",
	)

	cleanup := args.BoolLong(
		"cleanup",
		'C',
		"Clean up deastination database")

	args.Parse(os.Args)

	/* Show help */
	if *helpFlag == true {
		argsHelp(args)
		os.Exit(0)
	}

	/* Show version */
	if *versionFlag == true {
		argsVersion()
		os.Exit(0)
	}

	/* Config path */
	if args.IsSet("conf") == true {
		a.ConfigPath = *confPath
	} else {
		a.ConfigPath = confPathDefault
	}

	if args.IsSet("input") == true {

		f, err := os.Open(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open input file error: %s\n", err)
			os.Exit(1)
		}

		a.Input = f
	} else {
		a.Input = os.Stdin
	}

	if args.IsSet("output") == true {

		f, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open output file error: %s\n", err)
			os.Exit(1)
		}

		a.Output = f
	} else {
		a.Output = os.Stdout
	}

	a.Cleanup = *cleanup

	if args.IsSet("type") == true {
		a.DBType = DBType(*dbType)
	} else {
		fmt.Fprintf(os.Stderr, "'type' option must be specified\n")
		os.Exit(1)
	}

	return a
}

func argsHelp(args *getopt.Set) {

	additionalDescription := `
	
Additional description

  Just a sample.
`

	args.PrintUsage(os.Stdout)
	fmt.Println(additionalDescription)
}

func argsVersion() {
	fmt.Println("1.0")
}
