package helpers

import (
	"bytes"
	"fmt"
	cfg "github.com/JMURv/service-discovery/pkg/config"
	"github.com/goccy/go-json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var Conf *cfg.Config

func init() {
	Conf = cfg.MustLoad("../local.config.yaml")
}

func CleanDB(t *testing.T) {
	conn, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"postgres://%s:%s@%s:%v/%s",
				Conf.DB.User,
				Conf.DB.Password,
				Conf.DB.Host,
				Conf.DB.Port,
				Conf.DB.Database+"_test",
			),
		), &gorm.Config{},
	)
	if err != nil {
		t.Log(err)
	}

	sqlDB, err := conn.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	var tables []string
	if err := conn.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error; err != nil {
		t.Fatal(err)
	}

	for _, table := range tables {
		if err := conn.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Fatal(err)
		}
	}

	t.Log("Database cleaned")
}

func SendHttpRequest(t *testing.T, router *mux.Router, access string, method string, url string, body any) (map[string]any, *httptest.ResponseRecorder) {
	var err error
	var req *http.Request
	if body != nil {
		req, err = http.NewRequest(method, url, body.(*bytes.Buffer))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if access != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", access))
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return UnmarshallResponse(t, rr.Body), rr
}

func UnmarshallResponse(t *testing.T, rrBody *bytes.Buffer) (r map[string]any) {
	data, err := io.ReadAll(rrBody)
	if !assert.NoError(t, err) {
		t.Log(fmt.Sprintf("Error while reading: %v", err))
		t.Fatal(err)
	}

	if err = json.Unmarshal(data, &r); !assert.NoError(t, err) {
		t.Log(fmt.Sprintf("Error while unmarshalling: %v", err))
		t.Fatal(err)
	}
	return r
}
