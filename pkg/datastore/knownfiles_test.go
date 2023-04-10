package datastore_test

import (
	"io"
	"math/rand"
	"testing"

	"github.com/overlordtm/pmss/internal/testutils"
	"github.com/overlordtm/pmss/internal/utils"
	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/multihasher"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func LoadData(db *gorm.DB) {

	fixture := []datastore.KnownFile{
		{
			MD5:      utils.StringPtr("3e7705498e8be60520841409ebc69bc1"),
			SHA1:     utils.StringPtr("dba7673010f19a94af4345453005933fd511bea9"),
			SHA256:   utils.StringPtr("634b027b1b69e1242d40d53e312b3b4ac7710f55be81f289b549446ef6778bee"),
			Size:     utils.Int64Ptr(100),
			MimeType: utils.StringPtr("application/x-msdownload"),
			Status:   datastore.FileStatusMalicious,
		},
		{
			MD5:      utils.StringPtr("126a8a51b9d1bbd07fddc65819a542c3"),
			SHA1:     utils.StringPtr("9054fbe0b622c638224d50d20824d2ff6782e308"),
			SHA256:   utils.StringPtr("7d6fd7774f0d87624da6dcf16d0d3d104c3191e771fbe2f39c86aed4b2bf1a0f"),
			Size:     utils.Int64Ptr(100),
			MimeType: utils.StringPtr("application/x-msdownload"),
			Status:   datastore.FileStatusSafe,
		},
	}

	rnd := rand.New(rand.NewSource(42))

	for i := 0; i < 100; i++ {
		mHash, err := multihasher.Hash(io.LimitReader(rnd, 100))
		if err != nil {
			panic(err)
		}

		fixture = append(fixture, datastore.KnownFile{
			MD5:      utils.StringPtr(mHash.MD5),
			SHA1:     utils.StringPtr(mHash.SHA1),
			SHA256:   utils.StringPtr(mHash.SHA256),
			Size:     utils.Int64Ptr(100),
			MimeType: utils.StringPtr("application/x-msdownload"),
			Status:   datastore.FileStatusMalicious,
		})
	}

	testutils.MustExecute(datastore.KnownFiles().CreateInBatches(fixture), db)
}

func TestKnonFilesFindByScannedFile(t *testing.T) {
	db := datastore.MustOpen(testDbUrl).Begin()
	defer db.Rollback()
	datastore.AutoMigrate(db)
	LoadData(db)

	testCases := []struct {
		name           string
		scannedFile    *datastore.ScannedFile
		expectedErr    error
		expectedStatus datastore.FileStatus
	}{
		{
			name: "no match",
			scannedFile: &datastore.ScannedFile{
				MD5:    utils.StringPtr("b3af15b5431600e4b9ef45dc0dc41b92"),
				SHA1:   utils.StringPtr("d71afce9d88c842cb8638dab74a2f49e52a40c8d"),
				SHA256: utils.StringPtr("fef5e99fc83f08d7a0572fd2e8df48127e955cfef2059e5ece5aea6ab2482492"),
			},
			expectedErr: gorm.ErrRecordNotFound,
		},
		{
			name: "match by md5, file is malicious",
			scannedFile: &datastore.ScannedFile{
				MD5: utils.StringPtr("3e7705498e8be60520841409ebc69bc1"),
			},
			expectedStatus: datastore.FileStatusMalicious,
		},
		{
			name: "match by md5 and sha1, file is malicious",
			scannedFile: &datastore.ScannedFile{
				MD5:    utils.StringPtr("3e7705498e8be60520841409ebc69bc1"),
				SHA256: utils.StringPtr("634b027b1b69e1242d40d53e312b3b4ac7710f55be81f289b549446ef6778bee"),
			},
			expectedStatus: datastore.FileStatusMalicious,
		},
		{
			name: "match by md5 and sha256, md5 is safe, sha256 is malicious, file is malicious",
			scannedFile: &datastore.ScannedFile{
				MD5:    utils.StringPtr("126a8a51b9d1bbd07fddc65819a542c3"),
				SHA256: utils.StringPtr("634b027b1b69e1242d40d53e312b3b4ac7710f55be81f289b549446ef6778bee"),
			},
			expectedStatus: datastore.FileStatusMalicious,
		},
		{
			name: "match by sha1, file is safe",
			scannedFile: &datastore.ScannedFile{
				SHA1: utils.StringPtr("9054fbe0b622c638224d50d20824d2ff6782e308"),
			},
			expectedStatus: datastore.FileStatusSafe,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			knownFile := &datastore.KnownFile{}
			err := datastore.KnownFiles().FindByScannedFile(testCase.scannedFile, knownFile)(db)

			if testCase.expectedErr != nil {
				assert.ErrorIs(t, err, testCase.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedStatus, knownFile.Status)
			}
		})
	}

}
