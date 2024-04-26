package diagnose

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/printutil"
	"github.com/lgdd/lfr-cli/lfr/pkg/util/procutil"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:     "diagnose",
		Aliases: []string{"diag"},
		Run:     diagnose,
		Short:   `Run a diagnosis to verify your environment for Liferay development.`,
	}
)

func diagnose(cmd *cobra.Command, args []string) {
	verifyJava()
	verifyBlade()
	dockerInstalled := verifyDocker()
	fmt.Print("\n")
	verifyBundles()
	if dockerInstalled {
		verifyLiferayDockerImages()
		verifyElasticsearchDockerImages()
	}
	fmt.Println("\nMore information about compatibilities: https://www.liferay.com/compatibility-matrix")
}

func verifyJava() bool {
	_, javaVersionCmdErr, err := procutil.Exec("java", "-version")

	if err != nil {
		printutil.Danger("[✗] ")
		fmt.Printf("Liferay requires Java 8 or 11.\n")
		return false
	}

	javaVersionResult := javaVersionCmdErr.String()

	if strings.Contains(javaVersionResult, "build 1.8") ||
		strings.Contains(javaVersionResult, "build 11.") {
		printutil.Success("[✓] ")
		fmt.Printf("Java intalled (%s)\n", strings.Split(javaVersionResult, "\n")[0])
		printlnBulletPoint("Make sure that your Java edition is a Java Technical Compatibility Kit (TCK) compliant build.")
		printlnBulletPoint("JDK compatibility is for runtime and project compile time. DXP source compile is compatible with JDK 8 only.")
	} else {
		printutil.Warning("[!] ")
		fmt.Printf("Java (%s)\n", strings.Split(javaVersionResult, "\n")[0])
		printlnBulletPointWarning("Liferay supports Java 8 and 11 only.")
	}
	return true
}

func verifyBlade() bool {
	bladeVersionCmdOut, _, err := procutil.Exec("blade", "version")

	if err != nil {
		printutil.Danger("[✗] ")
		fmt.Printf("Blade is not installed.\n")
		printlnBulletPoint("You might like this tool, but Blade is still the official one with useful features.")
		printlnBulletPoint("Blade is supported by Liferay and used by Liferay IDE behind the scenes.")
		printlnBulletPoint("Checkout the documentation: https://learn.liferay.com/w/dxp/building-applications/tooling/blade-cli")
		return false
	}

	bladeVersionResult := bladeVersionCmdOut.String()
	printutil.Success("[✓] ")
	fmt.Printf("Blade installed (%s)\n", strings.Split(bladeVersionResult, "\n")[0])
	return true
}

func verifyDocker() bool {
	dockerVersionCmdOut, _, err := procutil.Exec("docker", "version", "--format", "json")

	if err != nil {
		printutil.Warning("[!] ")
		fmt.Printf("Docker is not installed.\n")
		printlnBulletPoint("Docker is not required, but it's a easy way to get started and try Liferay DXP.")
		printlnBulletPoint("Checkout official images: https://hub.docker.com/u/liferay")
		return false
	}

	var dockerVersion DockerVersion
	dockerVersionResult := dockerVersionCmdOut.String()
	json.Unmarshal([]byte(dockerVersionResult), &dockerVersion)
	printutil.Success("[✓] ")
	fmt.Printf("Docker installed (%s)\n", strings.Split(dockerVersion.Server.Platform.Name, "\n")[0])
	return true
}

func verifyBundles() {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	liferayHomeDir := filepath.Join(homeDir, ".liferay")
	liferayHomeBundlesDir := filepath.Join(liferayHomeDir, "bundles")

	bundlesDirSize, err := fileutil.DirSize(liferayHomeBundlesDir)

	if err == nil && bundlesDirSize > 0 {
		printutil.Info("[i] ")
		fmt.Printf("Downloaded bundles are using ~%s.\n", humanize.Bytes(uint64(bundlesDirSize)))
		printlnBulletPoint("They are stored under " + liferayHomeBundlesDir)
	}
}

func verifyLiferayDockerImages() {
	dxpImagesTotalSize := getDockerImagesSize("liferay/dxp")
	portalImagesTotalSize := getDockerImagesSize("liferay/portal")
	dockerImagesTotalSize := dxpImagesTotalSize + portalImagesTotalSize

	if dockerImagesTotalSize > 0 {
		printutil.Info("[i] ")
		fmt.Printf("Official Liferay Docker images are using ~%s.\n", humanize.Bytes(dockerImagesTotalSize))
		printlnBulletPoint("Run 'docker images liferay/dxp' to list DXP Images (EE)")
		printlnBulletPoint("Run 'docker images liferay/portal' to list Portal Images (CE)")
	}
}

func verifyElasticsearchDockerImages() {
	dockerHubTagName := "elasticsearch"
	elasticHubTagName := "docker.elastic.co/elasticsearch/elasticsearch"
	dockerHubImagesTotalSize := getDockerImagesSize(dockerHubTagName)
	elasticHubImagesTotalSize := getDockerImagesSize(elasticHubTagName)
	dockerImagesTotalSize := dockerHubImagesTotalSize + elasticHubImagesTotalSize

	if dockerImagesTotalSize > 0 {
		printutil.Info("[i] ")
		fmt.Printf("Official Elasticsearch Docker images are using ~%s.\n", humanize.Bytes(dockerImagesTotalSize))
		printlnBulletPoint(fmt.Sprintf("Run 'docker images %s' to list images from Docker Hub", dockerHubTagName))
		printlnBulletPoint(fmt.Sprintf("Run 'docker images %s' to list images from Elastic Hub", elasticHubTagName))
	}
}

func getDockerImagesSize(tag string) uint64 {
	var dockerImagesTotalSize uint64

	dockerImagesCmdOut, _, err := procutil.Exec("docker", "images", tag, "--format", "{{.Size}}")

	if err == nil {
		dockerImagesResult := dockerImagesCmdOut.String()
		dockerImagesSizes := strings.Split(dockerImagesResult, "\n")

		for _, dockerImageSize := range dockerImagesSizes {
			if len(dockerImageSize) > 0 {
				currentDockerImageSize, err := parseBytesFromString(dockerImageSize)
				if err == nil {
					dockerImagesTotalSize = dockerImagesTotalSize + currentDockerImageSize
				}
			}
		}

	}
	return dockerImagesTotalSize
}

func parseBytesFromString(size string) (uint64, error) {
	sizeSplit0 := strings.Split(size, "B")[0]
	unit := sizeSplit0[len(sizeSplit0)-1:] + "B"
	sizeParsed, err := strconv.ParseFloat(strings.Split(size, unit)[0], 64)

	if err != nil {
		return 0, err
	}

	switch unit {
	case "MB":
		return uint64(sizeParsed * math.Pow(1000, 2)), nil
	case "GB":
		return uint64(sizeParsed * math.Pow(1000, 3)), nil
	default:
		return 0, errors.New("Unexpected unit for a Liferay docker image")
	}
}

type DockerVersion struct {
	Server struct {
		Platform struct {
			Name string `json:"Name"`
		} `json:"Platform"`
	} `json:"Server"`
}

func printlnBulletPoint(msg string) {
	printutil.Bold("    • ")
	printutil.Bold(msg + "\n")
}

func printlnBulletPointWarning(msg string) {
	printutil.Bold("    ! ")
	printutil.Bold(msg + "\n")
}

func printlnBulletPointDanger(msg string) {
	printutil.Danger("    ✗ ")
	printutil.Bold(msg + "\n")
}
