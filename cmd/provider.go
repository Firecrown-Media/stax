package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/firecrown-media/stax/pkg/provider"
	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Manage hosting providers",
	Long: `Manage hosting providers for WordPress sites.

Stax supports multiple WordPress hosting providers through a pluggable
provider interface. Use these commands to list, inspect, and switch
between providers.`,
}

var providerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available providers",
	Long:  "List all registered hosting providers and their capabilities.",
	RunE:  runProviderList,
}

var providerShowCmd = &cobra.Command{
	Use:   "show <provider-name>",
	Short: "Show provider details",
	Long:  "Show detailed information about a specific provider including capabilities and requirements.",
	Args:  cobra.ExactArgs(1),
	RunE:  runProviderShow,
}

var providerSetCmd = &cobra.Command{
	Use:   "set <provider-name>",
	Short: "Set default provider for project",
	Long:  "Set the default provider for the current project in .stax.yml",
	Args:  cobra.ExactArgs(1),
	RunE:  runProviderSet,
}

var providerTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test current provider connection",
	Long:  "Test connectivity to the currently configured provider.",
	RunE:  runProviderTest,
}

var providerCompareCmd = &cobra.Command{
	Use:   "compare <provider1> <provider2>",
	Short: "Compare two providers",
	Long:  "Compare capabilities and features between two providers.",
	Args:  cobra.ExactArgs(2),
	RunE:  runProviderCompare,
}

var (
	providerOutputFormat string // json, yaml, table
	providerShowAll      bool
)

func init() {
	rootCmd.AddCommand(providerCmd)
	providerCmd.AddCommand(providerListCmd)
	providerCmd.AddCommand(providerShowCmd)
	providerCmd.AddCommand(providerSetCmd)
	providerCmd.AddCommand(providerTestCmd)
	providerCmd.AddCommand(providerCompareCmd)

	// Flags
	providerListCmd.Flags().StringVarP(&providerOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")
	providerShowCmd.Flags().StringVarP(&providerOutputFormat, "output", "o", "table", "Output format (table, json, yaml)")
	providerCompareCmd.Flags().StringVarP(&providerOutputFormat, "output", "o", "table", "Output format (table, json)")
}

func runProviderList(cmd *cobra.Command, args []string) error {
	infos, err := provider.GetAllProviderInfo()
	if err != nil {
		return fmt.Errorf("failed to get provider information: %w", err)
	}

	switch providerOutputFormat {
	case "json":
		return outputJSON(infos)
	case "yaml":
		return outputYAML(infos)
	default:
		return outputProviderListTable(infos)
	}
}

func outputProviderListTable(infos []*provider.ProviderInfo) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "PROVIDER\tDESCRIPTION\tDEFAULT\tCAPABILITIES")
	fmt.Fprintln(w, "--------\t-----------\t-------\t------------")

	for _, info := range infos {
		defaultMarker := ""
		if info.IsDefault {
			defaultMarker = "*"
		}

		// Count core capabilities
		caps := info.Capabilities
		coreCount := 0
		if caps.Authentication {
			coreCount++
		}
		if caps.SiteManagement {
			coreCount++
		}
		if caps.DatabaseExport {
			coreCount++
		}
		if caps.DatabaseImport {
			coreCount++
		}
		if caps.FileSync {
			coreCount++
		}

		// Count optional capabilities
		optionalCount := 0
		if caps.Deployment {
			optionalCount++
		}
		if caps.Environments {
			optionalCount++
		}
		if caps.Backups {
			optionalCount++
		}
		if caps.RemoteExecution {
			optionalCount++
		}
		if caps.MediaManagement {
			optionalCount++
		}

		capsSummary := fmt.Sprintf("%d core, %d optional", coreCount, optionalCount)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			info.Name,
			info.Description,
			defaultMarker,
			capsSummary,
		)
	}

	return nil
}

func runProviderShow(cmd *cobra.Command, args []string) error {
	providerName := args[0]

	info, err := provider.GetProviderInfo(providerName)
	if err != nil {
		return fmt.Errorf("failed to get provider info: %w", err)
	}

	switch providerOutputFormat {
	case "json":
		return outputJSON(info)
	case "yaml":
		return outputYAML(info)
	default:
		return outputProviderShowTable(info)
	}
}

func outputProviderShowTable(info *provider.ProviderInfo) error {
	fmt.Printf("Provider: %s\n", info.Name)
	fmt.Printf("Description: %s\n", info.Description)
	fmt.Printf("Default: %v\n\n", info.IsDefault)

	fmt.Println("Capabilities:")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	caps := info.Capabilities

	fmt.Fprintln(w, "  CAPABILITY\tSUPPORTED")
	fmt.Fprintln(w, "  ----------\t---------")
	fmt.Fprintf(w, "  Authentication\t%v\n", caps.Authentication)
	fmt.Fprintf(w, "  Site Management\t%v\n", caps.SiteManagement)
	fmt.Fprintf(w, "  Database Export\t%v\n", caps.DatabaseExport)
	fmt.Fprintf(w, "  Database Import\t%v\n", caps.DatabaseImport)
	fmt.Fprintf(w, "  File Sync\t%v\n", caps.FileSync)
	fmt.Fprintf(w, "  Deployment\t%v\n", caps.Deployment)
	fmt.Fprintf(w, "  Environments\t%v\n", caps.Environments)
	fmt.Fprintf(w, "  Backups\t%v\n", caps.Backups)
	fmt.Fprintf(w, "  Remote Execution\t%v\n", caps.RemoteExecution)
	fmt.Fprintf(w, "  Media Management\t%v\n", caps.MediaManagement)
	fmt.Fprintf(w, "  SSH Access\t%v\n", caps.SSHAccess)
	fmt.Fprintf(w, "  API Access\t%v\n", caps.APIAccess)
	fmt.Fprintf(w, "  Scaling\t%v\n", caps.Scaling)
	fmt.Fprintf(w, "  Monitoring\t%v\n", caps.Monitoring)
	fmt.Fprintf(w, "  Logging\t%v\n", caps.Logging)

	return nil
}

func runProviderSet(cmd *cobra.Command, args []string) error {
	providerName := args[0]

	// Verify provider exists
	if !provider.ProviderExists(providerName) {
		return fmt.Errorf("provider %s not found", providerName)
	}

	// TODO: Update .stax.yml with new provider
	// For now, just show what would be done
	fmt.Printf("Would set provider to: %s\n", providerName)
	fmt.Println("TODO: Implement .stax.yml update")

	return nil
}

func runProviderTest(cmd *cobra.Command, args []string) error {
	// TODO: Load configuration and create provider instance
	// For now, show placeholder
	fmt.Println("Testing provider connection...")
	fmt.Println("TODO: Implement provider connection test")

	return nil
}

func runProviderCompare(cmd *cobra.Command, args []string) error {
	provider1 := args[0]
	provider2 := args[1]

	comparison, err := provider.CompareProviders(provider1, provider2)
	if err != nil {
		return fmt.Errorf("failed to compare providers: %w", err)
	}

	switch providerOutputFormat {
	case "json":
		return outputJSON(comparison)
	default:
		return outputProviderComparisonTable(comparison)
	}
}

func outputProviderComparisonTable(comparison *provider.ProviderComparison) error {
	fmt.Printf("Comparing: %s vs %s\n\n", comparison.Provider1, comparison.Provider2)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "CAPABILITY\t"+comparison.Provider1+"\t"+comparison.Provider2)
	fmt.Fprintln(w, "----------\t----------\t----------")

	caps1 := comparison.Capabilities1
	caps2 := comparison.Capabilities2

	printCapabilityRow := func(name string, cap1, cap2 bool) {
		fmt.Fprintf(w, "%s\t%v\t%v\n", name, cap1, cap2)
	}

	printCapabilityRow("Authentication", caps1.Authentication, caps2.Authentication)
	printCapabilityRow("Site Management", caps1.SiteManagement, caps2.SiteManagement)
	printCapabilityRow("Database Export", caps1.DatabaseExport, caps2.DatabaseExport)
	printCapabilityRow("Database Import", caps1.DatabaseImport, caps2.DatabaseImport)
	printCapabilityRow("File Sync", caps1.FileSync, caps2.FileSync)
	printCapabilityRow("Deployment", caps1.Deployment, caps2.Deployment)
	printCapabilityRow("Environments", caps1.Environments, caps2.Environments)
	printCapabilityRow("Backups", caps1.Backups, caps2.Backups)
	printCapabilityRow("Remote Execution", caps1.RemoteExecution, caps2.RemoteExecution)
	printCapabilityRow("Media Management", caps1.MediaManagement, caps2.MediaManagement)
	printCapabilityRow("SSH Access", caps1.SSHAccess, caps2.SSHAccess)
	printCapabilityRow("API Access", caps1.APIAccess, caps2.APIAccess)

	fmt.Fprintf(w, "\nShared Features: %d\n", len(comparison.SharedFeatures))

	return nil
}

func outputJSON(v interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(v)
}

func outputYAML(v interface{}) error {
	// TODO: Implement YAML output
	return fmt.Errorf("YAML output not yet implemented")
}
