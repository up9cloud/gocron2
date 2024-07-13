<script setup>
import appPage from '../../layout/page.vue'
</script>
<template>
  <appPage>
    <el-row type="flex" justify="end">
      <el-button type="primary" icon="Edit" @click="toEdit(null)" v-if="$store.getters.user.isSuperAdmin">新增</el-button>
      <el-button type="info" icon="Refresh" @click="refresh">刷新</el-button>
    </el-row>
    <el-table :data="users" border row-key="id">
      <el-table-column
        prop="id"
        label="用户id"
        width="100">
      </el-table-column>
      <el-table-column
        prop="name"
        label="用户名">
      </el-table-column>
      <el-table-column
        prop="email"
        label="邮箱">
      </el-table-column>
      <el-table-column
        prop="is_admin"
        :formatter="formatRole"
        label="角色">
      </el-table-column>
      <el-table-column
        label="状态">
        <template #default="scope">
          <el-switch
            v-model="scope.row.status"
            :active-value="1"
            :inactive-value="0"
            @change="changeStatus(scope.row)">
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column
        align="center"
        header-align="left"
        label="操作"
        v-if="isSuperAdmin">
        <template #default="scope">
          <el-button-group>
            <el-button size="small" type="primary" @click="toEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="success" @click="editPassword(scope.row)">修改密码</el-button>
            <el-button size="small" type="danger" @click="remove(scope.row)">删除</el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <el-row type="flex" justify="end">
      <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="userTotal"
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
import userService from '../../api/user'
export default {
  data () {
    return {
      users: [],
      userTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1
      },
      isSuperAdmin: this.$store.getters.user.isSuperAdmin
    }
  },
  created () {
    this.search()
  },
  methods: {
    changeStatus (item) {
      if (item.status) {
        userService.enable(item.id)
      } else {
        userService.disable(item.id)
      }
    },
    formatRole (row, col) {
      if (row[col.property] === 2) {
        return '超级管理员'
      } else if (row[col.property] === 1) {
        return '管理员'
      } else {
        return '普通用户'
      }
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
      this.users = []
      userService.list(this.searchParams, (data) => {
        this.users = data.data
        this.userTotal = data.total
        if (callback) {
          callback()
        }
      })
    },
    remove (item) {
      this.$appConfirm(() => {
        userService.remove(item.id, () => {
          this.refresh()
        })
      })
    },
    toEdit (item) {
      let path = ''
      if (item === null) {
        path = '/user/create'
      } else {
        path = `/user/edit/${item.id}`
      }
      this.$router.push(path)
    },
    refresh () {
      this.search(() => {
        this.$message.success('刷新成功')
      })
    },
    editPassword (item) {
      this.$router.push(`/user/edit-password/${item.id}`)
    }
  }
}
</script>
