package main

import (
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "bytes"
    "regexp"
    "github.com/gorilla/mux"
)

func shell(w http.ResponseWriter, r *http.Request) {
    var out bytes.Buffer
    var envPath,cluster,deployment string
    workDir := "/Users/ichtar/git/helm-charts"
    // get last version of code
    cmd := exec.Command("git","pull","--rebase")
    cmd.Dir = workDir
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    // get all informations from pulled file
    myR := regexp.MustCompile(".*(environments/([^/]*)/(.*).yaml).*")
    searchResult := myR.FindAllStringSubmatch(out.String(),-1)
    // if it match nothing in format environments/{cluster}/{deployment}.yaml abort and send Notfound
    if searchResult != nil {
        for _,i := range searchResult {
            envPath=i[1]
            cluster=i[2]
            deployment=i[3]
        }
        // remove dry-run when want to push real update
        cmd = exec.Command("helm","--dry-run","--debug","--kube-context",cluster,"upgrade",deployment,"-f",envPath,"charts/bm-stack")
        cmd.Dir = workDir
        cmd.Stdout = &out
        err = cmd.Run()
        if err != nil {
            log.Fatal(err)
        }
        fmt.Fprintf(w, out.String())
      } else {
        http.Error(w, "Forbidden", http.StatusForbidden)
    }
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/trigger", shell ).Methods("POST")
    log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
    handleRequests()
}

