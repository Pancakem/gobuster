package cmd

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/OJ/gobuster/v3/cli"
	"github.com/OJ/gobuster/v3/gobustertld"
	"github.com/OJ/gobuster/v3/libgobuster"
	"github.com/spf13/cobra"
)

// nolint:gochecknoglobals
var cmdTLD *cobra.Command

func runTLD(cmd *cobra.Command, args []string) error {
	globalopts, pluginopts, err := parseTLDOptions()
	if err != nil {
		return fmt.Errorf("error on parsing arguments: %w", err)
	}

	plugin, err := gobustertld.NewGobusterTLD(globalopts, pluginopts)
	if err != nil {
		return fmt.Errorf("error on creating gobustertld: %w", err)
	}

	if err := cli.Gobuster(mainContext, globalopts, plugin); err != nil {
		var wErr *gobustertld.ErrWildcard
		if errors.As(err, &wErr) {
			return fmt.Errorf("%w. To force processing of Wildcard TLD, specify the '--wildcard' switch", wErr)
		}
		return fmt.Errorf("error on running gobuster: %w", err)
	}
	return nil
}

func parseTLDOptions() (*libgobuster.Options, *gobustertld.OptionsTLD, error) {
	globalopts, err := parseGlobalOptions()
	if err != nil {
		return nil, nil, err
	}
	plugin := gobustertld.NewOptionsTLD()

	plugin.Domain, err = cmdTLD.Flags().GetString("tld")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid value for top level domain: %w", err)
	}

	plugin.ShowIPs, err = cmdTLD.Flags().GetBool("show-ips")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid value for show-ips: %w", err)
	}

	plugin.ShowCNAME, err = cmdTLD.Flags().GetBool("show-cname")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid value for show-cname: %w", err)
	}

	plugin.WildcardForced, err = cmdTLD.Flags().GetBool("wildcard")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid value for wildcard: %w", err)
	}

	plugin.Timeout, err = cmdTLD.Flags().GetDuration("timeout")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid value for timeout: %w", err)
	}

	plugin.Resolver, err = cmdTLD.Flags().GetString("resolver")
	if err != nil {
		return nil, nil, fmt.Errorf("invalid value for resolver: %w", err)
	}

	if plugin.Resolver != "" && runtime.GOOS == "windows" {
		return nil, nil, fmt.Errorf("currently can not set custom dns resolver on windows. See https://golang.org/pkg/net/#hdr-Name_Resolution")
	}

	return globalopts, plugin, nil
}

// nolint:gochecknoinits
func init() {
	cmdTLD = &cobra.Command{
		Use:   "tld",
		Short: "Uses TLD domain enumeration mode",
		RunE:  runTLD,
	}

	cmdTLD.Flags().StringP("tld", "x", "", "The target top level domain")
	cmdTLD.Flags().BoolP("show-ips", "i", false, "Show IP addresses")
	cmdTLD.Flags().BoolP("show-cname", "c", false, "Show CNAME records (cannot be used with '-i' option)")
	cmdTLD.Flags().DurationP("timeout", "", time.Second, "TLD resolver timeout")
	cmdTLD.Flags().BoolP("wildcard", "", false, "Force continued operation when wildcard found")
	cmdTLD.Flags().StringP("resolver", "r", "", "Use custom TLD server (format server.com or server.com:port)")
	if err := cmdTLD.MarkFlagRequired("tld"); err != nil {
		log.Fatalf("error on marking flag as required: %v", err)
	}

	cmdTLD.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		configureGlobalOptions()
	}

	rootCmd.AddCommand(cmdTLD)
}
