package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"

    "github.com/blang/semver/v4"
    "golang.org/x/sys/windows/registry"
)

func banner() {
    fmt.Println(`
    _   _ _____ _____                      _                 
    | \ | | ____|_   _| __   _____ _ __ ___(_) ___  _ __  ___ 
    |  \| |  _|   | |   \ \ / / _ \ '__/ __| |/ _ \| '_ \/ __|
   _| |\  | |___  | |    \ V /  __/ |  \__ \ | (_) | | | \__ \
  (_)_| \_|_____| |_|     \_/ \___|_|  |___/_|\___/|_| |_|___/ ver. 0.3

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

    getCoreSdkRuntimes()

    fmt.Println()

    if !batchMode {
        fmt.Scanln()
    }
}

func writeDotnetClassicVersion(version string, spLevel string) {
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
            writeDotnetClassicVersion(name, "")
        } else {
            if sp != "" && install == "1" {
                writeDotnetClassicVersion(name, sp)
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
                    writeDotnetClassicVersion(name, "")
                } else {
                    if sp != "" && install == "1" {
                        writeDotnetClassicVersion(name, sp)
                    } else if install == "1" {
                        writeDotnetClassicVersion(name, "")
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
        writeDotnetClassicVersion(val, "")
    } else {
        rel, _, err := k.GetIntegerValue("Release")

        if err == nil {
            writeDotnetClassicVersion(checkFor45PlusVersion(int(rel)), "")
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

func writeDotnetCoreVersion(version string, location string) {
    version = strings.TrimSpace(version)
    if version == "" {
        return
    }

    fmt.Println(version, location)
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func getCoreSdkRuntimes() {
    const dotnet = "dotnet.exe"
    uniqueDotnetInstallDirs := make(map[string]interface{})

    // 1
    dotNetInstallDirEnv, exists := os.LookupEnv("DOTNET_INSTALL_DIR")
    if exists && dotNetInstallDirEnv != "" {
        if fileExists(filepath.Join(dotNetInstallDirEnv, dotnet)) {
            p1 := strings.TrimRight(strings.ToLower(dotNetInstallDirEnv), "\\/")
            uniqueDotnetInstallDirs[p1] = nil
        }
    }

    // 2
    pathEnv, exists := os.LookupEnv("PATH")
    if exists && pathEnv != "" {
        pathDirs := strings.Split(pathEnv, ";")

        for _, val := range pathDirs {
            if fileExists(filepath.Join(val, dotnet)) {
                p1 := strings.TrimRight(strings.ToLower(val), "\\/")
                uniqueDotnetInstallDirs[p1] = nil
            }
        }
    }

    // Enumerate SDK, RUNTIME versions
    var sdkList []string
    var runtimeList []string
    for val := range uniqueDotnetInstallDirs {
        sdkPath := filepath.Join(val, "sdk")
        runtimePath := filepath.Join(val, "shared")
        var line string

        dirs, err := os.ReadDir(sdkPath)
        if err == nil {
            for _, dir := range dirs {
                if dir.Type().IsDir() {
                    ver, err := semver.Make(dir.Name())
                    if err == nil {
                        line = "[SDK] " + ver.String() + " [" + filepath.Join(sdkPath, dir.Name()) + "]"
                        sdkList = append(sdkList, line)
                    }
                }
            }
        }

        dirs, err = os.ReadDir(runtimePath)
        if err == nil {
            for _, dir := range dirs {
                if dir.Type().IsDir() {
                    subs, err := os.ReadDir(filepath.Join(runtimePath, dir.Name()))
                    if err == nil {
                        for _, subd := range subs {
                            ver, err := semver.Make(subd.Name())
                            if err == nil {
                                line := "[" + dir.Name() + "] " + ver.String() + " [" + filepath.Join(runtimePath, dir.Name()) + "]"
                                runtimeList = append(runtimeList, line)
                            }
                        }
                    }
                }
            }
        }
    }

    if !batchMode {
        fmt.Println("\nCurrently installed .NET Core Versions in the system:")
    }

    for _, sdk := range sdkList {
        fmt.Println(sdk)
    }

    for _, runtime := range runtimeList {
        fmt.Println(runtime)
    }
}
