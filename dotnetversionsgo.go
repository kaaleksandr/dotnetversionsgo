package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func main() {
	batchMode := false

	args := os.Args[1:]

	if len(args) > 0 && (args[0] == "/help") {
		fmt.Println("Writes all the currently installed versions of \"classic\" .NET platform in the system.\r\nUse --b, -b or /b to use in a batch, showing only the installed versions, without any extra informational lines.")
	} else {
		if len(args) > 0 && (args[0] == "/b") {
			batchMode = true
		}

		if !batchMode {
			fmt.Println("Currently installed \"classic\" .NET Versions in the system:")
		}

		Get45PlusFromRegistry()
	}

	if !batchMode {
		fmt.Scanln()
	}
}

func WriteVersion(version string, spLevel string) {
	version = strings.Trim(version, " \t\n\r")
	if len(version) == 0 {
		return
	}

	spLevelString := ""
	if len(spLevelString) > 0 {
		spLevelString = " Service Pack " + spLevel
	}

	fmt.Println(version, spLevelString)
}

func Get1To45VersionFromRegistry() {
	const subkey string = "SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\"

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey, registry.READ)
	if err != nil {
		fmt.Println(err)

		return
	}

	defer k.Close()

	keys, err := k.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Println(err)

		return
	}

	for _, subkey1 := range keys {
		if subkey1 == "v4" {
			continue
		}

		if strings.HasPrefix(subkey1, "v") {
			versionKey, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey+subkey1, registry.READ)
			if err != nil {
				continue
			}

			defer versionKey.Close()

			name, _, err := versionKey.GetStringValue("Version")
			if err == nil {
				name = ""
			}

			sp, _, err := versionKey.GetStringValue("SP")
			if err == nil {
				sp = ""
			}

			install, _, err := versionKey.GetStringValue("Install")
			if err == nil || len(install) < 1 {
				WriteVersion(name, sp)
			} else {
				if len(sp) > 0 && install == "1" {
					WriteVersion(name, sp)
				}
			}

			if len(name) > 0 {
				continue
			}

			vers, err := versionKey.ReadSubKeyNames(-1)
			if err == nil {
				for _, subkey2 := range vers {
					subKeyVer, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey+subkey1+"\\"+subkey2, registry.READ)
					if err != nil {
						continue
					}

					defer subKeyVer.Close()

					name, _, err := subKeyVer.GetStringValue("Version")

					if err != nil {
						sp, _, err = subKeyVer.GetStringValue("SP")

						if err != nil {
							fmt.Println(err)
						}
					}

					install, _, err := subKeyVer.GetStringValue("Install")
					if err != nil {
						WriteVersion(name, sp)
					} else {
						if len(sp) > 0 && install == "1" {
							WriteVersion(name, sp)
						} else if install == "1" {
							WriteVersion(name, "")
						}
					}
				}
			}
		}
	}
}

func Get45PlusFromRegistry() {
	const subkey string = "SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v4\\Full\\"

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey, registry.READ)
	if err != nil {
		fmt.Println(err)

		return
	}

	defer k.Close()

	val, _, err := k.GetStringValue("Version")

	if err == nil {
		WriteVersion(val, "")
	} else {
		rel, _, err := k.GetIntegerValue("Release")

		if err == nil {
			WriteVersion(CheckFor45PlusVersion(int(rel)), "")
		}
	}
}

func CheckFor45PlusVersion(releaseKey int) string {
	switch {
	case releaseKey >= 533325:
		return "4.8.1"
	case releaseKey >= 528040:
		return "4.8"
	case releaseKey >= 461808:
		return "4.7.2"
	case releaseKey >= 461308:
		return "4.7.1"
	case releaseKey >= 460798:
		return "4.7"
	case releaseKey >= 394802:
		return "4.6.2"
	case releaseKey >= 394254:
		return "4.6.1"
	case releaseKey >= 393295:
		return "4.6"
	case releaseKey >= 379893:
		return "4.5.2"
	case releaseKey >= 378675:
		return "4.5.1"
	case releaseKey >= 378389:
		return "4.5"
	default:
		return ""
	}
}
