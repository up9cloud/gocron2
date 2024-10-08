package routers

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-macaron/binding"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/toolbox"
	"github.com/rakyll/statik/fs"
	"github.com/up9cloud/gocron2/internal/modules/app"
	"github.com/up9cloud/gocron2/internal/modules/logger"
	"github.com/up9cloud/gocron2/internal/modules/utils"
	"github.com/up9cloud/gocron2/internal/routers/host"
	"github.com/up9cloud/gocron2/internal/routers/install"
	"github.com/up9cloud/gocron2/internal/routers/loginlog"
	"github.com/up9cloud/gocron2/internal/routers/manage"
	"github.com/up9cloud/gocron2/internal/routers/task"
	"github.com/up9cloud/gocron2/internal/routers/tasklog"
	"github.com/up9cloud/gocron2/internal/routers/user"
	"gopkg.in/macaron.v1"

	_ "github.com/up9cloud/gocron2/internal/statik"
)

const (
	urlPrefix = "/api"
	staticDir = "assets"
)

var statikFS http.FileSystem

func init() {
	var err error
	statikFS, err = fs.New()
	if err != nil {
		log.Fatal(err)
	}
}

// Register 路由注册
func Register(m *macaron.Macaron) {
	m.SetURLPrefix(urlPrefix)
	// 所有GET方法，自动注册HEAD方法
	m.SetAutoHead(true)
	m.Get("/", func(ctx *macaron.Context) {
		file, err := statikFS.Open("/index.html")
		if err != nil {
			logger.Errorf("读取首页文件失败: %s", err)
			ctx.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Copy(ctx.Resp, file)
	})
	m.Get("/favicon.ico", func(ctx *macaron.Context) {
		file, err := statikFS.Open("/favicon.ico")
		if err != nil {
			logger.Errorf("读取首页文件失败: %s", err)
			ctx.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Copy(ctx.Resp, file)
	})
	m.Get("/robots.txt", func(ctx *macaron.Context) {
		file, err := statikFS.Open("/robots.txt")
		if err != nil {
			logger.Errorf("读取首页文件失败: %s", err)
			ctx.WriteHeader(http.StatusInternalServerError)
			return
		}
		io.Copy(ctx.Resp, file)
	})
	// 系统安装
	m.Group("/install", func() {
		m.Post("/store", binding.Bind(install.InstallForm{}), install.Store)
		m.Get("/status", func(ctx *macaron.Context) string {
			jsonResp := utils.JsonResponse{}
			return jsonResp.Success("", app.Installed)
		})
	})

	// 用户
	m.Group("/user", func() {
		m.Get("", user.Index)
		m.Get("/:id", user.Detail)
		m.Post("/store", binding.Bind(user.UserForm{}), user.Store)
		m.Post("/remove/:id", user.Remove)
		m.Post("/login", user.ValidateLogin)
		m.Post("/enable/:id", user.Enable)
		m.Post("/disable/:id", user.Disable)
		m.Post("/editMyPassword", user.UpdateMyPassword)
		m.Post("/editPassword/:id", user.UpdatePassword)
	})

	// 定时任务
	m.Group("/task", func() {
		m.Post("/store", binding.Bind(task.TaskForm{}), task.Store)
		m.Get("/:id", task.Detail)
		m.Get("", task.Index)
		m.Get("/dependency", task.Dependency)
		m.Get("/log", tasklog.Index)
		m.Post("/log/clear", tasklog.Clear)
		m.Post("/log/remove/:id", tasklog.Remove)
		m.Post("/log/remove/day/:id", tasklog.RemoveDay)
		m.Post("/log/stop", tasklog.Stop)
		m.Post("/remove/:id", task.Remove)
		m.Post("/enable/:id", task.Enable)
		m.Post("/disable/:id", task.Disable)
		m.Get("/run/:id", task.Run)
		m.Get("/tags", task.Tags)
	})

	// 主机
	m.Group("/host", func() {
		m.Get("/:id", host.Detail)
		m.Post("/store", binding.Bind(host.HostForm{}), host.Store)
		m.Get("", host.Index)
		m.Get("/all", host.All)
		m.Get("/ping/:id", host.Ping)
		m.Post("/remove/:id", host.Remove)
	})

	// 管理
	m.Group("/system", func() {
		m.Group("/slack", func() {
			m.Get("", manage.Slack)
			m.Post("/update", manage.UpdateSlack)
			m.Post("/channel", manage.CreateSlackChannel)
			m.Post("/channel/remove/:id", manage.RemoveSlackChannel)
		})
		m.Group("/mail", func() {
			m.Get("", manage.Mail)
			m.Post("/update", binding.Bind(manage.MailServerForm{}), manage.UpdateMail)
			m.Post("/user", manage.CreateMailUser)
			m.Post("/user/remove/:id", manage.RemoveMailUser)
		})
		m.Group("/webhook", func() {
			m.Get("", manage.WebHook)
			m.Post("/update", manage.UpdateWebHook)
		})
		m.Get("/login-log", loginlog.Index)
	})

	// API
	m.Group("/v1", func() {
		m.Post("/tasklog/remove/:id", tasklog.Remove)
		m.Post("/task/enable/:id", task.Enable)
		m.Post("/task/disable/:id", task.Disable)
	}, apiAuth)

	// 404错误
	m.NotFound(func(ctx *macaron.Context) string {
		jsonResp := utils.JsonResponse{}

		return jsonResp.Failure(utils.NotFound, "您访问的页面不存在")
	})
	// 50x错误
	m.InternalServerError(func(ctx *macaron.Context) string {
		jsonResp := utils.JsonResponse{}

		return jsonResp.Failure(utils.ServerError, "服务器内部错误, 请稍后再试")
	})
}

// 中间件注册
func RegisterMiddleware(m *macaron.Macaron) {
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	if macaron.Env != macaron.DEV {
		m.Use(gzip.Gziper())
	}
	m.Use(
		macaron.Static(
			staticDir,
			// https://pkg.go.dev/gopkg.in/macaron.v1?utm_source=godoc#StaticOptions
			macaron.StaticOptions{
				FileSystem: statikFS,
			},
		),
	)
	if macaron.Env == macaron.DEV {
		m.Use(toolbox.Toolboxer(m))
	}
	m.Use(macaron.Renderer())
	m.Use(appAuth)
}

// region 自定义中间件

func appAuth(ctx *macaron.Context) {
	// 检测应用是否已安装
	if !app.Installed {
		if strings.HasPrefix(ctx.Req.URL.Path, "/install") {
			return
		}
		jsonResp := utils.JsonResponse{}
		data := jsonResp.Failure(utils.AppNotInstall, "应用未安装")
		ctx.Write([]byte(data))
		return
	}
	// localhost bypass
	clientIp := ctx.RemoteAddr()
	isBypass := false
	if app.Setting.AllowLocalhostBypass {
		if clientIp == "localhost" || clientIp == "127.0.0.1" || clientIp == "[::1]" {
			isBypass = true
		}
	}
	ctx.Data["is_bypass"] = isBypass
	if isBypass {
		ctx.Data["uid"] = int(0)
		ctx.Data["username"] = ""
		ctx.Data["is_admin"] = int(2)
		return
	}
	// 用户认证
	user.RestoreToken(ctx)
	uri := strings.TrimRight(ctx.Req.URL.Path, "/")
	// 未登入允许访问的URL地址
	allowPaths := []string{
		"/user/login",
		"/install/status",
	}
	if utils.InStringSlice(allowPaths, uri) {
		return
	}
	if strings.HasPrefix(uri, "/v1") {
		return
	}
	if !user.IsLogin(ctx) {
		jsonResp := utils.JsonResponse{}
		data := jsonResp.Failure(utils.AuthError, "认证失败")
		ctx.Write([]byte(data))
		return
	}
	// IP验证, 通过反向代理访问gocron2，需设置Header X-Real-IP才能获取到客户端真实IP
	if app.Setting.AllowIps != "" {
		allowIps := strings.Split(app.Setting.AllowIps, ",")
		if !utils.InStringSlice(allowIps, clientIp) {
			logger.Warnf("非法IP访问-%s", clientIp)
			jsonResp := utils.JsonResponse{}
			data := jsonResp.Failure(utils.UnauthorizedError, "您无权限访问")
			ctx.Write([]byte(data))
			return
		}
	}
	if user.IsSuperAdmin(ctx) {
		return
	}
	if user.IsAdmin(ctx) {
		return
	}
	// 普通用户允许访问的URL地址
	allowPaths = []string{
		"/task",
		"/task/log",
		"/host",
		"/host/all",
		"/user/login",
		"/user/editMyPassword",
	}
	if utils.InStringSlice(allowPaths, uri) {
		return
	}
	jsonResp := utils.JsonResponse{}
	data := jsonResp.Failure(utils.UnauthorizedError, "您无权限访问")
	ctx.Write([]byte(data))
}

/** API接口签名验证 **/
func apiAuth(ctx *macaron.Context) {
	if !app.Installed {
		return
	}
	if !app.Setting.ApiSignEnable {
		return
	}
	apiKey := strings.TrimSpace(app.Setting.ApiKey)
	apiSecret := strings.TrimSpace(app.Setting.ApiSecret)
	json := utils.JsonResponse{}
	if apiKey == "" || apiSecret == "" {
		msg := json.CommonFailure("使用API前, 请先配置密钥")
		ctx.Write([]byte(msg))
		return
	}
	currentTimestamp := time.Now().Unix()
	time := ctx.QueryInt64("time")
	if time <= 0 {
		msg := json.CommonFailure("参数time不能为空")
		ctx.Write([]byte(msg))
		return
	}
	if time < (currentTimestamp - 1800) {
		msg := json.CommonFailure("time无效")
		ctx.Write([]byte(msg))
		return
	}
	sign := ctx.QueryTrim("sign")
	if sign == "" {
		msg := json.CommonFailure("参数sign不能为空")
		ctx.Write([]byte(msg))
		return
	}
	raw := apiKey + strconv.FormatInt(time, 10) + strings.TrimSpace(ctx.Req.URL.Path) + apiSecret
	realSign := utils.Md5(raw)
	if sign != realSign {
		msg := json.CommonFailure("签名验证失败")
		ctx.Write([]byte(msg))
		return
	}
}

// endregion
