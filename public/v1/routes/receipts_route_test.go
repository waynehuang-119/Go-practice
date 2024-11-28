package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type ReceiptsRouteTestSuite struct {
	suite.Suite
	router *mux.Router
}

func (suite *ReceiptsRouteTestSuite) SetupSuite() {
	suite.router = mux.NewRouter()
	Register(suite.router)
}

func (suite *ReceiptsRouteTestSuite) TestValidRoutes() {
	testCases := []struct {
		name       string
		method     string
		url        string
		body       string
		wantStatus int
	}{
		{
			name:       "Process Receipt",
			method:     http.MethodPost,
			url:        "/receipts/process",
			body:       `{"retailer":"Target","purchaseDate":"2022-01-01","purchaseTime":"13:01","items":[{"shortDescription":"Item 1","price":"5.00"}],"total":"5.00"}`,
			wantStatus: http.StatusOK, // Expecting success
		},
		{
			name:       "Get Points",
			method:     http.MethodGet,
			url:        "/receipts/12345/points",
			body:       "",
			wantStatus: http.StatusOK, // Expecting success
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.url, nil)
			}

			rec := httptest.NewRecorder()

			suite.router.ServeHTTP(rec, req)

			suite.Equal(tc.wantStatus, rec.Code, "Unexpected status code for %s", tc.url)
		})
	}
}

func (suite *ReceiptsRouteTestSuite) TestInvalidRoute() {
	req := httptest.NewRequest(http.MethodGet, "/invalid/url", nil)
	rec := httptest.NewRecorder()

	suite.router.ServeHTTP(rec, req)

	suite.Equal(http.StatusNotFound, rec.Code, "Invalid route should return 404")
}

// Run the test suite
func TestReceiptsRouteTestSuite(t *testing.T) {
	suite.Run(t, new(ReceiptsRouteTestSuite))
}
