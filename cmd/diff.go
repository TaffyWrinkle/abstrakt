package cmd

// visualise is a subcommand that constructs a graph representation of the yaml
// input file and renders this into GraphViz 'dot' notation.
// Initial version renders to dot syntax only, to graphically depict this the output
// has to be run through a graphviz visualisation tool/utiliyy

import (
	"bytes"
	"fmt"
	"github.com/microsoft/abstrakt/internal/platform/constellation"
	"github.com/microsoft/abstrakt/internal/tools/helpers"
	"github.com/microsoft/abstrakt/internal/tools/logger"
	"github.com/spf13/cobra"
	"strings"
)

type diffCmd struct {
	constellationFilePathOrg string
	constellationFilePathNew string
	showOriginal             *bool
	showNew                  *bool
	*baseCmd
}

func newDiffCmd() *diffCmd {
	cc := &diffCmd{}

	cc.baseCmd = newBaseCmd(&cobra.Command{
		Use:   "diff",
		Short: "Graphviz dot notation comparing two constellations",
		Long: `Diff is for producing a Graphviz dot notation representation of the difference between two constellations (line an old and new version)
	
Example: abstrakt diff -o [constellationFilePathOriginal] -n [constellationFilePathNew]`,

		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Debug("args: " + strings.Join(args, " "))
			logger.Debugf("constellationFilePathOrg: %v", cc.constellationFilePathOrg)
			logger.Debugf("constellationFilePathNew: %v", cc.constellationFilePathNew)

			logger.Debugf("showOriginalOutput: %t", *cc.showOriginal)
			logger.Debugf("showNewOutput: %t", *cc.showNew)

			if !helpers.FileExists(cc.constellationFilePathOrg) {
				return fmt.Errorf("Could not open original YAML input file for reading %v", cc.constellationFilePathOrg)
			}

			if !helpers.FileExists(cc.constellationFilePathNew) {
				return fmt.Errorf("Could not open new YAML input file for reading %v", cc.constellationFilePathNew)
			}

			dsGraphOrg := new(constellation.Config)
			err := dsGraphOrg.LoadFile(cc.constellationFilePathOrg)
			if err != nil {
				return fmt.Errorf("dagConfigService failed to load file %q: %s", cc.constellationFilePathOrg, err)
			}

			if *cc.showOriginal {
				out := &bytes.Buffer{}
				resStringOrg, err := dsGraphOrg.GenerateGraph(out)
				if err != nil {
					return err
				}
				logger.Output(resStringOrg)
			}

			dsGraphNew := new(constellation.Config)
			err = dsGraphNew.LoadFile(cc.constellationFilePathNew)
			if err != nil {
				return fmt.Errorf("dagConfigService failed to load file %q: %s", cc.constellationFilePathNew, err)
			}

			if *cc.showNew {
				out := &bytes.Buffer{}
				resStringNew, err := dsGraphNew.GenerateGraph(out)
				if err != nil {
					return err
				}
				logger.Output(resStringNew)
			}

			resStringDiff, err := constellation.CompareConstellations(dsGraphOrg, dsGraphNew)

			if err != nil {
				return err
			}

			logger.Output(resStringDiff)

			return nil
		},
	})

	cc.cmd.Flags().StringVarP(&cc.constellationFilePathOrg, "constellationFilePathOriginal", "o", "", "original or base constellation file path")
	cc.cmd.Flags().StringVarP(&cc.constellationFilePathNew, "constellationFilePathNew", "n", "", "new or changed constellation file path")
	cc.showOriginal = cc.cmd.Flags().Bool("showOriginalOutput", false, "will additionally produce dot notation for original constellation")
	cc.showNew = cc.cmd.Flags().Bool("showNewOutput", false, "will additionally produce dot notation for new constellation")
	_ = cc.cmd.MarkFlagRequired("constellationFilePathOriginal")
	_ = cc.cmd.MarkFlagRequired("constellationFilePathNew")

	return cc
}