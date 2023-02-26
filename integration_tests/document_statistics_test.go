package integrationtest

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
	"tryffel.net/go/virtualpaper/api"
)

type DocumentStatisticsSuite struct {
	ApiTestSuite
	docs map[string]*api.DocumentResponse
}

func TestDocumentStatistics(t *testing.T) {
	suite.Run(t, new(DocumentStatisticsSuite))
}

func (suite *DocumentStatisticsSuite) SetupTest() {
	suite.Init()
	clearDbMetadataTables(suite.T())
	clearDbDocumentTables(suite.T())

	testDocumentX86.Date = time.Date(2018, 06, 02, 0, 0, 0, 0, time.UTC)
	testDocumentX86Intel.Date = time.Date(2018, 06, 03, 0, 0, 0, 0, time.UTC)
	testDocumentJupiterMoons.Date = time.Date(2010, 06, 03, 12, 0, 0, 0, time.UTC)
	testDocumentYear1962.Date = time.Date(2010, 06, 03, 11, 59, 0, 0, time.UTC)
	testDocumentMetamorphosis.Date = time.Date(2011, 06, 03, 11, 59, 0, 0, time.UTC)
	testDocumentTransistorCount.Date = time.Date(2012, 06, 03, 11, 59, 0, 0, time.UTC)
	testDocumentTransistorCountAdminUser.Date = time.Date(2008, 06, 03, 11, 59, 0, 0, time.UTC)
	insertTestDocuments(suite.T())
}

func (suite *DocumentStatisticsSuite) TestBasicStats() {
	data := getDocumentStatistics(suite.T(), suite.userHttp, 200)

	assert.Equal(suite.T(), data.UserId, 1)
	assert.Equal(suite.T(), data.Indexing, false)
	assert.Equal(suite.T(), data.NumDocuments, 6)
	assert.Equal(suite.T(), data.NumMetadataKeys, 0)
	assert.Equal(suite.T(), data.NumMetadataValues, 0)

	assert.Len(suite.T(), data.LastDocumentsAdded, 6)
	assert.Equal(suite.T(), data.LastDocumentsAdded, testDocumentIdsUser, "last added")
	assert.Len(suite.T(), data.LastDocumentsUpdated, 6)
	assert.Equal(suite.T(), data.LastDocumentsUpdated, testDocumentIdsUser, "last updated")
	assert.Len(suite.T(), data.LastDocumentsViewed, 0)

	adminData := getDocumentStatistics(suite.T(), suite.adminHttp, 200)

	assert.Equal(suite.T(), adminData.UserId, 2)
	assert.Equal(suite.T(), adminData.Indexing, false)
	assert.Equal(suite.T(), adminData.NumDocuments, 1)
	assert.Equal(suite.T(), adminData.NumMetadataKeys, 0)
	assert.Equal(suite.T(), adminData.NumMetadataValues, 0)

	assert.Len(suite.T(), adminData.LastDocumentsAdded, 1)
	assert.Equal(suite.T(), adminData.LastDocumentsAdded, testDocumentIdsAdmin, "last added")
	assert.Len(suite.T(), adminData.LastDocumentsUpdated, 1)
	assert.Equal(suite.T(), adminData.LastDocumentsUpdated, testDocumentIdsAdmin, "last viewed")
	assert.Len(suite.T(), adminData.LastDocumentsViewed, 0)
}

func (suite *DocumentStatisticsSuite) TestUpdateDocument() {
	// update documents and check that their index in the list has moved up

	data := getDocumentStatistics(suite.T(), suite.userHttp, 200)
	adminData := getDocumentStatistics(suite.T(), suite.adminHttp, 200)
	assert.Len(suite.T(), adminData.LastDocumentsAdded, 1)
	assert.Len(suite.T(), adminData.LastDocumentsUpdated, 1)

	assert.Len(suite.T(), data.LastDocumentsAdded, 6)
	assert.Equal(suite.T(), data.LastDocumentsAdded, testDocumentIdsUser)
	assert.Len(suite.T(), data.LastDocumentsUpdated, 6)
	assert.Equal(suite.T(), data.LastDocumentsUpdated, testDocumentIdsUser)
	assert.Len(suite.T(), data.LastDocumentsViewed, 0)

	jupiter := getDocument(suite.T(), suite.userHttp, testDocumentJupiterMoons.Id, 200)
	jupiter.Description = "changed"
	updateDocument(suite.T(), suite.userHttp, jupiter, 200)

	intel := getDocument(suite.T(), suite.userHttp, testDocumentX86Intel.Id, 200)
	intel.Description = "changed"
	updateDocument(suite.T(), suite.userHttp, intel, 200)

	data = getDocumentStatistics(suite.T(), suite.userHttp, 200)
	adminData = getDocumentStatistics(suite.T(), suite.adminHttp, 200)
	assert.Len(suite.T(), adminData.LastDocumentsAdded, 1)
	assert.Len(suite.T(), adminData.LastDocumentsUpdated, 1)

	assert.Len(suite.T(), data.LastDocumentsAdded, 6)
	assert.Equal(suite.T(), data.LastDocumentsAdded, testDocumentIdsUser)
	assert.Len(suite.T(), data.LastDocumentsUpdated, 6)
	assert.Len(suite.T(), data.LastDocumentsViewed, 0)

	assert.Equal(suite.T(), data.LastDocumentsUpdated, []string{
		testDocumentX86Intel.Id,
		testDocumentJupiterMoons.Id,
		testDocumentTransistorCount.Id,
		testDocumentYear1962.Id,
		testDocumentMetamorphosis.Id,
		testDocumentX86.Id,
	})

	count := getDocument(suite.T(), suite.adminHttp, testDocumentTransistorCountAdminUser.Id, 200)
	count.Description = "changed"

	updateDocument(suite.T(), suite.adminHttp, count, 200)

	unChangedData := getDocumentStatistics(suite.T(), suite.userHttp, 200)
	// admin's data shouldn't be changed since there's only one document in a list
	unChangedAdminData := getDocumentStatistics(suite.T(), suite.adminHttp, 200)

	assert.Equal(suite.T(), data, unChangedData, "user's data not changed")
	assert.Equal(suite.T(), adminData, unChangedAdminData, "admin's data not changed")
}

func (suite *DocumentStatisticsSuite) TestViewDocument() {
	data := getDocumentStatistics(suite.T(), suite.userHttp, 200)
	assert.Len(suite.T(), data.LastDocumentsViewed, 0, "no views in the beginning")

	getDocument(suite.T(), suite.userHttp, testDocumentX86.Id, 200)

	data = getDocumentStatistics(suite.T(), suite.userHttp, 200)
	assert.Len(suite.T(), data.LastDocumentsViewed, 0, "no views without ?visit=1 param")

	// another user cannot add visit counter
	getDocumentWithVisit(suite.T(), suite.adminHttp, testDocumentX86.Id, 404)
	data = getDocumentStatistics(suite.T(), suite.userHttp, 200)
	assert.Len(suite.T(), data.LastDocumentsViewed, 0, "no views if another user tried access")
	adminData := getDocumentStatistics(suite.T(), suite.adminHttp, 200)
	assert.Len(suite.T(), adminData.LastDocumentsViewed, 0, "no views when accessing another user's document")

	getDocumentWithVisit(suite.T(), suite.userHttp, testDocumentX86.Id, 200)
	data = getDocumentStatistics(suite.T(), suite.userHttp, 200)
	assert.Len(suite.T(), data.LastDocumentsViewed, 1, "adds visit with ?visit=1 param")

}

func (suite *DocumentStatisticsSuite) TestAddMetdataKey() {}

func (suite *DocumentStatisticsSuite) TestAddMetdataValue() {}

func getDocumentStatistics(t *testing.T, client *httpClient, wantHttpStatus int) *api.UserDocumentStatistics {
	data := api.UserDocumentStatistics{}
	client.Get("/api/v1/documents/stats").ExpectName(t, "get document statistics", false).Json(t, &data).e.Status(wantHttpStatus).Done()
	return &data
}
