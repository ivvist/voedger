/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package apppartsctl

import (
	"context"

	"github.com/voedger/voedger/pkg/appparts"
)

type appPartitionsController struct {
	parts       appparts.IAppPartitions
	builtInApps []appparts.BuiltInApp
}

func newAppPartitionsController(parts appparts.IAppPartitions, apps []appparts.BuiltInApp) (ctl IAppPartitionsController, cleanup func(), err error) {
	apc := appPartitionsController{parts: parts, builtInApps: apps}

	return &apc, func() {}, err
}

func (ctl *appPartitionsController) Prepare() (err error) {
	return err
}

func (ctl *appPartitionsController) Run(ctx context.Context) {
	<-ctx.Done()
}
