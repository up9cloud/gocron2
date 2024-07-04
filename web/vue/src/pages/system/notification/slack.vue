<script setup>
import appPage from '../../../layout/page.vue'
</script>
<template>
  <appPage>
    <notification-tab></notification-tab>
    <el-form ref="form" :model="form" :rules="formRules" label-width="auto" style="max-width: 800px;">
      <el-form-item label="Slack Webhook URL" prop="url">
        <el-input v-model="form.url"></el-input>
      </el-form-item>
      <el-form-item label="模板" prop="template">
        <el-input
          type="textarea"
          :rows="10"
          v-model="form.template">
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submit">保存</el-button>
      </el-form-item>
      <h3>Channel &nbsp;&nbsp;&nbsp;<el-button type="primary" size="small" icon="Plus" plain @click="createChannel"></el-button></h3>
      <el-tag
        v-for="item in channels"
        :key="item.id"
        closable
        @close="deleteChannel(item)"
      >
        {{item.name}}
      </el-tag>
    </el-form>
    <el-dialog v-model="dialogVisible" width="30%">
      <el-form :model="form" label-width="auto">
        <el-form-item label="Channel名称" >
          <el-input v-model.trim="channel" v-focus></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="saveChannel">确 定</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </appPage>
</template>

<script>
import notificationTab from './tab.vue'
import notificationService from '../../../api/notification'
export default {
  data () {
    return {
      dialogVisible: false,
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
      },
      channels: [],
      channel: ''
    }
  },
  components: {notificationTab},
  created () {
    this.init()
  },
  methods: {
    createChannel () {
      this.dialogVisible = true
    },
    submit () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return false
        }
        this.save()
      })
    },
    save () {
      notificationService.updateSlack(this.form, () => {
        this.$message.success('更新成功')
        this.init()
      })
    },
    saveChannel () {
      if (this.channel === '') {
        this.$message.error('请输入Channel名称')
        return
      }
      notificationService.createSlackChannel(this.channel, () => {
        this.dialogVisible = false
        this.init()
      })
    },
    deleteChannel (item) {
      notificationService.removeSlackChannel(item.id, () => {
        this.init()
      })
    },
    init () {
      this.channel = ''
      notificationService.slack((data) => {
        this.form.url = data.url
        this.form.template = data.template
        this.channels = data.channels
      })
    }
  }
}
</script>
