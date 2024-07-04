<script setup>
import appPage from '../../layout/page.vue'
</script>
<template>
  <appPage>
    <el-form ref="form" :model="form" :rules="formRules" label-width="auto">
      <el-row>
        <el-col :span="8">
          <el-input v-model="form.id" type="hidden"></el-input>
          <el-form-item label="节点名" prop="alias">
            <el-input v-model="form.alias"></el-input>
          </el-form-item>
          <el-form-item label="主机名" prop="name">
            <el-input v-model="form.name"></el-input>
          </el-form-item>
          <el-form-item label="端口号" prop="port">
            <el-input v-model.number="form.port"></el-input>
          </el-form-item>
          <el-form-item label="备注" prop="remark">
            <el-input
              type="textarea"
              :rows="5"
              v-model="form.remark">
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submit()">保存</el-button>
            <el-button @click="cancel">取消</el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </appPage>
</template>

<script>
import hostService from '../../api/host'
export default {
  data: function () {
    return {
      form: {
        id: -1,
        name: '',
        port: 5921,
        alias: '',
        remark: ''
      },
      formRules: {
        name: [
          {required: true, message: '请输入主机名', trigger: 'blur'}
        ],
        port: [
          {required: true, message: '请输入端口', trigger: 'blur'},
          {type: 'number', message: '端口无效'}
        ],
        alias: [
          {required: true, message: '请输入节点名称', trigger: 'blur'}
        ]
      }
    }
  },
  created () {
    const id = this.$route.params.id
    if (!id) {
      return
    }
    hostService.detail(id, (data) => {
      if (!data) {
        this.$message.error('数据不存在')
        this.cancel()
        return
      }
      this.form.id = data.id
      this.form.name = data.name
      this.form.port = data.port
      this.form.alias = data.alias
      this.form.remark = data.remark
    })
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
      hostService.update(this.form, () => {
        this.$router.push('/host')
      })
    },
    cancel () {
      this.$router.push('/host')
    }
  }
}
</script>
