package cli

import (
	"context"

	"github.com/alecthomas/kingpin"
	"github.com/pkg/errors"

	"github.com/kopia/kopia/snapshot/policy"
)

type policyFilesFlags struct {
	// Ignore rules.
	policySetAddIgnore    []string
	policySetRemoveIgnore []string
	policySetClearIgnore  bool

	// Dot-ignore files to look at.
	policySetAddDotIgnore    []string
	policySetRemoveDotIgnore []string
	policySetClearDotIgnore  bool
	policySetMaxFileSize     string

	// Ignore other mounted filesystems.
	policyOneFileSystem string

	policyIgnoreCacheDirs string

	// Ignore extended attributes rules.
	policySetAddIgnoreExtendedAttributes    []string
	policySetRemoveIgnoreExtendedAttributes []string
	policySetClearIgnoreExtendedAttributes  bool
}

func (c *policyFilesFlags) setup(cmd *kingpin.CmdClause) {
	// Ignore rules.
	cmd.Flag("add-ignore", "List of paths to add to the ignore list").PlaceHolder("PATTERN").StringsVar(&c.policySetAddIgnore)
	cmd.Flag("remove-ignore", "List of paths to remove from the ignore list").PlaceHolder("PATTERN").StringsVar(&c.policySetRemoveIgnore)
	cmd.Flag("clear-ignore", "Clear list of paths in the ignore list").BoolVar(&c.policySetClearIgnore)

	// Dot-ignore files to look at.
	cmd.Flag("add-dot-ignore", "List of paths to add to the dot-ignore list").PlaceHolder("FILENAME").StringsVar(&c.policySetAddDotIgnore)
	cmd.Flag("remove-dot-ignore", "List of paths to remove from the dot-ignore list").PlaceHolder("FILENAME").StringsVar(&c.policySetRemoveDotIgnore)
	cmd.Flag("clear-dot-ignore", "Clear list of paths in the dot-ignore list").BoolVar(&c.policySetClearDotIgnore)
	cmd.Flag("max-file-size", "Exclude files above given size").PlaceHolder("N").StringVar(&c.policySetMaxFileSize)

	// Ignore other mounted filesystems.
	cmd.Flag("one-file-system", "Stay in parent filesystem when finding files ('true', 'false', 'inherit')").EnumVar(&c.policyOneFileSystem, booleanEnumValues...)

	cmd.Flag("ignore-cache-dirs", "Ignore cache directories ('true', 'false', 'inherit')").EnumVar(&c.policyIgnoreCacheDirs, booleanEnumValues...)

	// Ignore extended attribute rules.
	cmd.Flag("add-ignore-extended-attribute", "List of extended attribute names to add to the ignore extended attributes list").PlaceHolder("NAME").StringsVar(&c.policySetAddIgnoreExtendedAttributes)
	cmd.Flag("remove-ignore-extended-attribute", "List of extended attribute names to remove from the ignore extended attributes list").PlaceHolder("NAME").StringsVar(&c.policySetRemoveIgnoreExtendedAttributes)
	cmd.Flag("clear-ignore-extended-attributes", "Clear list of extended attribute names in the ignore attributes extended list").BoolVar(&c.policySetClearIgnoreExtendedAttributes)
}

func (c *policyFilesFlags) setFilesPolicyFromFlags(ctx context.Context, fp *policy.FilesPolicy, changeCount *int) error {
	if err := applyPolicyNumber64(ctx, "maximum file size", &fp.MaxFileSize, c.policySetMaxFileSize, changeCount); err != nil {
		return errors.Wrap(err, "maximum file size")
	}

	applyPolicyStringList(ctx, "dot-ignore filenames", &fp.DotIgnoreFiles, c.policySetAddDotIgnore, c.policySetRemoveDotIgnore, c.policySetClearDotIgnore, changeCount)
	applyPolicyStringList(ctx, "ignore rules", &fp.IgnoreRules, c.policySetAddIgnore, c.policySetRemoveIgnore, c.policySetClearIgnore, changeCount)

	if err := applyPolicyBoolPtr(ctx, "ignore cache dirs", &fp.IgnoreCacheDirectories, c.policyIgnoreCacheDirs, changeCount); err != nil {
		return err
	}

	applyPolicyStringList(ctx, "ignore extended attributes", &fp.IgnoreExtendedAttributes, c.policySetAddIgnoreExtendedAttributes, c.policySetRemoveIgnoreExtendedAttributes, c.policySetClearIgnoreExtendedAttributes, changeCount)

	return applyPolicyBoolPtr(ctx, "one filesystem", &fp.OneFileSystem, c.policyOneFileSystem, changeCount)
}
