package cmd

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"

    "github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
    Use:   "run",
    Short: "Generate and run a C++ program",
    Long:  `Generate a C++ program from the provided algorithm and main code, then compile and run it.`,
    Run: func(cmd *cobra.Command, args []string) {
        algorithmCode, _ := cmd.Flags().GetString("algorithm")
        mainCode, _ := cmd.Flags().GetString("main")

        if algorithmCode == "" || mainCode == "" {
            fmt.Println("Algorithm and main code must be provided.")
            return
        }

        cppCode := fmt.Sprintf(`#include <iostream>
using namespace std;

%s

int main() {
    %s
    return 0;
}
`, algorithmCode, mainCode)

        tmpFile, err := ioutil.TempFile("", "program-*.cpp")
        if err != nil {
            fmt.Println("Error creating temporary file:", err)
            return
        }
        defer os.Remove(tmpFile.Name())

        if _, err := tmpFile.Write([]byte(cppCode)); err != nil {
            fmt.Println("Error writing to temporary file:", err)
            return
        }
        tmpFile.Close()

        exeFile := strings.TrimSuffix(tmpFile.Name(), ".cpp")
        cmdCompile := exec.Command("g++", "-o", exeFile, tmpFile.Name())
        if output, err := cmdCompile.CombinedOutput(); err != nil {
            fmt.Printf("Error compiling C++ code: %s\n", string(output))
            return
        }

        cmdRun := exec.Command(exeFile)
        cmdRun.Stdout = os.Stdout
        cmdRun.Stderr = os.Stderr
        if err := cmdRun.Run(); err != nil {
            fmt.Println("Error running compiled program:", err)
        }
    },
}

func init() {
    rootCmd.AddCommand(runCmd)

    runCmd.Flags().StringP("algorithm", "a", "", "C++ algorithm code")
    runCmd.Flags().StringP("main", "m", "", "C++ main function code")
}
