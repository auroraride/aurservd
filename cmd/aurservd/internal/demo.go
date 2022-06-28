// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package internal

import (
    "context"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent/rider"
)

func Demo() {
    // res := esign.New().DocTemplate("4a4a554bc42e4deda22b5d503d0661f6")
    // for _, component := range res.StructComponents {
    //     fmt.Printf("%s\t\t%s\n", component.Key, component.Context.Label)
    // }
    // os.Exit(0)
    service.NewContract().Sign(ar.Ent.Rider.QueryNotDeleted().WithPerson().Where(rider.ID(68719476736)).FirstX(context.Background()), nil)
}
