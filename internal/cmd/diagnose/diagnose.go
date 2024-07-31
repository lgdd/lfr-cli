package diagnose

import (
	"encoding/json"
	"errors"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/lgdd/lfr-cli/pkg/util/fileutil"
	"github.com/lgdd/lfr-cli/pkg/util/logger"

	"github.com/lgdd/lfr-cli/pkg/util/procutil"
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
	verifyLCP()
	dockerInstalled := verifyDocker()
	logger.Print("\n")
	verifyBundles()
	if dockerInstalled {
		verifyLiferayDockerImages()
		verifyElasticsearchDockerImages()
	}
	logger.Println("\nMore information about compatibilities: https://www.liferay.com/compatibility-matrix")
}

func verifyJava() bool {
	major, version, err := procutil.GetCurrentJavaVersion()

	if err != nil {
		logger.PrintError("[✗] ")
		logger.Print("Liferay requires Java 8, 11, 17 or 21.\n")
		return false
	}

	if major == "8" || major == "11" {
		logger.PrintWarn("[!] ")
		logger.Printf("Java intalled (%s)\n", strings.Split(version, "\n")[0])
		if major == "8" {
			logger.PrintlnBold("\t! Liferay DXP 2024.Q1 and Liferay Portal 7.4 GA112 will be the last version to support Java 8.")
		}
		if major == "11" {
			logger.PrintlnBold("\t! Liferay DXP DXP 2024.Q2 and Liferay Portal 7.4 GA120 will be the last version to support Java 11.")
		}
	} else if major == "17" || major == "21" {
		logger.PrintSuccess("[✓] ")
		logger.Printf("Java intalled (%s)\n", strings.Split(version, "\n")[0])
		logger.PrintlnBold("\t• Liferay DXP 2024.Q2+ and Liferay Portal CE 7.4 GA120+ are fully certified to run on both Java JDK 17 and 21.")

	} else {
		logger.PrintWarn("[!] ")
		logger.Printf("Java installed (%s)\n", strings.Split(version, "\n")[0])
		logger.PrintlnBold("\t! Liferay supports Java 8 and 11 only.")
	}
	logger.PrintlnBold("\t• Make sure that your Java edition is a Java Technical Compatibility Kit (TCK) compliant build.")
	logger.PrintlnBold("\t• JDK compatibility is for runtime and project compile time.")
	return true
}

func verifyBlade() bool {
	bladeVersionCmdOut, _, err := procutil.Exec("blade", "version")

	if err != nil {
		logger.PrintError("[✗] ")
		logger.Printf("Blade is not installed.\n")
		logger.PrintlnBold("\t• You might like this tool, but Blade is still the official one with useful features.")
		logger.PrintlnBold("\t• Blade is supported by Liferay and used by Liferay IDE behind the scenes.")
		logger.PrintlnBold("\t• Checkout the documentation: https://learn.liferay.com/w/dxp/building-applications/tooling/blade-cli")
		return false
	}

	bladeVersionResult := bladeVersionCmdOut.String()
	logger.PrintSuccess("[✓] ")
	logger.Printf("Blade installed (%s)\n", strings.Split(bladeVersionResult, "\n")[0])
	return true
}

func verifyLCP() bool {
	lcpVersionCmdOut, _, err := procutil.Exec("lcp", "-v")

	if err != nil {
		logger.PrintWarn("[!] ")
		logger.Printf("LCP is not installed.\n")
		logger.PrintlnBold("\t• If you work on Liferay PaaS or Liferay SaaS, LCP can be used to view and manage your Liferay Cloud services.")
		logger.PrintlnBold("\t• Checkout the documentation: https://learn.liferay.com/w/liferay-cloud/reference/command-line-tool")
		return false
	}

	lcpVersionResult := lcpVersionCmdOut.String()
	logger.PrintSuccess("[✓] ")
	logger.Printf("LCP installed (%s)\n", strings.Split(lcpVersionResult, "\n")[0])
	return true
}

func verifyDocker() bool {
	dockerVersionCmdOut, _, err := procutil.Exec("docker", "version", "--format", "json")

	if err != nil {
		logger.PrintWarn("[!] ")
		logger.Printf("Docker is not installed.\n")
		logger.PrintlnBold("\t• Docker is not required, but it's a easy way to get started and try Liferay DXP.")
		logger.PrintlnBold("\t• Checkout official images: https://hub.docker.com/u/liferay")
		return false
	}

	var dockerVersion DockerVersion
	dockerVersionResult := dockerVersionCmdOut.String()
	json.Unmarshal([]byte(dockerVersionResult), &dockerVersion)
	logger.PrintSuccess("[✓] ")
	logger.Printf("Docker installed (%s)\n", strings.Split(dockerVersion.Server.Platform.Name, "\n")[0])
	return true
}

func verifyBundles() {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		logger.Fatal(err.Error())
	}

	liferayHomeDir := filepath.Join(homeDir, ".liferay")
	liferayHomeBundlesDir := filepath.Join(liferayHomeDir, "bundles")

	bundlesDirSize, err := fileutil.DirSize(liferayHomeBundlesDir)

	if err == nil && bundlesDirSize > 0 {
		logger.PrintWarn("[!] ")
		logger.Printf("Downloaded bundles are using ~%s.\n", humanize.Bytes(uint64(bundlesDirSize)))
		logger.PrintlnBold("\t• They are stored under " + liferayHomeBundlesDir)
	}
}

func verifyLiferayDockerImages() {
	dxpImagesTotalSize := getDockerImagesSize("liferay/dxp")
	portalImagesTotalSize := getDockerImagesSize("liferay/portal")
	dockerImagesTotalSize := dxpImagesTotalSize + portalImagesTotalSize

	if dockerImagesTotalSize > 0 {
		logger.PrintWarn("[!] ")
		logger.Printf("Official Liferay Docker images are using ~%s.\n", humanize.Bytes(dockerImagesTotalSize))
		logger.PrintlnBold("\t• Run 'docker images liferay/dxp' to list DXP Images (EE)")
		logger.PrintlnBold("\t• Run 'docker images liferay/portal' to list Portal Images (CE)")
	}
}

func verifyElasticsearchDockerImages() {
	dockerHubTagName := "elasticsearch"
	elasticHubTagName := "docker.elastic.co/elasticsearch/elasticsearch"
	dockerHubImagesTotalSize := getDockerImagesSize(dockerHubTagName)
	elasticHubImagesTotalSize := getDockerImagesSize(elasticHubTagName)
	dockerImagesTotalSize := dockerHubImagesTotalSize + elasticHubImagesTotalSize

	if dockerImagesTotalSize > 0 {
		logger.PrintInfo("[i] ")
		logger.Printf("Official Elasticsearch Docker images are using ~%s.\n", humanize.Bytes(dockerImagesTotalSize))
		logger.PrintfBold("\t• Run 'docker images %s' to list images from Docker Hub\n", dockerHubTagName)
		logger.PrintfBold("\t• Run 'docker images %s' to list images from Elastic Hub\n", elasticHubTagName)
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
		return 0, errors.New("unexpected unit for a Liferay docker image")
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
	logger.Print("    • ")
	logger.Print(msg)
}

func printlnBulletPointWarning(msg string) {
	logger.Print("    ! ")
	logger.Print(msg)
}
