package main

import (
    "flag"
    "fmt"
    "os"
    "strconv"
    "strings"

    "golang.org/x/sys/windows/registry"
)

func banner() {
    fmt.Println(`
    _   _ _____ _____                      _                 
    | \ | | ____|_   _| __   _____ _ __ ___(_) ___  _ __  ___ 
    |  \| |  _|   | |   \ \ / / _ \ '__/ __| |/ _ \| '_ \/ __|
   _| |\  | |___  | |    \ V /  __/ |  \__ \ | (_) | | | \__ \
  (_)_| \_|_____| |_|     \_/ \___|_|  |___/_|\___/|_| |_|___/ ver. 0.2

  `)
}

var noBanner bool
var batchMode bool

func init() {
    flag.BoolVar(&noBanner, "nobanner", false, "No banner")
    flag.BoolVar(&batchMode, "batch", false, "Enable batch mode")
}

func main() {
    flag.Parse()

    if !noBanner {
        banner()
    }

    if !batchMode {
        fmt.Println("Currently installed \"classic\" .NET Versions in the system:")
    }

    get1To45VersionFromRegistry()
    get45PlusFromRegistry()

    fmt.Println()

    if !batchMode {
        fmt.Scanln()
    }
}

func writeVersion(version string, spLevel string) {
    version = strings.TrimSpace(version)
    if version == "" {
        return
    }

    spLevelString := ""
    if spLevel != "" {
        spLevelString = " Service Pack " + spLevel
    }

    fmt.Println(version, spLevelString)
}

func get1To45VersionFromRegistry() {
    const subkey = `SOFTWARE\Microsoft\NET Framework Setup\NDP\`

    k, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey, registry.READ)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }

    defer k.Close()

    keys, err := k.ReadSubKeyNames(-1)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }

    for _, versionKeyName := range keys {
        if versionKeyName == "v4" {
            continue
        }

        if !strings.HasPrefix(versionKeyName, "v") {
            continue
        }

        versionKey, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey+versionKeyName, registry.READ)
        defer versionKey.Close()

        var name string
        var install string
        var sp string

        name, _, err = versionKey.GetStringValue("Version")
        spUint, _, err := versionKey.GetIntegerValue("SP")
        if err == nil {
            sp = strconv.FormatUint(spUint, 10)
        }

        installUint, _, err := versionKey.GetIntegerValue("Install")
        if err == nil {
            install = strconv.FormatUint(installUint, 10)
        }

        if install == "" {
            writeVersion(name, "")
        } else {
            if sp != "" && install == "1" {
                writeVersion(name, sp)
            }
        }

        if name != "" {
            continue
        }

        vers, err := versionKey.ReadSubKeyNames(-1)
        if err == nil {
            for _, subKeyName := range vers {
                subKey, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey+versionKeyName+"\\"+subKeyName, registry.READ)
                if err != nil {
                    continue
                }
                defer subKey.Close()

                sp = ""
                install = ""
                name, _, err = subKey.GetStringValue("Version")

                if name != "" {
                    spUint, _, err := versionKey.GetIntegerValue("SP")
                    if err == nil {
                        sp = strconv.FormatUint(spUint, 10)
                    }
                }

                installUint, _, err := versionKey.GetIntegerValue("Install")
                if err == nil {
                    install = strconv.FormatUint(installUint, 10)
                }

                if install == "" {
                    writeVersion(name, "")
                } else {
                    if sp != "" && install == "1" {
                        writeVersion(name, sp)
                    } else if install == "1" {
                        writeVersion(name, "")
                    }
                }
            }
        }
    }
}

func get45PlusFromRegistry() {
    const subkey string = "SOFTWARE\\Microsoft\\NET Framework Setup\\NDP\\v4\\Full\\"

    k, err := registry.OpenKey(registry.LOCAL_MACHINE, subkey, registry.READ)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }

    defer k.Close()

    val, _, err := k.GetStringValue("Version")

    if err == nil {
        writeVersion(val, "")
    } else {
        rel, _, err := k.GetIntegerValue("Release")

        if err == nil {
            writeVersion(checkFor45PlusVersion(int(rel)), "")
        }
    }
}

func checkFor45PlusVersion(releaseKey int) string {
    switch {
    case releaseKey >= 533320:
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
