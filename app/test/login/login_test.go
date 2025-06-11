package login

import (
	"encoding/json"
	"govote/app/logic"
	"govote/app/model"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// 测试响应结构
type LoginResponse struct {
	Message string `json:"message"`
}

func TestLoginSuccess(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.POST("/login", logic.DoLogin)

	// 测试用例
	testCases := []struct {
		name     string
		password string
		expected string
	}{
		{"admin", "123456", "登录成功"},
		{"test", "test123", "登录成功"},
		{"user1", "password1", "登录成功"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 准备表单数据
			formData := url.Values{}
			formData.Set("name", tc.name)
			formData.Set("password", tc.password)

			// 创建请求
			req, err := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
			if err != nil {
				t.Fatalf("创建请求失败: %s", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 检查状态码
			if w.Code != http.StatusOK {
				t.Errorf("期望状态码 %d, 实际得到 %d", http.StatusOK, w.Code)
			}

			// 读取响应
			body, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("读取响应失败: %s", err)
			}

			// 解析JSON响应
			var response LoginResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("解析JSON失败: %s", err)
			}

			// 检查响应消息
			if response.Message != tc.expected {
				t.Errorf("期望消息 '%s', 实际得到 '%s'", tc.expected, response.Message)
			}

			t.Logf("用户 %s 登录成功: %s", tc.name, response.Message)
		})
	}
}

func TestLoginFailure(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.POST("/login", logic.DoLogin)

	// 测试用例
	testCases := []struct {
		name     string
		password string
		expected string
	}{
		{"admin", "wrong", "账号或者密码有误"},
		{"nonexistent", "123456", "账号或者密码有误"},
		{"test", "wrong", "账号或者密码有误"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 准备表单数据
			formData := url.Values{}
			formData.Set("name", tc.name)
			formData.Set("password", tc.password)

			// 创建请求
			req, err := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
			if err != nil {
				t.Fatalf("创建请求失败: %s", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 检查状态码
			if w.Code != http.StatusOK {
				t.Errorf("期望状态码 %d, 实际得到 %d", http.StatusOK, w.Code)
			}

			// 读取响应
			body, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("读取响应失败: %s", err)
			}

			// 解析JSON响应
			var response LoginResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("解析JSON失败: %s", err)
			}

			// 检查响应消息
			if response.Message != tc.expected {
				t.Errorf("期望消息 '%s', 实际得到 '%s'", tc.expected, response.Message)
			}

			t.Logf("用户 %s 登录失败: %s", tc.name, response.Message)
		})
	}
}

func TestLoginInvalidData(t *testing.T) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.POST("/login", logic.DoLogin)

	// 测试用例
	testCases := []struct {
		name        string
		contentType string
		body        string
		expected    string
	}{
		{"空用户名", "application/x-www-form-urlencoded", "name=&password=123456", "参数绑定失败"},
		{"空密码", "application/x-www-form-urlencoded", "name=admin&password=", "参数绑定失败"},
		{"无效JSON", "application/json", `{"invalid": "json"}`, "参数绑定失败"},
		{"空请求体", "application/x-www-form-urlencoded", "", "参数绑定失败"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建请求
			req, err := http.NewRequest("POST", "/login", strings.NewReader(tc.body))
			if err != nil {
				t.Fatalf("创建请求失败: %s", err)
			}
			req.Header.Set("Content-Type", tc.contentType)

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			router.ServeHTTP(w, req)

			// 检查状态码
			if w.Code != http.StatusOK {
				t.Errorf("期望状态码 %d, 实际得到 %d", http.StatusOK, w.Code)
			}

			// 读取响应
			body, err := io.ReadAll(w.Body)
			if err != nil {
				t.Fatalf("读取响应失败: %s", err)
			}

			// 解析JSON响应
			var response LoginResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("解析JSON失败: %s", err)
			}

			// 检查响应消息
			if response.Message != tc.expected {
				t.Errorf("期望消息 '%s', 实际得到 '%s'", tc.expected, response.Message)
			}

			t.Logf("测试 %s 通过: %s", tc.name, response.Message)
		})
	}
}

func TestGetLoginPage(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.GET("/login", logic.GetLogin)

	// 创建请求
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %s", err)
	}

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 检查状态码
	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 实际得到 %d", http.StatusOK, w.Code)
	}

	// 读取响应
	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("读取响应失败: %s", err)
	}

	// 检查响应内容是否包含HTML
	responseBody := string(body)
	if !strings.Contains(responseBody, "<html") {
		t.Error("响应应该包含HTML内容")
	}

	if !strings.Contains(responseBody, "用户登录") {
		t.Error("响应应该包含登录页面标题")
	}

	t.Logf("登录页面加载成功，响应长度: %d", len(responseBody))
}

func TestLogout(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.GET("/logout", logic.Logout)

	// 创建请求
	req, err := http.NewRequest("GET", "/logout", nil)
	if err != nil {
		t.Fatalf("创建请求失败: %s", err)
	}

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 执行请求
	router.ServeHTTP(w, req)

	// 检查状态码（应该是重定向）
	if w.Code != http.StatusFound {
		t.Errorf("期望状态码 %d, 实际得到 %d", http.StatusFound, w.Code)
	}

	// 检查重定向位置
	location := w.Header().Get("Location")
	if location != "/index" {
		t.Errorf("期望重定向到 '/index', 实际重定向到 '%s'", location)
	}

	t.Logf("退出登录成功，重定向到: %s", location)
}

func BenchmarkLogin(t *testing.B) {
	// 初始化数据库
	err := model.Init()
	if err != nil {
		t.Fatalf("数据库初始化失败: %s", err)
	}
	defer model.Close()

	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	router := gin.New()
	router.POST("/login", logic.DoLogin)

	// 准备表单数据
	formData := url.Values{}
	formData.Set("name", "admin")
	formData.Set("password", "123456")

	for i := 0; i < t.N; i++ {
		// 创建请求
		req, err := http.NewRequest("POST", "/login", strings.NewReader(formData.Encode()))
		if err != nil {
			t.Fatalf("创建请求失败: %s", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// 创建响应记录器
		w := httptest.NewRecorder()

		// 执行请求
		router.ServeHTTP(w, req)

		// 检查状态码
		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d, 实际得到 %d", http.StatusOK, w.Code)
		}
	}
}
