package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestJSONResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the JSON helper
	JSON(c, 200, true, "ok", map[string]string{"hello": "world"}, "")

	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if resp["success"] != true {
		t.Fatalf("expected success=true, got %v", resp["success"])
	}
	if resp["message"] != "ok" {
		t.Fatalf("expected message=ok, got %v", resp["message"])
	}
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data to be an object, got %T", resp["data"])
	}
	if data["hello"] != "world" {
		t.Fatalf("expected data.hello=world, got %v", data["hello"])
	}
}
