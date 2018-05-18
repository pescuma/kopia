package cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"

	"github.com/kopia/kopia/internal/packindex"

	"github.com/kopia/kopia/block"
	"github.com/kopia/kopia/repo"
)

var (
	blockIndexShowCommand = blockIndexCommands.Command("show", "List block indexes").Alias("cat")
	blockIndexShowIDs     = blockIndexShowCommand.Arg("id", "IDs of index blocks to show").Required().Strings()
	blockIndexShowRaw     = blockIndexShowCommand.Flag("raw", "Show raw block data").Bool()
)

func getIndexBlocksToShow(ctx context.Context, rep *repo.Repository) ([]block.PhysicalBlockID, error) {
	var blockIDs []block.PhysicalBlockID
	for _, id := range *blockIndexShowIDs {
		blockIDs = append(blockIDs, block.PhysicalBlockID(id))
	}

	if len(blockIDs) == 1 && blockIDs[0] == "active" {
		b, err := rep.Blocks.IndexBlocks(ctx)
		if err != nil {
			return nil, err
		}

		sort.Slice(b, func(i, j int) bool {
			return b[i].Timestamp.Before(b[j].Timestamp)
		})

		blockIDs = nil
		for _, bi := range b {
			blockIDs = append(blockIDs, bi.BlockID)
		}
	}

	return blockIDs, nil
}

func runShowBlockIndexesAction(ctx context.Context, rep *repo.Repository) error {
	blockIDs, err := getIndexBlocksToShow(ctx, rep)
	if err != nil {
		return err
	}

	for _, blockID := range blockIDs {
		data, err := rep.Blocks.GetIndexBlock(ctx, blockID)
		if err != nil {
			return fmt.Errorf("can't read block %q: %v", blockID, err)
		}

		if *blockIndexShowRaw {
			os.Stdout.Write(data) //nolint:errcheck
		} else {
			ndx, err := packindex.Open(bytes.NewReader(data))
			if err != nil {
				return err
			}

			_ = ndx.Iterate("", func(l block.Info) error {
				if l.Deleted {
					fmt.Printf("  del %v %v\n", l.BlockID, l.Timestamp().Format(timeFormat))
				} else {
					fmt.Printf("  add %v %v %v %v+%v\n", l.BlockID, l.Timestamp().Format(timeFormat), l.PackBlockID, l.PackOffset, l.Length)
				}
				return nil
			})
		}
	}

	return nil
}

func init() {
	blockIndexShowCommand.Action(repositoryAction(runShowBlockIndexesAction))
}
