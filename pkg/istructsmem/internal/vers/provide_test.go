/*
 * Copyright (c) 2021-present Sigma-Soft, Ltd.
 * @author: Nikolay Nikitin
 */

package vers

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/voedger/voedger/pkg/istorage"
	"github.com/voedger/voedger/pkg/istorageimpl"
	"github.com/voedger/voedger/pkg/istructs"
)

func Test_BasicUsage(t *testing.T) {
	sp := istorageimpl.Provide(istorage.ProvideMem())
	storage, _ := sp.AppStorage(istructs.AppQName_test1_app1)

	versions := NewVersions()
	if err := versions.Prepare(storage); err != nil {
		panic(err)
	}

	t.Run("basic Versions methods", func(t *testing.T) {
		require := require.New(t)

		const (
			key VersionKey   = 55
			ver VersionValue = 77
		)

		require.Equal(UnknownVersion, versions.GetVersion(key))

		versions.PutVersion(key, ver)
		require.Equal(ver, versions.GetVersion(key))

		t.Run("must be able to load early stored versions", func(t *testing.T) {
			otherVersions := NewVersions()
			if err := otherVersions.Prepare(storage); err != nil {
				panic(err)
			}
			require.Equal(ver, otherVersions.GetVersion(key))
		})
	})

}