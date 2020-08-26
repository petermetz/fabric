/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tests

import (
	"fmt"
	"testing"

	"github.com/petermetz/fabric/core/ledger"
	"github.com/petermetz/fabric/core/ledger/kvledger"
	"github.com/stretchr/testify/assert"

	"github.com/petermetz/fabric/common/ledger/testutil"
	"github.com/petermetz/fabric/protos/common"
	protopeer "github.com/petermetz/fabric/protos/peer"
)

// Test data used in the tests in this file was generated by v1.1 code https://gerrit.hyperledger.org/r/#/c/22749/6/core/ledger/kvledger/tests/v11_generate_test.go@22
// Folder, "testdata/v11/sample_ledgers" contains the data that was generated before commit hash feature was added.
// Folder, "testdata/v11/sample_ledgers_with_commit_hashes" contains the data that was generated after commit hash feature was added.

// TestV11 tests that a ledgersData folder created by v1.1 can be used with future releases in a backward compatible way
func TestV11(t *testing.T) {
	fsPath := defaultConfig["peer.fileSystemPath"].(string)
	testutil.CopyDir("testdata/v11/sample_ledgers/ledgersData", fsPath)
	env := newEnv(defaultConfig, t)
	defer env.cleanup()

	h1, h2 := newTestHelperOpenLgr("ledger1", t), newTestHelperOpenLgr("ledger2", t)
	dataHelper := &v11SampleDataHelper{}

	dataHelper.verifyBeforeStateRebuild(h1)
	dataHelper.verifyBeforeStateRebuild(h2)

	env.closeAllLedgersAndDrop(rebuildableStatedb)
	h1, h2 = newTestHelperOpenLgr("ledger1", t), newTestHelperOpenLgr("ledger2", t)
	dataHelper.verifyAfterStateRebuild(h1)
	dataHelper.verifyAfterStateRebuild(h2)

	env.closeAllLedgersAndDrop(rebuildableStatedb + rebuildableBlockIndex + rebuildableConfigHistory)
	h1, h2 = newTestHelperOpenLgr("ledger1", t), newTestHelperOpenLgr("ledger2", t)
	dataHelper.verifyAfterStateRebuild(h1)
	dataHelper.verifyAfterStateRebuild(h2)

	h1.verifyCommitHashNotExists()
	h2.verifyCommitHashNotExists()
	h1.simulateDataTx("txid1_with_new_binary", func(s *simulator) {
		s.setState("cc1", "new_key", "new_value")
	})

	// add a new block and the new block should not contain a commit hash
	// because the previously committed block from 1.1 code did not contain commit hash
	h1.cutBlockAndCommitWithPvtdata()
	h1.verifyCommitHashNotExists()
}

func TestV11CommitHashes(t *testing.T) {
	testCases := []struct {
		description               string
		v11SampleDataPath         string
		preResetCommitHashExists  bool
		resetFunc                 func(h *testhelper)
		postResetCommitHashExists bool
	}{
		{
			"Reset (no existing CommitHash)",
			"testdata/v11/sample_ledgers/ledgersData",
			false,
			func(h *testhelper) {
				assert.NoError(t, kvledger.ResetAllKVLedgers())
			},
			true,
		},

		{
			"Rollback to genesis block (no existing CommitHash)",
			"testdata/v11/sample_ledgers/ledgersData",
			false,
			func(h *testhelper) {
				assert.NoError(t, kvledger.RollbackKVLedger(h.lgrid, 0))
			},
			true,
		},

		{
			"Rollback to block other than genesis block (no existing CommitHash)",
			"testdata/v11/sample_ledgers/ledgersData",
			false,
			func(h *testhelper) {
				assert.NoError(t, kvledger.RollbackKVLedger(h.lgrid, h.currentHeight()/2+1))
			},
			false,
		},

		{
			"Reset (existing CommitHash)",
			"testdata/v11/sample_ledgers_with_commit_hashes/ledgersData",
			true,
			func(h *testhelper) {
				assert.NoError(t, kvledger.ResetAllKVLedgers())
			},
			true,
		},

		{
			"Rollback to genesis block (existing CommitHash)",
			"testdata/v11/sample_ledgers_with_commit_hashes/ledgersData",
			true,
			func(h *testhelper) {
				assert.NoError(t, kvledger.RollbackKVLedger(h.lgrid, 0))
			},
			true,
		},

		{
			"Rollback to block other than genesis block (existing CommitHash)",
			"testdata/v11/sample_ledgers_with_commit_hashes/ledgersData",
			true,
			func(h *testhelper) {
				assert.NoError(t, kvledger.RollbackKVLedger(h.lgrid, h.currentHeight()/2+1))
			},
			true,
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.description,
			func(t *testing.T) {
				testV11CommitHashes(
					t,
					testCase.v11SampleDataPath,
					testCase.preResetCommitHashExists,
					testCase.resetFunc,
					testCase.postResetCommitHashExists,
				)
			})
	}
}

func testV11CommitHashes(t *testing.T,
	v11DataPath string,
	preResetCommitHashExists bool,
	resetFunc func(*testhelper),
	postResetCommitHashExists bool,
) {
	fsPath := defaultConfig["peer.fileSystemPath"].(string)
	testutil.CopyDir(v11DataPath, fsPath)
	env := newEnv(defaultConfig, t)
	defer env.cleanup()

	h := newTestHelperOpenLgr("ledger1", t)
	blocksAndPvtData := h.retrieveCommittedBlocksAndPvtdata(0, h.currentHeight()-1)
	if preResetCommitHashExists {
		h.verifyCommitHashExists()
	} else {
		h.verifyCommitHashNotExists()
	}

	closeLedgerMgmt()
	resetFunc(h)
	initLedgerMgmt()

	h = newTestHelperOpenLgr("ledger1", t)
	for i := int(h.currentHeight()); i < len(blocksAndPvtData); i++ {
		d := blocksAndPvtData[i]
		// add metadata slot for commit hash, as this would have be missing in the blocks from 1.1 prior to this feature
		for len(d.Block.Metadata.Metadata) < int(common.BlockMetadataIndex_COMMIT_HASH)+1 {
			d.Block.Metadata.Metadata = append(d.Block.Metadata.Metadata, []byte{})
		}
		// set previous block hash, as this is not present in the test blocks from 1.1
		d.Block.Header.PreviousHash = blocksAndPvtData[i-1].Block.Header.Hash()
		assert.NoError(t, h.lgr.CommitWithPvtData(d, &ledger.CommitOptions{FetchPvtDataFromLedger: true}))
	}

	if postResetCommitHashExists {
		h.verifyCommitHashExists()
	} else {
		h.verifyCommitHashNotExists()
	}

	h.closeAndReopenLgr()
	h.simulateDataTx("txid1_with_new_binary", func(s *simulator) {
		s.setState("cc1", "new_key", "new_value")
	})
	h.cutBlockAndCommitWithPvtdata()

	if postResetCommitHashExists {
		h.verifyCommitHashExists()
	} else {
		h.verifyCommitHashNotExists()
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// v11SampleDataHelper provides a set of functions to verify the ledger. This verifies the ledger under the assumption that the
// ledger is generated by this code from v1.1 (https://gerrit.hyperledger.org/r/#/c/22749/1/core/ledger/kvledger/tests/v11_generate_test.go@22).
// In summary, the above generate function, constructs two ledgers and populates the ledgers uses this code
// (https://gerrit.hyperledger.org/r/#/c/22749/1/core/ledger/kvledger/tests/util_sample_data.go@55)
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type v11SampleDataHelper struct {
}

func (d *v11SampleDataHelper) verifyBeforeStateRebuild(h *testhelper) {
	dataHelper := &v11SampleDataHelper{}
	dataHelper.verifyState(h)
	dataHelper.verifyBlockAndPvtdata(h)
	dataHelper.verifyGetTransactionByID(h)
	dataHelper.verifyConfigHistoryDoesNotExist(h)
}

func (d *v11SampleDataHelper) verifyAfterStateRebuild(h *testhelper) {
	dataHelper := &v11SampleDataHelper{}
	dataHelper.verifyState(h)
	dataHelper.verifyBlockAndPvtdata(h)
	dataHelper.verifyGetTransactionByID(h)
	dataHelper.verifyConfigHistory(h)
}

func (d *v11SampleDataHelper) verifyState(h *testhelper) {
	lgrid := h.lgrid
	h.verifyPubState("cc1", "key1", d.sampleVal("value13", lgrid))
	h.verifyPubState("cc1", "key2", "")
	h.verifyPvtState("cc1", "coll1", "key3", d.sampleVal("value14", lgrid))
	h.verifyPvtState("cc1", "coll1", "key4", "")
	h.verifyPvtState("cc1", "coll2", "key3", d.sampleVal("value09", lgrid))
	h.verifyPvtState("cc1", "coll2", "key4", d.sampleVal("value10", lgrid))

	h.verifyPubState("cc2", "key1", d.sampleVal("value03", lgrid))
	h.verifyPubState("cc2", "key2", d.sampleVal("value04", lgrid))
	h.verifyPvtState("cc2", "coll1", "key3", d.sampleVal("value07", lgrid))
	h.verifyPvtState("cc2", "coll1", "key4", d.sampleVal("value08", lgrid))
	h.verifyPvtState("cc2", "coll2", "key3", d.sampleVal("value11", lgrid))
	h.verifyPvtState("cc2", "coll2", "key4", d.sampleVal("value12", lgrid))
}

func (d *v11SampleDataHelper) verifyConfigHistory(h *testhelper) {
	lgrid := h.lgrid
	h.verifyMostRecentCollectionConfigBelow(10, "cc1",
		&expectedCollConfInfo{5, d.sampleCollConf2(lgrid, "cc1")})

	h.verifyMostRecentCollectionConfigBelow(5, "cc1",
		&expectedCollConfInfo{3, d.sampleCollConf1(lgrid, "cc1")})

	h.verifyMostRecentCollectionConfigBelow(10, "cc2",
		&expectedCollConfInfo{5, d.sampleCollConf2(lgrid, "cc2")})

	h.verifyMostRecentCollectionConfigBelow(5, "cc2",
		&expectedCollConfInfo{3, d.sampleCollConf1(lgrid, "cc2")})
}

func (d *v11SampleDataHelper) verifyConfigHistoryDoesNotExist(h *testhelper) {
	h.verifyMostRecentCollectionConfigBelow(10, "cc1", nil)
	h.verifyMostRecentCollectionConfigBelow(10, "cc2", nil)
}

func (d *v11SampleDataHelper) verifyBlockAndPvtdata(h *testhelper) {
	lgrid := h.lgrid
	h.verifyBlockAndPvtData(2, nil, func(r *retrievedBlockAndPvtdata) {
		r.hasNumTx(2)
		r.hasNoPvtdata()
	})

	h.verifyBlockAndPvtData(4, nil, func(r *retrievedBlockAndPvtdata) {
		r.hasNumTx(2)
		r.pvtdataShouldContain(0, "cc1", "coll1", "key3", d.sampleVal("value05", lgrid))
		r.pvtdataShouldContain(1, "cc2", "coll1", "key3", d.sampleVal("value07", lgrid))
	})
}

func (d *v11SampleDataHelper) verifyGetTransactionByID(h *testhelper) {
	h.verifyTxValidationCode("txid7", protopeer.TxValidationCode_VALID)
	h.verifyTxValidationCode("txid8", protopeer.TxValidationCode_MVCC_READ_CONFLICT)
}

func (d *v11SampleDataHelper) sampleVal(val, ledgerid string) string {
	return fmt.Sprintf("%s:%s", val, ledgerid)
}

func (d *v11SampleDataHelper) sampleCollConf1(ledgerid, ccName string) []*collConf {
	return []*collConf{
		{name: "coll1", members: []string{"org1", "org2"}},
		{name: ledgerid, members: []string{"org1", "org2"}},
		{name: ccName, members: []string{"org1", "org2"}},
	}
}

func (d *v11SampleDataHelper) sampleCollConf2(ledgerid string, ccName string) []*collConf {
	return []*collConf{
		{name: "coll1", members: []string{"org1", "org2"}},
		{name: "coll2", members: []string{"org1", "org2"}},
		{name: ledgerid, members: []string{"org1", "org2"}},
		{name: ccName, members: []string{"org1", "org2"}},
	}
}
