<script setup>
import appPage from '../../layout/page.vue'
</script>
<template>
  <appPage>
    <el-table :data="logs" border row-key="id">
      <el-table-column
        prop="id"
        label="ID">
      </el-table-column>
      <el-table-column
        prop="username"
        label="用户名">
      </el-table-column>
      <el-table-column
        prop="ip"
        label="登录IP">
      </el-table-column>
      <el-table-column
        label="登录时间">
        <template #default="scope">
          {{$filters.formatTime(scope.row.created)}}
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end">
      <el-pagination
        background
        layout="prev, pager, next, sizes, total"
        :total="logTotal"
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
import systemService from '../../api/system'
export default {
  data () {
    return {
      logs: [],
      logTotal: 0,
      searchParams: {
        page_size: 20,
        page: 1
      }
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
    search () {
      this.logs = []
      systemService.loginLogList(this.searchParams, (data) => {
        this.logs = data.data
        this.logTotal = data.total
      })
    }
  }
}
</script>
