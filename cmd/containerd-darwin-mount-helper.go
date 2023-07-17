package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func main() {
	var fstype string
	var options []string
	var cmd = &cobra.Command{
		Args:         cobra.ExactArgs(2),
		ArgAliases:   []string{"source", "target"},
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var commandName string
			if fstype == "bind" {
				// macOS doesn't natively support bindfs/nullfs
				// The way to emulate it is via FUSE fs named "bindfs"
				commandName = "bindfs"
			} else {
				return fmt.Errorf("unsupported mount type: %s", fstype)
				// This could be used in the future to support other mount types
				// commandName = fmt.Sprintf("mount_%s", m.Type)
			}

			var bindargs []string
			for _, option := range options {
				if option == "rbind" {
					// On one side, rbind is not supported by macOS mounting tools
					// On the other, bindfs works as if rbind is enabled anyway
					continue
				}

				bindargs = append(bindargs, "-o", option)
			}
			bindargs = append(bindargs, args...)

			bindcmd := exec.Command(commandName, bindargs...)

			output, err := bindcmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("%w: %s", err, string(output))
			}

			return nil
		},
	}
	cmd.Flags().StringVarP(&fstype, "type", "t", "", "mount type")
	cmd.Flags().StringArrayVarP(&options, "option", "o", nil, "mount options")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
