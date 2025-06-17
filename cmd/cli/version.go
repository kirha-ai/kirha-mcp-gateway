package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const Version = "0.0.2"

type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit,omitempty"`
	BuildDate string `json:"build_date,omitempty"`
	GoVersion string `json:"go_version,omitempty"`
}

type NPMRegistryResponse struct {
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
}

func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "Display the current version of Kirha MCP Gateway",
		RunE:  runVersion,
	}

	return cmd
}

func runVersion(cmd *cobra.Command, args []string) error {
	versionInfo := VersionInfo{
		Version: Version,
	}

	output, err := json.MarshalIndent(versionInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal version info: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func checkForUpdates() {
	go func() {
		latest, err := getLatestVersion()
		if err != nil {
			return
		}

		if latest != Version {
			fmt.Fprintf(os.Stderr, "\n🔄 Update available: %s → %s\n", Version, latest)
			fmt.Fprintf(os.Stderr, "Run 'npx @kirha/mcp-gateway@latest update' to update\n\n")
		}
	}()
}

func getLatestVersion() (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get("https://registry.npmjs.org/@kirha/mcp-gateway")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("registry returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var registryResp NPMRegistryResponse
	if err := json.Unmarshal(body, &registryResp); err != nil {
		return "", err
	}

	return registryResp.DistTags.Latest, nil
}
