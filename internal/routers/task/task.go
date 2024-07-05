package task

import (
	"strconv"
	"strings"

	"github.com/ouqiang/goutil"

	"github.com/go-macaron/binding"
	"github.com/jakecoffman/cron"
	"github.com/up9cloud/gocron2/internal/models"
	"github.com/up9cloud/gocron2/internal/modules/logger"
	"github.com/up9cloud/gocron2/internal/modules/utils"
	"github.com/up9cloud/gocron2/internal/routers/base"
	"github.com/up9cloud/gocron2/internal/routers/user"
	"github.com/up9cloud/gocron2/internal/service"
	"gopkg.in/macaron.v1"
)

type TaskForm struct {
	Id               int
	Level            models.TaskLevel `binding:"Required;In(1,2)"`
	DependencyStatus models.TaskDependencyStatus
	DependencyTaskId string
	Name             string `binding:"Required;MaxSize(32)"`
	Spec             string
	Protocol         models.TaskProtocol   `binding:"In(1,2)"`
	Command          string                `binding:"Required;MaxSize(256)"`
	HttpMethod       models.TaskHTTPMethod `binding:"In(1,2)"`
	Timeout          int                   `binding:"Range(0,86400)"`
	Multi            int8                  `binding:"In(1,2)"`
	RetryTimes       int8
	RetryInterval    int16
	HostId           string
	Tag              string
	Remark           string
	NotifyStatus     int8 `binding:"In(1,2,3,4)"`
	NotifyType       int8 `binding:"In(1,2,3,4)"`
	NotifyReceiverId string
	NotifyKeyword    string
	Status           models.Status `binding:"In(0,1)"`
}

func (f TaskForm) Error(ctx *macaron.Context, errs binding.Errors) {
	if len(errs) == 0 {
		return
	}
	json := utils.JsonResponse{}
	content := json.CommonFailure("表单验证失败, 请检测输入")

	ctx.Resp.Write([]byte(content))
}

// 首页
func Index(ctx *macaron.Context) string {
	taskModel := new(models.Task)
	queryParams := parseQueryParams(ctx)
	total, err := taskModel.Total(queryParams)
	if err != nil {
		logger.Error(err)
	}
	tasks, err := taskModel.List(queryParams)
	if err != nil {
		logger.Error(err)
	}
	for i, item := range tasks {
		tasks[i].NextRunTime = service.ServiceTask.NextRunTime(item)
	}
	jsonResp := utils.JsonResponse{}

	return jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  tasks,
	})
}

func Dependency(ctx *macaron.Context) string {
	taskModel := new(models.Task)
	queryParams := parseQueryParams(ctx)
	total, err := taskModel.Total(queryParams)
	if err != nil {
		logger.Error(err)
	}
	tasks, err := taskModel.DependencyList(queryParams)
	if err != nil {
		logger.Error(err)
	}
	for i, item := range tasks {
		tasks[i].NextRunTime = service.ServiceTask.NextRunTime(item)
	}
	jsonResp := utils.JsonResponse{}

	return jsonResp.Success(utils.SuccessContent, map[string]interface{}{
		"total": total,
		"data":  tasks,
	})
}

// Detail 任务详情
func Detail(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	taskModel := new(models.Task)
	task, err := taskModel.Detail(id)
	jsonResp := utils.JsonResponse{}
	if err != nil || task.Id == 0 {
		logger.Errorf("编辑任务#获取任务详情失败#任务ID-%d", id)
		return jsonResp.Success(utils.SuccessContent, nil)
	}

	return jsonResp.Success(utils.SuccessContent, task)
}

// 标签列表
func Tags(ctx *macaron.Context) string {
	taskModel := new(models.Task)
	tagsList, err := taskModel.TagsList()
	jsonResp := utils.JsonResponse{}
	if err != nil {
		logger.Errorf("编辑任务#获取任务Tags列表出错")
		return jsonResp.Success(utils.SuccessContent, nil)
	}
	return jsonResp.Success(utils.SuccessContent, tagsList)
}

// 保存任务  todo 拆分为多个方法
func Store(ctx *macaron.Context, form TaskForm) string {
	json := utils.JsonResponse{}
	taskModel := models.Task{}
	var id = form.Id
	nameExists, err := taskModel.NameExist(form.Name, form.Id)
	if err != nil {
		return json.CommonFailure(utils.FailureContent, err)
	}
	if nameExists {
		return json.CommonFailure("任务名称已存在")
	}

	if form.DependencyTaskId != "" {
		childrenIds := strings.Split(form.DependencyTaskId, ",")
		for _, childrenId := range childrenIds {
			intChildrenId, _ := strconv.Atoi(childrenId)
			childrenUsed := taskModel.ChildrenExist(id, intChildrenId)
			if childrenUsed {
				return json.CommonFailure("子任务不能被重复使用")
			}
		}
	}

	if form.Protocol == models.TaskRPC && form.HostId == "" {
		return json.CommonFailure("请选择主机名")
	}

	// 编辑时校验权限
	if id > 0 {
		task, errMsg := taskModel.Detail(id)
		if errMsg != nil || task.Id <= 0 {
			return json.CommonFailure("获取任务详情失败", errMsg)
		}

		if !validateWriteAuth(ctx, task.Creator) {
			return json.CommonFailure("该任务无权限操作")
		}
	}

	taskModel.Name = form.Name
	taskModel.Protocol = form.Protocol
	taskModel.Command = strings.TrimSpace(form.Command)
	taskModel.Timeout = form.Timeout
	taskModel.Tag = form.Tag
	taskModel.Remark = form.Remark
	taskModel.Multi = form.Multi
	taskModel.RetryTimes = form.RetryTimes
	taskModel.RetryInterval = form.RetryInterval
	if taskModel.Multi != 1 {
		taskModel.Multi = 0
	}
	taskModel.NotifyStatus = form.NotifyStatus - 1
	taskModel.NotifyType = form.NotifyType - 1
	taskModel.NotifyReceiverId = form.NotifyReceiverId
	taskModel.NotifyKeyword = form.NotifyKeyword
	taskModel.Spec = form.Spec
	taskModel.Level = form.Level
	taskModel.Status = form.Status
	taskModel.DependencyStatus = form.DependencyStatus
	taskModel.DependencyTaskId = strings.TrimSpace(form.DependencyTaskId)
	if taskModel.NotifyStatus > 0 && taskModel.NotifyType != 3 && taskModel.NotifyReceiverId == "" {
		return json.CommonFailure("至少选择一个通知接收者")
	}
	taskModel.HttpMethod = form.HttpMethod
	if taskModel.Protocol == models.TaskHTTP {
		command := strings.ToLower(taskModel.Command)
		if !strings.HasPrefix(command, "http://") && !strings.HasPrefix(command, "https://") {
			return json.CommonFailure("请输入正确的URL地址")
		}
		if taskModel.Timeout > 300 {
			return json.CommonFailure("HTTP任务超时时间不能超过300秒")
		}
	}

	if taskModel.RetryTimes > 10 || taskModel.RetryTimes < 0 {
		return json.CommonFailure("任务重试次数取值0-10")
	}

	if taskModel.RetryInterval > 3600 || taskModel.RetryInterval < 0 {
		return json.CommonFailure("任务重试间隔时间取值0-3600")
	}

	if taskModel.DependencyStatus != models.TaskDependencyStatusStrong &&
		taskModel.DependencyStatus != models.TaskDependencyStatusWeak {
		return json.CommonFailure("请选择依赖关系")
	}

	if taskModel.Level == models.TaskLevelParent {
		err = goutil.PanicToError(func() {
			cron.Parse(form.Spec)
		})
		if err != nil {
			return json.CommonFailure("crontab表达式解析失败", err)
		}
	} else {
		//taskModel.DependencyTaskId = ""
		taskModel.Spec = ""
	}

	if id > 0 && taskModel.DependencyTaskId != "" {
		dependencyTaskIds := strings.Split(taskModel.DependencyTaskId, ",")
		if utils.InStringSlice(dependencyTaskIds, strconv.Itoa(id)) {
			return json.CommonFailure("不允许设置当前任务为子任务")
		}
	}

	taskModel.Updater = user.Uid(ctx)

	if id == 0 {
		taskModel.Creator = user.Uid(ctx)
		id, err = taskModel.Create()
	} else {
		_, err = taskModel.UpdateBean(id)
	}

	if err != nil {
		return json.CommonFailure("保存失败", err)
	}

	taskHostModel := new(models.TaskHost)
	if form.Protocol == models.TaskRPC {
		hostIdStrList := strings.Split(form.HostId, ",")
		hostIds := make([]int, len(hostIdStrList))
		for i, hostIdStr := range hostIdStrList {
			hostIds[i], _ = strconv.Atoi(hostIdStr)
		}
		taskHostModel.Add(id, hostIds)
	} else {
		taskHostModel.Remove(id)
	}

	status, _ := taskModel.GetStatus(id)
	if status == models.Enabled {
		if taskModel.Level == models.TaskLevelParent {
			addTaskToTimer(id)
		}
	} else {
		service.ServiceTask.Remove(id)
	}

	return json.Success("保存成功", nil)
}

// 删除任务
func Remove(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	json := utils.JsonResponse{}
	taskModel := new(models.Task)

	task, errMsg := taskModel.Detail(id)
	if errMsg != nil || task.Id <= 0 {
		return json.CommonFailure("获取任务详情失败", errMsg)
	}

	if !validateWriteAuth(ctx, task.Creator) {
		return json.CommonFailure("该任务无权限操作")
	}

	_, err := taskModel.Delete(id)
	if err != nil {
		return json.CommonFailure(utils.FailureContent, err)
	}

	taskHostModel := new(models.TaskHost)
	taskHostModel.Remove(id)

	service.ServiceTask.Remove(id)

	return json.Success(utils.SuccessContent, nil)
}

// 激活任务
func Enable(ctx *macaron.Context) string {
	return changeStatus(ctx, models.Enabled)
}

// 暂停任务
func Disable(ctx *macaron.Context) string {
	return changeStatus(ctx, models.Disabled)
}

// 手动运行任务
func Run(ctx *macaron.Context) string {
	id := ctx.ParamsInt(":id")
	json := utils.JsonResponse{}
	taskModel := new(models.Task)
	task, err := taskModel.Detail(id)
	if err != nil || task.Id <= 0 {
		return json.CommonFailure("获取任务详情失败", err)
	}

	if !validateWriteAuth(ctx, task.Creator) {
		return json.CommonFailure("该任务无权限操作", err)
	}

	task.Spec = "手动运行"
	service.ServiceTask.Run(task)

	return json.Success("任务已开始运行, 请到任务日志中查看结果", nil)
}

// 改变任务状态
func changeStatus(ctx *macaron.Context, status models.Status) string {
	id := ctx.ParamsInt(":id")
	json := utils.JsonResponse{}
	taskModel := new(models.Task)

	task, errMsg := taskModel.Detail(id)
	if errMsg != nil || task.Id <= 0 {
		return json.CommonFailure("获取任务详情失败", errMsg)
	}

	if !validateWriteAuth(ctx, task.Creator) {
		return json.CommonFailure("该任务无权限操作")
	}

	_, err := taskModel.Update(id, models.CommonMap{
		"status": status,
	})
	if err != nil {
		return json.CommonFailure(utils.FailureContent, err)
	}

	if status == models.Enabled {
		addTaskToTimer(id)
	} else {
		service.ServiceTask.Remove(id)
	}

	return json.Success(utils.SuccessContent, nil)
}

// 添加任务到定时器
func addTaskToTimer(id int) {
	taskModel := new(models.Task)
	task, err := taskModel.Detail(id)
	if err != nil {
		logger.Error(err)
		return
	}

	service.ServiceTask.RemoveAndAdd(task)
}

// 检查任务是否有写权限
func validateWriteAuth(ctx *macaron.Context, uid int) bool {
	// 兼容历史
	if uid == 0 {
		return true
	}

	// 超级管理员
	if user.IsSuperAdmin(ctx) {
		return true
	}

	// 管理员用户自己的任务
	if uid == user.Uid(ctx) && user.IsAdmin(ctx) {
		return true
	}

	return false
}

// 解析查询参数
func parseQueryParams(ctx *macaron.Context) models.CommonMap {
	var params models.CommonMap = models.CommonMap{}
	params["Id"] = ctx.QueryInt("id")
	params["HostId"] = ctx.QueryInt("host_id")
	params["Name"] = ctx.QueryTrim("name")
	params["Protocol"] = ctx.QueryInt("protocol")
	params["Tag"] = ctx.QueryTrim("tag")
	params["Command"] = ctx.QueryTrim("command")
	status := ctx.QueryInt("status")
	if status >= 0 {
		status -= 1
	}
	params["Status"] = status
	base.ParsePageAndPageSize(ctx, params)

	return params
}
