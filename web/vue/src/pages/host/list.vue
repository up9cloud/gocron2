<script setup>
import appPage from '../../layout/page.vue'
</script>
<template>
  <appPage>
    <el-form :inline="true">
      <el-row>
        <el-form-item>
          <el-input placeholder="请输入节点ID" v-model.trim="searchParams.id"></el-input>
        </el-form-item>
        <el-form-item>
          <el-input placeholder="请输入主机名" v-model.trim="searchParams.name"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search()">搜索</el-button>
        </el-form-item>
      </el-row>
    </el-form>
    <el-row type="flex" justify="end">
      <el-button type="primary" icon="Edit" @click="toEdit(null)" v-if="$store.getters.user.isAdmin">新增</el-button>
      <el-button type="info" icon="Refresh" @click="refresh">刷新</el-button>
    </el-row>
    <el-table :data="hosts" border row-key="id">
      <el-table-column
        prop="id"
        label="节点ID">
      </el-table-column>
      <el-table-column
        prop="alias"
        label="节点名">
      </el-table-column>
      <el-table-column
        prop="name"
        label="主机名">
      </el-table-column>
      <el-table-column
        prop="port"
        label="端口">
      </el-table-column>
      <el-table-column
        prop="remark"
        label="备注">
      </el-table-column>
      <el-table-column
        align="center"
        header-align="left"
        label="操作"
        v-if="isAdmin">
        <template #default="scope">
          <el-button-group>
            <el-button size="small" type="primary" @click="toEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="success" @click="toTasks(scope.row)">查看任务</el-button>
            <el-button size="small" type="danger" @click="remove(scope.row)">删除</el-button>
            <el-button size="small" type="info" @click="ping(scope.row)">测试连接</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <el-row type="flex" justify="end">
      <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="hostTotal"
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
import hostService from '../../api/host'
export default {
  data () {
    return {
      hosts: [],
      hostTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1,
        id: '',
        name: '',
        alias: ''
      },
      isAdmin: this.$store.getters.user.isAdmin
    }
  },
  created () {
    this.search()
  },
  methods: {
    changePage (page) {
      this.searchParams.page = page
      this.search()
    },
    changePageSize (pageSize) {
      this.searchParams.page_size = pageSize
      this.search()
    },
    search (callback = null) {
      this.hosts = []
      hostService.list(this.searchParams, (data) => {
        this.hosts = data.data
        this.hostTotal = data.total
        if (callback) {
          callback()
        }
      })
    },
    remove (item) {
      this.$appConfirm(() => {
        hostService.remove(item.id, () => this.refresh())
      })
    },
    ping (item) {
      hostService.ping(item.id, () => {
        this.$message.success('连接成功')
      })
    },
    toEdit (item) {
      let path = ''
      if (item === null) {
        path = '/host/create'
      } else {
        path = `/host/edit/${item.id}`
      }
      this.$router.push(path)
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    toTasks (item) {
      this.$router.push(
        {
          path: '/task',
          query: {
            host_id: item.id
          }
        })
    }
  }
}
</script>
