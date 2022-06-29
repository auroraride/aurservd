// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-29
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "github.com/auroraride/aurservd/pkg/utils"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
    "regexp"
)

func deleteAllMutations() {

}

func SplitMutation() {
    re := regexp.MustCompile(`(?m)/\*\* (.*)? START \*/\n([\s\S]*?)/\*\* .*? END \*/`)
    p := "internal/ent"
    b, err := ioutil.ReadFile(filepath.Join(p, "mutation.go"))
    if err != nil {
        log.Fatal(err)
    }

    for _, matches := range re.FindAllSubmatch(b, -1) {
        log.Println(utils.StrToSnakeCase(string(matches[1])))
    }

    log.Println(len(b), len(re.FindAll(b, -1)))
    os.Exit(0)
}
