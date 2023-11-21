<template>
  <el-container>
    <task-sidebar></task-sidebar>
    <el-main>
      <el-form :inline="true">
        <el-row>
          <el-form-item label="任务ID">
            <el-input v-model.trim="searchParams.id"></el-input>
          </el-form-item>
          <el-form-item label="任务名称">
            <el-input v-model.trim="searchParams.name"></el-input>
          </el-form-item>
          <el-form-item label="执行方式">
            <el-select v-model.trim="searchParams.protocol">
              <el-option label="全部" value=""></el-option>
              <el-option
                v-for="item in protocolList"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="任务节点">
            <el-select v-model.trim="searchParams.host_id">
              <el-option label="全部" value=""></el-option>
              <el-option
                v-for="item in hosts"
                :key="item.id"
                :label="item.alias + ' - ' + item.name + ':' + item.port "
                :value="item.id">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model.trim="searchParams.status">
              <el-option label="全部" value=""></el-option>
              <el-option
                v-for="item in statusList"
                :key="item.value"
                :label="item.label"
                :value="item.value">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="search()">搜索</el-button>
          </el-form-item>
        </el-row>
      </el-form>

      <el-row type="flex" justify="end">
        <el-col :span="2">
          <el-button type="primary" @click="toEdit(null)" v-if="this.$store.getters.user.isAdmin">新增</el-button>
        </el-col>
        <el-col :span="2">
          <el-button type="info" @click="refresh">刷新</el-button>
        </el-col>
      </el-row>
      <el-divider content-position="left">任务依赖关系</el-divider>
      <el-table
        :data="tasks"
        style="width: 100%;margin-bottom: 20px;"
        border
        row-key="id"
      >
        <el-table-column
          prop="id"
          label="任务ID"
          sortable>
        </el-table-column>
        <el-table-column
          prop="name"
          label="任务名称"
          sortable>
        </el-table-column>

        <el-table-column
          prop="tag"
          label="标签">
        </el-table-column>
        <el-table-column
          prop="spec"
          label="cron表达式"
          width="150">
        </el-table-column>

        <el-table-column label="下次执行时间" width="170">
          <template slot-scope="scope">
            {{scope.row.next_run_time | formatTime}}
          </template>
        </el-table-column>
        <el-table-column label="上次时间" width="170">
          <template slot-scope="scope">
            {{scope.row.last_run_info.end_time | formatTime}}
          </template>
        </el-table-column>
        <el-table-column label="上次结果" width="80">
          <template slot-scope="scope">
            <span v-if="scope.row.last_run_info.id === 0">-</span>
            <span style="color:red" v-else-if="scope.row.last_run_info.status === 0">失败</span>
            <span style="color:green" v-else-if="scope.row.last_run_info.status === 1">执行中</span>
            <span v-else-if="scope.row.last_run_info.status === 2">成功</span>
            <span style="color:#4499EE" v-else-if="scope.row.last_run_info.status === 3">取消</span>
          </template>
        </el-table-column>
        <el-table-column
          label="状态" v-if="this.isAdmin" width="80">
          <template slot-scope="scope">
            <el-switch
              v-if="scope.row.level === 1"
              v-model="scope.row.status"
              :active-value="1"
              :inactive-vlaue="0"
              active-color="#13ce66"
              @change="changeStatus(scope.row)"
              inactive-color="#ff4949">
            </el-switch>
          </template>
        </el-table-column>

        <el-table-column label="状态" v-else width="80">
          <template slot-scope="scope">
            <el-switch
              v-if="scope.row.level === 1"
              v-model="scope.row.status"
              :active-value="1"
              :inactive-vlaue="0"
              active-color="#13ce66"
              :disabled="true"
              inactive-color="#ff4949">
            </el-switch>
          </template>
        </el-table-column>
        <el-table-column
          prop="protocol"
          :formatter="formatProtocol"
          label="执行方式" width="80">
        </el-table-column>
        <el-table-column label="操作" width="220" v-if="this.isAdmin">
          <template slot-scope="scope">
            <el-row>
              <el-col span="12">
                <el-button type="success" @click="runTask(scope.row)" size="mini">手动执行</el-button>
              </el-col>
              <el-col span="12">
                <el-button type="info" @click="jumpToLog(scope.row)" size="mini">查看日志</el-button>
              </el-col>
            </el-row>
          </template>
        </el-table-column>
      </el-table>

    </el-main>
  </el-container>
</template>

<script>
import taskSidebar from './sidebar'
import taskService from '../../api/task'

export default {
  name: 'task-dependency',
  data() {
    return {
      tasks: [],
      hosts: [],
      taskTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        id: '',
        protocol: '',
        name: '',
        tag: '',
        host_id: '',
        status: ''
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
  components: {taskSidebar},
  created() {
    const hostId = this.$route.query.host_id
    if (hostId) {
      this.searchParams.host_id = hostId
    }

    this.search()
  },
  filters: {
    formatLevel(value) {
      if (value === 1) {
        return '主任务'
      }
      return '子任务'
    },
    formatTimeout(value) {
      if (value > 0) {
        return value + '秒'
      }
      return '不限制'
    },
    formatRetryTimesInterval(value) {
      if (value > 0) {
        return value + '秒'
      }
      return '系统默认'
    },
    formatMulti(value) {
      if (value > 0) {
        return '否'
      }
      return '是'
    }
  },
  methods: {
    changeStatus(item) {
      if (item.status) {
        taskService.enable(item.id)
      } else {
        taskService.disable(item.id)
      }
    },
    formatProtocol(row, col) {
      if (row[col.property] === 2) {
        return 'shell'
      }
      if (row.http_method === 1) {
        return 'http-get'
      }
      return 'http-post'
    },
    changePage(page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize(pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search(callback = null) {
      taskService.dependencyList(this.searchParams, (tasks, hosts) => {
        this.tasks = tasks.data
        this.taskTotal = tasks.total
        this.hosts = hosts
        if (callback) {
          callback()
        }
      })
    },
    runTask(item) {
      this.$appConfirm(() => {
        taskService.run(item.id, () => {
          this.$message.success('任务已开始执行')
        })
      }, true)
    },
    remove(item) {
      this.$appConfirm(() => {
        taskService.remove(item.id, () => {
          this.refresh()
        })
      })
    },
    jumpToLog(item) {
      this.$router.push(`/task/log?task_id=${item.id}`)
    },
    refresh() {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    toEdit(item) {
      let path = ''
      if (item === null) {
        path = '/task/create'
      } else {
        path = `/task/edit/${item.id}`
      }
      this.$router.push(path)
    }
  }
}
</script>
<style>
.custom-tree-node {
  flex: 1;
  justify-content: space-between;
  font-size: 14px;
}
</style>
<style scoped>
.demo-table-expand {
  font-size: 0;
}

.demo-table-expand label {
  width: 90px;
  color: #99a9bf;
}

.demo-table-expand .el-form-item {
  margin-right: 0;
  margin-bottom: 0;
  width: 50%;
}
</style>
