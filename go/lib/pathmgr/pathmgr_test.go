// Copyright 2017 ETH Zurich
// Copyright 2018 ETH Zurich, Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pathmgr

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/scionproto/scion/go/lib/ctrl/path_mgmt"
	"github.com/scionproto/scion/go/lib/infra"
	"github.com/scionproto/scion/go/lib/pathpol"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/sciond/mock_sciond"
	"github.com/scionproto/scion/go/lib/spath/spathmeta"
	"github.com/scionproto/scion/go/lib/xtest"
)

const timeUnitDuration time.Duration = 10 * time.Millisecond

func getDuration(units time.Duration) time.Duration {
	return units * timeUnitDuration
}

func TestQuery(t *testing.T) {
	t.Log("Query")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sd := mock_sciond.NewMockConnector(ctrl)
	pm := New(sd, Timers{})

	srcIA, dstIA := xtest.MustParseIA("1-ff00:0:133"), xtest.MustParseIA("1-ff00:0:131")

	paths := []string{}

	f := func(p string) {
		t.Log("get")
		paths = append(paths, p)
		sd.EXPECT().Paths(gomock.Any(), dstIA, srcIA, gomock.Any(), gomock.Any()).Return(
			buildSDAnswer(paths...), nil,
		)
		aps := pm.Query(context.Background(), srcIA, dstIA, sciond.PathReqFlags{})
		assert.Len(t, aps, len(paths), fmt.Sprintf("get %d paths", len(paths)))
		assert.ElementsMatch(t, getPathStrings(aps), paths)
	}

	pathOne := fmt.Sprintf("%s#1019 1-ff00:0:132#1910 1-ff00:0:132#1916 %s#1619", srcIA, dstIA)
	f(pathOne)

	pathTwo := fmt.Sprintf("%s#101902 1-ff00:0:132#191002 1-ff00:0:132#1916 %s#1619", srcIA, dstIA)
	f(pathTwo)
}

var allowEntry = &pathpol.ACLEntry{Action: pathpol.Allow, Rule: pathpol.NewHopPredicate()}
var denyEntry = &pathpol.ACLEntry{Action: pathpol.Deny, Rule: pathpol.NewHopPredicate()}

func TestQueryFilter(t *testing.T) {
	t.Log("Query with policy filter")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sd := mock_sciond.NewMockConnector(ctrl)
	pm := New(sd, Timers{})

	srcIA := xtest.MustParseIA("1-ff00:0:133")
	dstIA := xtest.MustParseIA("1-ff00:0:131")

	pathOne := fmt.Sprintf("%s#1019 1-ff00:0:132#1910 1-ff00:0:132#1916 %s#1619", srcIA, dstIA)
	paths := []string{pathOne}

	sd.EXPECT().Paths(gomock.Any(), dstIA, srcIA, gomock.Any(), gomock.Any()).Return(
		buildSDAnswer(paths...), nil,
	).AnyTimes()

	t.Run("Hop does not exist in paths, default deny", func(t *testing.T) {
		pp, err := pathpol.HopPredicateFromString("0-0#0")
		require.NoError(t, err)
		policy := &pathpol.Policy{ACL: &pathpol.ACL{Entries: []*pathpol.ACLEntry{
			{Action: pathpol.Allow, Rule: pp},
			denyEntry,
		}}}

		aps := pm.QueryFilter(context.Background(), srcIA, dstIA, policy)
		assert.Len(t, aps, 1, "only one path should remain")
		assert.ElementsMatch(t, getPathStrings(aps), paths)
	})

	t.Run("Hop does not exist paths, default allow", func(t *testing.T) {
		pp, err := pathpol.HopPredicateFromString("1-ff00:0:134#1910")
		require.NoError(t, err)
		policy := &pathpol.Policy{ACL: &pathpol.ACL{Entries: []*pathpol.ACLEntry{
			{Action: pathpol.Allow, Rule: pp},
			allowEntry,
		}}}
		aps := pm.QueryFilter(context.Background(), srcIA, dstIA, policy)
		assert.Len(t, aps, 1, "only one path should remain")
		assert.ElementsMatch(t, getPathStrings(aps), paths)
	})

	t.Run("Hop exists in paths, default deny", func(t *testing.T) {
		pp, err := pathpol.HopPredicateFromString("1-ff00:0:132#1910")
		require.NoError(t, err)
		policy := &pathpol.Policy{ACL: &pathpol.ACL{Entries: []*pathpol.ACLEntry{
			{Action: pathpol.Deny, Rule: pp},
			denyEntry,
		}}}
		aps := pm.QueryFilter(context.Background(), srcIA, dstIA, policy)
		assert.Len(t, aps, 0, "no path should remain")
	})
}

func TestACLPolicyFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sd := mock_sciond.NewMockConnector(ctrl)
	pm := New(sd, Timers{})

	srcIA, dstIA := xtest.MustParseIA("2-ff00:0:222"), xtest.MustParseIA("1-ff00:0:131")
	hop1, hop2 := "1-ff00:0:121", "2-ff00:0:211"

	paths := []string{
		fmt.Sprintf("%s#1019 1-ff00:0:122#1910 1-ff00:0:122#1916 %s#1619", srcIA, dstIA),
		fmt.Sprintf("%s#1019 %s#1912 %s#2328 %s#1619", srcIA, hop1, hop1, dstIA),
		fmt.Sprintf("%s#1019 %s#1911 %s#2327 %s#1619", srcIA, hop2, hop2, dstIA),
	}

	sd.EXPECT().Paths(gomock.Any(), dstIA, srcIA, gomock.Any(), gomock.Any()).Return(
		buildSDAnswer(paths...), nil,
	).AnyTimes()

	pp1, err := pathpol.HopPredicateFromString(hop1 + "#0")
	require.NoError(t, err)
	pp2, err := pathpol.HopPredicateFromString(hop2 + "#2327")
	require.NoError(t, err)

	t.Run("Query with ACL policy filter (hop1 denied)", func(t *testing.T) {
		policy := &pathpol.Policy{ACL: &pathpol.ACL{Entries: []*pathpol.ACLEntry{
			{Action: pathpol.Deny, Rule: pp1},
			allowEntry,
		}}}
		aps := pm.QueryFilter(context.Background(), srcIA, dstIA, policy)
		assert.Len(t, aps, 2)
		assert.NotContains(t, strings.Join(getPathStrings(aps), " "), hop1)
		assert.Contains(t, strings.Join(getPathStrings(aps), " "), hop2)
	})

	t.Run("Query with longer ACL policy filter (hop1,hop2 denied)", func(t *testing.T) {
		policy := &pathpol.Policy{ACL: &pathpol.ACL{Entries: []*pathpol.ACLEntry{
			{Action: pathpol.Deny, Rule: pp1},
			{Action: pathpol.Deny, Rule: pp2},
			allowEntry,
		}}}
		aps := pm.QueryFilter(context.Background(), srcIA, dstIA, policy)
		assert.Len(t, aps, 1)
		assert.NotContains(t, strings.Join(getPathStrings(aps), " "), hop1)
		assert.NotContains(t, strings.Join(getPathStrings(aps), " "), hop2)
	})
}

func TestWatchCount(t *testing.T) {
	t.Log("Given a path manager and adding a watch")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sd := mock_sciond.NewMockConnector(ctrl)
	pr := New(sd, Timers{})

	src := xtest.MustParseIA("1-ff00:0:111")
	dst := xtest.MustParseIA("1-ff00:0:110")

	sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(), gomock.Any()).Return(
		buildSDAnswer(), nil,
	).AnyTimes()

	assert.Equal(t, pr.WatchCount(), 0, " the count is initially 0")
	sp, err := pr.Watch(context.Background(), src, dst)
	require.NoError(t, err)
	assert.Equal(t, pr.WatchCount(), 1, "the number of watches increases to 1")
	sp.Destroy()
	assert.Equal(t, pr.WatchCount(), 0, "the number of watches decreases to 0")
}

func TestWatchPolling(t *testing.T) {
	t.Log("Given a path manager and adding a watch that retrieves zero paths")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sd := mock_sciond.NewMockConnector(ctrl)
	pr := New(sd, Timers{ErrorRefire: getDuration(1)})

	src := xtest.MustParseIA("1-ff00:0:111")
	dst := xtest.MustParseIA("1-ff00:0:110")
	gomock.InOrder(
		sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(), gomock.Any()).Return(
			buildSDAnswer(), nil,
		),
		sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(), gomock.Any()).Return(
			buildSDAnswer(
				"1-ff00:0:111#105 1-ff00:0:130#1002 1-ff00:0:130#1004 1-ff00:0:110#2",
			), nil,
		).MinTimes(1),
	)

	sp, err := pr.Watch(context.Background(), src, dst)
	require.NoError(t, err)
	assert.Len(t, sp.Load().APS, 0, "there are 0 paths currently available")
	time.Sleep(getDuration(4))
	assert.Len(t, sp.Load().APS, 1, "and after waiting, we get new paths")
}

func TestWatchFilter(t *testing.T) {
	t.Log("Given a path manager and adding a watch that should retrieve 1 path")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sd := mock_sciond.NewMockConnector(ctrl)
	pr := New(sd, Timers{ErrorRefire: getDuration(1)})

	src := xtest.MustParseIA("1-ff00:0:111")
	dst := xtest.MustParseIA("1-ff00:0:110")
	gomock.InOrder(
		sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(), gomock.Any()).Return(
			buildSDAnswer(
				"1-ff00:0:111#104 1-ff00:0:120#5 1-ff00:0:120#6 1-ff00:0:110#1",
			), nil,
		),
		sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(), gomock.Any()).Return(
			buildSDAnswer(
				"1-ff00:0:111#105 1-ff00:0:130#1002 1-ff00:0:130#1004 1-ff00:0:110#2",
				"1-ff00:0:111#104 1-ff00:0:120#5 1-ff00:0:120#6 1-ff00:0:110#1",
			), nil,
		).AnyTimes(),
	)

	seq, err := pathpol.NewSequence("1-ff00:0:111#105 0 0")
	require.NoError(t, err)
	filter := pathpol.NewPolicy("test-1-ff00:0:111#105", nil, seq, nil)

	sp, err := pr.WatchFilter(context.Background(), src, dst, filter)
	require.NoError(t, err)
	assert.Len(t, sp.Load().APS, 0, "there are 0 paths due to filtering")
	time.Sleep(getDuration(4))
	assert.Len(t, sp.Load().APS, 1, "and after waiting, we get 1 path that is not filtered")
}

func TestRevokeFastRecovery(t *testing.T) {
	t.Log("Given a path manager with a long normal timer and very small error timer")
	t.Log("A revocation that deletes everything triggers an immediate requery")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	src := xtest.MustParseIA("1-ff00:0:111")
	dst := xtest.MustParseIA("1-ff00:0:110")

	sd := mock_sciond.NewMockConnector(ctrl)
	pr := New(sd, Timers{NormalRefire: getDuration(100), ErrorRefire: getDuration(1)})

	sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(), gomock.Any()).Return(
		buildSDAnswer(
			"1-ff00:0:111#105 1-ff00:0:130#1002 1-ff00:0:130#1004 1-ff00:0:110#2",
		), nil,
	)

	_, err := pr.Watch(context.Background(), src, dst)
	require.NoError(t, err)

	// Once everything is revoked a fast request is immediately
	// triggered. We check for at least 2 iterations to make sure we
	// are in error recovery mode, and the aggressive timer is used.
	// We actually test that the mock .{Revnotifications,Paths} functions are
	// being called within a 5 time units. It will fail with "missing
	// call(s)" error message
	gomock.InOrder(
		sd.EXPECT().RevNotification(gomock.Any(), gomock.Any()).Return(
			&sciond.RevReply{Result: sciond.RevValid}, nil,
		),
		sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(),
			gomock.Any()).Return(
			buildSDAnswer(), nil,
		).MinTimes(2),
	)
	pr.Revoke(context.Background(), newTestRev(t, "1-ff00:0:130#1002"))
	time.Sleep(getDuration(5))
}

func TestRevoke(t *testing.T) {
	t.Log("Given a path manager and a watch that")

	src := xtest.MustParseIA("1-ff00:0:111")
	dst := xtest.MustParseIA("1-ff00:0:110")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	paths := []string{
		"1-ff00:0:111#105 1-ff00:0:130#1002 1-ff00:0:130#1004 1-ff00:0:110#2",
		"1-ff00:0:111#104 1-ff00:0:120#5 1-ff00:0:120#6 1-ff00:0:110#1",
	}

	tests := map[string]struct {
		Paths         []string
		RevReply      *sciond.RevReply
		RevReplyError error
		Revocation    *path_mgmt.SignedRevInfo
		Remaining     int
	}{
		"retrieves one path, revokes an IFID that matches the path": {
			Paths:      paths[:1],
			RevReply:   &sciond.RevReply{Result: sciond.RevValid},
			Revocation: newTestRev(t, "1-ff00:0:130#1002"),
			Remaining:  0,
		},

		"retrieves one path, revokes an IFID that does not match the path": {
			Paths:      paths[:1],
			RevReply:   &sciond.RevReply{Result: sciond.RevValid},
			Revocation: newTestRev(t, "2-ff00:0:1#1"),
			Remaining:  1,
		},
		"tries to revoke an IFID, but SCIOND encounters an error": {
			Paths:         paths[:1],
			RevReplyError: errors.New("some error"),
			Revocation:    newTestRev(t, "1-ff00:0:130#1002"),
			Remaining:     1,
		},
		"tries to revoke an IFID, but the revocation is invalid": {
			Paths:      paths[:1],
			RevReply:   &sciond.RevReply{Result: sciond.RevInvalid},
			Revocation: newTestRev(t, "1-ff00:0:130#1002"),
			Remaining:  1,
		},
		"tries to revoke an IFID, but the revocation is stale": {
			Paths:      paths[:1],
			RevReply:   &sciond.RevReply{Result: sciond.RevStale},
			Revocation: newTestRev(t, "1-ff00:0:130#1002"),
			Remaining:  1,
		},
		"tries to revoke an IFID, but the revocation is unknown": {
			Paths:      paths[:1],
			RevReply:   &sciond.RevReply{Result: sciond.RevUnknown},
			Revocation: newTestRev(t, "1-ff00:0:130#1002"),
			Remaining:  0,
		},
		"retrieves two paths, revokes an IFID that matches one path": {
			Paths:      paths,
			RevReply:   &sciond.RevReply{Result: sciond.RevValid},
			Revocation: newTestRev(t, "1-ff00:0:130#1002"),
			Remaining:  1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sd := mock_sciond.NewMockConnector(ctrl)
			pr := New(sd, Timers{})

			sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(), gomock.Any()).Return(
				buildSDAnswer(test.Paths...), nil,
			)
			sd.EXPECT().Paths(gomock.Any(), dst, src, gomock.Any(), gomock.Any()).Return(
				buildSDAnswer(), nil,
			).AnyTimes()
			sp, err := pr.Watch(context.Background(), src, dst)
			require.NoError(t, err)
			sd.EXPECT().RevNotification(gomock.Any(), gomock.Any()).Return(
				test.RevReply, test.RevReplyError,
			)
			pr.Revoke(context.Background(), test.Revocation)
			assert.Len(t, sp.Load().APS, test.Remaining)
		})
	}

}

func newTestRev(t *testing.T, rev string) *path_mgmt.SignedRevInfo {
	pi := mustParsePI(rev)
	signedRevInfo, err := path_mgmt.NewSignedRevInfo(
		&path_mgmt.RevInfo{
			IfID:     pi.IfID,
			RawIsdas: pi.RawIsdas,
		}, infra.NullSigner)
	require.NoError(t, err)
	return signedRevInfo
}

func getPathStrings(aps spathmeta.AppPathSet) (ss []string) {
	for _, v := range aps {
		ss = append(ss, strings.Trim(fmt.Sprintf("%v", v.Entry.Path.Interfaces), "[]"))
	}
	return
}
