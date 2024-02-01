package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var (
		cmd               string
		cmdDuration       int
		cmdStatus         int
		threshold         int
		currFgAppAsn      string
		prevFgAppAsn      string
		prevFgAppBundleID string
	)

	flag.StringVar(&cmd, "cmd", "", "Command to run")
	flag.IntVar(&cmdDuration, "duration", -1, "Duration of command")
	flag.IntVar(&cmdStatus, "status", -1, "Status of command")
	flag.IntVar(&threshold, "threshold", 5000, "Threshold for command duration in ms")
	flag.StringVar(&prevFgAppAsn, "prev_fg_app_asn", "", "Foreground App ASN before command")
	flag.StringVar(&prevFgAppBundleID, "prev_fg_app_bundleid", "", "Foreground App bundleID before command (optional)")
	flag.StringVar(&currFgAppAsn, "curr_fg_app_asn", "", "Foreground App ASN after command")
	flag.Parse()

	missing := []string{}
	if cmd == "" {
		missing = append(missing, "cmd")
	}
	if cmdDuration == -1 {
		missing = append(missing, "duration")
	}
	if cmdStatus == -1 {
		missing = append(missing, "status")
	}
	if prevFgAppAsn == "" {
		missing = append(missing, "prev_fg_app_asn")
	}
	if currFgAppAsn == "" {
		missing = append(missing, "curr_fg_app_asn")
	}
	if len(missing) > 0 {
		fmt.Println("missing:", strings.Join(missing, ", "))
		flag.Usage()
		os.Exit(1)
		return
	}

	// don't show a notification for these commands
	excluded := []string{"bash", "less", "man", "more", "ssh", "nvim", "vim", "webpack-dev-server", "tmux"}
	for _, e := range excluded {
		if strings.HasPrefix(cmd, e) {
			return
		}
	}
	// don't show a notification for commands that run faster than the threshold
	if cmdDuration < threshold {
		return
	}
	// don't show a notification if the foreground app didn't change
	if currFgAppAsn == prevFgAppAsn {
		return
	}

	if prevFgAppBundleID == "" {
		bundleInfo, err := execCmd("lsappinfo", "info", "-only", "bundleid", prevFgAppAsn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		bundleInfoParts := strings.Split(strings.TrimSpace(bundleInfo), "\"")
		if len(bundleInfoParts) == 5 {
			prevFgAppBundleID = bundleInfoParts[3]
		} else {
			fmt.Printf("warning: invalid bundle info: '%s'\n", bundleInfo)
		}

	}

	durationSeconds := float64(cmdDuration) / 1000.0
	var msg string
	var sound string
	if cmdStatus != 0 {
		msg = fmt.Sprintf("Failed in %.2fs (exit code %d)", durationSeconds, cmdStatus)
		sound = "Sosumi.aiff"
	} else {
		msg = fmt.Sprintf("Finished in %.2fs", durationSeconds)
		sound = "Hero.aiff"
	}

	args := []string{
		"-group", "fish-slow-commands",
		"-title", cmd,
		"-message", msg,
		"-sound", sound,
		"-ignoreDnD",
	}
	if prevFgAppBundleID != "" {
		args = append(args, "-sender", prevFgAppBundleID)
	}

	_, err := execCmd("terminal-notifier", args...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}

func execCmd(bin string, args ...string) (string, error) {
	cmd := exec.Command(bin, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, stderr.String())
	}
	return stdout.String(), nil
}
