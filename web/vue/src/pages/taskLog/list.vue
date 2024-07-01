<template>
  <el-container>
    <task-sidebar></task-sidebar>
    <el-main>
      <el-breadcrumb separator-icon="ArrowRight" style="margin-bottom:20px">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ path: '/task' }">任务管理</el-breadcrumb-item>
        <el-breadcrumb-item>任务日志</el-breadcrumb-item>
      </el-breadcrumb>
      <el-form :inline="true" label-width="auto">
        <el-form-item>
          <el-input placeholder="请输入任务ID" v-model.trim="searchParams.task_id"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model.trim="searchParams.protocol" placeholder="请选择执行方式">
            <el-option value=""></el-option>
            <el-option
            v-for="item in protocolList"
            :key="item.value"
            :label="item.label"
            :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model.trim="searchParams.status" placeholder="请选择任务状态">
            <el-option value=""></el-option>
            <el-option
              v-for="item in statusList"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="Search" @click="search()">搜索</el-button>
        </el-form-item>
      </el-form>
      <el-row type="flex" justify="end">
          <el-button type="danger" icon="Delete" v-if="this.$store.getters.user.isSuperAdmin" @click="clearLog">清空日志</el-button>
          <el-button type="info" icon="Refresh" @click="refresh">刷新</el-button>
      </el-row>
      <el-table
        :data="logs"
        border
        ref="table"
        style="margin: 20px 0;">
        <el-table-column type="expand">
          <template #default="scope">
            <el-form label-width="auto" label-position="left">
              <el-form-item>
                  重试次数: {{scope.row.retry_times}} <br>
                  cron表达式: {{scope.row.spec}} <br>
                  命令: {{scope.row.command}}
              </el-form-item>
            </el-form>
          </template>
        </el-table-column>
        <el-table-column
          prop="id"
          label="ID"
          width="100">
        </el-table-column>
        <el-table-column
          prop="task_id"
          label="任务ID"
          width="100">
        </el-table-column>
        <el-table-column
          prop="name"
          label="任务名称">
        </el-table-column>
        <el-table-column
          prop="protocol"
          label="执行方式"
          :formatter="formatProtocol"
          width="100">
        </el-table-column>
        <el-table-column
          label="任务节点"
          width="150">
          <template #default="scope">
            <div>{{scope.row.hostname}}</div>
          </template>
        </el-table-column>
        <el-table-column
          label="执行时长"
          width="250">
          <template #default="scope">
            执行时长: {{scope.row.total_time > 0 ? scope.row.total_time : 1}}秒<br>
            开始时间: {{$filters.formatTime(scope.row.start_time)}}<br>
            <span v-if="scope.row.status !== 1">结束时间: {{$filters.formatTime(scope.row.end_time)}}</span>
          </template>
        </el-table-column>
        <el-table-column
          label="状态"
          width="100"
          align="center"
          >
          <template #default="scope">
            <span v-if="scope.row.status === 0">
              <el-tag type="danger">失败</el-tag>
            </span>
            <span v-else-if="scope.row.status === 1">
              <el-tag>执行中</el-tag>
            </span>
            <span v-else-if="scope.row.status === 2">
              <el-tag type="success">成功</el-tag>
            </span>
            <span v-else-if="scope.row.status === 3">
              <el-tag type="info">取消</el-tag>
            </span>
          </template>
        </el-table-column>
        <el-table-column
          label="操作"
          align="center"
          header-align="left"
          width="110" v-if="this.isAdmin">
          <template #default="scope">
            <el-button size="small" type="success"
                       v-if="scope.row.status === 2"
                       @click="showTaskResult(scope.row)">查看结果</el-button>
            <el-button size="small" type="warning"
                       v-if="scope.row.status === 0"
                       @click="showTaskResult(scope.row)" >查看结果</el-button>
            <el-button size="small" type="danger"
                       v-if="scope.row.status === 1 && scope.row.protocol === 2"
                       @click="stopTask(scope.row)">停止任务
            </el-button>
          </template>
        </el-table-column>
        <el-table-column
          label="执行结果"
          width="102" v-else>
          <template #default="scope">
            <el-button size="small" type="success"
                       v-if="scope.row.status === 2"
                       @click="showTaskResult(scope.row)">查看结果</el-button>
            <el-button size="small" type="warning"
                       v-if="scope.row.status === 0"
                       @click="showTaskResult(scope.row)" >查看结果</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-row type="flex" justify="end">
        <el-pagination
          background
          layout="prev, pager, next, sizes, total"
          :total="logTotal"
          :page-size="20"
          @size-change="changePageSize"
          @current-change="changePage"
          @prev-click="changePage"
          @next-click="changePage">
        </el-pagination>
      </el-row>
      <el-dialog v-model="dialogVisible" width="70%">
        <div>
          <pre>{{currentTaskResult.command}}</pre>
        </div>
        <div>
          <pre>{{currentTaskResult.result}}</pre>
        </div>
      </el-dialog>
    </el-main>
  </el-container>
</template>

<script>
import taskSidebar from '../task/sidebar.vue'
import taskLogService from '../../api/taskLog'

export default {
  name: 'task-log',
  data () {
    return {
      logs: [],
      logTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        task_id: '',
        protocol: '',
        status: ''
      },
      isAdmin: this.$store.getters.user.isAdmin,
      dialogVisible: false,
      currentTaskResult: {
        command: '',
        result: ''
      },
      protocolList: [
        {
          value: '1',
          label: 'http'
        },
        {
          value: '2',
          label: 'shell'
        }
      ],
      statusList: [
        {
          value: '1',
          label: '失败'
        },
        {
          value: '2',
          label: '执行中'
        },
        {
          value: '3',
          label: '成功'
        },
        {
          value: '4',
          label: '取消'
        }
      ]
    }
  },
  components: {taskSidebar},
  created () {
    if (this.$route.query.task_id) {
      this.searchParams.task_id = this.$route.query.task_id
    }
    this.search()
  },
  methods: {
    formatProtocol (row, col) {
      if (row[col.property] === 1) {
        return 'http'
      }
      return 'shell'
    },
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search (callback = null) {
      taskLogService.list(this.searchParams, (data) => {
        this.logs = data.data
        this.logTotal = data.total

        if (callback) {
          callback()
        }
      })
    },
    clearLog () {
      this.$appConfirm(() => {
        taskLogService.clear(() => {
          this.searchParams.page = 1
          this.search()
        })
      })
    },
    stopTask (item) {
      taskLogService.stop(item.id, item.task_id, () => {
        this.search()
      })
    },
    showTaskResult (item) {
      this.dialogVisible = true
      this.currentTaskResult.command = item.command
      this.currentTaskResult.result = item.result
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    }
  }
}
</script>
<style scoped>
  pre {
    white-space: pre-wrap;
    word-wrap: break-word;
    padding: 10px;
    background-color: #334157;
    color: white;
    border-radius: 2px;
  }
</style>
