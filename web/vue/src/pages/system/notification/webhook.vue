<script setup>
import appPage from '../../../layout/page.vue'
</script>
<template>
  <appPage>
    <notification-tab></notification-tab>
    <el-form ref="form" :model="form" :rules="formRules" label-width="auto" style="max-width: 800px;">
      <el-form-item label="URL" prop="url">
        <template #label>
          URL
          <el-tooltip placement="top-start">
            <template #content>
              通知内容推送到指定URL, POST请求, 设置Header [Content-Type: application/json]
            </template>
            <el-icon><QuestionFilled /></el-icon>
          </el-tooltip>
        </template>
        <el-input v-model.trim="form.url"></el-input>
      </el-form-item>
      <el-form-item label="模板" prop="template">
        <el-input
          type="textarea"
          :rows="10"
          v-model="form.template">
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submit()">保存</el-button>
      </el-form-item>
    </el-form>
  </appPage>
</template>

<script>
import notificationTab from './tab.vue'
import notificationService from '../../../api/notification'
export default {
  data () {
    return {
      form: {
        url: '',
        template: ''
      },
      formRules: {
        url: [
          {type: 'url', required: true, message: '请输入有效的通知URL', trigger: 'blur'}
        ],
        template: [
          {required: true, message: '请输入通知模板', trigger: 'blur'}
        ]
      }
    }
  },
  components: {notificationTab},
  created () {
    this.init()
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
      notificationService.updateWebHook(this.form, () => {
        this.$message.success('更新成功')
        this.init()
      })
    },
    init () {
      notificationService.webhook((data) => {
        this.form.url = data.url
        this.form.template = data.template
      })
    }
  }
}
</script>
