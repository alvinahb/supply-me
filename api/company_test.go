package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/alvinahb/supply-me/db/mock"
	db "github.com/alvinahb/supply-me/db/sqlc"
	"github.com/alvinahb/supply-me/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetCompanyAPI(t *testing.T) {
	company := randomCompany()

	testCases := []struct {
		name          string
		companyID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(company, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchCompany(t, recorder.Body, company)
			},
		},
		{
			name:      "NotFound",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(db.Company{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			companyID: company.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Eq(company.ID)).
					Times(1).
					Return(db.Company{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "InvalidID",
			companyID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetCompany(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// Start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/company/%d", tc.companyID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomCompany() db.Company {
	return db.Company{
		ID:          util.RandomInt(1, 1000),
		CompanyType: util.RandomCompanyType(),
		CompanyName: util.RandomString(20),
	}
}

func requiredBodyMatchCompany(t *testing.T, body *bytes.Buffer, company db.Company) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotCompany db.Company
	err = json.Unmarshal(data, &gotCompany)
	require.NoError(t, err)
	require.Equal(t, company, gotCompany)
}
