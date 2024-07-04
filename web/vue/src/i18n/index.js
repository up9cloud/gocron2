import { createI18n } from 'vue-i18n'
const i18n = createI18n({
    locale: 'zh-CN',
    fallbackLocale: 'zh-CN',
    messages: {
        'zh-CN': {
            routes: {
                '--': '首页',
                'task': '任务管理',
                'task--list': '定时任务',
                'task--create': '新增',
                'task--edit': '编辑',
                'task--dependency': '任务依赖',
                'task--log': '任务日志',
                'host': '任务节点',
                'host--list': '节点列表',
                'host--create': '新增',
                'host--edit': '编辑',
                'user': '用户管理',
                'user--list': '用户列表',
                'user--create': '新增',
                'user--edit': '编辑',
                'user--edit-my-password': '修改密码',
                'user--edit-password': '修改密码',
                'system': '系统管理',
                'system--login-log': '登录日志',
                'system--notification': '通知配置'
            },
            task: {
                level: {
                    tooltip: `强依赖: 主任务执行成功，才会运行子任务
弱依赖: 无论主任务执行是否成功，都会运行子任务`
                }
            }
        },
    }
})
export default i18n
