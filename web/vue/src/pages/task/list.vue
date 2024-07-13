<script setup>
import appPage from '../../layout/page.vue'
</script>
<template>
  <appPage>
    <el-form :inline="true">
      <el-row>
        <el-form-item>
          <el-input placeholder="请输入任务ID" v-model.trim="searchParams.id"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input placeholder="请输入任务名称" v-model.trim="searchParams.name"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input placeholder="请输入任务命令" v-model.trim="searchParams.command"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input placeholder="请输入标签" v-model.trim="searchParams.tag"></el-input>
        </el-form-item>
        <el-form-item>
          <el-select v-model.trim="searchParams.protocol" placeholder="请选择执行方式" clearable>
            <el-option
              v-for="item in protocolList"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model.trim="searchParams.host_id" placeholder="请选择任务节点" clearable>
            <el-option
              v-for="item in hosts"
              :key="item.id"
              :label="item.alias + ' - ' + item.name + ':' + item.port"
              :value="item.id">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-select v-model.trim="searchParams.status" placeholder="请选择状态" clearable>
            <el-option
              v-for="item in statusList"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
      </el-row>
      <el-row>
        <el-button type="primary" icon="Search" @click="search()">搜索</el-button>
        <el-button type="info" icon="CircleClose" @click="resetSearch()">重置</el-button>
      </el-row>
    </el-form>
    <el-row type="flex" justify="end">
      <el-button type="primary" icon="Edit" @click="toEdit(null)" v-if="$store.getters.user.isAdmin">新增</el-button>
      <el-button type="info" icon="Refresh" @click="refresh">刷新</el-button>
    </el-row>
    <el-container>
      <el-table :data="tasks" border row-key="id">
        <el-table-column type="expand">
          <template #default="scope">
            <el-descriptions border>
              <el-descriptions-item>
                <template #label>任务创建时间</template>
                {{$filters.formatTime(scope.row.created)}}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>任务类型</template>
                {{formatLevel(scope.row.level)}}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>单实例运行</template>
                {{formatMulti(scope.row.multi)}}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>超时时间</template>
                {{formatTimeout(scope.row.timeout)}}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>重试次数</template>
                {{scope.row.retry_times}}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>重试间隔</template>
                {{formatRetryTimesInterval(scope.row.retry_interval)}}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>任务节点</template>
                <el-tag v-for="item in scope.row.hosts" :key="item.host_id">
                  {{item.alias}} - {{item.name}}:{{item.port}}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>命令</template>
                <el-input
                  v-model="scope.row.command"
                  type="textarea"
                  disabled
                  autosize
                />
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>备注</template>
                <el-input
                  v-model="scope.row.remark"
                  type="textarea"
                  disabled
                  autosize
                />
              </el-descriptions-item>
            </el-descriptions>
          </template>
        </el-table-column>
        <el-table-column
          prop="id"
          label="任务ID" min-width="60">
        </el-table-column>
        <el-table-column
          prop="name"
          label="任务名称" min-width="100">
        </el-table-column>
        <el-table-column
          prop="tag"
          label="标签">
          <template #default="scope">
            <el-button size="small" type="success" plain @click="toTasksByTag(scope.row)" v-if="scope.row.tag">{{scope.row.tag}}</el-button>
          </template>
        </el-table-column>
        <el-table-column
          prop="spec"
          label="cron表达式"
          width="120">
        </el-table-column>
        <el-table-column label="下次执行时间" width="150">
          <template #default="scope">
            {{$filters.formatTime(scope.row.next_run_time)}}
          </template>
        </el-table-column>
        <el-table-column
          prop="protocol"
          :formatter="formatProtocol"
          label="执行方式" width="80">
        </el-table-column>
        <el-table-column
          label="状态" width="80" v-if="isAdmin" >
            <template #default="scope">
              <el-switch
                v-if="scope.row.level === 1"
                v-model="scope.row.status"
                :active-value="1"
                :inactive-value="0"
                @change="changeStatus(scope.row)"
                :disabled="!checkAuth(scope.row)">
              </el-switch>
            </template>
        </el-table-column>
        <el-table-column label="状态" width="100" v-else>
          <template #default="scope">
            <el-switch
              v-if="scope.row.level === 1"
              v-model="scope.row.status"
              :active-value="1"
              :inactive-value="0"
              disabled>
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column
          align="center"
          header-align="left"
          label="操作"
          width="150"
          v-if="isAdmin">
          <template #default="scope">
            <el-button-group>
              <el-button type="primary" size="small" @click="toEdit(scope.row)" :disabled="!checkAuth(scope.row)">编辑</el-button>
              <el-button type="success" size="small" @click="runTask(scope.row)" :disabled="!checkAuth(scope.row)">手动执行</el-button>
              <el-button type="danger" size="small" @click="remove(scope.row)" :disabled="!checkAuth(scope.row)">删除</el-button>
              <el-button type="info" size="small" @click="jumpToLog(scope.row)">查看日志</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-container>
    <el-row type="flex" justify="end">
      <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="taskTotal"
        :current-page="searchParams.page"
        :page-size="searchParams.page_size"
        @size-change="changePageSize"
        @current-change="changePage"
        @prev-click="changePage"
        @next-click="changePage">
      </el-pagination>
    </el-row>
  </appPage>
</template>

<script>
import taskService from '../../api/task'

export default {
  data () {
    return {
      tasks: [],
      hosts: [],
      tagList: [],
      taskTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        id: sessionStorage.getItem('searchParams.id') || '',
        protocol: sessionStorage.getItem('searchParams.protocol') || '',
        name: sessionStorage.getItem('searchParams.name') || '',
        tag: sessionStorage.getItem('searchParams.tag') || '',
        host_id: Number(sessionStorage.getItem('searchParams.host_id')) || '',
        status: sessionStorage.getItem('searchParams.status') || '',
        command: sessionStorage.getItem('searchParams.command') || '',
      },
      isAdmin: this.$store.getters.user.isAdmin,
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
          value: '2',
          label: '激活'
        },
        {
          value: '1',
          label: '停止'
        }
      ]
    }
  },
  created () {
    const hostId = this.$route.query.host_id
    if (hostId) {
      this.searchParams.host_id = hostId
    }
    this.initTagList()
    this.search()
  },
  methods: {
    changeStatus (item) {
      if (item.status) {
        taskService.enable(item.id)
      } else {
        taskService.disable(item.id)
      }
    },
    formatProtocol (row, col) {
      if (row[col.property] === 2) {
        return 'shell'
      }
      if (row.http_method === 1) {
        return 'http-get'
      }
      return 'http-post'
    },
    formatLevel (value) {
      if (value === 1) {
        return '主任务'
      }
      return '子任务'
    },
    formatTimeout (value) {
      if (value > 0) {
        return value + '秒'
      }
      return '不限制'
    },
    formatRetryTimesInterval (value) {
      if (value > 0) {
        return value + '秒'
      }
      return '系统默认'
    },
    formatMulti (value) {
      if (value > 0) {
        return '否'
      }
      return '是'
    },
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    resetSearch () {
      this.searchParams = this.$options.data.call(this).searchParams
      this.search()
    },
    search (callback = null) {
      this.tasks = []
      this.hosts = []
      taskService.list(this.searchParams, (tasks, hosts) => {
        this.tasks = tasks.data
        this.taskTotal = tasks.total
        this.hosts = hosts
        if (callback) {
          callback()
        }
      })
    },
    initTagList(){
      taskService.tagList((list) => {
        this.tagList = list
      })
    },
    runTask (item) {
      this.$appConfirm(() => {
        taskService.run(item.id, () => {
          this.$message.success('任务已开始执行')
        })
      }, true)
    },
    remove (item) {
      this.$appConfirm(() => {
        taskService.remove(item.id, () => {
          this.refresh()
        })
      })
    },
    jumpToLog (item) {
      this.$router.push(`/task/log?task_id=${item.id}`)
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    checkAuth (item) {
      return item.creator === 0 || this.$store.getters.user.isSuperAdmin || item.creator === parseInt(this.$store.getters.user.uid)
    },
    toEdit (item) {
      let path = ''
      if (item === null) {
        path = '/task/create'
      } else {
        path = `/task/edit/${item.id}`
      }
      this.$router.push(path)
    },
    toTasksByTag (item) {
      this.searchParams.tag = item.tag
      this.search()
    }
  }
}
</script>
