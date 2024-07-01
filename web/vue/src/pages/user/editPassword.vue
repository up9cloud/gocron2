<template>
  <el-container>
    <user-sidebar></user-sidebar>
    <el-main>
      <el-breadcrumb separator-icon="ArrowRight" style="margin-bottom:20px">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/user' }">用户管理</el-breadcrumb-item>
          <el-breadcrumb-item>修改密码</el-breadcrumb-item>
      </el-breadcrumb>
      <el-form ref="form" class="page-form" :model="form" :rules="formRules" label-width="auto">
        <el-row>
          <el-col :span="8">
            <el-form-item label="新密码" prop="new_password">
              <el-input v-model="form.new_password" type="password"></el-input>
            </el-form-item>
            <el-form-item label="确认新密码" prop="confirm_new_password">
              <el-input v-model="form.confirm_new_password" type="password"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submit()">保存</el-button>
              <el-button @click="cancel">取消</el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-main>
  </el-container>
</template>

<script>
import userSidebar from './sidebar.vue'
import userService from '../../api/user'
export default {
  name: 'user-edit-password',
  data: function () {
    return {
      form: {
        id: '',
        new_password: '',
        confirm_new_password: ''
      },
      formRules: {
        new_password: [
          {required: true, message: '请输入新密码', trigger: 'blur'}
        ],
        confirm_new_password: [
          {required: true, message: '请再次输入新密码', trigger: 'blur'}
        ]
      }
    }
  },
  components: {userSidebar},
  created () {
    const id = this.$route.params.id
    if (!id) {
      return
    }
    this.form.id = id
  },
  methods: {
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      userService.editPassword(this.form, () => {
        this.$router.push('/user')
      })
    },
    cancel () {
      this.$router.push('/user')
    }
  }
}
</script>
