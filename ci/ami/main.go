package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type AMIBuildConfig struct {
	K8sReleases map[string]string `json:"k8s_releases"`
}

type ReleaseVersion struct {
	Major int
	Minor int
	Patch int
}

func (r *ReleaseVersion) toString() string {
	return "v" + strconv.Itoa(r.Major) + "." + strconv.Itoa(r.Minor) + "." + strconv.Itoa(r.Patch)
}

func BuildReleaseVersion(ver string) ReleaseVersion {
	verSplit := strings.Split(string(ver), ".")
	major, err := strconv.Atoi(strings.ReplaceAll(verSplit[0], "v", ""))
	check(err)
	minor, err := strconv.Atoi(verSplit[1])
	check(err)
	patch, err := strconv.Atoi(verSplit[2])
	check(err)

	return ReleaseVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func main() {
	var m2, m3 string
	var updated bool
	url := "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
	k8sReleaseResponse, err := http.Get(url)
	check(err)

	min1, err := ioutil.ReadAll(k8sReleaseResponse.Body)
	check(err)

	min1Release := BuildReleaseVersion(string(min1))

	log.Print("Info: min1Release: ReleaseVersion ", min1Release.toString())
	log.Print("Info: min1Release: Major ", min1Release.Major, ", Minor ", min1Release.Minor, ", Patch ", min1Release.Patch)
	fmt.Println()

	if min1Release.Minor >= 2 {
		m2 = strconv.Itoa(min1Release.Major) + "." + strconv.Itoa(min1Release.Minor-1)
		m3 = strconv.Itoa(min1Release.Major) + "." + strconv.Itoa(min1Release.Minor-2)
	}

	url = fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/stable-%s.txt", m2)
	k8sReleaseResponse, err = http.Get(url)
	check(err)

	min2, err := ioutil.ReadAll(k8sReleaseResponse.Body)
	check(err)

	min2Release := BuildReleaseVersion(string(min2))

	log.Print("Info: min2Release: ReleaseVersion ", min2Release.toString())
	log.Print("Info: min2Release: Major ", min2Release.Major, ", Minor ", min2Release.Minor, ", Patch ", min2Release.Patch)
	fmt.Println()

	url = fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/stable-%s.txt", m3)
	k8sReleaseResponse, err = http.Get(url)
	check(err)

	min3, err := ioutil.ReadAll(k8sReleaseResponse.Body)
	check(err)

	min3Release := BuildReleaseVersion(string(min3))

	log.Print("Info: min2Release: ReleaseVersion ", min3Release.toString())
	log.Print("Info: min2Release: Major ", min3Release.Major, ", Minor ", min3Release.Minor, ", Patch ", min3Release.Patch)
	fmt.Println()

	latestAMIBuildConfig := &AMIBuildConfig{
		K8sReleases: map[string]string{
			"min1": string(min1),
			"min2": string(min2),
			"min3": string(min3),
		},
	}

	fmt.Println(*latestAMIBuildConfig)

	latestAMIBuildConfigFileBytes, err := json.MarshalIndent(latestAMIBuildConfig, "", "  ")
	check(err)

	AMIBuildConfigFilepath := "AMIBuildConfig.json"

	dat, err := os.ReadFile(AMIBuildConfigFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.WriteFile("AMIBuildConfig.json", latestAMIBuildConfigFileBytes, 0666)
			check(err)
			log.Printf("Info: Created \"AMIBuildConfig.json\" K8s versions \"%s\"", latestAMIBuildConfig.K8sReleases)

			os.Setenv("K8S_UPDATED", "true")

			return
		} else {
			os.Setenv("K8S_UPDATED", "false")

			log.Fatal(err)
		}
	}

	currentAMIBuildConfig := new(AMIBuildConfig)
	err = json.Unmarshal(dat, currentAMIBuildConfig)
	check(err)
	if !cmp.Equal(currentAMIBuildConfig, latestAMIBuildConfig) {
		err = os.WriteFile(AMIBuildConfigFilepath, latestAMIBuildConfigFileBytes, 0666)
		check(err)

		log.Printf("Info: Updated \"%s\" with K8s versions from \"%s\" to \"%s\"", AMIBuildConfigFilepath, currentAMIBuildConfig.K8sReleases, latestAMIBuildConfig.K8sReleases)
		updated = true
	} else {
		log.Printf("Info: \"%s\" is up-to-date.", AMIBuildConfigFilepath)
		updated = false
	}

	if updated {
		os.Setenv("K8S_UPDATED", "true")
	} else {
		os.Setenv("K8S_UPDATED", "false")
	}
}
